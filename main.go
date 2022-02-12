package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
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

	var tty *os.File
	var winSize *pty.Winsize
	piped := false

	cmd := exec.Command(name, args[1:]...)
	if _, err := pty.GetsizeFull(os.Stdin); err != nil {
		cmd.Stdin = os.Stdin
	}

	if size, err := pty.GetsizeFull(os.Stdout); err == nil {
		tty = os.Stdout
		winSize = size
	} else {
		cmd.Stdout = os.Stdout
		piped = piped || isPipe(os.Stdout)
	}

	if size, err := pty.GetsizeFull(os.Stderr); err == nil {
		tty = os.Stderr
		winSize = size
	} else {
		cmd.Stderr = os.Stderr
		piped = piped || isPipe(os.Stdout)
	}

	// ap should only work under tty, otherwise fall back to doing nothing.
	if tty == nil || piped {
		err = syscall.Exec(name, args, os.Environ())
		fmt.Fprintf(os.Stderr, "Can't exec %v: %v\n", args, err)
		os.Exit(1)
	}

	exitCode := run(cmd, tty, winSize)

	os.Exit(exitCode)
}

func run(cmd *exec.Cmd, tty *os.File, winSize *pty.Winsize) int {
	p, err := pty.StartWithAttrs(cmd, winSize, &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true,
		Ctty:    int(tty.Fd()),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't exec %v: %v\n", cmd.Args, err)
		return 1
	}

	state, err := term.MakeRaw(0)

	output := new(bytes.Buffer)

	go func() {
		checkAppType := time.Tick(1 * time.Second)
		copyStdin := time.Tick(50 * time.Millisecond)
		signalCh := make(chan os.Signal, 10)
		signal.Notify(signalCh, syscall.SIGCHLD)
		signal.Notify(signalCh, os.Interrupt)
		signal.Notify(signalCh, syscall.SIGWINCH)
		for {
			select {
			case sig := <-signalCh:
				switch sig {
				case os.Interrupt:
					cmd.Process.Signal(os.Interrupt)
					break
				case syscall.SIGWINCH:
					var err error
					winSize, err = pty.GetsizeFull(os.Stdout)
					if err == nil {
						pty.Setsize(p, winSize)
					}
				case syscall.SIGCHLD:
					return
				}

			case <-checkAppType:
				if bytes.Contains(output.Bytes(), []byte("\x1b[?1049h")) {
					syscall.SetNonblock(0, false)
					go relayTTY(p, os.Stdin)
					return
				}

			case <-copyStdin:
				syscall.SetNonblock(0, true)
				io.Copy(p, os.Stdin)
				syscall.SetNonblock(0, false)
			}
		}
	}()

	relayTTY(io.MultiWriter(output, tty), p)
	cmd.Wait()

	if state != nil {
		term.Restore(0, state)
	}

	if optHeight == 0 {
		optHeight = int(winSize.Rows) * 80 / 100
	} else if optHeight < 0 {
		optHeight = int(winSize.Rows) * -optHeight / 100
	}

	if bytes.Count(output.Bytes(), []byte("\n")) > optHeight &&
		!bytes.Contains(output.Bytes(), []byte("\x1b[?1049h")) {
		paging(output, tty)
	}

	return cmd.ProcessState.ExitCode()
}

func paging(input io.Reader, output io.Writer) {
	var pager string
	if optPager != "" {
		pager = optPager
	} else if s := os.Getenv("AP_PAGER"); s != "" {
		pager = s
	} else if s := os.Getenv("PAGER"); s != "" {
		pager = s
	} else {
		pager = "less -Fr"
	}

	args := strings.Fields(pager)
	c := exec.Command(args[0], args[1:]...)
	c.Stdout = output
	c.Stderr = output
	c.Stdin = input

	c.Run()
}

func isPipe(file *os.File) bool {
	stat := &unix.Stat_t{}
	unix.Fstat(int(file.Fd()), stat)
	return stat.Mode&unix.S_IFIFO == 1
}

func parseOptions() []string {
	var (
		bash bool
		fish bool
		zsh  bool
	)

	flag.StringVar(&optPager, "pager", "", "what pager to be used, defaults to `less -Fr'")
	flag.IntVar(&optHeight, "height", -80, "enable paging when the number of lines exceeds this height. negative numbers means percentages. defaults to -80(means 80%)")
	flag.BoolVar(&bash, "bash", false, "output bash script")
	flag.BoolVar(&fish, "fish", false, "output fish script")
	flag.BoolVar(&zsh, "zsh", false, "output zsh script")
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
