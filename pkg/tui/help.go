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
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	GoHelp   key.Binding
	Help     key.Binding
	Quit     key.Binding
	Programs key.Binding
	Rooms    key.Binding
	Devices  key.Binding
	Info     key.Binding
	Auth     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Programs, k.Rooms, k.Devices, k.Auth, k.Info}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},          // first column
		{k.Help, k.Quit, k.Info},                 // second column
		{k.Rooms, k.Programs, k.Devices, k.Auth}, // second column
	}
}

var Keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	GoHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Show help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+q", "esc", "ctrl+c"),
		key.WithHelp("ctrl+q", "return"),
	),
	Programs: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("ctrl+p", "programs"),
	),
	Rooms: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "rooms"),
	),
	Devices: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "devices"),
	),
	Auth: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("ctrl+a", "auth"),
	),
	Info: key.NewBinding(
		key.WithKeys("ctrl+i"),
		key.WithHelp("ctrl+i", "info"),
	),
}

type HelpModel struct {
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	quitting   bool
}

func NewHelpModel() HelpModel {
	return HelpModel{
		keys:       Keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.lastKey = "↑\n\nMoves the menu up one item"
		case key.Matches(msg, m.keys.Down):
			m.lastKey = "↓\n\nMoves the menu down one item"
		case key.Matches(msg, m.keys.Left):
			m.lastKey = "←\n\nMoves the menu to the left"
		case key.Matches(msg, m.keys.Right):
			m.lastKey = "→\n\nMoves the menu to the right"

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return ReturnToHomeModel(helpState), tea.Batch(tick, DeviceInfoCommand)
		}

		switch msg.String() {
		case "p", "P":
			m.lastKey = "ctrl+p\n\nView and manage loaded program files"
		case "r", "R":
			m.lastKey = "ctrl+r\n\nView and manage active rooms"
		case "d", "D":
			m.lastKey = "ctrl+d\n\nView device maps and communication status"
		case "a", "A":
			m.lastKey = "ctrl+a\n\nCreate and access API tokens"
		case "i", "I":
			m.lastKey = "ctrl+i\n\nRefresh and view device information"
		}
	}

	return m, nil
}

func (m HelpModel) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var status string
	if m.lastKey == "" {
		status = "Enter key below for extended help information..."
	} else {
		status = "You chose: " + m.inputStyle.Render(m.lastKey)
	}

	helpView := m.help.View(m.keys)
	height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")

	return "\n" + status + strings.Repeat("\n", height) + helpView
}

func (m HelpModel) renderHelpInfo() string {

	helpView := m.help.View(m.keys)
	return "\n" + helpView
}
