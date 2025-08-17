package app

import (
	"math/rand"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func (a *App) generateText(length int) string {
	// Grab a random starting word
	keys := make([]string, 0, len(a.Words))
	for key := range a.Words {
		keys = append(keys, key)
	}
	idx := int(rand.Float64() * float64(len(keys)-1))
	lastWord := keys[idx]
	words := make([]string, 0, length)
	words = append(words, lastWord)
	for range length - 1 {
		word := a.getWord(lastWord)
		words = append(words, word)
		lastWord = word
	}

	return strings.Join(words, " ")
}

func (a *App) getWord(word string) string {
	if len(a.Words[word]) == 0 {
		panic("WTF")
	}
	idx := int(rand.Float64() * float64(len(a.Words[word])-1))
	idx = min(len(a.Words)-1, idx)
	return a.Words[word][idx]
}

func stringWidth(text []rune) int {
	var x int = 0
	for _, r := range text {
		x += runewidth.RuneWidth(r)
	}
	return x
}

func clearRow(row int) {
	w, _ := termbox.Size()
	for i := range w {
		termbox.SetCell(i, row, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
}
