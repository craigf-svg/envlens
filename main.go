package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

const (
	modeNormal = "Normal"
	modeSearch = "Search"
)

type model struct {
	variables []string
	cursor    int
	choices   []string
	selected  map[int]struct{}
	mode      string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No local .env found")
	}

	var envList []string = os.Environ()
	if os.Getenv("DEBUG") != "" {
		printList(envList)
	}

	p := tea.NewProgram(initialModel(envList, modeNormal))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func printList(list []string) {
	for count := 0; count < len(list); count++ {
		fmt.Printf("Environment Variable #%v: %v\n", count, list[count])
	}
	fmt.Println("Total Number of Environment Variables", len(list))
}

func initialModel(prop []string, initMode string) model {
	return model{
		variables: prop,
		choices:   prop,
		selected:  map[int]struct{}{},
		mode:      initMode,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		switch m.mode {
		case modeNormal:
			// Cool, what was the actual key pressed?
			switch msg.String() {

			// These keys should exit the program.
			case "ctrl+c", "q":
				return m, tea.Quit

			// The "up" and "k" keys move the cursor up
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			// The "down" and "j" keys move the cursor down
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case "enter", " ":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}

			case "s":
				m.mode = modeSearch
			}
		case modeSearch:
			switch msg.String() {
			case "esc":
				m.mode = modeNormal
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	s := "Select environment variables:\n\n"
	s += renderList(m)
	s += renderFooter(m)
	return s
}

func renderList(model model) string {
	var renderedList string
	renderedList = ""

	for index, choice := range model.choices {
		cursor := " "
		if model.cursor == index {
			cursor = ">"
		}

		checked := " "
		if _, ok := model.selected[index]; ok {
			checked = "x"
		}

		renderedList += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	return renderedList
}

func renderFooter(m model) string {
	footer := "\nPress q to quit.\n"
	switch m.mode {
	case modeNormal:
		footer += "Normal mode - press s to search."
	case modeSearch:
		footer += "Search mode - press esc to exit."
	default:
		footer += "Unknown mode."
	}

	return footer
}
