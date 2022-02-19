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
	buf := make([]byte, 1024)
	fd := int(r.tty.Fd())
	fds := unix.FdSet{}
	fds.Set(fd)

	for !r.quit {
		timeout := unix.Timeval{Sec: 0, Usec: 30000}
		rs := fds
		n, _ := unix.Select(fd+1, &rs, nil, nil, &timeout)
		if n > 0 {
			if n, err := r.tty.Read(buf); err == nil && n > 0 {
				r.pty.Write(buf[0:n])
			}
		}
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

func printVersion() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Can't get build info.")
		return
	}

	version := info.Main.Version
	fmt.Printf("%v version %s\n", filepath.Base(info.Path), version)

	l := len(version)
	if l < len("vX.0.0-yyyymmddhhmmss-abcdefabcdef") {
		return
	}

	tag := version[0 : l-30]
	typ := version[l-30 : l-29]
	time, _ := time.Parse("20060102150405", version[l-27:l-13])
	commit := version[l-12 : l]
	timeStr := time.Local().Format("2006-01-02 15:04:05 MST")

	if version[0:7] == "v0.0.0-" {
		tag = "untagged branch"
	} else if typ == "-" {
		parts := strings.Split(tag, ".")
		patch, _ := strconv.Atoi(parts[2])
		if patch > 0 {
			patch = patch - 1
		}
		tag = parts[0] + "." + parts[1] + "." + strconv.Itoa(patch)
	}

	fmt.Printf("base on %s, commit at %s, commit ID is %s\n",
		tag, timeStr, commit)
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
		printVersion()
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
