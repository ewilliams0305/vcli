package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	vc "github.com/ewilliams0305/VC4-CLI/vc"
)

type RoomsTableModel struct {
	table        table.Model
	rooms        vc.Rooms
	selectedRoom vc.Room
	err          error
	help         RoomsHelpModel
	busy         busy
	cursor       int
	width        int
	height       int
}

func InitialRoomsModel() *RoomsTableModel {
	return &RoomsTableModel{
		rooms:        vc.Rooms{},
		selectedRoom: vc.Room{},
		cursor:       0,
		help:         NewRoomsHelpModel(),
	}
}

func BusyRoomsModel(b busy, rooms vc.Rooms) *RoomsTableModel {
	return &RoomsTableModel{
		busy:         b,
		rooms:        rooms,
		selectedRoom: vc.Room{},
		cursor:       0,
	}
}

func (m RoomsTableModel) Init() tea.Cmd {
	return DeviceInfoCommand
}

func (m RoomsTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case int:
		m.cursor = msg
		m.selectedRoom = m.rooms[msg]
		return m, nil

	case busy:
		m.busy = msg
		return m, RefreshRoomData

	case vc.Rooms:
		m.busy = busy{flag: false}
		m.rooms = msg
		m.table = newRoomsTable(msg, m.cursor)
		m.selectedRoom = msg[m.cursor]
		return m, nil
		//return NewRoomsTable(msg, m.cursor), nil

	case error:
		m.err = msg
		return NewRoomsErrorTable(msg), nil

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c", "esc":
			return InitialModel(), DeviceInfoCommand
		case "down":
			m.table.SetCursor(m.table.Cursor() + 1)
			return m, cmdCursor(m.table.Cursor())
		case "up":
			m.table.SetCursor(m.table.Cursor() - 1)
			return m, cmdCursor(m.table.Cursor())

		case "ctrl+s":
			if m.selectedRoom.Status == string(vc.Running) {
				return m, cmdRoomStop(m.selectedRoom.ID)
			} else if m.selectedRoom.Status == string(vc.Starting) {
				return m, cmdRoomStop(m.selectedRoom.ID)
			} else if m.selectedRoom.Status == string(vc.Stopped) {
				return m, cmdRoomStart(m.selectedRoom.ID)
			} else if m.selectedRoom.Status == string(vc.Stopping) {
				return m, cmdRoomStart(m.selectedRoom.ID)
			} else if m.selectedRoom.Status == string(vc.Aborted) {
				return m, cmdRoomStart(m.selectedRoom.ID)
			}

		case "ctrl+r":
			return m, cmdRoomRestart(m.selectedRoom.ID)

		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m RoomsTableModel) View() string {
	s := DisplayLogo(m.width)
	s += BaseStyle.Render(m.table.View()) + "\n\n"

	if m.busy.flag {
		s += RenderMessageBox(m.width).Render(m.busy.message)
	} else {
		room := fmt.Sprintf("❓ %s %s\n", m.selectedRoom.ID, m.selectedRoom.CompileDateTime)
		s += RenderMessageBox(m.width).Render(room)
	}

	s += m.help.renderHelpInfo()
	return s
}

// func NewRoomsTable(rooms vc.Rooms, cursor int) RoomsTableModel {
// 	columns := []table.Column{
// 		{Title: "ID", Width: 20},
// 		{Title: "NAME", Width: 20},
// 		{Title: "PROGRAM", Width: 30},
// 		{Title: "NOTES", Width: 30},
// 		{Title: "TYPE", Width: 8},
// 		{Title: "STATUS", Width: 8},
// 		{Title: "DEBUG", Width: 8},
// 	}

// 	rows := []table.Row{}

// 	for _, room := range rooms {
// 		rows = append(rows, table.Row{room.ID, room.Name, room.ProgramName, room.Notes, room.ProgramType, GetStatus(room.Status), CheckMark(room.Debugging)})
// 	}

// 	t := table.New(
// 		table.WithColumns(columns),
// 		table.WithRows(rows),
// 		table.WithFocused(false),
// 		table.WithHeight(9),
// 	)

