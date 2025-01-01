package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// Define gradient colors and styles
var (
	gradientStart, _ = colorful.Hex("#F096DD")
	gradientEnd, _   = colorful.Hex("#BC52F1")
	subtextColor     = lipgloss.Color("#F095DD")
	checkboxColor    = lipgloss.Color("#FFFFFF")
	checkboxChecked  = lipgloss.Color("#BC52F1")
)

// keyMap defines keybindings for navigation and actions
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the compact help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings to be shown in the expanded help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},             // First column
		{k.Select, k.Help, k.Quit}, // Second column
	}
}

// Initialize keybindings
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
		key.WithHelp("enter", "select option"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// Model for the Bubble Tea program
type model struct {
	Choice int // Store selected choice (0 or 1)
	keys   keyMap
	help   help.Model
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
		case key.Matches(msg, m.keys.Up):
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case key.Matches(msg, m.keys.Down):
			m.Choice++
			if m.Choice > 1 {
				m.Choice = 1
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}

func (m model) View() string {
	// Gradient title
	title := gradientText("❒ block-vision", gradientStart, gradientEnd)

	// Subtitle
	subtext := lipgloss.NewStyle().
		Foreground(subtextColor).
		Render("• v1.0.0  ✦︎ Lighting up your crypto journey ✦︎")

	// Checkbox options
	checkboxOptions := checkboxPicker([]string{"Trending", "Search"}, m.Choice)

	// Help view
	helpView := m.help.View(m.keys)

	// Combine all parts with additional newlines
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

// Helper to create gradient text
func gradientText(text string, start, end colorful.Color) string {
	gradient := make([]string, len(text))
	for i, char := range text {
		color := start.BlendLuv(end, float64(i)/float64(len(text)))
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color.Hex()))
		gradient[i] = style.Render(string(char))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, gradient...)
}

// Create the checkbox options
func checkboxPicker(options []string, choice int) string {
	var renderedOptions []string
	for i, option := range options {
		// Determine if the option is selected
		checked := "[ ]"
		if i == choice {
			checked = fmt.Sprintf("[%s]", lipgloss.NewStyle().Foreground(checkboxChecked).Render("x"))
		}

		// Style the option
		optionText := lipgloss.NewStyle().Foreground(checkboxColor).Render(option)
		renderedOptions = append(renderedOptions, fmt.Sprintf("%s %s", checked, optionText))
	}
	return lipgloss.JoinVertical(lipgloss.Left, renderedOptions...)
}

func main() {
	p := tea.NewProgram(model{keys: keys, help: help.New()}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
