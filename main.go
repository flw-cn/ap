package main

import (
	"bytes"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/creack/pty"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

//go:embed ap.bash
var bashScript string

//go:embed ap.fish
var fishScript string

//go:embed ap.zsh
var zshScript string

var optPager string
var optHeight int

func main() {
	args := parseOptions()
	if len(args) == 0 {
		os.Exit(1)
	}

	name, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't exec %v: %v\n", args, err)
		os.Exit(1)
	}

	var winSize *pty.Winsize
	piped := false

	cmd := exec.Command(name, args[1:]...)

	if _, err := pty.GetsizeFull(os.Stdin); err != nil {
		cmd.Stdin = os.Stdin
	}

	if size, err := pty.GetsizeFull(os.Stdout); err == nil {
		winSize = size
	} else {
		cmd.Stdout = os.Stdout
		piped = piped || isPipe(os.Stdout)
	}

	if size, err := pty.GetsizeFull(os.Stderr); err == nil {
		winSize = size
	} else {
		cmd.Stderr = os.Stderr
		piped = piped || isPipe(os.Stderr)
	}

	// ap should only work under tty, otherwise fall back to doing nothing.
	if winSize == nil || piped {
		err = syscall.Exec(name, args, os.Environ())
		fmt.Fprintf(os.Stderr, "Can't exec %v: %v\n", args, err)
		os.Exit(1)
	}

	runner := &Runner{
		cmd:     cmd,
		winSize: winSize,
	}

	exitCode := runner.Run()

	if optHeight == 0 {
		optHeight = int(runner.winSize.Rows) * 80 / 100
	} else if optHeight < 0 {
		optHeight = int(runner.winSize.Rows) * -optHeight / 100
	}

	if strings.Count(runner.output.String(), "\n") > optHeight {
		paging(&runner.output, runner.tty)
	}

	os.Exit(exitCode)
}

type Runner struct {
	cmd      *exec.Cmd    // the command to be run
	tty      *os.File     // local TTY device file
	ttyState *term.State  // the old state of local TTY
	pty      *os.File     // PTY master for run cmd
	output   ScreenBuffer // the command TTY output
	winSize  *pty.Winsize // the window size of local TTY & PTY master
	quit     bool         // indicates whether the child process has exited
}

func (r *Runner) Run() int {
	var err error

	r.tty, err = os.OpenFile("/dev/tty", os.O_RDWR, 0644)
	if err != nil {
		r.tty = os.Stdin
	}

	r.ttyState, err = term.MakeRaw(int(r.tty.Fd()))
	if err == nil {
		defer term.Restore(int(r.tty.Fd()), r.ttyState)
	}

	err = r.StartProcess()
	if err != nil {
		fmt.Fprintf(r.tty, "Can't exec %v: %v\r\n", r.cmd.Args, err)
		return 1
	}

	go r.relaySignal()
	go r.relayInput()
	r.relayOutput()
	r.cmd.Wait()

	return r.cmd.ProcessState.ExitCode()
}

func (r *Runner) StartProcess() (err error) {
	var tty *os.File

	r.pty, tty, err = pty.Open()
	if err != nil {
		return err
	}
	defer tty.Close()

	pty.Setsize(r.pty, r.winSize)

	if r.cmd.Stdout == nil {
		r.cmd.Stdout = tty
	}
	if r.cmd.Stderr == nil {
		r.cmd.Stderr = tty
	}
	if r.cmd.Stdin == nil {
		r.cmd.Stdin = tty
	}

	// NOTE: the index of `tty' here is 0
	r.cmd.ExtraFiles = []*os.File{tty}
	r.cmd.SysProcAttr = &syscall.SysProcAttr{
		// Setsid lets the child process to create a new session
		Setsid: true,
		// Setctty & Ctty lets child process connects to a controlling terminal
		Setctty: true,
		// NOTE: Golang requires us to predict the TTY file descriptor in the
		// child processes.
		// `3' is reserved for stdio, `0' is the index of `tty' in `ExtraFiles'
		Ctty: 3 + 0,
	}

	if err = r.cmd.Start(); err != nil {
		r.pty.Close()
		return err
	}

	return err
}

