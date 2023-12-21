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
type serviceKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Help    key.Binding
	Quit    key.Binding
	Start   key.Binding
	Stop    key.Binding
	Restart key.Binding
	Logs    key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k serviceKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Up, k.Down, k.Start, k.Stop, k.Restart, k.Logs}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k serviceKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Up, k.Down},    // first column
		{k.Start, k.Stop, k.Logs}, // second column
	}
}

var serviceKeys = serviceKeyMap{
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
	Stop: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "stop service"),
	),
	Start: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "start service"),
	),
	Restart: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "restart room"),
	),
	Logs: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "display logs"),
	),
}

type SystemsHelpModel struct {
	keys       serviceKeyMap
	help       help.Model
	inputStyle lipgloss.Style
}

func NewSystensHelpModel() SystemsHelpModel {
	return SystemsHelpModel{
		keys:       serviceKeys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (m SystemsHelpModel) Init() tea.Cmd {
	return nil
}

func (m SystemsHelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return VirtualControlServiceModel{}, nil
		}
	}
	return m, nil
}

func (m SystemsHelpModel) View() string {
	var status string

	helpView := m.help.View(m.keys)
	height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")

	return "\n" + status + strings.Repeat("\n", height) + helpView
}

func (m SystemsHelpModel) renderHelpInfo() string {
	helpView := m.help.View(m.keys)
	return "\n" + helpView
}
