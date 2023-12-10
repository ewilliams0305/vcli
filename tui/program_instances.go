package tui

import (
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	vc "github.com/ewilliams0305/VC4-CLI/vc"
)

type RoomsTableModel struct {
	table table.Model
	rooms vc.ProgramInstanceLibrary
	err   error
	help  HelpModel
	row   string
	busy  busy
}

func NewRoomsModel() *RoomsTableModel {
	return &RoomsTableModel{}
}

func BusyRoomsModel(b busy, rooms vc.ProgramInstanceLibrary) *RoomsTableModel {
	return &RoomsTableModel{
		busy:  b,
		rooms: rooms,
	}
}

func (m RoomsTableModel) Init() tea.Cmd {
	return DeviceInfoCommand
}

func (m RoomsTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case busy:
		m.busy = msg
		return m, RefreshRoomData

	case vc.ProgramInstanceLibrary:
		m.rooms = msg
		return NewRoomsTable(msg), nil

	case error:
		m.err = msg
		return NewRoomsErrorTable(msg), nil

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c", "esc":
			return InitialModel(), DeviceInfoCommand
		case "down":
			m.table.SetCursor(m.table.Cursor() + 1)
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		case "up":
			m.table.SetCursor(m.table.Cursor() - 1)
			m.table, cmd = m.table.Update(msg)
			return m, cmd

		case "ctrl+s":

			room, err := getRoomFromCurson(m.rooms, m.table.Cursor())
			if err != nil {
				return NewRoomsErrorTable(err), nil
			}
			if room.Status == "Running" {
				return m, cmdRoomStop(room.ProgramInstanceID)
			} else if room.Status == "Stopped" {
				return m, cmdRoomStart(room.ProgramInstanceID)
			} else if room.Status == "Aborted" {
				return m, cmdRoomStart(room.ProgramInstanceID)
			}

		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	//m.row = m.table.SelectedRow()[0] + ": " + m.table.SelectedRow()[1]
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func getRoomFromCurson(rooms vc.ProgramInstanceLibrary, cursor int) (vc.ProgramInstance, error) {
	var i int = 0
	if rooms == nil {
		return vc.ProgramInstance{}, errors.New("OUT OF RANGE")
	}
	for _, r := range rooms {
		if i == cursor {
			return r, nil
		}
		i++
	}
	return vc.ProgramInstance{}, errors.New("OUT OF RANGE")
}

func (m RoomsTableModel) View() string {
	s := DisplayLogo()
	s += BaseStyle.Render(m.table.View()) + "\n\n"

	if len(m.row) > 0 {
		s += SelectText.Render(m.row)
	}

	if m.busy.flag {
		s += HighlightedText.Render(m.busy.message)
	} else {
		s += m.help.renderHelpInfo()
	}

	return s
}

func NewRoomsTable(rooms vc.ProgramInstanceLibrary) RoomsTableModel {
	columns := []table.Column{
		{Title: "ID", Width: 20}, {Title: "NAME", Width: 20}, {Title: "NOTES", Width: 30}, {Title: "STATUS", Width: 8}, {Title: "DEBUG", Width: 8},
	}

	//items := len(rooms)
	rows := []table.Row{}

	for key, room := range rooms {
		rows = append(rows, table.Row{key, room.Name, room.Notes, GetStatus(room.Status), CheckMark(room.DebuggingEnabled)})
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
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return RoomsTableModel{
		table: t,
		rooms: rooms,
		help:  NewHelpModel()}
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

	return RoomsTableModel{table: t, help: NewHelpModel()}
}

func RoomCommand() tea.Msg {

	info, err := server.ProgramInstances()
	if err != nil {
		return err
	}
	return info
}

func RefreshRoomData() tea.Msg {

	var rooms vc.ProgramInstanceLibrary
	var err error
	_, err = server.ProgramInstances()
	if err != nil {
		return err
	}

	//for i := 0; i < 3; i++ {
	time.Sleep(3 * time.Second)

	rooms, err = server.ProgramInstances()
	if err != nil {
		return err
	}

	//}

	return rooms
}

func RoomRestart(id string) tea.Msg {

	_, err := server.RestartRoom(id)
	if err != nil {
		return err
	}
	return busy{flag: true, message: fmt.Sprintf("Restarting room %s", id)}
}
func RoomStop(id string) tea.Msg {

	_, err := server.StopRoom(id)
	if err != nil {
		return err
	}
	return busy{flag: true, message: fmt.Sprintf("Stopping room %s", id)}
}
func RoomStart(id string) tea.Msg {

	_, err := server.StartRoom(id)
	if err != nil {
		return err
	}
	return busy{flag: true, message: fmt.Sprintf("Starting room %s", id)}
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