func (r *Runner) relaySignal() {
	signalCh := make(chan os.Signal, 10)
	signal.Notify(signalCh,
		syscall.SIGCHLD,
		syscall.SIGWINCH,
		os.Interrupt,
	)

LOOP:
	for sig := range signalCh {
		switch sig {
		case os.Interrupt:
			r.cmd.Process.Signal(os.Interrupt)
		case syscall.SIGWINCH:
			var err error
			r.winSize, err = pty.GetsizeFull(os.Stdout)
			if err == nil {
				pty.Setsize(r.pty, r.winSize)
			}
		case syscall.SIGCHLD:
			if r.cmd.Process.Signal(syscall.SIGCONT) != nil {
				r.quit = true
				break LOOP
			}
		}
	}
}

func (r *Runner) relayInput() {
	// waiting for the child process to start
	time.Sleep(20 * time.Millisecond)
	fd := int(r.tty.Fd())
	for !r.quit {
		syscall.SetNonblock(fd, true)
		io.Copy(r.pty, r.tty)
		syscall.SetNonblock(fd, false)
		time.Sleep(20 * time.Millisecond)
	}
}

func (r *Runner) relayOutput() {
	var perr *fs.PathError

	for !r.quit {
		_, err := io.Copy(io.MultiWriter(&r.output, r.tty), r.pty)
		if err == nil {
			break
		} else if errors.As(err, &perr) && perr.Err == syscall.EIO {
			break
		} else {
			time.Sleep(1 * time.Millisecond)
		}
	}
}

type ScreenBuffer struct {
	buf       bytes.Buffer
	altScreen bool // is currently writing to the alternate screen?
}

func (b *ScreenBuffer) Write(p []byte) (int, error) {
	n := len(p)

	if !b.altScreen {
		flag := []byte("\x1b[?1049h")
		if i := bytes.Index(p, flag); i > -1 {
			p = p[0:i]
			b.altScreen = true
		}
	} else {
		flag := []byte("\x1b[?1049l")
		if i := bytes.Index(p, flag); i > -1 {
			p = p[i+len(flag):]
			b.altScreen = false
		} else {
			p = nil
		}
	}

	b.buf.Write(p)

	return n, nil
}

func (b *ScreenBuffer) Read(p []byte) (int, error) {
	return b.buf.Read(p)
}

func (b *ScreenBuffer) String() string {
	return b.buf.String()
}

func (b *ScreenBuffer) Bytes() []byte {
	return b.buf.Bytes()
}

func paging(input io.Reader, tty *os.File) {
	args := strings.Fields(optPager)
	c := exec.Command(args[0], args[1:]...)
	c.Stdout = tty
	c.Stderr = tty
	c.Stdin = input

	c.Run()
}

func isPipe(file *os.File) bool {
	stat := &unix.Stat_t{}
	unix.Fstat(int(file.Fd()), stat)
	return stat.Mode&unix.S_IFIFO != 0
}

func printVersion(w io.Writer) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Fprintln(w, "Can't get build info.")
		return
	}

	version := info.Main.Version
	fmt.Fprintf(w, "%v version %s, built with %v\n",
		filepath.Base(info.Path), version, info.GoVersion)

	vcs := "unknown"
	vcsRev := "unknown"
	vcsTime := "unknown"

	typ, tag, rev, t := ParseVersion(version)
	switch typ {
	case Release, PreRelease:
		// info.Settings can't contains any valid VCS information. just return
		return
	case ErrorVersion:
		tag = "unknown branch"
	case Devel:
		tag = "clean working copy"
	case PseudoBaseNoTag, PseudoBaseRelease, PseudoBasePreRelease:
		if typ == PseudoBaseNoTag {
			tag = "untagged branch"
		} else {
			tag = "branch base on tag " + tag
		}
		vcsRev = rev
		vcsTime = t.Local().Format("2006-01-02 15:04:05 MST")
	}

	for _, s := range info.Settings {
		switch s.Key {
		case "vcs":
			vcs = s.Value
		case "vcs.revision":
			vcsRev = s.Value
		case "vcs.time":
			t, e := time.Parse(time.RFC3339, s.Value)
			if e == nil {
				vcsTime = t.Local().Format("2006-01-02 15:04:05 MST")
			}
		case "vcs.modified":
			if s.Value == "true" {
				tag = "dirty working copy"
			}
		}
	}

	fmt.Fprintf(w, `WARNING! This is not a release version, it's built from a %s.

VCS information:
VCS:         %v
Module path: %v
Commit time: %v
Revision id: %v

Please visit %v to get updates.
`,
		tag, vcs, info.Main.Path, vcsTime, vcsRev, info.Main.Path,
	)
}

