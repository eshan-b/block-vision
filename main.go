package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/lucasb-eyer/go-colorful"
)

// ViewMode represents the different application states or "views"
type ViewMode int

// Enum definitions for view modes
const (
	MenuView     ViewMode = iota // main menu view
	TrendingView                 // view displaying trending data
	SearchView                   // view for searching cryptocurrencies
	ViewingList                  // view for displaying search results
	ViewingInfo                  // view for displaying detailed coin information
)

// The main application model holds the state of the program
type model struct {
	viewMode    ViewMode        // current view mode
	menuChoice  int             // selected option in the menu
	keys        keyMap          // key bindings
	help        help.Model      // help menu model
	table       table.Model     // table to display data
	err         error           // error tracking for data operations
	textInput   textinput.Model // input for search
	list        list.Model      // list for displaying search results
	selected    Coin            // currently selected result
	coinDetails CoinDetails     // detailed info of selected result
}

// Color definitions for gradients and UI elements
var (
	gradientStart, _ = colorful.Hex("#F096DD")   // start color of the gradient
	gradientEnd, _   = colorful.Hex("#BC52F1")   // end color of the gradient
	subtextColor     = lipgloss.Color("#F095DD") // subtext color
	subtleColor      = lipgloss.Color("#898989") // subtle color
	checkboxColor    = lipgloss.Color("#FFFFFF") // checkbox outline color
	checkboxChecked  = lipgloss.Color("#BC52F1") // checkbox checked color
	baseStyle        = lipgloss.NewStyle().      // base style for table rendering
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240"))
	docStyle = lipgloss.NewStyle().Margin(1, 2) // style for the list view
)

// keyMap defines bindings for user input and their associated actions
type keyMap struct {
	Up     key.Binding // navigate up
	Down   key.Binding // navigate down
	Select key.Binding // select an item
	Help   key.Binding // toggle help menu
	Quit   key.Binding // quit the application
	Back   key.Binding // navigate back
}

// ShortHelp returns a minimal list of key bindings for display
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp provides a comprehensive list of key bindings
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},                     // navigation keys
		{k.Select, k.Help, k.Back, k.Quit}, // action keys
	}
}

// Default key bindings for the application
var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"), // moving up
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"), // moving down
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"), // selecting an option
		key.WithHelp("enter", "select"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"), // toggle help
		key.WithHelp("?", "toggle help"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"), // navigate back
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"), // quit the program
		key.WithHelp("q", "quit"),
	),
}

// Initialization method for the model; no initial commands are required
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles user input and updates the program's state accordingly
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit and back key presses globally
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			if m.viewMode != MenuView {
				m.viewMode = MenuView
				return m, nil
			}
		}

		// Handle view-specific input
		switch m.viewMode {
		case MenuView:
			return m.updateMenu(msg) // update logic for the menu view
		case TrendingView:
			return m.updateTrending(msg) // update logic for the trending view
		case SearchView:
			return m.updateSearch(msg) // update logic for the search view
		case ViewingList:
			return m.updateViewingList(msg) // update logic for the viewing list view
		case ViewingInfo:
			return m.updateViewingInfo(msg) // update logic for the viewing info view
		}
	case errMsg:
		// Handle errors and update the model's error field
		m.err = msg.err
		return m, nil
	case trendingDataMsg:
		// Populate table with fetched trending data
		m.table.SetRows(msg.rows)
		return m, nil
	}
	return m, nil
}

// Menu-specific update logic
func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		// Navigate up through menu options
		m.menuChoice--
		if m.menuChoice < 0 {
			m.menuChoice = 0
		}
	case key.Matches(msg, m.keys.Down):
		// Navigate down through menu options
		m.menuChoice++
		if m.menuChoice > 1 {
			m.menuChoice = 1
		}
	case key.Matches(msg, m.keys.Select):
		// Select an option and change view accordingly
		if m.menuChoice == 0 {
			m.viewMode = TrendingView
			return m, fetchTrendingCmd
		} else if m.menuChoice == 1 {
			m.viewMode = SearchView
			return m, nil
		}
	case key.Matches(msg, m.keys.Help):
		// Toggle help menu visibility
		m.help.ShowAll = !m.help.ShowAll
	}
	return m, nil
}

// Logic for updating the trending view
func (m model) updateTrending(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Help):
		// Toggle help menu in trending view
		m.help.ShowAll = !m.help.ShowAll
		return m, nil
	case key.Matches(msg, m.keys.Up), key.Matches(msg, m.keys.Down):
		// Handle table navigation
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}
	return m, nil
}

// Logic for updating the search view
func (m model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		query := m.textInput.Value()
		items, err := fetchCoins(query)
		if err != nil {
			m.err = err
			return m, nil
		}
		m.list.SetItems(items)
		m.viewMode = ViewingList
	case "esc":
		m.viewMode = MenuView
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// Logic for updating the viewing list view
func (m model) updateViewingList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		selectedItem := m.list.SelectedItem().(item)
		m.selected = selectedItem.data
		coinDetails, err := fetchCoinDetails(m.selected.ID)
		if err != nil {
			m.err = err
			return m, nil
		}
		m.coinDetails = coinDetails
		m.viewMode = ViewingInfo
	case "esc":
		m.viewMode = SearchView
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// Logic for updating the viewing info view
func (m model) updateViewingInfo(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		m.viewMode = ViewingList
	}
	return m, nil
}

