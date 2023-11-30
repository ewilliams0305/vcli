package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc4"
)

func deviceInfoCommand() tea.Msg {

	info, err := vc4.GetDeviceInfo()
	if err != nil {
		return info
	}

	return info
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type DeviceTableModel struct {
	table table.Model
}

func (m DeviceTableModel) Init() tea.Cmd {
	return nil
}

func (m DeviceTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return InitialModel(), nil
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m DeviceTableModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func NewDeviceTable(info vc4.DeviceInfo) DeviceTableModel {
	columns := []table.Column{
		{Title: "Device Id", Width: 4},
		{Title: "MAC", Width: 4},
		{Title: "Build Data", Width: 10},
		{Title: "Mono Version", Width: 10},
		{Title: "Python Version", Width: 10},
	}

	rows := []table.Row{
		{info.ID, info.MacAddress, info.BuildDate, info.MonoVersion, info.PythonVersion},
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
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return DeviceTableModel{t}
}

func NewDeviceErrorTable(msg error) DeviceTableModel {
	columns := []table.Column{
		{Title: "ERROR", Width: 10},
		{Title: "MESSAGE", Width: 10},
		{Title: "CODE", Width: 10},
	}

	rows := []table.Row{
		{msg.Error(), msg.Error(), msg.Error()},
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
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return DeviceTableModel{t}
}
