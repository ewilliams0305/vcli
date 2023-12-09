package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	vc "github.com/ewilliams0305/VC4-CLI/vc"
)

var server vc.VirtualControl

type appState int

const (
	// HOME VIEW, Displays dev info and main menu with text input for commands
	// ALL otger const states will make our main menu
	home appState = 0
	// PROGAM VIEW, displays all programs loaded
	programs appState = 1
	// ROOM VIEW, display all program instances
	rooms appState = 2
	// INFO VIEW, displays all hardware and system information
	info appState = 3
	// DEVICES VIEW, displays all the device IP Tables and maps
	devices appState = 4
	// AUTH VIEW, displays all auth and api tokens
	auth appState = 4
)

type MainModel struct {
	state    appState
	device   vc.DeviceInfo
	err      error
	actions  []string
	cursor   int
	selected map[int]struct{}
	help     HelpModel
}

func InitialModel() MainModel {
	return MainModel{
		device:   vc.DeviceInfo{},
		actions:  []string{"Manage Programs", "Manage Rooms", "View Logs"},
		selected: make(map[int]struct{}),
		help:     NewHelpModel(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return DeviceInfoCommand
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case *vc.DeviceInfo:
		m.device = *msg

	case vc.DeviceInfo:
		m.device = msg

	case error:
		m.err = msg

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
	s := Logo + "\n\n"

	if m.err != nil {
		info := NewDeviceErrorTable(m.err)
		s += BaseStyle.Render(info.Table.View()) + "\n"
	} else {
		info := NewDeviceTable(m.device)
		s += BaseStyle.Render(info.Table.View()) + "\n"
	}

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
	s += m.help.renderHelpInfo()
	return s
}

func Run() {

	server = initServer()

	// TODO: Process addtional flags to send instant actions to the device.
	p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("VC4 CLI failed to start, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initServer() vc.VirtualControl {
	// TODO: Possible create an error if token and host are invlid
	if (len(Hostname) > 0) && (len(Token) > 0) {
		return vc.NewRemoteVC(Hostname, Token)
	}
	return vc.NewLocalVC()
}
