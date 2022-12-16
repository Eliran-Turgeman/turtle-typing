package main

import (
	"fmt"
	"os"
	"strings"
	word_sets "turtle-typing/word_sets"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	_words, _      = word_sets.ShuffleWordSet(word_sets.EnPop200)
	wordsToDisplay = strings.Join(_words, " ")

	Red   = "\033[31m"
	Green = "\033[32m"
)

type model struct {
	words          string // words that appear on the screen
	redHighlight   []int  // word indexes that are highlighted with red
	greenHighlight []int  // word indexes that are highlighted with green
	// cursor           int              // which to-do list item our cursor is pointing at
	// selected         map[int]struct{} // which to-do items are selected
	keyStrokeCounter int // how many keystorkes player made
}

func initialModel() model {
	return model{
		words:            wordsToDisplay,
		redHighlight:     []int{},
		greenHighlight:   []int{},
		keyStrokeCounter: 0,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		// selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.keyStrokeCounter++
	var correctChar = string(m.words[m.keyStrokeCounter])

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case correctChar:
			m.greenHighlight = append(m.greenHighlight, m.keyStrokeCounter)

		default:
			m.redHighlight = append(m.redHighlight, m.keyStrokeCounter)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Turtle Typing"

	for i, ch := range m.words {
		if word_sets.Contains(m.greenHighlight, i) {
			s += fmt.Sprintf(Green + string(ch) + Green)
		} else {
			s += fmt.Sprintf(Red + string(ch) + Red)
		}

	}

	// The footer
	s += "\nPress ctrl+c to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	fmt.Println("HELLO")
	fmt.Println(wordsToDisplay)
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
