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

//go:embed ap.zsh
var zshScript string

var optPager string
var optHeight int

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: %v [<option> [<option args>]] -- <command> <args>\n", os.Args[0])
		os.Exit(1)
	}

	var zsh bool

	flag.StringVar(&optPager, "pager", "", "what pager to be used, defaults to `less -FR'")
	flag.IntVar(&optHeight, "height", 0, "enable paging when the number of lines exceeds this height. negative numbers means percentages. defaults to -80(means 80%)")
	flag.BoolVar(&zsh, "zsh", false, "output zsh script")
	flag.Parse()

	if zsh {
		fmt.Println(zshScript)
		return
	}

	// ap should only work under tty, otherwise fall back to doing nothing.
	isTTY := true
	for fd := 0; fd < 3; fd++ {
		if _, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ); err != nil {
			isTTY = false
			break
		}
	}

	if !isTTY {
		args := flag.Args()
		name, err := exec.LookPath(args[0])
		if err != nil {
			os.Exit(1)
		}
		err = syscall.Exec(name, args, os.Environ())
		os.Exit(1)
	}

	run()
}

func run() {
	args := flag.Args()
	c := exec.Command(args[0], args[1:]...)

	p, err := pty.StartWithSize(c, getSize(os.Stdout))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't exec %v: %v\n", args, err)
		os.Exit(1)
	}

	state, _ := term.MakeRaw(0)

	go func() {
		signalCh := make(chan os.Signal, 10)
		signal.Notify(signalCh, os.Interrupt)
		signal.Notify(signalCh, syscall.SIGWINCH)

		for sig := range signalCh {
			switch sig {
			case os.Interrupt:
				c.Process.Signal(os.Interrupt)
				break
			case syscall.SIGWINCH:
				pty.Setsize(p, getSize(os.Stdout))
			}
		}
	}()

	output := new(bytes.Buffer)

	go func() {
		signalCh := make(chan os.Signal, 10)
		signal.Notify(signalCh, syscall.SIGCHLD)
		waitStdin := time.After(1 * time.Second)
		stdinReady := false
		for {
			select {
			case <-signalCh:
				return

			case <-waitStdin:
				if bytes.Contains(output.Bytes(), []byte("\x1b[?1049h")) {
					io.Copy(p, os.Stdin)
					return
				} else {
					stdinReady = true
				}

			default:
				if stdinReady {
					syscall.SetNonblock(0, true)
					os.Stdin.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
					io.Copy(p, os.Stdin)
					syscall.SetNonblock(0, false)
				}
			}
		}
	}()

	io.Copy(output, io.TeeReader(p, os.Stderr))
	c.Wait()

	winSize := getSize(os.Stdout)
	if optHeight == 0 {
		optHeight = int(winSize.Rows) * 80 / 100
	} else if optHeight < 0 {
		optHeight = int(winSize.Rows) * optHeight / 100
	}

	if bytes.Count(output.Bytes(), []byte("\n")) > optHeight &&
		!bytes.Contains(output.Bytes(), []byte("\x1b[?1049h")) {
		paging(output)
	}

	term.Restore(0, state)
	os.Exit(c.ProcessState.ExitCode())
}

func paging(output io.Reader) {
	var pager string
	if optPager != "" {
		pager = optPager
	} else if s := os.Getenv("AP_PAGER"); s != "" {
		pager = s
	} else if s := os.Getenv("PAGER"); s != "" {
		pager = s
	} else {
		pager = "less -FR"
	}

	args := strings.Fields(pager)

	c := exec.Command(args[0], args[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = output

	c.Run()
}

func getSize(p *os.File) *pty.Winsize {
	rows, cols, _ := pty.Getsize(p)
	return &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols), X: 0, Y: 0}
}
