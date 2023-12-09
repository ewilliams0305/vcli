package tui

import (
	"flag"
)

// // The hostname passed into the application as a command line argument flag
// var HostnameFlag = flag.String("host", "127.0.0.1", "The IP or hostname of the virtual control service")

// // The API token generated from the VC4 webpage, this is required to control an external appliance
// var TokenFlag = flag.String("token", "", "The API token generated from the VC4 webpage, this is required to control an external appliance")

var Hostname string
var Token string

func InitFlags() {
	const (
		defaultHost  = "127.0.0.1"
		defaultToken = ""
		hostUsage    = "The IP or hostname of the virtual control service"
		tokenUsage   = "The API token generated from the VC4 webpage, this is required to control an external appliance"
	)
	flag.StringVar(&Hostname, "host", defaultHost, hostUsage)
	flag.StringVar(&Hostname, "h", defaultHost, hostUsage+" (shorthand)")
	flag.StringVar(&Token, "token", defaultToken, tokenUsage)
	flag.StringVar(&Token, "t", defaultToken, tokenUsage+" (shorthand)")
}
