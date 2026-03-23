package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <Foundation/Foundation.h>
#include <stdlib.h>

// Read a symbolic hotkey's "enabled" state. Returns 1 if enabled, 0 if disabled, -1 if not found.
static int getHotkeyEnabled(int keyID) {
    NSUserDefaults *prefs = [[NSUserDefaults alloc] initWithSuiteName:@"com.apple.symbolichotkeys"];
    NSDictionary *hotkeys = [prefs objectForKey:@"AppleSymbolicHotKeys"];
    if (!hotkeys) return -1;

    NSString *key = [NSString stringWithFormat:@"%d", keyID];
    NSDictionary *entry = hotkeys[key];
    if (!entry) return -1;

    NSNumber *enabled = entry[@"enabled"];
    if (!enabled) return -1;
    return [enabled boolValue] ? 1 : 0;
}

// Set a symbolic hotkey's "enabled" state.
static void setHotkeyEnabled(int keyID, int enabled) {
    NSUserDefaults *prefs = [[NSUserDefaults alloc] initWithSuiteName:@"com.apple.symbolichotkeys"];
    NSMutableDictionary *hotkeys = [[prefs objectForKey:@"AppleSymbolicHotKeys"] mutableCopy];
    if (!hotkeys) {
        hotkeys = [NSMutableDictionary dictionary];
    }

    NSString *key = [NSString stringWithFormat:@"%d", keyID];
    NSMutableDictionary *entry = [hotkeys[key] mutableCopy];
    if (!entry) {
        entry = [NSMutableDictionary dictionary];
    }
    entry[@"enabled"] = @(enabled ? YES : NO);
    hotkeys[key] = entry;

    [prefs setObject:hotkeys forKey:@"AppleSymbolicHotKeys"];
    [prefs synchronize];
}
*/
import "C"

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Symbolic hotkey IDs for Mission Control family
var symbolicHotKeyIDs = []int{32, 34, 36, 62, 79, 80, 160}

// savedStates maps hotkey ID -> was it enabled before we disabled it
var savedStates = map[int]bool{}

var breadcrumbPath string

func init() {
	home, _ := os.UserHomeDir()
	breadcrumbPath = filepath.Join(home, ".kidmode-active")
}

func disableSymbolicHotkeys() {
	for _, id := range symbolicHotKeyIDs {
		state := int(C.getHotkeyEnabled(C.int(id)))
		if state == -1 {
			// Not found — assume was enabled
			savedStates[id] = true
		} else {
			savedStates[id] = state == 1
		}
		C.setHotkeyEnabled(C.int(id), 0)
	}
	activateSettings()

	// Write breadcrumb
	os.WriteFile(breadcrumbPath, []byte("active"), 0644)
}

func restoreSymbolicHotkeys() {
	for id, wasEnabled := range savedStates {
		if wasEnabled {
			C.setHotkeyEnabled(C.int(id), 1)
		}
	}
	activateSettings()

	// Remove breadcrumb
	os.Remove(breadcrumbPath)
}

func activateSettings() {
	cmd := exec.Command(
		"/System/Library/PrivateFrameworks/SystemAdministration.framework/Resources/activateSettings",
		"-u",
	)
	cmd.Run()
}

func startKeyBlocker() {
	disableSymbolicHotkeys()
	go startEventTap()
}

func stopKeyBlocker() {
	restoreSymbolicHotkeys()
}
