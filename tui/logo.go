package tui

import "github.com/charmbracelet/lipgloss"

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

var LogoStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(1).
	PaddingLeft(1).
	MarginBottom(1).
	Width(76).Align(lipgloss.Center).
	Height(3)

func DisplayLogo() string {
	return LogoStyle.Render(Logo) + "\n\n"
}
