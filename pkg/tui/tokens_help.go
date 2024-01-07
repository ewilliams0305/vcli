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
type tokensKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Quit   key.Binding
	Create key.Binding
	Delete key.Binding
	Edit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k tokensKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Up, k.Down, k.Create, k.Edit, k.Delete}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k tokensKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},               // first column
		{k.Create, k.Edit, k.Delete}, // second column
	}
}

var tokenKeys = tokensKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),

	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Create: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "create room"),
	),
	Edit: key.NewBinding(
		key.WithKeys("ctrl+e", "enter"),
		key.WithHelp("enter", "edit room"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete"),
		key.WithHelp("delete", "delete room"),
	),
}

type TokensHelpModel struct {
	keys       tokensKeyMap
	help       help.Model
	inputStyle lipgloss.Style
}

func NewtokensHelpModel() TokensHelpModel {
	return TokensHelpModel{
		keys:       tokenKeys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (m TokensHelpModel) Init() tea.Cmd {
	return nil
}

func (m TokensHelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	}
	return m, nil
}

func (m TokensHelpModel) View() string {
	var status string

	helpView := m.help.View(m.keys)
	height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")

	return "\n" + status + strings.Repeat("\n", height) + helpView
}

func (m TokensHelpModel) renderHelpInfo() string {
	helpView := m.help.View(m.keys)
	return "\n" + helpView
}
