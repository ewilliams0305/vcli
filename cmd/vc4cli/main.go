package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ewilliams0305/VC4-CLI/cmd/vc4cli/views"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc4"
)

type errMsg struct{ err error }

type MainModel struct {
	device   vc4.DeviceInfo
	err      string
	actions  []string
	cursor   int
	selected map[int]struct{}
}

func InitialModel() MainModel {
	return MainModel{
		device:   vc4.DeviceInfo{},
		actions:  []string{"Manage Programs", "Manage Rooms", "View Logs"},
		selected: make(map[int]struct{}),
	}
}

func (m MainModel) Init() tea.Cmd {
	return deviceInfoCommand
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case int:
		m.err = "GOT YOU"

	case *vc4.DeviceInfo:
		m.device = *msg

	case vc4.DeviceInfo:
		m.device = msg

	case errMsg:
		m.err = msg.err.Error()

	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.actions)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

			return NewHelpModel(), nil
		case "i":
			return NewDeviceTable(m.device), nil
		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m MainModel) View() string {
	// The header
	s := views.Logo + "\n\n"
	info := NewDeviceTable(m.device)
	s += baseStyle.Render(info.table.View()) + "\n"

	// Iterate over our choices
	for i, choice := range m.actions {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"
	s += "\n" + m.err

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
