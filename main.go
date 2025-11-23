package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

const (
	modeNormal = "Normal"
	modeSearch = "Search"
	modeDetail = "Detail"
)

type model struct {
	osEnvVars     SelectionModel
	localEnvVars  SelectionModel
	mode          string
	searchTerm    string
	statusMessage string
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

	envSlice := make([]string, 0, len(localEnv))
	for k, v := range localEnv {
		envSlice = append(envSlice, k+"="+v)
	}

	envList := os.Environ()
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

	osEnvVars := SelectionModel{
		variables: envList,
		choices:   envList,
		selected:  map[int]struct{}{},
		cursor:    0,
	}

	localEnvVars := SelectionModel{
		variables: localEnv,
		choices:   localEnv,
		selected:  map[int]struct{}{},
		cursor:    0,
	}

	return model{
		osEnvVars:     osEnvVars,
		localEnvVars:  localEnvVars,
		mode:          initMode,
		searchTerm:    "",
		statusMessage: "",
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
		m.statusMessage = ""

		switch m.mode {
		case modeNormal:
			switch msg.String() {
			// These keys should exit the program.
			case "ctrl+c", "q":
				return m, tea.Quit
			// The "up" and "k" keys move the cursor up
			case "up", "k":
				if m.osEnvVars.cursor > 0 {
					m.osEnvVars.cursor--
				}
			// The "down" and "j" keys move the cursor down
			case "down", "j":
				if m.osEnvVars.cursor < len(m.osEnvVars.choices)-1 {
					m.osEnvVars.cursor++
				}

			case "enter", " ":
				_, ok := m.osEnvVars.selected[m.osEnvVars.cursor]
				if ok {
					delete(m.osEnvVars.selected, m.osEnvVars.cursor)
				} else {
					m.osEnvVars.selected[m.osEnvVars.cursor] = struct{}{}
				}

			case "y":
				v := m.osEnvVars.variables[m.osEnvVars.cursor]
				err := clipboard.WriteAll(v)
				if err != nil {
					fmt.Println("Failed to copy to clipboard:", err)
				} else {
					m.statusMessage = "Successfully copied to clipboard"
				}
				return m, nil

			case "s":
				m.mode = modeSearch

			case "d":
				m.mode = modeDetail
			}
		case modeSearch:
			switch msg.Type {
			case tea.KeyEsc:
				m.mode = modeNormal
				m.searchTerm = ""
			case tea.KeyBackspace:
				if len(m.searchTerm) > 0 {
					m.searchTerm = m.searchTerm[:len(m.searchTerm)-1]
				}
			default:
				if msg.Type == tea.KeyRunes {
					m.searchTerm += msg.String()
				}
			}
		case modeDetail:
			switch msg.String() {
			case "esc":
				m.mode = modeNormal
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.localEnvVars.cursor > 0 {
					m.localEnvVars.cursor--
				}
			case "down", "j":
				if m.localEnvVars.cursor < len(m.localEnvVars.choices)-1 {
					m.localEnvVars.cursor++
				}
			case "enter", " ":
				_, ok := m.localEnvVars.selected[m.localEnvVars.cursor]
				if ok {
					delete(m.localEnvVars.selected, m.localEnvVars.cursor)
				} else {
					m.localEnvVars.selected[m.localEnvVars.cursor] = struct{}{}
				}
			case "y":
				v := m.localEnvVars.variables[m.localEnvVars.cursor]
				err := clipboard.WriteAll(v)
				if err != nil {
					fmt.Println("Failed to copy to clipboard:", err)
				} else {
					m.statusMessage = "Successfully copied to clipboard"
				}
				return m, nil
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
		for index, choice := range model.localEnvVars.variables {
			cursorSymbol := " "
			if model.localEnvVars.cursor == index {
				cursorSymbol = ">"
			}

			checked := " "
			if _, ok := model.localEnvVars.selected[index]; ok {
				checked = "x"
			}

			renderedList += fmt.Sprintf("%s [%s] %s\n", cursorSymbol, checked, choice)
		}
		return renderedList
	default:
		for index, choice := range model.osEnvVars.choices {
			if model.mode == modeSearch && !strings.Contains(strings.ToLower(choice), strings.ToLower(model.searchTerm)) {
				continue
			}
			cursor := " "
			if model.osEnvVars.cursor == index {
				cursor = ">"
			}

			checked := " "
			if _, ok := model.osEnvVars.selected[index]; ok {
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
		footer += "\nPress y to copy value to your clipboard. Press q to quit.  \n"
		footer += "Normal mode - press s to search, press d for details."
	case modeSearch:
		footer += "\nSearch Query: " + m.searchTerm
		footer += "\nSearch mode - press esc for normal mode."
	case modeDetail:
		footer += "\nPress y to copy value to your clipboard. Press v to toggle details"
		footer += "\nDetail mode - press esc for normal mode."
	default:
		footer += "Unknown mode."
	}

	if m.statusMessage != "" {
		footer += "\n" + m.statusMessage
	}

	return footer
}