// 	s := table.DefaultStyles()
// 	s.Header = s.Header.
// 		BorderStyle(lipgloss.NormalBorder()).
// 		BorderForeground(lipgloss.Color("240")).
// 		BorderBottom(true).
// 		Background(lipgloss.Color(AccentColor)).
// 		Foreground(lipgloss.Color(AccentColor)).
// 		Bold(true)
// 	s.Selected = s.Selected.
// 		Foreground(lipgloss.Color("229")).
// 		Background(lipgloss.Color(AccentColor)).
// 		Bold(false)
// 	t.SetStyles(s)

// 	t.SetCursor(cursor)

// 	return RoomsTableModel{
// 		table:        t,
// 		rooms:        rooms,
// 		selectedRoom: rooms[cursor],
// 		help:         NewRoomsHelpModel(),
// 	}
// }

func newRoomsTable(rooms vc.Rooms, cursor int) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "NAME", Width: 20},
		{Title: "PROGRAM", Width: 30},
		{Title: "NOTES", Width: 30},
		{Title: "TYPE", Width: 8},
		{Title: "STATUS", Width: 8},
		{Title: "DEBUG", Width: 8},
	}

	rows := []table.Row{}

	for _, room := range rooms {
		rows = append(rows, table.Row{room.ID, room.Name, room.ProgramName, room.Notes, room.ProgramType, GetStatus(room.Status), CheckMark(room.Debugging)})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(9),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Foreground(lipgloss.Color(AccentColor)).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(PrimaryLight)).
		Background(lipgloss.Color(PrimaryDark)).
		Bold(false)
	t.SetStyles(s)

	t.SetCursor(cursor)

	return t
}

func NewRoomsErrorTable(msg vc.VirtualControlError) RoomsTableModel {
	columns := []table.Column{
		{Title: "SERVER ERROR", Width: 20},
		{Title: "", Width: 100},
	}

	rows := []table.Row{
		{"ERROR", msg.Error()},
		{"", ""},
		{"MESSAGE", "There was an error connecting to the VC4 service"},
		{"", "Please verify your IP address and token"},
		{"", "Please veriify the virtualcontrol service is enabled and running."},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("#FF0000")).
		Bold(true).Italic(true)
	t.SetStyles(s)

	return RoomsTableModel{table: t, help: NewRoomsHelpModel()}
}

func RoomCommand() tea.Msg {

	info, err := server.GetRooms()
	if err != nil {
		return err
	}
	return info
}

func RefreshRoomData() tea.Msg {

	var rooms vc.Rooms
	var err error
	_, err = server.GetRooms()
	if err != nil {
		return err
	}

	//for i := 0; i < 3; i++ {
	time.Sleep(3 * time.Second)
	rooms, err = server.GetRooms()
	if err != nil {
		return err
	}
	return rooms
}

func RoomRestart(id string) tea.Msg {

	_, err := server.RestartRoom(id)
	if err != nil {
		return err
	}
	return busy{flag: true, message: fmt.Sprintf("restarting room %s, please wait...", id)}
}
func RoomStop(id string) tea.Msg {

	_, err := server.StopRoom(id)
	if err != nil {
		return err
	}
	return busy{flag: true, message: fmt.Sprintf("stopping room %s, please wait...", id)}
}
func RoomStart(id string) tea.Msg {

	_, err := server.StartRoom(id)
	if err != nil {
		return err
	}
	return busy{flag: true, message: fmt.Sprintf("starting room %s, please wait...", id)}
}

func cmdRoomStop(id string) tea.Cmd {
	return func() tea.Msg {
		return RoomStop(id)
	}
}

func cmdRoomStart(id string) tea.Cmd {
	return func() tea.Msg {
		return RoomStart(id)
	}
}

type busy struct {
	flag    bool
	message string
}

func cmdRoomRestart(id string) tea.Cmd {
	return func() tea.Msg {
		return RoomRestart(id)
	}
}

func cmdCursor(cursor int) tea.Cmd {
	return func() tea.Msg {
		return cursor
	}
}
