package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	BannerNormalState BannerState = iota
	BannerErrorState  BannerState = 1
)

type BannerModel struct {
	message string
	state   BannerState
	width   int
}

type BannerState int

func (m BannerModel) Init() tea.Cmd {
	return nil
}

func NewBanner(message string, state BannerState, width int) *BannerModel {
	return &BannerModel{
		message: message,
		state:   state,
		width:   width,
	}
}

func (m BannerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case BannerState:
		m.state = msg
	case string:
		m.message = msg
	case error:
		m.state = BannerErrorState
		m.message = msg.Error()
	}
	return m, nil
}

func (m BannerModel) View() string {
	return renderBanner(m)
}

func renderBanner(model BannerModel) string {
	var bg lipgloss.Color
	if model.state == BannerNormalState {
		bg = lipgloss.Color(PrimaryDark)
	} else {
		bg = lipgloss.Color(ErrorColor)
	}

	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(PrimaryLight)).
		Background(lipgloss.Color(bg)).
		PaddingTop(1).
		PaddingLeft(1).
		MarginBottom(1).
		Width(model.width).
		Align(lipgloss.Center).
		Render(model.message + "\n")
}
