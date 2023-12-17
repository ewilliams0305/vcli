package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	vc "github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

var programsView *ProgramsModel

type ProgramsModel struct {
	table         table.Model
	Programs      vc.Programs
	selected      vc.ProgramEntry
	err           error
	help          programsHelpModel
	busy          busy
	cursor        int
	width, height int
}

func InitialProgramsModel(width, height int) *ProgramsModel {
	programsView = &ProgramsModel{
		table:    newProgramsTable(make(vc.Programs, 0), 0, width),
		Programs: vc.Programs{},
		selected: vc.ProgramEntry{},
		cursor:   0,
		help:     NewProgramsHelpModel(),
		width:    width,
		height:   height,
	}
	return programsView
}

func ReturnToPrograms() ProgramsModel {

	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	programsView.width = w
	programsView.height = h

	return *programsView
}

func BusyProgramsModel(b busy, Programs vc.Programs) *ProgramsModel {
	return &ProgramsModel{
		busy:     b,
		Programs: Programs,
		selected: vc.ProgramEntry{},
		cursor:   0,
	}
}

func (m ProgramsModel) Init() tea.Cmd {
	return ProgramsQuery
}

func (m ProgramsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tickMsg:
		w, h, _ := term.GetSize(int(os.Stdout.Fd()))
		if w != m.width || h != m.height {
			m.width = w
			m.height = h
		}

		// TODO: handle the errors and stop the polling
		return m, tea.Batch(ProgramsQuery, tick)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table, cmd = m.table.Update(msg)
		return m, cmd

	case vc.ProgramDeleteResult:
		return m, tea.Batch(ProgramsQuery, tick)

	case int:
		if msg <= len(m.Programs) {
			m.cursor = msg
			m.selected = m.Programs[msg]
		}
		_, cmd := m.table.Update(msg)
		return m, cmd

	case busy:
		m.busy = msg
		return m, nil

	case vc.Programs:
		m.busy = busy{flag: false}
		m.Programs = msg
		if len(msg) > 0 {
			m.table = newProgramsTable(msg, m.cursor, m.width)
			m.selected = msg[m.cursor]
			return m, nil
		}
		return m, nil

	case error:
		m.err = msg
		return NewProgramsErrorTable(msg), nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+q", "q", "ctrl+c", "esc":
			return ReturnToHomeModel(programs), tea.Batch(tick, DeviceInfoCommand)
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

		case "ctrl+n":
			if m.err == nil {
				form := NewProgramFormModel()
				return form, form.Init()
			}

		case "ctrl+e", "enter":
			if m.err == nil && len(m.selected.AppFile) > 0 {
				form := EditProgramFormModel(&m.selected)
				return form, form.Init()
			}

		case "ctrl+d":
			if m.err == nil {
				if m.cursor == len(m.Programs) {
					m.cursor = m.cursor - 1
				}
				if len(m.Programs) > 0 {
					return DeleteProgramFormModel(&m.selected), nil
				}
				return m, nil
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ProgramsModel) View() string {
	s := DisplayLogo(m.width)
	s += BaseStyle.Render(m.table.View()) + "\n\n"

	if m.busy.flag {
		s += RenderMessageBox(m.width).Render(m.busy.message)
	} else if m.err == nil {
		prog := fmt.Sprintf("\u2192 use keyboard actions to manage %s %s (ctrl+s, ctrl+d...)\n", m.selected.FriendlyName, m.selected.AppFile)
		s += RenderMessageBox(m.width).Render(prog)
	} else {
		if m.err != nil {
			s += RenderErrorBox("error performing program operation", m.err)
			return s
		}
	}

	s += m.help.renderHelpInfo()
	return s
}

func newProgramsTable(Programs vc.Programs, cursor int, width int) table.Model {

	columns := getProgramColumns(width)
	rows := getProgramRows(width, cursor, Programs)

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

func getProgramColumns(width int) []table.Column {

	if width < 200 {

		return []table.Column{
			{Title: "", Width: 1},
			{Title: "Name", Width: 20},
			{Title: "App File", Width: 35},
			{Title: "Notes", Width: width - 80},
			{Title: "Type", Width: 12},
		}
	}
	return []table.Column{
		{Title: "", Width: 1},
		{Title: "Name", Width: 20},
		{Title: "App File", Width: 35},
		{Title: "Notes", Width: width - 154},
		{Title: "Type", Width: 16},
		{Title: "Compiled", Width: 32},
		{Title: "Crestron DB", Width: 16},
		{Title: "Device DB", Width: 16},
	}
}

func getProgramRows(width int, cursor int, Programs vc.Programs) []table.Row {
	rows := []table.Row{}
	small := width < 200

	if len(Programs) == 0 {
		if small {
			rows = append(rows, table.Row{"", "No programs loaded to system, press ctl+n to a new program", "", "", "ctrl+n"})
		} else {
			rows = append(rows, table.Row{"", "No programs loaded to system, press ctl+n to a new program", "", "", "", "", "", "ctrl+n"})
		}
	}

	for i, prog := range Programs {
		marker := ""
		if cursor == i {
			marker = "\u2192"
		}
		if small {
			rows = append(rows, table.Row{marker, prog.FriendlyName, prog.AppFile, prog.Notes, prog.ProgramType})
		} else {
			rows = append(rows, table.Row{marker, prog.FriendlyName, prog.AppFile, prog.Notes, prog.ProgramType, prog.CompileDateTime, prog.CresDBVersion, prog.DeviceDBVersion})
		}
	}
	return rows
}

func NewProgramsErrorTable(msg vc.VirtualControlError) ProgramsModel {
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

	return ProgramsModel{table: t, err: msg, help: NewProgramsHelpModel()}
}

func ProgramsQuery() tea.Msg {

	programs, err := server.GetPrograms()
	if err != nil {
		return err
	}
	return programs
}

func CreateNewProgram(options vc.ProgramOptions) tea.Msg {

	result, err := server.CreateProgram(options)
	if err != nil {
		return err
	}
	return result
}

func EditProgram(options vc.ProgramOptions) tea.Msg {

	result, err := server.EditProgram(options)
	if err != nil {
		return err
	}
	return result
}

func DeleteProgram(id int) tea.Cmd {

	return func() tea.Msg {
		result, err := server.DeleteProgram(id)
		if err != nil {
			return err
		}
		return result
	}
}
