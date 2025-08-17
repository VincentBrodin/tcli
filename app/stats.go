package app

import (
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type Stat struct {
	Title  string
	Value  float64
	Prefix string
}

const (
	STATS_COUNT int = 3
)

func (a *App) renderStats() {
	a.clearStats()

	const AVG_WORD_LEN float64 = 4.7
	acc := float64(a.validHits) / float64(a.totalHits) * 100

	dur := a.done.Sub(a.start)
	chars := float64(len(a.text) - a.Length - 1) // a.Lenght is the amount of words, so the amount of words -1 = the amount of spaces
	wpm := (chars / AVG_WORD_LEN) / dur.Minutes()

	stats := [STATS_COUNT]Stat{
		{
			Title:  "wpm",
			Value:  wpm,
			Prefix: "",
		},
		{
			Title:  "accuracy",
			Value:  acc,
			Prefix: "%",
		},
		{
			Title:  "time",
			Value:  dur.Seconds(),
			Prefix: "s",
		},
	}

	offset := len(stats) / 2
	for i, stat := range stats {
		x := a.width/2 - stringWidth([]rune(fmt.Sprintf("%s: %.1f%s", stat.Title, stat.Value, stat.Prefix)))/2
		y := a.height/2 + (i - offset)
		for _, r := range fmt.Sprintf("%s: ", stat.Title) {
			termbox.SetCell(x, y, r, termbox.ColorDefault|termbox.AttrDim, termbox.ColorDefault)
			x += runewidth.RuneWidth(r)
		}
		for _, r := range fmt.Sprintf("%.1f", stat.Value) {
			termbox.SetCell(x, y, r, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault)
			x += runewidth.RuneWidth(r)
		}
		for _, r := range fmt.Sprintf("%s", stat.Prefix) {
			termbox.SetCell(x, y, r, termbox.ColorDefault|termbox.AttrDim, termbox.ColorDefault)
			x += runewidth.RuneWidth(r)
		}
	}
}

func (a *App) clearStats() {
	offset := STATS_COUNT / 2
	for i := range STATS_COUNT {
		y := a.height/2 + (i - offset)
		clearRow(y)
	}
}
