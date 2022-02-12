//go:build !linux

package main

import (
	"io"
	"os"
)

func relayTTY(dst io.Writer, tty *os.File) {
	for {
		_, err := io.Copy(dst, tty)
		if err == nil {
			break
		}
	}
}
