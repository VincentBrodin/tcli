package app

import (
	"time"

	"github.com/nsf/termbox-go"
)

type State int

const (
	StateIdle State = iota
	StateRunning
	StateDone
)

type App struct {
	state State

	start time.Time
	done  time.Time

	totalHits int
	validHits int

	text  []rune
	input []rune

	tools []Tool

	width  int
	height int

	Length int
	Words  map[string][]string
}

type Tool struct {
	Key         string
	Description string
}

func (a *App) Run() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	defer termbox.Close()

	a.width, a.height = termbox.Size()
	a.state = StateIdle

	a.tools = []Tool{
		{Key: "ctrl+q", Description: "quit"},
		{Key: "ctrl+r", Description: "restart"},
		{Key: "ctrl+n", Description: "new"},
	}

	a.text = []rune(a.generateText(a.Length))
	a.input = make([]rune, 0, len(a.text))

	// Event loop
	for {
		a.render()
		if err := termbox.Flush(); err != nil {
			return err
		}

		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventKey:
			if e.Ch == rune(0) {
				if isExitEvent(e.Key) {
					return nil
				}
				a.handleKey(e.Key)
			} else {
				a.handleRune(e.Ch)
			}
		case termbox.EventResize:
			if err := a.handleResize(e.Width, e.Height); err != nil {
				return err
			}
		case termbox.EventInterrupt:
		case termbox.EventMouse:
		case termbox.EventError:
			return e.Err
		case termbox.EventNone:
		}
		if len(a.text) == len(a.input) && a.state != StateDone {
			a.state = StateDone
			a.done = time.Now()
			a.clearTypingTest()
		}
	}
}

func (a *App) tryStart() {
	if a.state == StateIdle {
		a.start = time.Now()
		a.state = StateRunning
		a.totalHits = 0
		a.validHits = 0
	}
}

func isExitEvent(key termbox.Key) bool {
	return key == termbox.KeyCtrlQ
}

func (a *App) handleKey(key termbox.Key) {
	switch key {
	case termbox.KeySpace:
		a.tryStart()
		a.input = append(a.input, ' ')
		if a.state != StateDone {
			a.totalHits++
			i := len(a.input) - 1
			if a.text[i] == a.input[i] {
				a.validHits++
			}
		}
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		a.tryStart()
		if len(a.input) != 0 {
			a.input = a.input[:len(a.input)-1]
		}
	case termbox.KeyCtrlR:
		a.clearStats()
		a.state = StateIdle
		a.input = make([]rune, 0, len(a.text))
	case termbox.KeyCtrlN:
		a.clearStats()
		a.state = StateIdle
		a.text = []rune(a.generateText(a.Length))
		a.input = make([]rune, 0, len(a.text))
	}
}

func (a *App) handleRune(r rune) {
	a.tryStart()
	a.input = append(a.input, r)
	if a.state != StateDone {
		a.totalHits++
		i := len(a.input) - 1
		if a.text[i] == a.input[i] {
			a.validHits++
		}
	}
}

func (a *App) handleResize(width, height int) error {
	a.clearTypingTest()
	a.clearToolbar()
	a.clearStats()

	a.width = width
	a.height = height
	a.render()
	return termbox.Flush()
}

func (a *App) render() {
	switch a.state {

	case StateIdle, StateRunning:
		a.renderTypingTest()
	case StateDone:
		a.renderStats()
	}
	a.renderToolbar()
}
