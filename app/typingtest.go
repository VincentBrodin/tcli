package app

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func (a *App) renderTypingTest() {
	a.clearTypingTest()
	rows := a.rows()

	offset := len(rows) / 2
	for idy, row := range rows {
		x := a.width/2 - stringWidth(row)/2
		y := a.height/2 + (idy - offset)
		for idx := range row {
			r, fg := a.getRuneAndAttr(posToI(idx, idy, rows))
			termbox.SetCell(x, y, r, fg, termbox.ColorDefault)
			x += runewidth.RuneWidth(r)
		}
	}
}

func (a *App) getRuneAndAttr(i int) (rune, termbox.Attribute) {
	var fg termbox.Attribute
	var r rune
	if i > len(a.input)-1 { // All the text left that the user has not written
		r = a.text[i]
		fg = termbox.ColorDefault | termbox.AttrDim
	} else if i > len(a.text)-1 { // All the input that overflows the text
		r = a.input[i]
		fg = termbox.ColorRed
	} else if a.input[i] == a.text[i] { // All the right characters
		r = a.input[i]
		fg = termbox.ColorDefault | termbox.AttrBold
	} else if a.text[i] == ' ' { // If a users missed a space
		r = '•'
		fg = termbox.ColorRed
	} else { // All the wrong characters
		r = a.text[i]
		fg = termbox.ColorRed
	}

	if i == len(a.input) { // Cursors
		fg |= termbox.AttrReverse
	}
	return r, fg
}

// Converts the string into rows if they are out of bounds
func (a *App) rows() [][]rune {
	rows := make([][]rune, 0, 16)
	rows = append(rows, make([]rune, 0, len(a.text)))
	width := int(float64(a.width) * 0.5) // Adds 25% margin to each side
	var x, y int

	words := strings.SplitSeq(string(a.text), " ")

	for word := range words {
		runes := []rune(word)
		rw := stringWidth(runes)
		if x+rw+1 > width {
			rows[y] = append(rows[y], []rune(" ")...)
			rows = append(rows, make([]rune, 0, len(a.text)))
			x = 0
			y++
		} else if x != 0 {
			rows[y] = append(rows[y], []rune(" ")...)
		}
		rows[y] = append(rows[y], runes...)
		x += rw
	}

	return rows
}

func posToI(x, y int, rows [][]rune) int {
	if y < 0 || y >= len(rows) || x < 0 || x >= len(rows[y]) {
		return -1
	}

	index := 0
	for i := 0; i < y; i++ {
		index += len(rows[i])
	}
	index += x
	return index
}

func (a *App) clearTypingTest() {
	for i := range a.height - 1 {
		clearRow(i)
	}
}
