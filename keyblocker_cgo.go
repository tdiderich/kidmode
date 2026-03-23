package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework Foundation -framework ApplicationServices

#include <CoreGraphics/CoreGraphics.h>
#include <CoreFoundation/CoreFoundation.h>

// Forward declaration of the Go callback
extern CGEventRef goEventTapCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *userInfo);

// C shim that CGEventTapCreate can use as a function pointer
static CGEventRef eventTapShim(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *userInfo) {
    return goEventTapCallback(proxy, type, event, userInfo);
}

// Wrapper to create the event tap using our shim
static CFMachPortRef createEventTap(CGEventMask mask) {
    return CGEventTapCreate(
        kCGHIDEventTap,
        kCGHeadInsertEventTap,
        kCGEventTapOptionDefault,
        mask,
        eventTapShim,
        NULL
    );
}

// Wrapper to run the current CFRunLoop
static void runCFRunLoop(void) {
    CFRunLoopRun();
}

// Wrapper to add event tap to run loop
static void addTapToRunLoop(CFMachPortRef tap) {
    CFRunLoopSourceRef source = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, tap, 0);
    CFRunLoopAddSource(CFRunLoopGetCurrent(), source, kCFRunLoopCommonModes);
    CGEventTapEnable(tap, true);
    CFRelease(source);
}

// Re-enable a disabled tap
static void reenableTap(CFMachPortRef tap) {
    CGEventTapEnable(tap, true);
}
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// Shortcut defines a key combination to block
type Shortcut struct {
	keyCode      uint16
	requiredFlag uint64
	anyModifiers bool
}

var shortcutsToBlock = []Shortcut{
	// Cmd+Space (Spotlight)
	{49, C.kCGEventFlagMaskCommand, false},
	// Cmd+Tab (App Switcher)
	{48, C.kCGEventFlagMaskCommand, false},
	// Cmd+Q (Quit)
	{12, C.kCGEventFlagMaskCommand, false},
	// Cmd+W (Close window)
	{13, C.kCGEventFlagMaskCommand, false},
	// Cmd+H (Hide)
	{4, C.kCGEventFlagMaskCommand, false},
	// Cmd+M (Minimize)
	{41, C.kCGEventFlagMaskCommand, false},
	// F3 / Mission Control
	{99, 0, true},
	// F4 / Launchpad
	{118, 0, true},
	// Mission Control key (newer Apple keyboards)
	{160, 0, true},
	// Launchpad key (newer Apple keyboards)
	{131, 0, true},
	// Control+Up (Mission Control)
	{126, C.kCGEventFlagMaskControl, false},
	// Control+Down (App Exposé)
	{125, C.kCGEventFlagMaskControl, false},
	// Control+Left (Switch Spaces)
	{123, C.kCGEventFlagMaskControl, false},
	// Control+Right (Switch Spaces)
	{124, C.kCGEventFlagMaskControl, false},
	// Globe/fn key
	{179, 0, true},
	// Cmd+Shift+3 (Screenshot full screen)
	{20, C.kCGEventFlagMaskCommand | 0x00020000, false},
	// Cmd+Shift+4 (Screenshot selection)
	{21, C.kCGEventFlagMaskCommand | 0x00020000, false},
	// Cmd+Shift+5 (Screenshot/recording panel)
	{23, C.kCGEventFlagMaskCommand | 0x00020000, false},
}

var eventTapRef C.CFMachPortRef

func shouldBlockKeyEvent(eventType C.CGEventType, event C.CGEventRef) bool {
	keyCode := uint16(C.CGEventGetIntegerValueField(event, C.kCGKeyboardEventKeycode))
	flags := uint64(C.CGEventGetFlags(event))

	for _, s := range shortcutsToBlock {
		if keyCode == s.keyCode {
			if s.anyModifiers {
				return true
			}
			if flags&s.requiredFlag == s.requiredFlag {
				return true
			}
		}
	}
	return false
}

//export goEventTapCallback
func goEventTapCallback(proxy C.CGEventTapProxy, eventType C.CGEventType, event C.CGEventRef, userInfo unsafe.Pointer) C.CGEventRef {
	// Re-enable tap if macOS disabled it
	if eventType == C.kCGEventTapDisabledByTimeout || eventType == C.kCGEventTapDisabledByUserInput {
		if eventTapRef != 0 {
			C.reenableTap(eventTapRef)
		}
		return event
	}

	// Block all mouse events
	switch eventType {
	case C.kCGEventLeftMouseDown, C.kCGEventLeftMouseUp, C.kCGEventLeftMouseDragged,
		C.kCGEventRightMouseDown, C.kCGEventRightMouseUp, C.kCGEventRightMouseDragged,
		C.kCGEventMouseMoved, C.kCGEventScrollWheel:
		return 0
	case C.kCGEventOtherMouseDown, C.kCGEventOtherMouseUp, C.kCGEventOtherMouseDragged:
		return 0
	}

	// Block key events
	if eventType == C.kCGEventKeyDown || eventType == C.kCGEventKeyUp {
		if shouldBlockKeyEvent(eventType, event) {
			return 0 // nil = block the event
		}
	}

	// Block all system-defined events (media keys, volume, brightness, Mission Control, etc.)
	if uint32(eventType) == 14 {
		return 0
	}

	return event
}

// startEventTap runs on a dedicated OS thread and blocks indefinitely.
func startEventTap() {
	runtime.LockOSThread()

	// Event mask: keyDown | keyUp | systemDefined | all mouse events
	mask := C.CGEventMask(
		(1 << C.kCGEventKeyDown) | (1 << C.kCGEventKeyUp) | (1 << 14) |
			(1 << C.kCGEventLeftMouseDown) | (1 << C.kCGEventLeftMouseUp) | (1 << C.kCGEventLeftMouseDragged) |
			(1 << C.kCGEventRightMouseDown) | (1 << C.kCGEventRightMouseUp) | (1 << C.kCGEventRightMouseDragged) |
			(1 << C.kCGEventMouseMoved) | (1 << C.kCGEventScrollWheel) |
			(1 << C.kCGEventOtherMouseDown) | (1 << C.kCGEventOtherMouseUp) | (1 << C.kCGEventOtherMouseDragged),
	)

	tap := C.createEventTap(mask)
	if tap == 0 {
		write("Failed to create event tap. Grant Accessibility permissions to your terminal.\n")
		write("System Settings → Privacy & Security → Accessibility → enable your terminal app\n")
		return
	}

	eventTapRef = tap
	C.addTapToRunLoop(tap)
	C.runCFRunLoop() // blocks forever
}
