package main

import (
	"fmt"
	"os"
	"time"
	word_sets "turtle-typing/word_sets"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var (
	_words, _ = word_sets.ShuffleWordSet(word_sets.EnPop200)
)

const (
	timeout         = time.Second * 60
	wordsWindowSize = word_sets.WordsPerLineLimit * 2
	commonWordSize  = 5 // To calculate WPM
)

type model struct {
	Words    []rune      // chars that appear on the screen
	Typed    []rune      // chars player has typed
	Mistakes int         // number of mistakes player made
	Score    int         // player score
	Timer    timer.Model // timer
	Quit     bool        // should quit game
}

func initialModel() model {
	return model{
		Words:    _words,
		Typed:    []rune{},
		Mistakes: 0,
		Score:    0,
		Timer:    timer.NewWithInterval(timeout, time.Second),
	}
}

func (m model) Init() tea.Cmd {
	return m.Timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.Quit = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case tea.KeyBackspace.String():
			if len(m.Typed) > 0 {
				m.Typed = m.Typed[:len(m.Typed)-1]
			}

		}

		if msg.Type != tea.KeyRunes && msg.Type != tea.KeySpace {
			return m, nil
		}

		ch := msg.Runes[0]
		next := rune(m.Words[len(m.Typed)])

		if next == '\n' {
			m.Typed = append(m.Typed, next)

			// Since we need to perform a line break
			// if the user types a space we should simply ignore it.
			if ch == ' ' {
				return m, nil
			}
		}

		m.Typed = append(m.Typed, ch)

		if len(m.Typed) == 1 {
			m.Timer.Timeout = timeout
			m.Timer.Toggle()
		}

		if ch == next {
			m.Score += 1
		} else {
			m.Mistakes += 1
		}
	}

	return m, nil
}

func (m model) View() string {
	header := "🐢 Turtle Typing "
	footer := "\nPress ctrl+c to quit.\n"

	var typed string
	remaining := m.Words[len(m.Typed):][:wordsWindowSize]

	for i, c := range m.Typed {
		if c == m.Words[i] {
			output := termenv.String(string(c))
			output = output.Foreground(termenv.RGBColor("#90EE90"))
			typed += output.String()
		} else {
			output := termenv.String(string(c))
			output = output.Foreground(termenv.RGBColor("#DC143C"))
			typed += output.String()
		}
	}

	timerView := m.Timer.View()

	if m.Timer.Timedout() {
		timerView = "Time's Up!"
		WPM := m.Score / commonWordSize
		var errorRate = 0
		if len(m.Typed) > 0 {
			errorRate = m.Mistakes * 100 / len(m.Typed)
		}

		s := fmt.Sprintf(
			"\n  %s\n\nWPM = %s, Error Rate = %s%%\n\n%s",
			header,
			fmt.Sprint(WPM),
			fmt.Sprint(errorRate),
			footer,
		)
		return s
	}
	timerView += "\n"
	if !m.Quit {
		timerView = "Remaining... " + timerView
	}

	s := fmt.Sprintf(
		"\n  %s\n\n%s%s\n\n%s\n\nscore=%s,mistakes=%s\n\n%s",
		header,
		typed,
		string(remaining),
		timerView,
		fmt.Sprint(m.Score),
		fmt.Sprint(m.Mistakes),
		footer,
	)

	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
