package tui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc"
	"golang.org/x/term"
)

const (
	PrimaryColor string = "#3F51B5"
	PrimaryLight string = "#C5CAE9"
	PrimaryDark  string = "#001F5F"
	AccentColor  string = "#00796B"
	ErrorColor   string = "#8A0B29"
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

func RenderErrorBox(header string, err error) string {
	w, _, _ := term.GetSize(int(os.Stdout.Fd()))

	s := GreyedOutText.Width(w).
		Render("⚠ " + header + "\n")

	e := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#CCCCCC")).
		Background(lipgloss.Color(ErrorColor)).
		PaddingTop(1).
		PaddingLeft(1).
		MarginBottom(1).
		Width(w).Align(lipgloss.Top).
		Height(3).
		Render("\n" + err.Error() + "\n")
	return s + e
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

func GetOnlineIcon(status string) string {
	switch status {
	case "ONLINE":
		return "✅"
	case "OFFLINE":
		return "❌"
	}
	return "❌"
}

func CheckMark(status bool) string {
	if status {
		return " \u2713"
	}
	return "\u274C"
}

func GetIcons() string {
	return "🎏🍔🍒🍥🎮📦🦁🐶🐸🍕🥐🧲"
}
