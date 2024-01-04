package tui

import (
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	vc "github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

type IpTableModel struct {
	roomId        string
	table         table.Model
	entries       []vc.IpTableEntry
	selected      vc.IpTableEntry
	err           error
	help          HelpModel
	cursor        int
	width, height int
	banner        *BannerModel
}

func InitialIpTableModel(width, height int) *IpTableModel {
	return &IpTableModel{
		entries:  make([]vc.IpTableEntry, 0),
		selected: vc.IpTableEntry{},
		cursor:   0,
		help:     NewHelpModel(),
		width:    width,
		height:   height,
		banner:   NewBanner("VIEW PROGRAM IP TABLES", BannerNormalState, width),
	}
}

func (m IpTableModel) Init() tea.Cmd {
	return RoomsQuery
}

func (m IpTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tickMsg:
		w, h, _ := term.GetSize(int(os.Stdout.Fd()))
		if w != m.width || h != m.height {
			m.width = w
			m.height = h
		}
		return m, tea.Batch(IpTableQuery(m.roomId), tick)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table, cmd = m.table.Update(msg)
		return m, cmd

	case int:
		m.cursor = msg
		return m, nil

	case []vc.IpTableEntry:
		t := newIpTableDeviceTable(msg, m.cursor, m.width)
		m.table = t
		return m, nil

	case error:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+q", "q", "ctrl+c", "esc":
			return ReturnToHomeModel(devices), tea.Batch(tick, DeviceInfoCommand)
		case "down":
			if m.err == nil {
				m.table.SetCursor(m.table.Cursor() + 1)
				return m, cmdCursor(m.table.Cursor())
			}
		case "up":
			if m.err == nil {
				m.table.SetCursor(m.table.Cursor() - 1)
				return m, cmdCursor(m.table.Cursor())
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m IpTableModel) View() string {
	s := m.banner.View() + "\n"
	s += BaseStyle.Render(m.table.View()) + "\n\n"

	// if m.busy.flag {
	// 	s += RenderMessageBox(m.width).Render(m.busy.message)
	// } else {
	// 	room := fmt.Sprintf("\u2192 use keyboard actions to manage %s %s (ctrl+s, ctrl+d...)\n", m.selectedRoom.ID, m.selectedRoom.ProgramName)
	// 	s += RenderMessageBox(m.width).Render(room)
	// }

	s += m.help.renderHelpInfo()
	return s
}

func newIpTableDeviceTable(entries []vc.IpTableEntry, cursor int, width int) table.Model {

	columns := getIpTableColumns(width)
	rows := getIpTableRows(width, cursor, entries)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(16),
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

func getIpTableColumns(width int) []table.Column {

	if width < 120 {

		return []table.Column{
			{Title: "", Width: 1},
			{Title: "IPID", Width: 8},
			{Title: "MODEL", Width: 28},
			{Title: "DESCRIPTION", Width: width - 57},
			{Title: "STATUS", Width: 8},
		}
	}
	return []table.Column{
		{Title: "", Width: 1},
		{Title: "ID", Width: 20},
		{Title: "MODEL", Width: 35},
		{Title: "DESCRIPTION", Width: width - 113},
		{Title: "IP ADDRESS", Width: 35},
		{Title: "STATUS", Width: 8},
	}
}

func getIpTableRows(width int, cursor int, entries []vc.IpTableEntry) []table.Row {
	rows := []table.Row{}
	small := width < 120

	for i, ipt := range entries {
		marker := ""
		if cursor == i {
			marker = "\u2192"
		}
		if small {
			rows = append(rows, table.Row{marker, strconv.FormatInt(int64(ipt.ProgramIPID), 16), ipt.Model, ipt.Description, GetOnlineIcon(ipt.Status)})
		} else {
			rows = append(rows, table.Row{marker, strconv.FormatInt(int64(ipt.ProgramIPID), 16), ipt.Model, ipt.Description, ipt.RemoteIP, GetOnlineIcon(ipt.Status)})
		}
	}
	return rows
}

func NewIpTableErrorTable(msg vc.VirtualControlError) IpTableModel {
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

	return IpTableModel{
		table: t,
		err:   msg,
	}
}

func IpTableQuery(id string) tea.Cmd {

	return func() tea.Msg {
		ipTable, err := server.GetIpTable(id)
		if err != nil {
			return err
		}
		return ipTable
	}
}