// View method to render the appropriate screen based on viewMode
func (m model) View() string {
	switch m.viewMode {
	case MenuView:
		return m.menuView() // render menu view
	case TrendingView:
		return m.trendingView() // render trending view
	case SearchView:
		return m.searchView() // render search view
	case ViewingList:
		return m.viewingListView() // render viewing list view
	case ViewingInfo:
		return m.viewingInfoView() // render viewing info view
	default:
		return "Unknown view" // handle unexpected states
	}
}

// Menu view rendering logic
func (m model) menuView() string {
	// Title with gradient text
	title := gradientText("❒ block-vision", gradientStart, gradientEnd)

	// Subtitle
	subtext := lipgloss.NewStyle().
		Foreground(subtextColor).
		Render("• v1.0.0  ✦︎ Lighting up your crypto journey ✦︎")

	// Menu options
	checkboxOptions := checkboxPicker([]string{"Trending", "Search"}, m.menuChoice)

	// Help menu
	helpView := m.help.View(m.keys)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtext,
		"", // new line
		checkboxOptions,
		"", // new line
		helpView,
	)
}

// Trending view rendering logic
func (m model) trendingView() string {
	// Display error message if present
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	// Render the table and help menu
	helpView := m.help.View(m.keys)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		baseStyle.Render(m.table.View()),
		"", // new line
		helpView,
	)
}

// Search view rendering logic
func (m model) searchView() string {
	return fmt.Sprintf(
		"Search for a cryptocurrency:\n\n%s\n\n%s",
		m.textInput.View(),
		"(Press Enter to search, Esc to quit)",
	)
}

// Viewing list view rendering logic
func (m model) viewingListView() string {
	return docStyle.Render(m.list.View())
}

// Viewing info view rendering logic
func (m model) viewingInfoView() string {
	// Parse and format price
	price := m.coinDetails.MarketData.CurrentPrice["usd"]
	priceFormatted := humanize.Commaf(price)

	// Parse sentiment values
	upPercentage := fmt.Sprintf("%.2f%%", m.coinDetails.SentimentVotesUpPercentage)
	downPercentage := fmt.Sprintf("%.2f%%", m.coinDetails.SentimentVotesDownPercentage)

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(subtextColor).Render(fmt.Sprintf("%s (%s)", m.coinDetails.Name, strings.ToUpper(m.coinDetails.Symbol))),
		"", // new line
		fmt.Sprintf("• $%s", priceFormatted),
		fmt.Sprintf(
			"• Sentiment: %s ↑ %s ↓",
			lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(upPercentage),   // green
			lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(downPercentage), // red
		),
		fmt.Sprintf(
			"• Learn More ↗: %s",
			lipgloss.NewStyle().Foreground(subtleColor).Render(m.coinDetails.Links.Whitepaper),
		),
		"", // new line
		lipgloss.NewStyle().Foreground(subtleColor).Render("(press esc to go back)"),
	)
}

// Message type to handle rows of trending data
type trendingDataMsg struct {
	rows []table.Row
}

// Message type to handle errors during data fetching
type errMsg struct {
	err error
}

// Command to fetch trending data asynchronously
func fetchTrendingCmd() tea.Msg {
	rows, err := fetchTrendingCryptos()
	if err != nil {
		return errMsg{err}
	}
	return trendingDataMsg{rows}
}

// Renders text with a gradient effect (uses go-colorful library)
func gradientText(text string, start, end colorful.Color) string {
	gradient := make([]string, len(text))
	for i, char := range text {
		// Compute gradient color
		color := start.BlendLuv(end, float64(i)/float64(len(text)))

		// Apply current computed color
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color.Hex()))
		gradient[i] = style.Render(string(char))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, gradient...)
}

// Renders a checkbox menu with a visual indicator for the selected item
func checkboxPicker(options []string, choice int) string {
	var renderedOptions []string
	for i, option := range options {
		// Default checkbox state
		checked := "[ ]"

		// Checked state
		if i == choice {
			checked = fmt.Sprintf("[%s]", lipgloss.NewStyle().Foreground(checkboxChecked).Render("x"))
		}

		// Option styling
		optionText := lipgloss.NewStyle().Foreground(checkboxColor).Render(option)
		renderedOptions = append(renderedOptions, fmt.Sprintf("%s %s", checked, optionText))
	}
	return lipgloss.JoinVertical(lipgloss.Left, renderedOptions...)
}

// Creates and initializes the application model with default values
func initialModel() model {
	// Initialize the text input for cryptocurrency search
	ti := textinput.New()
	ti.Placeholder = "Search for a cryptocurrency"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	// Initialize the list to display search results
	ls := list.New([]list.Item{}, list.NewDefaultDelegate(), 30, 20)
	ls.Title = "Search Results"

	return model{
		viewMode:  MenuView,                  // start in the menu view
		keys:      keys,                      // initialize key bindings
		help:      help.New(),                // initialize help model
		table:     InitializeTrendingTable(), // set up the trending data table
		textInput: ti,                        // text input for search
		list:      ls,                        // list for displaying search results
	}
}

// Main function to run the application
func main() {
	// Enable alternate screen mode (fullscreen) for clean UI
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	// Run the program and handle errors
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
