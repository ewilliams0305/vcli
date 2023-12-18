package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	vc "github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

var roomsModel *RoomsTableModel

type RoomsTableModel struct {
	table         table.Model
	rooms         vc.Rooms
	selectedRoom  vc.Room
	err           error
	help          RoomsHelpModel
	busy          busy
	cursor        int
	width, height int
}

func InitialRoomsModel(width, height int) *RoomsTableModel {
	roomsModel = &RoomsTableModel{
		rooms:        vc.Rooms{},
		selectedRoom: vc.Room{},
		cursor:       0,
		help:         NewRoomsHelpModel(),
		width:        width,
		height:       height,
	}
	return roomsModel
}

func BusyRoomsModel(b busy, rooms vc.Rooms) *RoomsTableModel {
	roomsModel.busy = b
	roomsModel.rooms = rooms
	roomsModel.selectedRoom = vc.Room{}
	roomsModel.cursor = 0
	return roomsModel
}

func (m RoomsTableModel) Init() tea.Cmd {
	return RoomCommand
}

func (m RoomsTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tickMsg:
		w, h, _ := term.GetSize(int(os.Stdout.Fd()))
		if w != m.width || h != m.height {
			roomsModel.width = w
			roomsModel.height = h
		}
		return roomsModel, tea.Batch(RoomCommand, tick)

	case tea.WindowSizeMsg:
		roomsModel.width = msg.Width
		roomsModel.height = msg.Height
		roomsModel.table, cmd = roomsModel.table.Update(msg)
		return roomsModel, cmd

	case int:
		roomsModel.cursor = msg
		roomsModel.selectedRoom = roomsModel.rooms[msg]
		return roomsModel, nil

	case busy:
		roomsModel.busy = msg
		return roomsModel, nil

	case vc.Rooms:
		if len(msg) > 0 {
			roomsModel.busy = busy{flag: false}
			roomsModel.rooms = msg
			roomsModel.table = newRoomsTable(msg, roomsModel.cursor, roomsModel.width)
			roomsModel.selectedRoom = msg[m.cursor]
			return roomsModel, nil
		}
		return NewRoomsErrorTable(fmt.Errorf("THERE ARE NO ROOMS LOADED TO THE SYSTEM, PRESS CTRL+N TO CREATE NEW ROOM")), nil

	case error:
		roomsModel.err = msg
		return NewRoomsErrorTable(msg), nil

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c", "esc":
			return ReturnToHomeModel(rooms), tea.Batch(tick, DeviceInfoCommand)
		case "down":
			if roomsModel.err == nil {
				roomsModel.table.SetCursor(roomsModel.table.Cursor() + 1)
				return roomsModel, cmdCursor(roomsModel.table.Cursor())
			}
		case "up":
			if roomsModel.err == nil {
				roomsModel.table.SetCursor(roomsModel.table.Cursor() - 1)
				return roomsModel, cmdCursor(roomsModel.table.Cursor())
			}

		case "ctrl+s":
			if roomsModel.err == nil {
				if m.selectedRoom.Status == string(vc.Running) {
					return roomsModel, cmdRoomStop(roomsModel.selectedRoom.ID)
				} else if roomsModel.selectedRoom.Status == string(vc.Starting) {
					return roomsModel, cmdRoomStop(roomsModel.selectedRoom.ID)
				} else if roomsModel.selectedRoom.Status == string(vc.Stopped) {
					return roomsModel, cmdRoomStart(roomsModel.selectedRoom.ID)
				} else if roomsModel.selectedRoom.Status == string(vc.Stopping) {
					return roomsModel, cmdRoomStart(roomsModel.selectedRoom.ID)
				} else if roomsModel.selectedRoom.Status == string(vc.Aborted) {
					return roomsModel, cmdRoomStart(roomsModel.selectedRoom.ID)
				}
			}

		case "ctrl+r":
			if roomsModel.err == nil {
				return roomsModel, cmdRoomRestart(roomsModel.selectedRoom.ID)
			}

		case "ctrl+d":
			if roomsModel.err == nil {
				return roomsModel, cmdRoomDebug(roomsModel.selectedRoom.ID, !roomsModel.selectedRoom.Debugging)
			}

		case "ctrl+n":
			if roomsModel.err == nil {
				// Query the programs to populate the option list for the programs!
				return NewRoomFormModel(), ProgramsQuery
			}

		case "ctrl+e", "enter":
			if roomsModel.err == nil {
				// Query the programs to populate the option list for the programs!
				return NewRoomFormModel(), ProgramsQuery
			}
		}
	}
	roomsModel.table, cmd = roomsModel.table.Update(msg)
	return roomsModel, cmd
}

