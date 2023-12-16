package tui

import (
	"flag"
)

var (

	// The hostname passed into the application as a command line argument flag
	Hostname string
	// The API token generated from the VC4 webpage, this is required to control an external appliance
	Token string

	// The room id passed into the application. providing a room id creates addtional options that can be leveraged to process room actions
	RoomID string
	// Program file path passed into the application used to create a new program entry
	ProgramFile string
	// Actions include GET PUT DEL and can be used with a combination of the Room ID and program file flags
	Action string
)

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

	flag.StringVar(&RoomID, "room", "", "An optional room id used from CRUD operations")
	flag.StringVar(&RoomID, "r", "", " An optional room id used for CRUD operations (shorthand)")
}
