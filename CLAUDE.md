# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
make build        # builds ./kidmode binary
./kidmode         # run (must be in a terminal, requires Accessibility permissions)
```

Exit by typing the password (default: `adulttakeover`, override with `KIDMODE_PASSWORD` env var).

## What This Is

A macOS terminal app (Go + CGO/Objective-C) that turns the keyboard into a toy for kids. Every keypress triggers random visual effects and sounds. The app locks down the system so kids can't escape to other apps.

## Architecture

All files are in a single `main` package. There are no tests.

**Input pipeline:** `main.go` → `input.go` (raw mode, password detection, control char filtering) → `effects.go` / `sound.go`

**Three-layer lockdown system (all macOS-specific via CGO):**
1. `keyblocker.go` — Disables macOS symbolic hotkeys (Mission Control, Exposé, etc.) via `NSUserDefaults` and restores them on exit. Writes `~/.kidmode-active` breadcrumb for crash recovery.
2. `keyblocker_cgo.go` — CoreGraphics event tap that intercepts and blocks system keyboard shortcuts (Cmd+Tab, Cmd+Q, etc.) and all mouse events at the session level.
3. `input.go` — Terminal-level filtering of Ctrl+C/D/Z/\ and escape sequences.

`keyblocker.go:startKeyBlocker()` orchestrates layers 1+2. Layer 3 is handled inline during `readInput()`.

**Display:** `display.go` (ANSI escape helpers, screen size), `effects.go` (6 weighted visual effects), `art.go` (ASCII art + block letter data).

## CGO Constraints

The CGO code links CoreGraphics, CoreFoundation, Foundation, and ApplicationServices. The event tap callback (`goEventTapCallback`) is exported to C and called from a dedicated OS thread — it must not block or allocate heavily. Returning `0` (nil) from the callback suppresses an event; returning the event passes it through.
