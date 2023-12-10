package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ewilliams0305/VC4-CLI/vc"
)

const (
	PrimaryColor string = "#3F51B5"
	PrimaryLight string = "#C5CAE9"
	PrimaryDark  string = "#303F9F"
	AccentColor  string = "#00796B"
)

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

var HighlightedText = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF75B7")).
	Bold(true)

var GreyedOutText = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#555555")).
	Bold(false)

func RenderMessageBox(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(PrimaryLight)).
		Background(lipgloss.Color(PrimaryDark)).
		PaddingTop(1).
		PaddingLeft(1).
		MarginBottom(1).
		Width(width).Align(lipgloss.Top).
		Height(3)
}

func GetStatus(status string) string {
	switch status {
	case string(vc.Starting):
		return "💨"
	case string(vc.Running):
		return "🚀"
	case string(vc.Stopping):
		return "🛑"
	case string(vc.Stopped):
		return "🤚"
	case string(vc.Aborted):
		return "😈"
	}
	return "🤚"
}

func CheckMark(status bool) string {
	if status {
		return " " + "\u2713"
	}
	return "\u274C"
}

func GetIcons() string {
	return "🎏🍔🍒🍥🎮📦🦁🐶🐸🍕🥐🧲"
}
