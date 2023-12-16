package main

import (
	"flag"
	"fmt"

	"github.com/ewilliams0305/VC4-CLI/pkg/tui"
)

func main() {

	tui.InitFlags()
	flag.Parse()

	fmt.Printf("\n\nCLI Started with host flag: %s %s\n\n", tui.Hostname, tui.Token)
	tui.Run()
}
