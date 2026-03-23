package main

import (
	"math"
	"math/rand"
	"strings"
)

var effectChars = []string{
	"★", "✦", "●", "◆", "♦", "♥", "♠", "♣", "☀", "☁",
	"♪", "♫", "▲", "△", "▼", "◇", "○", "◎", "✿", "❀",
	"✻", "✼", "❋", "⚡", "⭐",
}

func randomInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + rand.Intn(max-min+1)
}

func pickRandom(arr []string) string {
	return arr[rand.Intn(len(arr))]
}

// charBurst: burst of random colorful characters
func charBurst() {
	rows, cols := screenSize()
	var out strings.Builder
	count := randomInt(5, 15)
	for i := 0; i < count; i++ {
		r := randomInt(1, rows)
		c := randomInt(1, cols-1)
		ch := pickRandom(effectChars)
		out.WriteString(moveTo(r, c) + randomVividFg() + bold() + ch + resetStyle())
	}
	write(out.String())
}

// bigLetter: display the pressed key as a big block letter
func bigLetter(key byte) {
	lines := getBlockLetter(key)
	if lines == nil {
		charBurst()
		return
	}
	rows, cols := screenSize()
	startRow := randomInt(1, max(1, rows-6))
	startCol := randomInt(1, max(1, cols-6))
	color := randomVividFg()
	var out strings.Builder
	out.WriteString(color + bold())
	for i, line := range lines {
		out.WriteString(moveTo(startRow+i, startCol) + line)
	}
	out.WriteString(resetStyle())
	write(out.String())
}

// asciiArt: display random ASCII art
func asciiArtEffect() {
	art := getRandomArt()
	artLines := strings.Split(art, "\n")
	rows, cols := screenSize()
	maxWidth := 0
	for _, l := range artLines {
		if len(l) > maxWidth {
			maxWidth = len(l)
		}
	}
	startRow := randomInt(1, max(1, rows-len(artLines)-1))
	startCol := randomInt(1, max(1, cols-maxWidth-1))
	color := randomVividFg()
	var out strings.Builder
	out.WriteString(color + bold())
	for i, line := range artLines {
		out.WriteString(moveTo(startRow+i, startCol) + line)
	}
	out.WriteString(resetStyle())
	write(out.String())
}

// colorFlood: flood a region with a random background color
func colorFlood() {
	rows, cols := screenSize()
	floodRows := randomInt(3, min(10, rows))
	floodCols := randomInt(10, min(30, cols))
	startRow := randomInt(1, max(1, rows-floodRows))
	startCol := randomInt(1, max(1, cols-floodCols))
	bg := randomVividBg()
	var out strings.Builder
	out.WriteString(bg)
	spaces := strings.Repeat(" ", floodCols)
	for i := 0; i < floodRows; i++ {
		out.WriteString(moveTo(startRow+i, startCol) + spaces)
	}
	out.WriteString(resetStyle())
	write(out.String())
}

// sparkles: sparkles scattered across the screen
func sparkles() {
	rows, cols := screenSize()
	sparkleChars := []string{"✦", "✧", "⋆", "✫", "✬", "✶", "✷", "✸", "✹", "✺", "·", "˚", "°"}
	var out strings.Builder
	count := randomInt(15, 40)
	for i := 0; i < count; i++ {
		r := randomInt(1, rows)
		c := randomInt(1, cols-1)
		out.WriteString(moveTo(r, c) + randomVividFg() + sparkleChars[rand.Intn(len(sparkleChars))])
	}
	out.WriteString(resetStyle())
	write(out.String())
}

// rainbow: rainbow horizontal stripe
func rainbow() {
	rows, cols := screenSize()
	row := randomInt(1, rows)
	rainbowColors := []int{196, 208, 226, 46, 21, 93, 201}
	var out strings.Builder
	for c := 1; c <= cols; c++ {
		colorIdx := int(math.Floor(float64(c-1) / float64(cols) * float64(len(rainbowColors))))
		if colorIdx >= len(rainbowColors) {
			colorIdx = len(rainbowColors) - 1
		}
		out.WriteString(moveTo(row, c) + fg256(rainbowColors[colorIdx]) + "█")
	}
	out.WriteString(resetStyle())
	write(out.String())
}

// runRandomEffect picks an effect with weighted probability
// 35% char burst, 20% big letter, 15% ASCII art, 12% color flood, 10% sparkles, 8% rainbow
func runRandomEffect(key byte) {
	roll := rand.Intn(100)
	switch {
	case roll < 35:
		charBurst()
	case roll < 55:
		bigLetter(key)
	case roll < 70:
		asciiArtEffect()
	case roll < 82:
		colorFlood()
	case roll < 92:
		sparkles()
	default:
		rainbow()
	}
}
