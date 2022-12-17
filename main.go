package main

import (
	"fmt"
	"os"
	word_sets "turtle-typing/word_sets"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var (
	_words, _ = word_sets.ShuffleWordSet(word_sets.EnPop200)
)

type model struct {
	Words            []rune // words that appear on the screen
	redHighlight     []int  // word indexes that are highlighted with red
	greenHighlight   []int  // word indexes that are highlighted with green
	keyStrokeCounter int    // how many keystorkes player made
	Typed            []rune // words player has typed
	Mistakes         int    // number of mistakes player made
	Score            int    // player score
}

func initialModel() model {
	return model{
		Words:            _words,
		redHighlight:     []int{},
		greenHighlight:   []int{},
		keyStrokeCounter: -1,
		Typed:            []rune{},
		Mistakes:         0,
		Score:            0,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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

		if ch == next {
			m.Score += 1
		} else {
			m.Mistakes += 1
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	header := "Turtle Typing\t\t"
	var typed string
	remaining := m.Words[len(m.Typed):]

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

	// The footer
	footer := "\nPress ctrl+c to quit.\n"

	s := fmt.Sprintf(
		"\n  %s\n\n%s%s\n\nscore=%s,mistakes=%s\n\n%s",
		header,
		typed,
		string(remaining),
		fmt.Sprint(m.Score),
		fmt.Sprint(m.Mistakes),
		footer,
	)

	// Send the UI for rendering
	return s
}

func main() {
	fmt.Println(string(_words[:10]))
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
