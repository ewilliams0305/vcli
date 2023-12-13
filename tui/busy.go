package tui

import tea "github.com/charmbracelet/bubbletea"

type busy struct {
	flag    bool
	message string
}

func ShowBusyMessage(message string) tea.Cmd {
	b := busy{
		flag:    true,
		message: message,
	}
	return func() tea.Msg {
		return b
	}
}

func HideBusyMessage() tea.Msg {
	return busy{
		flag: false,
	}
}
