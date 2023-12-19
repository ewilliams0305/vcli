package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	loadProgram   initialAction = 0
	createRoom    initialAction = 1
	loadAndCreate initialAction = 2
)

type initialAction int

type ActionsModel struct {
	state    appState
	err      error
	message  string
	action   initialAction
	progress progress.Model
}

func InitialActionModel(message string, action initialAction) *ActionsModel {

	var state appState

	if action == loadProgram {
		state = programs
	}

	if action == createRoom {
		state = rooms
	}

	if action == loadAndCreate {
		state = rooms
	}

	return &ActionsModel{
		state:    state,
		message:  message,
		action:   action,
		err:      nil,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m ActionsModel) Init() tea.Cmd {

	var cmd tea.Cmd

	// if m.action == loadProgram {
	// 	cmd = CreateNewProgram()
	// }

	// if m.action == createRoom {
	// 	state = rooms
	// }

	// if m.action == loadAndCreate {
	// 	state = rooms
	// }

	return tea.Batch(cmd, actionsTickCmd())
}

func (m ActionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		m.err = msg
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case progressTick:
		if m.progress.Percent() == 1.0 {
			return m, nil
		}
		return m, tea.Batch(actionsTickCmd(), m.progress.IncrPercent(0.20))

	}
	return m, nil
}

func (m ActionsModel) View() string {
	s := m.message

	if m.progress.Percent() != 0.0 {
		s += "\n" + m.progress.View() + "\n\n"
	}
	s += GreyedOutText.Render("this is the actions model")
	return s
}

func actionsTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return progressTick(t)
	})
}
