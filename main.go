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
	modeDetail = "Detail"
)

type model struct {
	variables     []string
	cursor        int
	choices       []string
	selected      map[int]struct{}
	mode          string
	searchTerm    string
	localVars     []string
	localChoices  []string
	localSelected map[int]struct{}
	localCursor   int
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No local .env found")
	}

	localEnv, err := godotenv.Read()
	if err != nil {
		fmt.Println("err", err)
	} else {
		fmt.Println("localEnv", localEnv)
	}

	var envSlice []string
	for k, v := range localEnv {
		envSlice = append(envSlice, k+"="+v)
	}

	var envList []string = os.Environ()
	if os.Getenv("DEBUG") != "" {
		printList(envList)
	}

	p := tea.NewProgram(initialModel(envList, modeNormal, envSlice))
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

func initialModel(envList []string, initMode string, localEnv []string) model {
	return model{
		variables:     envList,
		choices:       envList,
		selected:      map[int]struct{}{},
		mode:          initMode,
		searchTerm:    "",
		localVars:     localEnv,
		localChoices:  localEnv,
		localSelected: map[int]struct{}{},
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

			case "enter", " ":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}

			case "s":
				m.mode = modeSearch

			case "d":
				m.mode = modeDetail
			}
		case modeSearch:
			switch msg.String() {
			case "esc":
				m.mode = modeNormal
			case "o":
				m.searchTerm += "o"
			}
		case modeDetail:
			switch msg.String() {
			case "esc":
				m.mode = modeNormal
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.localCursor > 0 {
					m.localCursor--
				}
			case "down", "j":
				if m.localCursor < len(m.localChoices)-1 {
					m.localCursor++
				}
			case "enter", " ":
				_, ok := m.localSelected[m.localCursor]
				if ok {
					delete(m.localSelected, m.localCursor)
				} else {
					m.localSelected[m.localCursor] = struct{}{}
				}
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

	switch model.mode {

	case modeDetail:
		for index, choice := range model.localVars {
			localCursor := " "
			if model.localCursor == index {
				localCursor = ">"
			}

			checked := " "
			if _, ok := model.localSelected[index]; ok {
				checked = "x"
			}

			renderedList += fmt.Sprintf("%s [%s] %s\n", localCursor, checked, choice)
		}
		return renderedList
	default:
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
}

func renderFooter(m model) string {
	footer := ""
	switch m.mode {
	case modeNormal:
		footer += "\nPress q to quit.\n"
		footer += "Normal mode - press s to search, press d for details."
	case modeSearch:
		footer += "\nSearch Query: " + m.searchTerm
		footer += "\nSearch mode - press esc for normal mode."
	case modeDetail:
		footer += "\nPress v to toggle details."
		footer += "\nDetail mode - press esc for normal mode."
	default:
		footer += "Unknown mode."
	}

	return footer
}
