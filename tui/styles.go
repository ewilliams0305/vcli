package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ewilliams0305/VC4-CLI/vc"
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

func GetStatus(status string) string {
	switch status {
	case string(vc.Running):
		return "\u2713"
	case string(vc.Stopped):
		return "\u274C"
	case string(vc.Aborted):
		return "\u1F643"
	}
	return "\u1F641"
}

func CheckMark(status bool) string {
	if status {
		return " " + "\u2713"
	}
	return "\u274C"
}

func GetIcons() string {
  return "🎏🍔🍒🍥🎮📦🦁🐶🐸🍕🥐🧲";
}
