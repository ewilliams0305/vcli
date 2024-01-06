package tui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	vc "github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

var iptable *IpTableModel

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

func InitialIpTableModel(width, height int, roomid string) *IpTableModel {
	iptable = &IpTableModel{
		roomId:   roomid,
		entries:  make([]vc.IpTableEntry, 0),
		selected: vc.IpTableEntry{},
		cursor:   0,
		help:     NewHelpModel(),
		width:    width,
		height:   height,
		banner:   NewBanner(fmt.Sprintf("VIEWING %s PROGRAM IP TABLES", roomid), BannerNormalState, width),
	}
	return iptable
}

func (m IpTableModel) Init() tea.Cmd {
	return IpTableQuery(m.roomId)
}

func (m IpTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tickMsg:
		w, h, _ := term.GetSize(int(os.Stdout.Fd()))
		if w != m.width || h != m.height {
			iptable.width = w
			iptable.height = h
		}
		return iptable, tea.Batch(tick, IpTableQuery(iptable.roomId))

	case tea.WindowSizeMsg:
		iptable.width = msg.Width
		iptable.height = msg.Height
		iptable.table, cmd = iptable.table.Update(msg)
		return iptable, cmd

	case int:
		iptable.cursor = msg
		return iptable, nil

	case []vc.IpTableEntry:
		iptable.err = nil
		t := newIpTableDeviceTable(msg, iptable.cursor, iptable.width)
		iptable.table = t
		iptable.entries = msg
		iptable.banner = NewBanner(fmt.Sprintf("VIEWING %s PROGRAM IP TABLES", iptable.roomId), BannerNormalState, iptable.width)
		return iptable, nil

	case error:
		iptable.err = msg
		iptable.banner = NewBanner(msg.Error(), BannerErrorState, iptable.width)
		return iptable, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+q", "q", "ctrl+c", "esc":
			return ReturnRoomsModel(), tea.Batch(tick, RoomsQuery)
		case "down":
			if iptable.err == nil {
				iptable.table.SetCursor(iptable.table.Cursor() + 1)
				return iptable, cmdCursor(iptable.table.Cursor())
			}
		case "up":
			if iptable.err == nil {
				iptable.table.SetCursor(iptable.table.Cursor() - 1)
				return iptable, cmdCursor(iptable.table.Cursor())
			}
		}
	}
	iptable.table, cmd = iptable.table.Update(msg)
	return iptable, cmd
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

func IpTableQuery(id string) tea.Cmd {

	return func() tea.Msg {
		ipTable, err := server.GetIpTable(id)
		if err != nil {
			return err
		}
		return ipTable
	}
}
