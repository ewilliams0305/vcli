package tui

import (
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type VirtualControlServiceModel struct {
	altscreenActive bool
	err             error
	help            SystemsHelpModel
	list            list.Model
	progress        progress.Model
	banner          *BannerModel
}

type serviceOption struct {
	title, desc string
}

type journalClosedMessage struct{ err error }

type serviceClosedMessage struct{ err error }

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (i serviceOption) Title() string       { return i.title }
func (i serviceOption) Description() string { return i.desc }
func (i serviceOption) FilterValue() string { return i.title }

func InitialSystemModel() VirtualControlServiceModel {

	items := []list.Item{
		serviceOption{title: "Stop", desc: "stops the virtual control systemd service"},
		serviceOption{title: "Start", desc: "starts the virtual control systemd service"},
		serviceOption{title: "Restart", desc: "restarts the virtual control systemd service"},
		serviceOption{title: "Logs", desc: "views the virtual control service journal"},
	}

	prog := progress.New(progress.WithDefaultGradient())
	prog.Width = app.width

	m := VirtualControlServiceModel{
		help:            NewSystensHelpModel(),
		altscreenActive: true,
		list:            list.New(items, list.NewDefaultDelegate(), 100, 20),
		progress:        prog,
		banner:          NewBanner("MANAGE VIRTUAL CONTROL SERVICE", BannerNormalState, app.width),
	}

	m.list.Title = "Virtual Control Service Actions"
	return m
}

func (m VirtualControlServiceModel) Init() tea.Cmd {
	return nil
}

func (m VirtualControlServiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "s", "ctrl+s":
			if m.list.FilterState() != list.Filtering {
				return m, tea.Batch(stopService(), systemTickCmd())
			}
		case "n", "ctrl+n":
			if m.list.FilterState() != list.Filtering {
				return m, tea.Batch(startService(), systemTickCmd())
			}
		case "r", "ctrl+r":
			if m.list.FilterState() != list.Filtering {
				return m, tea.Batch(restartService(), systemTickCmd())
			}
		case "l", "ctrl+l":
			if m.list.FilterState() != list.Filtering {
				return m, openJournal()
				// return m, tea.Batch(openJournal(), tea.EnterAltScreen)
			}
		case "esc", "ctrl+q", "q":
			if m.list.FilterState() != list.Filtering {
				return ReturnToHomeModel(systemd), DeviceInfoCommand
			}

		case "enter":

			switch m.list.Cursor() {
			case 0:
				return m, tea.Batch(stopService(), systemTickCmd())
			case 1:
				return m, tea.Batch(startService(), systemTickCmd())
			case 2:
				return m, tea.Batch(restartService(), systemTickCmd())
			case 3:
				return m, openJournal()
				// return m, tea.Batch(openJournal(), tea.EnterAltScreen)

			}
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

	case progress.Model:
		m.progress = msg

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case progressTick:
		if m.progress.Percent() == 1.0 {
			m.progress.SetPercent(0)
			// return m, tea.Batch(openJournal(), tea.EnterAltScreen)
			return m, openJournal()
		}
		return m, tea.Batch(systemTickCmd(), m.progress.IncrPercent(0.20))
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m VirtualControlServiceModel) View() string {
	s := m.banner.View() + "\n"
	s += docStyle.Render(m.list.View())
	s += "\n\n\n"

	if m.err != nil {
		s += RenderErrorBox("Failed interacting with the virtual control service", m.err)
	}

	if m.progress.Percent() != 0.0 {
		s += "\n" + m.progress.View() + "\n"
	} else {
		s += "\n\n\n"
	}

	s += m.help.renderHelpInfo()
	return s
}

func systemTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*300, func(t time.Time) tea.Msg {
		return progressTick(t)
	})
}
func openJournal() tea.Cmd {

	c := exec.Command("journalctl", "-u", "virtualcontrol.service", "-f")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return journalClosedMessage{err}
	})
}
func stopService() tea.Cmd {

	c := exec.Command("systemctl", "stop", "virtualcontrol.service")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return serviceClosedMessage{err}
	})
}
func startService() tea.Cmd {

	c := exec.Command("systemctl", "start", "virtualcontrol.service")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return serviceClosedMessage{err}
	})
}
func restartService() tea.Cmd {

	c := exec.Command("systemctl", "restart", "virtualcontrol.service")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return serviceClosedMessage{err}
	})
}
