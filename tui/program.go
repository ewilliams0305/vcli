package tui

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"

	vc "github.com/ewilliams0305/VC4-CLI/vc"
)

var server vc.VirtualControl

type appState int

const (
	initializing appState = iota
	// HOME VIEW, Displays dev info and main menu with text input for commands
	// ALL otger const states will make our main menu
	home appState = 10
	// PROGAM VIEW, displays all programs loaded
	programs appState = 1
	// ROOM VIEW, display all program instances
	rooms appState = 2
	// INFO VIEW, displays all hardware and system information
	info appState = 3
	// DEVICES VIEW, displays all the device IP Tables and maps
	devices appState = 4
	// AUTH VIEW, displays all auth and api tokens
	auth appState = 5
	// HELP VIEW
	helpState appState = 6
)

var program *MainModel

type MainModel struct {
	state         appState
	device        vc.DeviceInfo
	err           error
	actions       []string
	cursor        int
	help          HelpModel
	width, height int
}

func InitialModel() MainModel {

	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	program = &MainModel{
		device:  vc.DeviceInfo{},
		actions: []string{"Refresh", "Manage Programs", "Manage Rooms", "Device Information", "Devices", "Authorization", "Help"},
		help:    NewHelpModel(),
		width:   w,
		height:  h,
	}

	return *program
}

func ReturnToHomeModel(state appState) MainModel {

	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	program.width = w
	program.height = h
	program.state = state
	program.cursor = int(state)

	return *program
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch(tick, DeviceInfoCommand)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		w, h, _ := term.GetSize(int(os.Stdout.Fd()))
		if w != m.width || h != m.height {
			m.updateSize(w, h)
		}
		return m, tea.Batch(tick, func() tea.Msg { return tea.WindowSizeMsg{Width: w, Height: h} })

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.state = home
		return m, nil

	case vc.DeviceInfo:
		m.device = msg
		m.state = home
		return m, nil
	case error:
		m.err = msg

	case tea.KeyMsg:

		//TODO: Change these to match the keys in the help.go file

		// THE MESSAGE IS A KEYPRESS
		switch msg.String() {

		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.actions)-1 {
				m.cursor++
			}

		case "?", "":
			return NewHelpModel(), nil

		case "enter", " ":
			return arrowSelected(&m)

		case "i":
			m.state = info
			return NewDeviceInfo(m.width, m.height), DeviceInfoCommand

		case "r", "ctrl+r":
			m.state = rooms
			return InitialRoomsModel(m.width, m.height), RoomCommand

		case "p":
			m.state = programs
			return NewDeviceInfo(m.width, m.height), DeviceInfoCommand
		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m MainModel) View() string {
	var s string

	if m.state == initializing {
		s = RenderMessageBox(1000).Render("Initializing application...")
		return s
	}

	s = DisplayLogo(m.width)

	if m.err != nil {
		info := NewDeviceErrorTable(m.err, m.width)
		s += BaseStyle.Render(info.Table.View()) + "\n"
	} else {
		info := HomeDeviceInfo(m.device, m.width)
		s += BaseStyle.Render(info.Table.View()) + "\n"
	}

	// Iterate over our choices
	for i, choice := range m.actions {

		// Is the cursor pointing at this choice?
		cursor := GreyedOutText.Render("  "+" "+choice) + "\n"
		if m.cursor == i {
			cursor = HighlightedText.Render("\u2192"+"  "+choice) + "\n"
		}
		s += cursor

	}

	// The footer
	s += m.help.renderHelpInfo()
	return s
}

func arrowSelected(m *MainModel) (tea.Model, tea.Cmd) {

	switch m.cursor {
	case int(programs):
		m.state = programs
		return NewDeviceInfo(m.width, m.height), DeviceInfoCommand

	case int(rooms):
		m.state = rooms
		return InitialRoomsModel(m.width, m.height), RoomCommand

	case int(info):
		m.state = info
		return NewDeviceInfo(m.width, m.height), DeviceInfoCommand

	case int(auth):
	case int(devices):
	case int(helpState):
		return NewHelpModel(), nil

	}

	return program, nil
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

/********************************************************
*
* VIEWPORT BUGGY WINDOW SIZE HACK
*
*********************************************************/

// Pointless type to trigger the update function
type tickMsg int

// Updates the entire view if the size changed
func (m *MainModel) updateSize(w, h int) {
	m.width = w
	m.width = h

	//m.View()
}

// Sends a message back to the update function to start the tick over again.
func tick() tea.Msg {
	time.Sleep(time.Second + 1)
	return tickMsg(1)
}
