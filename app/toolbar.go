package app

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func (a *App) renderToolbar() {
	a.clearToolbar()
	text := toolsToRunes(a.tools)
	msgWidth := stringWidth(text)
	start := a.width/2 - msgWidth/2
	x := start
	for _, r := range text {
		termbox.SetCell(x, a.height-1, r, termbox.ColorDefault|termbox.AttrDim, termbox.ColorDefault)
		x += runewidth.RuneWidth(r)
	}
}

func toolsToRunes(tools []Tool) []rune {
	sb := strings.Builder{}
	size := len(tools)
	i := 0
	for _, tool := range tools {
		sb.WriteString(tool.Key)
		sb.WriteString(":")
		sb.WriteString(tool.Description)
		if size-1 > i {
			sb.WriteString(" - ")
		}
		i++
	}

	return []rune(sb.String())
}

func (a *App) clearToolbar() {
	clearRow(a.height - 1)
}
