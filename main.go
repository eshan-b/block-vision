// main.go
package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

type ViewMode int

const (
	MenuView ViewMode = iota
	TrendingView
)

type model struct {
	viewMode   ViewMode
	menuChoice int
	keys       keyMap
	help       help.Model
	table      table.Model
	err        error
}

var (
	gradientStart, _ = colorful.Hex("#F096DD")
	gradientEnd, _   = colorful.Hex("#BC52F1")
	subtextColor     = lipgloss.Color("#F095DD")
	checkboxColor    = lipgloss.Color("#FFFFFF")
	checkboxChecked  = lipgloss.Color("#BC52F1")
	baseStyle        = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240"))
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding
	Back   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Select, k.Help, k.Back, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			if m.viewMode != MenuView {
				m.viewMode = MenuView
				return m, nil
			}
		}

		switch m.viewMode {
		case MenuView:
			return m.updateMenu(msg)
		case TrendingView:
			return m.updateTrending(msg)
		}
	case errMsg:
		m.err = msg.err
		return m, nil
	case trendingDataMsg:
		m.table.SetRows(msg.rows)
		return m, nil
	}
	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		m.menuChoice--
		if m.menuChoice < 0 {
			m.menuChoice = 0
		}
	case key.Matches(msg, m.keys.Down):
		m.menuChoice++
		if m.menuChoice > 1 {
			m.menuChoice = 1
		}
	case key.Matches(msg, m.keys.Select):
		if m.menuChoice == 0 {
			m.viewMode = TrendingView
			return m, fetchTrendingCmd
		}
	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
	}
	return m, nil
}

func (m model) updateTrending(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
		return m, nil
	case key.Matches(msg, m.keys.Up), key.Matches(msg, m.keys.Down):
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	switch m.viewMode {
	case MenuView:
		return m.menuView()
	case TrendingView:
		return m.trendingView()
	default:
		return "Unknown view"
	}
}

func (m model) menuView() string {
	title := gradientText("❒ block-vision", gradientStart, gradientEnd)
	subtext := lipgloss.NewStyle().
		Foreground(subtextColor).
		Render("• v1.0.0  ✦︎ Lighting up your crypto journey ✦︎")
	checkboxOptions := checkboxPicker([]string{"Trending", "Search"}, m.menuChoice)
	helpView := m.help.View(m.keys)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtext,
		"",
		checkboxOptions,
		"",
		helpView,
	)
}

func (m model) trendingView() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	// Combine the table view with the help view
	helpView := m.help.View(m.keys)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		baseStyle.Render(m.table.View()),
		"", // Add a blank line between table and help
		helpView,
	)
}

// Message types for handling trending data
type trendingDataMsg struct {
	rows []table.Row
}

type errMsg struct {
	err error
}

func fetchTrendingCmd() tea.Msg {
	rows, err := fetchTrendingCryptos()
	if err != nil {
		return errMsg{err}
	}
	return trendingDataMsg{rows}
}

// Helper functions
func gradientText(text string, start, end colorful.Color) string {
	gradient := make([]string, len(text))
	for i, char := range text {
		color := start.BlendLuv(end, float64(i)/float64(len(text)))
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color.Hex()))
		gradient[i] = style.Render(string(char))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, gradient...)
}

func checkboxPicker(options []string, choice int) string {
	var renderedOptions []string
	for i, option := range options {
		checked := "[ ]"
		if i == choice {
			checked = fmt.Sprintf("[%s]", lipgloss.NewStyle().Foreground(checkboxChecked).Render("x"))
		}
		optionText := lipgloss.NewStyle().Foreground(checkboxColor).Render(option)
		renderedOptions = append(renderedOptions, fmt.Sprintf("%s %s", checked, optionText))
	}
	return lipgloss.JoinVertical(lipgloss.Left, renderedOptions...)
}

func initialModel() model {
	return model{
		viewMode: MenuView,
		keys:     keys,
		help:     help.New(),
		table:    InitializeTrendingTable(),
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
