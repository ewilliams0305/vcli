package main

import (
	"flag"
	"fmt"

	"github.com/ewilliams0305/VC4-CLI/tui"
)

func main() {

	// ADD OUR COMMAND LINE ARGS
	tui.InitFlags()
	flag.Parse()

	fmt.Printf("CLI Started with host flag: %s and token: %s", tui.Hostname, tui.Token)
	tui.Run()

	// p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("VC4 CLI failed to start, there's been an error: %v", err)
	// 	os.Exit(1)
	// }
}
