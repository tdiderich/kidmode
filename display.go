package main

import (
	"fmt"
	"math/rand"
	"os"

	"golang.org/x/term"
)

const esc = "\x1b["

func clear() string {
	return esc + "2J" + esc + "H"
}

func moveTo(row, col int) string {
	return fmt.Sprintf("%s%d;%dH", esc, row, col)
}

func hideCursor() string {
	return esc + "?25l"
}

func showCursor() string {
	return esc + "?25h"
}

func resetStyle() string {
	return esc + "0m"
}

func fg256(n int) string {
	return fmt.Sprintf("%s38;5;%dm", esc, n)
}

func bg256(n int) string {
	return fmt.Sprintf("%s48;5;%dm", esc, n)
}

func bold() string {
	return esc + "1m"
}

func randomVividFg() string {
	n := 196 + rand.Intn(36)
	return fg256(n)
}

func randomVividBg() string {
	n := 196 + rand.Intn(36)
	return bg256(n)
}

func screenSize() (rows, cols int) {
	cols, rows, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 24, 80
	}
	return rows, cols
}

func write(s string) {
	os.Stdout.WriteString(s)
}