type VersionType int

const (
	Devel VersionType = iota
	Release
	PreRelease
	PseudoBaseNoTag
	PseudoBaseRelease
	PseudoBasePreRelease
	ErrorVersion
)

// ParseVersion parses Go Module Version string to three parts:
// vcs tag, vcs revision ID, and commit time
//
// A Go Module Version string layout is one of follow formats:
//	* dirty vcs work directory: (devel)
//	* release version: vX.Y.Z
//	* pre-release version: v1.2.3-RC1
//	* pseudo version:
//		- untagged branch: v0.0.0-YYYYmmddHHMMSS-aabbccddeeff
//		- base on release version: vX.Y.(Z+1)-0.YYYYmmddHHMMSS-aabbccddeeff
//		- base on pre-release version: vX.Y.Z-RC1.0.YYYYmmddHHMMSS-aabbccddeeff
//
// see also: https://go.dev/ref/mod#glossary
//
func ParseVersion(version string) (typ VersionType, tag, rev string, t time.Time) {
	parts := strings.Split(version, "-")
	tag = parts[0]
	n := len(parts)
	if n < 3 { // this is not a pseudo version
		if tag == "(devel)" {
			typ = Devel
		} else if strings.Contains(tag, "-") {
			typ = PreRelease
		} else {
			typ = Release
		}
		return
	}

	rev = parts[n-1]
	timeStr := parts[n-2]
	actualLen := len(timeStr)
	expectLen := len("YYYYmmddHHMMSS")
	if actualLen < expectLen {
		return ErrorVersion, "", "", t
	}

	t, err := time.Parse("20060102150405", timeStr[actualLen-expectLen:actualLen])
	if err != nil {
		return ErrorVersion, "", "", t
	}

	if actualLen == expectLen {
		return PseudoBaseNoTag, "", rev, t
	}

	if actualLen == expectLen+2 {
		parts := strings.Split(tag, ".")
		patch, _ := strconv.Atoi(parts[2])
		if patch > 0 {
			patch = patch - 1
		}
		tag = parts[0] + "." + parts[1] + "." + strconv.Itoa(patch)
		return PseudoBaseRelease, tag, rev, t
	}

	tagLen := len(version) - len(".0.yyyymmddhhmmss-aabbccddeeff")
	tag = version[0:tagLen]
	return PseudoBasePreRelease, tag, rev, t
}

func parseOptions() []string {
	var (
		bash bool
		fish bool
		zsh  bool
		ver  bool
	)

	flag.StringVar(&optPager, "pager", "", "what pager to be used, defaults to `less -Fr'")
	flag.IntVar(&optHeight, "height", -80, "enable paging when the number of lines exceeds this height. negative numbers means percentages. defaults to -80(means 80%)")
	flag.BoolVar(&bash, "bash", false, "output bash script")
	flag.BoolVar(&fish, "fish", false, "output fish script")
	flag.BoolVar(&zsh, "zsh", false, "output zsh script")
	flag.BoolVar(&ver, "version", false, "print version information")
	flag.Parse()

	if bash {
		fmt.Println(bashScript)
		return nil
	}

	if fish {
		fmt.Println(fishScript)
		return nil
	}

	if zsh {
		fmt.Println(zshScript)
		return nil
	}

	if optPager == "" {
		if s := os.Getenv("AP_PAGER"); s != "" {
			optPager = s
		} else if s := os.Getenv("PAGER"); s != "" {
			optPager = s
		} else {
			optPager = "less -Fr"
		}
	}

	if ver {
		printVersion(os.Stderr)
		return nil
	}

	args := flag.Args()
	if len(args) == 0 {
		usage := `
Usage: %v [<option> [<option args>]] -- <command> [<args>]
       %v --help for more information
`
		fmt.Fprintf(os.Stderr, usage[1:], os.Args[0], os.Args[0])
	}

	return args
}
