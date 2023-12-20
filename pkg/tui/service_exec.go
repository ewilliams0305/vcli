package tui

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type VirtualControlServiceModel struct {
	altscreenActive bool
	err             error
	help            SystemsHelpModel
}

func InitialSystemModel() VirtualControlServiceModel {
	return VirtualControlServiceModel{
		help:            NewSystensHelpModel(),
		altscreenActive: true,
	}
}

func (m VirtualControlServiceModel) Init() tea.Cmd {
	return nil
}

func (m VirtualControlServiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "l", "ctrl+l":
			return m, tea.Batch(openJournal(), tea.EnterAltScreen)

		case "s", "ctrl+s":
			return m, tea.Batch(stopService(), tea.EnterAltScreen)

		case "n", "ctrl+n":
			return m, tea.Batch(startService(), tea.EnterAltScreen)

		case "r", "ctrl+r":
			return m, tea.Batch(restartService(), tea.EnterAltScreen)

		case "esc", "ctrl+q", "q":
			return ReturnToHomeModel(systemd), DeviceInfoCommand
		}

	case error:
		m.err = msg
		return m, nil

	case journalClosedMessage:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
	case serviceClosedMessage:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
	}
	return m, nil
}

func (m VirtualControlServiceModel) View() string {
	s := HighlightedText.Render("\nManage the Virtual Control Service\n\n")

	if m.err != nil {
		s += RenderErrorBox("Failed interacting with the virtual control service", m.err)
	}

	s += m.help.renderHelpInfo()
	return s
}

func openJournal() tea.Cmd {

	c := exec.Command("journalctl", "-u", "virtualcontrol.service", "-f")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return journalClosedMessage{err}
	})
}
func stopService() tea.Cmd {

	c := exec.Command("systemd", "stop", "virtualcontrol.service")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return serviceClosedMessage{err}
	})
}
func startService() tea.Cmd {

	c := exec.Command("systemd", "start", "virtualcontrol.service")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return serviceClosedMessage{err}
	})
}
func restartService() tea.Cmd {

	c := exec.Command("systemd", "restart", "virtualcontrol.service")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return serviceClosedMessage{err}
	})
}

type journalClosedMessage struct{ err error }

type serviceClosedMessage struct{ err error }
