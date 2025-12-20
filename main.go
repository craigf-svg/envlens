package main

import (
	"flag"
	"fmt"
	"os"

	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

const (
	modeNormal            = "Normal"
	modeSearch            = "Search"
	modeEdit              = "Edit"
	modeLocalEnv          = "LocalEnv"
	footerRightPadding    = 5
	ctrlY                 = rune(0x19)
	errMsgClipboardFailed = "Failed to copy to clipboard:"
)

func supportsModernTerminal() bool {
	return os.Getenv("WT_SESSION") != "" || os.Getenv("TERM_PROGRAM") != "" || os.Getenv("COLORTERM") != ""
}

func icon(emoji, fallback string) string {
	if supportsModernTerminal() {
		return emoji
	}
	return fallback
}

var (
	cursorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	checkStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styledCheck       = checkStyle.Render("x")
	searchCursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
)

type model struct {
	osEnvVars     SelectionModel
	localEnvVars  SelectionModel
	mode          string
	searchTerm    string
	searchCursor  int
	editBuffer    string
	editCursor    int
	editingIndex  int
	statusMessage string
	hideValues    bool
	hasLocalEnv   bool
	height        int
	width         int
}

func main() {
	showVersion := flag.Bool("version", false, "Print version and exit")
	demoMode := flag.Bool("demo", false, "Run with test data")
	flag.Parse()

	if *showVersion {
		fmt.Println("envlens", getVersion())
		return
	}

	var osEnvList []string
	var localEnvList []string
	var hasLocalEnv bool

	if *demoMode {
		osEnvList = demoEnvVars()
		localEnvList = demoLocalEnvVars()
		hasLocalEnv = true
	} else {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("No local .env found:", err)
		}

		localEnv, err := godotenv.Read()
		if err != nil {
			fmt.Println("Error reading .env:", err)
		}

		localEnvList = make([]string, 0, len(localEnv))
		for k, v := range localEnv {
			localEnvList = append(localEnvList, k+"="+v)
		}

		osEnvList = os.Environ()
		hasLocalEnv = (err == nil && len(localEnv) > 0)
	}

	osHiddenList := autoHideFilter(osEnvList)
	localHiddenList := autoHideFilter(localEnvList)

	hideValuesDefault := false
	p := tea.NewProgram(
		initialModel(osEnvList, modeNormal, localEnvList, hideValuesDefault, hasLocalEnv, osHiddenList, localHiddenList),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func autoHideFilter(envList []string) map[int]struct{} {
	hiddenList := make(map[int]struct{})
	for i, env := range envList {
		env = strings.ToLower(env)
		key, _, found := strings.Cut(env, "=")

		if !found {
			continue
		}
		if strings.Contains(key, "key") || strings.Contains(key, "private") || strings.Contains(key, "secret") {
			hiddenList[i] = struct{}{}
		}
	}

	return hiddenList
}

func initialModel(envList []string, initMode string, localEnv []string, hideValues bool, hasLocalEnv bool, hidden map[int]struct{}, localHidden map[int]struct{}) model {
	return model{
		osEnvVars: SelectionModel{
			variables: envList,
			choices:   envList,
			selected:  map[int]struct{}{},
			hidden:    hidden,
		},
		localEnvVars: SelectionModel{
			variables: localEnv,
			choices:   localEnv,
			selected:  map[int]struct{}{},
			hidden:    localHidden,
		},
		mode:        initMode,
		hideValues:  hideValues,
		hasLocalEnv: hasLocalEnv,
		height:      20,
		width:       80,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *model) copyItemToClipboard(text string) {
	err := clipboard.WriteAll(text)
	if err != nil {
		fmt.Println(errMsgClipboardFailed, err)
	} else {
		m.statusMessage = "Successfully copied to clipboard"
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle terminal resize
	case tea.WindowSizeMsg:
		m.height = msg.Height - 5
		m.width = msg.Width
		if m.height < 3 {
			m.height = 3
		}
		return m, nil

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

			case "e":
				// Enter edit mode
				if len(m.osEnvVars.choices) > 0 && m.osEnvVars.cursor < len(m.osEnvVars.choices) {
					m.mode = modeEdit
					m.editingIndex = m.osEnvVars.cursor
					m.editBuffer = m.osEnvVars.choices[m.osEnvVars.cursor]
					// place cursor at end of buffer
					m.editCursor = len([]rune(m.editBuffer))
				}
			case "enter", " ":
				_, ok := m.osEnvVars.selected[m.osEnvVars.cursor]
				if ok {
					delete(m.osEnvVars.selected, m.osEnvVars.cursor)
				} else {
					m.osEnvVars.selected[m.osEnvVars.cursor] = struct{}{}
				}

			case "tab":
				m.hideValues = !m.hideValues

			case "y":
				v := m.osEnvVars.variables[m.osEnvVars.cursor]
				status, err := copySingleVarToClipboard(v)
				if err != nil {
					fmt.Println(errMsgClipboardFailed, err)
				} else {
					m.statusMessage = status
				}
				return m, nil

			case "Y":
				status, err := copySelectedVarsToClipboard(m.osEnvVars.selected, m.osEnvVars.variables)
				if err != nil {
					fmt.Println(errMsgClipboardFailed, err)
				} else {
					m.statusMessage = status
				}
				return m, nil

			case "s":
				m.mode = modeSearch

			case "d":
				if !m.hasLocalEnv {
					m.statusMessage = "Failed to load .env file at runtime"
					return m, nil
				}
				m.mode = modeLocalEnv
			}

		case modeEdit:
			// Separate edit mode - editBuffer holds the current text being edited
			switch msg.Type {
			case tea.KeyEsc:
				m.mode = modeNormal
				m.editBuffer = ""
				m.editCursor = 0
				return m, nil
			case tea.KeyBackspace:
				// delete rune before cursor
				r := []rune(m.editBuffer)
				if m.editCursor > 0 && len(r) > 0 {
					idx := m.editCursor
					if idx > len(r) {
						idx = len(r)
					}
					r = append(r[:idx-1], r[idx:]...)
					m.editBuffer = string(r)
					m.editCursor--
				}
				return m, nil
			case tea.KeyLeft:
				if m.editCursor > 0 {
					m.editCursor--
				}
				return m, nil
			case tea.KeyRight:
				r := []rune(m.editBuffer)
				if m.editCursor < len(r) {
					m.editCursor++
				}
				return m, nil
			case tea.KeyTab:
				m.hideValues = !m.hideValues
				return m, nil
			case tea.KeyEnter, tea.KeySpace:
				// save to memory and persist to .env file
				if m.editingIndex >= 0 && m.editingIndex < len(m.osEnvVars.variables) {
					key, val, found := strings.Cut(m.editBuffer, "=")
					if !found {
						m.statusMessage = "Edit must be in KEY=VALUE format; not saved to .env"
						m.editBuffer = ""
					} else {
						m.osEnvVars.variables[m.editingIndex] = m.editBuffer
						m.osEnvVars.choices[m.editingIndex] = m.editBuffer
						m.statusMessage = "Edited variable"

						// set in process environment
						_ = os.Setenv(key, val)

						// attempt to persist to system environment as well
						if perr := persistPlatformEnv(key, val); perr != nil {
							m.statusMessage = "Failed to persist system env: " + perr.Error()
						} else {
							m.statusMessage = "Saved to system environment; restart terminal to apply"
						}

					}
				}
				m.editBuffer = ""
				m.editCursor = 0
				m.mode = modeNormal
				return m, nil
			default:
				if msg.Type != tea.KeyRunes {
					return m, nil
				}
				// Ignore alt and control
				runes := msg.Runes
				if msg.Alt || (len(runes) == 1 && runes[0] < 32) {
					return m, nil
				}
				// insert runes at cursor
				r := []rune(m.editBuffer)
				runesToInsert := runes
				idx := m.editCursor
				if idx > len(r) {
					idx = len(r)
				}
				r = append(r[:idx], append(runesToInsert, r[idx:]...)...)
				m.editBuffer = string(r)
				m.editCursor += len(runesToInsert)
				return m, nil
			}
		case modeSearch:
			items, indices := filterChoices(m.osEnvVars.choices, m.searchTerm)
			switch msg.Type {
			case tea.KeyEsc:
				m.mode = modeNormal
				m.searchTerm = ""
				m.searchCursor = 0
				return m, nil
			case tea.KeyBackspace:
				if len(m.searchTerm) > 0 {
					m.searchTerm = m.searchTerm[:len(m.searchTerm)-1]
					m.searchCursor = 0
				}
			case tea.KeyTab:
				m.hideValues = !m.hideValues
			case tea.KeyUp:
				if m.searchCursor > 0 {
					m.searchCursor--
				}
			case tea.KeyDown:
				if m.searchCursor < len(items)-1 {
					m.searchCursor++
				}
			case tea.KeyEnter, tea.KeySpace:
				if len(indices) > 0 {
					idx := indices[m.searchCursor]
					if _, ok := m.osEnvVars.selected[idx]; ok {
						delete(m.osEnvVars.selected, idx)
					} else {
						m.osEnvVars.selected[idx] = struct{}{}
					}
				}
			case tea.KeyCtrlY:
				if len(items) > 0 && m.searchCursor < len(items) {
					m.copyItemToClipboard(items[m.searchCursor])
				}
				return m, nil
			default:
				if msg.Type != tea.KeyRunes {
					return m, nil
				}
				runes := msg.Runes
				// Handle ctrl+y
				if len(runes) == 1 && runes[0] == ctrlY {
					if len(items) > 0 && m.searchCursor < len(items) {
						m.copyItemToClipboard(items[m.searchCursor])
					}
					return m, nil
				}

				// Ignore alt and control
				if msg.Alt || (len(runes) == 1 && runes[0] < 32) {
					return m, nil
				}

				m.searchTerm += msg.String()
				m.searchCursor = 0
			}
		case modeLocalEnv:
			switch msg.String() {
			case "esc", "d":
				m.mode = modeNormal
			case "tab":
				m.hideValues = !m.hideValues
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
				status, err := copySingleVarToClipboard(v)
				if err != nil {
					fmt.Println(errMsgClipboardFailed, err)
				} else {
					m.statusMessage = status
				}
				return m, nil
			case "Y":
				status, err := copySelectedVarsToClipboard(m.localEnvVars.selected, m.localEnvVars.variables)
				if err != nil {
					fmt.Println(errMsgClipboardFailed, err)
				} else {
					m.statusMessage = status
				}
				return m, nil
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func filterChoices(choices []string, term string) ([]string, []int) {
	term = strings.ToLower(term)
	var items []string
	var indices []int
	for i, choice := range choices {
		if strings.Contains(strings.ToLower(choice), term) {
			items = append(items, choice)
			indices = append(indices, i)
		}
	}
	return items, indices
}

func visibleRange(cursor, total, height int) (start, end int) {
	if total <= height {
		return 0, total
	}
	start = cursor - height/2
	if start < 0 {
		start = 0
	}
	end = start + height
	if end > total {
		end = total
		start = end - height
	}
	return start, end
}

func (m model) View() string {
	var header string
	switch m.mode {
	case modeNormal:
		header = icon("📋", "[ENV]") + " Environment Variables:"
	case modeLocalEnv:
		header = icon("📁", "[.env]") + " Local .env file:"
	case modeSearch:
		header = icon("🔍", "[S]") + " Search Results:"
	default:
		header = "Environment Variables:"
	}
	s := header + "\n\n"
	s += renderList(m)
	s += renderFooter(m)
	return s
}

func renderList(m model) string {
	var items []string
	var cursor int
	var selected map[int]struct{}
	var hiddenList map[int]struct{}

	switch m.mode {
	case modeLocalEnv:
		items = m.localEnvVars.variables
		cursor = m.localEnvVars.cursor
		selected = m.localEnvVars.selected
		hiddenList = m.localEnvVars.hidden
	case modeSearch:
		var indices []int
		items, indices = filterChoices(m.osEnvVars.choices, m.searchTerm)
		cursor = m.searchCursor
		if cursor >= len(items) {
			cursor = len(items) - 1
		}
		if cursor < 0 {
			cursor = 0
		}
		selected = make(map[int]struct{})
		hiddenList = make(map[int]struct{})
		for i, origIdx := range indices {
			if _, ok := m.osEnvVars.selected[origIdx]; ok {
				selected[i] = struct{}{}
			}
			if _, ok := m.osEnvVars.hidden[origIdx]; ok {
				hiddenList[i] = struct{}{}
			}
		}
		if len(items) == 0 {
			return "No results found\n"
		}
	default:
		items = m.osEnvVars.choices
		cursor = m.osEnvVars.cursor
		selected = m.osEnvVars.selected
		hiddenList = m.osEnvVars.hidden
	}

	listHeight := m.height
	if m.mode == modeSearch {
		listHeight = m.height - 3
	}
	start, end := visibleRange(cursor, len(items), listHeight)

	var output string
	for i := start; i < end; i++ {
		drawCursor := i == cursor
		symbol := " "
		if drawCursor {
			symbol = ">"
		}
		check := " "
		if _, ok := selected[i]; ok {
			if drawCursor {
				check = "x"
			} else {
				check = styledCheck
			}
		}
		line := "%s [%s] %s"
		var formatLine string

		if _, hidden := hiddenList[i]; hidden {
			formatLine = fmt.Sprintf(line, symbol, check, maskEnvVar(items[i], m.hideValues, true))
		} else {
			formatLine = fmt.Sprintf(line, symbol, check, maskEnvVar(items[i], m.hideValues, false))
		}

		if drawCursor {
			formatLine = cursorStyle.Render(formatLine)
		}
		output += formatLine + "\n"
	}
	return output
}

func renderFooter(m model) string {
	footer := ""
	switch m.mode {
	case modeNormal:
		cursorPosition := fmt.Sprintf("<%d-%d>", m.osEnvVars.cursor+1, len(m.osEnvVars.variables))
		footer += lipgloss.PlaceHorizontal(m.width-footerRightPadding, lipgloss.Right, cursorPosition)
		footer += "\n[↑/↓] Navigate [↵] Select [e] Edit  [y/Y] Copy (one/all)  [tab] Toggle  [s] Search  [d] Local  [q] Quit"

	case modeEdit:
		cursorBlock := searchCursorStyle.Render("█")
		width := m.width - 1
		if width < 0 {
			width = 0
		}
		editBorderStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("245")).
			Padding(0, 1).
			Width(width)

		// render edit buffer with cursor inserted at editCursor
		r := []rune(m.editBuffer)
		if m.editCursor < 0 {
			m.editCursor = 0
		}
		if m.editCursor > len(r) {
			m.editCursor = len(r)
		}
		var display string
		// insert cursor block in between runes
		display = string(r[:m.editCursor]) + cursorBlock + string(r[m.editCursor:])

		editWithBorder := editBorderStyle.Render(fmt.Sprintf("Edit: %s", display))
		footer += fmt.Sprintf("\n%s\n[←/→] Move [↵] Save [tab] Toggle [esc] Cancel", editWithBorder)

	case modeSearch:
		cursorBlock := searchCursorStyle.Render("█")
		items, _ := filterChoices(m.osEnvVars.choices, m.searchTerm)
		cursorPosition := fmt.Sprintf("<%d-%d>", m.searchCursor+1, len(items))
		footer = lipgloss.PlaceHorizontal(m.width-footerRightPadding, lipgloss.Right, cursorPosition)
		width := m.width - 1
		if width < 0 {
			width = 0
		}
		searchBorderStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("245")).
			Padding(0, 1).
			Width(width)

		content := fmt.Sprintf("Search mode: %s%s", m.searchTerm, cursorBlock)
		searchWithBorder := searchBorderStyle.Render(content)
		footer += fmt.Sprintf("\n%s\n[↑/↓] Navigate [↵] Select [e] Edit [ctrl+y] Copy (one) [tab] Toggle [esc] Back", searchWithBorder)
	case modeLocalEnv:
		footer += "\n[↑/↓] Navigate [↵] Select [e] Edit [y/Y] Copy (one/all) [tab] Toggle [d] Global  [q] Quit"
	default:
		footer += "\n[?] Unknown mode"
	}

	footer += "\n" + m.statusMessage

	return footer
}

func maskEnvVar(line string, hide bool, hidden bool) string {
	if hide || hidden {
		key, val, found := strings.Cut(line, "=")

		if !found {
			return line
		}

		return key + "=" + strings.Repeat("*", len(val))
	}

	return line
}

// persistPlatformEnv attempts to persist an environment variable across sessions on the host platform.
// On Windows it calls `setx`to set the variable for the user,
// on Unix like systems it updates (or appends) an `export KEY="VALUE"` line in a shell rc file.

// NOTEE-  user have to restart their terminal or source the rc file for changes to take effect.
func persistPlatformEnv(key, val string) error {
	// setx on Windows will persist the variable for the user
	if runtime.GOOS == "windows" {
		cmd := exec.Command("setx", key, val)
		return cmd.Run()
	}

	// For Unix like systems, update the first existing rc file or ~/.profile
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	candidates := []string{".profile", ".bash_profile", ".bashrc", ".zshrc"}
	var target string
	for _, f := range candidates {
		p := filepath.Join(home, f)
		if _, err := os.Stat(p); err == nil {
			target = p
			break
		}
	}
	if target == "" {
		target = filepath.Join(home, ".profile")
	}

	// read existing content (ignore error -> treat as empty)
	data, _ := os.ReadFile(target)
	content := string(data)

	// prepare export line - escape double quotes
	esc := strings.ReplaceAll(val, "\"", "\\\"")
	newLine := fmt.Sprintf("export %s=\"%s\"", key, esc)

	// replace existing export KEY=... line if present
	re := regexp.MustCompile(`(?m)^\s*export\s+` + regexp.QuoteMeta(key) + `=.*$`)
	if re.MatchString(content) {
		content = re.ReplaceAllString(content, newLine)
	} else {
		if len(content) > 0 && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += newLine + "\n"
	}

	// write back to file

	return os.WriteFile(target, []byte(content), 0644)
}
