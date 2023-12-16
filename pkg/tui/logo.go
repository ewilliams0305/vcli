package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	Logo string = `
__      _______ _  _      _____ _      _____ 
\ \    / / ____| || |    / ____| |    |_   _|
 \ \  / / |    | || |_  | |    | |      | |  
  \ \/ /| |    |__   _| | |    | |      | |  
   \  / | |____   | |   | |____| |____ _| |_ 
    \/   \_____|  |_|    \_____|______|_____|

`
)

func DisplayLogo(width int) string {

	var LogoStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(PrimaryLight)).
		Background(lipgloss.Color(PrimaryDark)).
		PaddingTop(1).
		PaddingLeft(1).
		MarginBottom(1).
		Width(width).
		Align(lipgloss.Center)
	//Height(20)

	//box := fmt.Sprintf("%s\n", logo())
	return LogoStyle.Render(Logo) + "\n\n"
}
