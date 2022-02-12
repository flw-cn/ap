//go:build linux

package main

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

// relayTTY ensures that the output of tty is copied in full to dst.
//
// The pty slave close(2) under linux does not cause the pty master to receive
// an EOF, so io.Copy() cannot use this to determine that the read is finished.
//
// The ioctl(2) system call is used here to get the current status of the tty.
func relayTTY(dst io.Writer, tty *os.File) {
	for {
		_, err := io.Copy(dst, tty)
		if err == nil {
			break
		} else if err != nil && !getSID(tty) {
			break
		}
	}
}

func getSID(tty *os.File) bool {
	var pid int
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, tty.Fd(), syscall.TIOCGSID, uintptr(unsafe.Pointer(&pid)))
	return err != syscall.ENOTTY
}
