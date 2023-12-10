package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type roomsKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Help    key.Binding
	Quit    key.Binding
	Start   key.Binding
	Stop    key.Binding
	Restart key.Binding
	Create  key.Binding
	Delete  key.Binding
	Edit    key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k roomsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Start, k.Stop, k.Delete}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k roomsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Up, k.Down},      // first column
		{k.Start, k.Stop, k.Delete}, // second column
	}
}

var roomKeys = roomsKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?", "h"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Start: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "star room"),
	),
	Stop: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "stop room"),
	),
}

type RoomsHelpModel struct {
	keys       roomsKeyMap
	help       help.Model
	inputStyle lipgloss.Style
}

func NewRoomsHelpModel() RoomsHelpModel {
	return RoomsHelpModel{
		keys:       roomKeys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (m RoomsHelpModel) Init() tea.Cmd {
	return nil
}

func (m RoomsHelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		// case key.Matches(msg, m.keys.Up):
		// 	m.lastKey = "↑\n\nMoves the menu up one item"
		// case key.Matches(msg, m.keys.Down):
		// 	m.lastKey = "↓\n\nMoves the menu down one item"
		// case key.Matches(msg, m.keys.Left):
		// 	m.lastKey = "←\n\nMoves the menu to the left"
		// case key.Matches(msg, m.keys.Right):
		// 	m.lastKey = "→\n\nMoves the menu to the right"

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return InitialModel(), DeviceInfoCommand
		}

		// switch msg.String() {
		// case "p", "P":
		// 	m.lastKey = "ctrl+p\n\nView and manage loaded program files"
		// case "r", "R":
		// 	m.lastKey = "ctrl+r\n\nView and manage active rooms"
		// case "d", "D":
		// 	m.lastKey = "ctrl+d\n\nView device maps and communication status"
		// case "a", "A":
		// 	m.lastKey = "ctrl+a\n\nCreate and access API tokens"
		// case "i", "I":
		// 	m.lastKey = "ctrl+i\n\nRefresh and view device information"
		// }
	}

	return m, nil
}

func (m RoomsHelpModel) View() string {
	var status string
	// if m.lastKey == "" {
	// 	status = "Enter key below for extended help information..."
	// } else {
	// 	status = "You chose: " + m.inputStyle.Render(m.lastKey)
	// }

	helpView := m.help.View(m.keys)
	height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")

	return "\n" + status + strings.Repeat("\n", height) + helpView
}

func (m RoomsHelpModel) renderHelpInfo() string {
	helpView := m.help.View(m.keys)
	return "\n" + helpView
}