func (m RoomsTableModel) View() string {
	s := DisplayLogo(m.width)
	s += BaseStyle.Render(m.table.View()) + "\n\n"

	if m.busy.flag {
		s += RenderMessageBox(m.width).Render(m.busy.message)
	} else {
		room := fmt.Sprintf("\u2192 use keyboard actions to manage %s %s (ctrl+s, ctrl+d...)\n", m.selectedRoom.ID, m.selectedRoom.ProgramName)
		s += RenderMessageBox(m.width).Render(room)
	}

	s += m.help.renderHelpInfo()
	return s
}

func newRoomsTable(rooms vc.Rooms, cursor int, width int) table.Model {

	columns := getRoomsColumns(width)
	rows := getRoomsRows(width, cursor, rooms)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(9),
		table.WithWidth(width),
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

func getRoomsColumns(width int) []table.Column {

	if width < 120 {

		return []table.Column{
			{Title: "", Width: 1},
			{Title: "ID", Width: 20},
			{Title: "NAME", Width: width - 49},
			{Title: "STATUS", Width: 8},
			{Title: "DEBUG", Width: 8},
		}
	}
	return []table.Column{
		{Title: "", Width: 1},
		{Title: "ID", Width: 20},
		{Title: "NAME", Width: 35},
		{Title: "PROGRAM", Width: 35},
		{Title: "NOTES", Width: width - 141},
		{Title: "TYPE", Width: 16},
		{Title: "STATUS", Width: 8},
		{Title: "DEBUG", Width: 8},
	}
}

func getRoomsRows(width int, cursor int, rooms vc.Rooms) []table.Row {
	rows := []table.Row{}
	small := width < 120

	for i, room := range rooms {
		marker := ""
		if cursor == i {
			marker = "\u2192"
		}
		if small {
			rows = append(rows, table.Row{marker, room.ID, room.Name, GetStatus(room.Status), CheckMark(room.Debugging)})
		} else {
			rows = append(rows, table.Row{marker, room.ID, room.Name, room.ProgramName, room.Notes, room.ProgramType, GetStatus(room.Status), CheckMark(room.Debugging)})
		}
	}
	return rows
}

func NewRoomsErrorTable(msg vc.VirtualControlError) *RoomsTableModel {
	columns := []table.Column{
		{Title: "SERVER ERROR", Width: 50},
		{Title: "", Width: roomsModel.width - 56},
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

	roomsModel.table = t
	roomsModel.err = msg

	return roomsModel
}

func RoomCommand() tea.Msg {

	info, err := server.GetRooms()
	if err != nil {
		return err
	}
	return info
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

func RoomDebug(id string, enable bool) tea.Msg {

	_, err := server.DebugRoom(id, enable)
	if err != nil {
		return err
	}
	if enable {
		return busy{flag: true, message: fmt.Sprintf("enable debugging on room %s, please wait...", id)}
	}
	return busy{flag: true, message: fmt.Sprintf("disabling debugging on room %s, please wait...", id)}
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

func cmdRoomDebug(id string, enable bool) tea.Cmd {
	return func() tea.Msg {
		return RoomDebug(id, enable)
	}
}

func cmdRoomRestart(id string) tea.Cmd {
	return func() tea.Msg {
		// VC4 DOES NOT WORK WHEN RESTART IS ISSUED => HACK
		RoomStop(id)
		time.Sleep(time.Second + 3)

		return RoomStart(id)
	}
}

func cmdCursor(cursor int) tea.Cmd {
	return func() tea.Msg {
		return cursor
	}
}

func CreateRoom(options *vc.RoomOptions) tea.Cmd {

	return func() tea.Msg {
		result, err := server.CreateRoom(options)
		if err != nil {
			return err
		}
		return result
	}
}
func EditRoom(options *vc.RoomOptions) tea.Cmd {

	return func() tea.Msg {
		result, err := server.EditRoom(options)
		if err != nil {
			return err
		}
		return result
	}
}
