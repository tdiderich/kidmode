package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
)

const (
	clearInterval = 40  // clear screen every N keypresses
	soundChance   = 0.6 // 60% chance of sound per keypress
)

var version = "dev"

func main() {
	// TTY check
	if !isTTY() {
		os.Stderr.WriteString("kidmode: must be run in a terminal (TTY required)\n")
		os.Exit(1)
	}

	// Check for unclean exit from last run
	if _, err := os.Stat(breadcrumbPath); err == nil {
		restoreSymbolicHotkeys()
	}

	initInput()

	// Cleanup function
	cleanup := func() {
		write(resetStyle() + showCursor() + clear())
		restoreTerminal()
		stopKeyBlocker()
	}

	// Signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-sigCh
		cleanup()
		os.Exit(0)
	}()

	// Start key blocker (symbolic hotkeys + event tap)
	startKeyBlocker()

	// Set up raw terminal mode
	setupRawMode()

	// Initialize screen
	write(hideCursor() + clear())

	// Welcome splash
	rows, cols := screenSize()
	msg := "*** KIDMODE *** Bang on the keyboard!"
	msgRow := rows / 2
	msgCol := max(1, (cols-len(msg))/2)
	write(moveTo(msgRow, msgCol) + randomVividFg() + bold() + msg + resetStyle())

	// Main input loop
	keypressCount := 0
	readInput(func(key byte) {
		keypressCount++

		// Periodic screen clear
		if keypressCount%clearInterval == 0 {
			write(clear())
		}

		// Visual effect
		runRandomEffect(key)

		// Sound (60% chance)
		if rand.Float64() < soundChance {
			playRandomSound()
		}
	})

	// Exit password was typed
	cleanup()
	fmt.Println("Bye! Kidmode exited.")
}

func isTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
