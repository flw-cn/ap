package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"

	"github.com/creack/pty"
	"golang.org/x/sys/unix"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: %v <command> <args>\n", os.Args[0])
		os.Exit(1)
	}

	c := exec.Command(os.Args[1], os.Args[2:]...)

	isTTY := true

	files := []*os.File{os.Stdin, os.Stdout, os.Stderr}

	for _, file := range files {
		if _, err := unix.IoctlGetWinsize(int(file.Fd()), unix.TIOCGWINSZ); err != nil {
			isTTY = false
			break
		}
	}

	if !isTTY {
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
		c.Wait()
		os.Exit(c.ProcessState.ExitCode())
	}

	p, err := pty.Start(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't exec %v: %v\n", os.Args[1:], err)
		os.Exit(1)
	}

	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt)
		<-signalCh
		c.Process.Signal(os.Interrupt)
	}()

	buf, _ := ioutil.ReadAll(io.TeeReader(p, os.Stderr))
	c.Wait()

	if bytes.Count(buf, []byte("\n")) > 30 {
		paging(bytes.NewBuffer(buf))
		os.Exit(c.ProcessState.ExitCode())
	}

	os.Exit(0)
}

func paging(file io.Reader) {
	c := exec.Command("/usr/bin/less", "-FR")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = file

	c.Run()
	c.Wait()
}
