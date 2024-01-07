package tui

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"

	vc "github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

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
	// SYSTEM SERVICE VIEW, displays logs and service status
	systemd appState = 6
	// HELP VIEW
	helpState appState = 7
)

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
	app = &MainModel{
		device:  vc.DeviceInfo{},
		actions: []string{"Refresh", "Manage Programs", "Manage Rooms", "Device Information", "Devices", "Authorization", "System Service", "Help"},
		help:    NewHelpModel(),
		width:   w,
		height:  h,
	}

	return *app
}

func ReturnToHomeModel(state appState) MainModel {

	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	app.width = w
	app.height = h
	app.state = state
	app.cursor = int(state)

	return *app
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

		case "i", "ctrl+i":
			m.state = info
			return NewDeviceInfo(m.width, m.height), DeviceInfoCommand

		case "r", "ctrl+r":
			m.state = rooms
			return InitialRoomsModel(m.width, m.height), RoomsQuery

		case "p", "ctrl+p":
			m.state = programs
			return InitialProgramsModel(m.width, m.height), DeviceInfoCommand

		case "s", "ctrl+s":
			m.state = programs
			return InitialSystemModel(), nil
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

	for i, choice := range m.actions {
		cursor := GreyedOutText.Render("  "+" "+choice) + "\n"
		if m.cursor == i {
			cursor = HighlightedText.Render("\u2192"+"  "+choice) + "\n"
		}
		s += cursor
	}
	s += m.help.renderHelpInfo()
	return s
}

func arrowSelected(m *MainModel) (tea.Model, tea.Cmd) {

	switch m.cursor {
	case 0:
		return m, DeviceInfoCommand
	case int(programs):
		m.state = programs
		return InitialProgramsModel(m.width, m.height), ProgramsQuery
	case int(rooms):
		m.state = rooms
		return InitialRoomsModel(m.width, m.height), RoomsQuery
	case int(info):
		m.state = info
		return NewDeviceInfo(m.width, m.height), DeviceInfoCommand
	case int(auth):
	case int(devices):
	case int(systemd):
		return InitialSystemModel(), nil
	case int(helpState):
		return NewHelpModel(), nil

	}

	return app, nil
}

/********************************************************
*
* INITIALIZE THE APP WITH FLAGS
*
*********************************************************/

func Run() {

	server = initServer()
	initialView, err := initActions()
	if err != nil {
		fmt.Printf("VC4 CLI failed execute intial actions, there's been an error: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialView, tea.WithAltScreen())
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

func initActions() (tea.Model, error) {

	if len(RoomID) > 0 && len(ProgramFile) > 0 && len(ProgramName) > 0 {
		return InitialActionModel(fmt.Sprintf("Uploading and creating new room %s", RoomID), loadAndCreate), nil
	}

	if len(ProgramFile) > 0 && len(ProgramName) > 0 {
		return InitialActionModel(fmt.Sprintf("Loading new program %s", ProgramFile), loadProgram), nil
	}
	return InitialModel(), nil
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
