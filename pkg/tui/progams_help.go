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
type programsKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Help   key.Binding
	Quit   key.Binding
	New    key.Binding
	Delete key.Binding
	Edit   key.Binding
	Room   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k programsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Up, k.Down, k.New, k.Delete, k.Edit, k.Room}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k programsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Up, k.Down},            // first column
		{k.New, k.Delete, k.Edit, k.Room}, // second column
	}
}

var programKeys = programsKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?", "h"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	New: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "new program"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete"),
	),
	Edit: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "edit"),
	),
	Room: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "create room"),
	),
}

type programsHelpModel struct {
	keys       programsKeyMap
	help       help.Model
	inputStyle lipgloss.Style
}

func NewProgramsHelpModel() programsHelpModel {
	return programsHelpModel{
		keys:       programKeys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (m programsHelpModel) Init() tea.Cmd {
	return nil
}

func (m programsHelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return InitialProgramsModel(200, 200), ProgramsQuery
		}
	}

	return m, nil
}

func (m programsHelpModel) View() string {
	var status string

	helpView := m.help.View(m.keys)
	height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")

	return "\n" + status + strings.Repeat("\n", height) + helpView
}

func (m programsHelpModel) renderHelpInfo() string {
	helpView := m.help.View(m.keys)
	return "\n" + helpView
}
