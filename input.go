package main

import (
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	defaultPassword = "adulttakeover"
	bufferSize      = 50
)

var (
	password  string
	charBuf   string
	oldState  *term.State
)

func initInput() {
	password = os.Getenv("KIDMODE_PASSWORD")
	if password == "" {
		password = defaultPassword
	}
}

func setupRawMode() {
	var err error
	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		os.Stderr.WriteString("kidmode: failed to set raw mode\n")
		os.Exit(1)
	}
}

func restoreTerminal() {
	if oldState != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
}

func checkPassword(ch byte) bool {
	charBuf += string(ch)
	if len(charBuf) > bufferSize {
		charBuf = charBuf[len(charBuf)-bufferSize:]
	}
	return strings.Contains(charBuf, password)
}

// readInput reads raw stdin and calls onKey for each keypress.
// Returns when the exit password is typed.
func readInput(onKey func(key byte)) {
	buf := make([]byte, 256)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return
		}
		for i := 0; i < n; i++ {
			b := buf[i]

			// Swallow dangerous control chars
			switch b {
			case 0x03, 0x04, 0x1A, 0x1C: // Ctrl+C, Ctrl+D, Ctrl+Z, Ctrl+\
				onKey('*')
				continue
			}

			// Skip escape sequences (arrow keys, function keys, etc.)
			if b == 0x1B {
				if i+1 < n && buf[i+1] == 0x5B {
					i += 2 // skip ESC [
					for i < n && buf[i] < 0x40 {
						i++ // skip parameter bytes
					}
					// i now at final byte; loop will increment past it
				}
				onKey('*')
				continue
			}

			// Printable ASCII
			if b >= 32 && b <= 126 {
				if checkPassword(b) {
					return
				}
				onKey(b)
			} else {
				onKey('*')
			}
		}
	}
}
