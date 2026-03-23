package main

import (
	"math/rand"
	"os/exec"
	"sync"
	"time"
)

var systemSounds = []string{
	"Basso", "Blow", "Bottle", "Frog", "Funk", "Glass",
	"Hero", "Morse", "Ping", "Pop", "Purr", "Sosumi",
	"Submarine", "Tink",
}

var (
	lastPlayTime time.Time
	soundMu      sync.Mutex
	throttleMS   = 150 * time.Millisecond
)

func playRandomSound() {
	soundMu.Lock()
	now := time.Now()
	if now.Sub(lastPlayTime) < throttleMS {
		soundMu.Unlock()
		return
	}
	lastPlayTime = now
	soundMu.Unlock()

	sound := systemSounds[rand.Intn(len(systemSounds))]
	path := "/System/Library/Sounds/" + sound + ".aiff"

	cmd := exec.Command("afplay", path)
	if err := cmd.Start(); err != nil {
		return
	}
	// Reap in background to avoid zombies
	go cmd.Wait()
}
