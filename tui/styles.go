package tui

import "github.com/charmbracelet/lipgloss"

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var SelectText = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(1).
	PaddingLeft(1).
	MarginBottom(1).
	Width(76).Align(lipgloss.Top).
	Height(3)
