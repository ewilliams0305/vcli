package tui

import (
	"flag"
)

var (

	// The hostname passed into the application as a command line argument flag
	Hostname string
	// The API token generated from the VC4 webpage, this is required to control an external appliance
	Token string
	// Program file path passed into the application used to create a new program entry
	ProgramFile string
	// The programs and room name used to whip up a new room with a simply command line arg
	ProgramName string
	// The room id passed into the application. providing a room id creates additional options that can be leveraged to process room actions
	RoomID string
	// Actions include GET PUT DEL and can be used with a combination of the Room ID and program file flags
	Action string
	// When the -o flag is provided the application will override the provided room.
	OverrideFile bool
)

func InitFlags() {
	const (
		defaultHost       = "127.0.0.1"
		defaultToken      = ""
		hostFlagUsage     = "The IP or hostname of the virtual control service"
		tokenFlagUsage    = "The API token generated from the VC4 webpage, this is required to control an external appliance"
		progFlagUsage     = "An optional flag to load a program file"
		nameFagUsage      = "An optional glag used to name the loaded program flag"
		roomFlagUsage     = "An optional room ID used to spin up a new room"
		overrideFlagUsage = "An option flag to let the application know to override the provided program file"
	)

	flag.StringVar(&Hostname, "host", defaultHost, hostFlagUsage)
	flag.StringVar(&Hostname, "h", defaultHost, hostFlagUsage+" (shorthand)")
	flag.StringVar(&Token, "token", defaultToken, tokenFlagUsage)
	flag.StringVar(&Token, "t", defaultToken, tokenFlagUsage+" (shorthand)")

	flag.StringVar(&RoomID, "room", "", roomFlagUsage)
	flag.StringVar(&RoomID, "r", "", roomFlagUsage+" (shorthand)")

	flag.StringVar(&ProgramFile, "file", "", progFlagUsage)
	flag.StringVar(&ProgramFile, "f", "", progFlagUsage+" (shorthand)")

	flag.StringVar(&ProgramName, "name", "", nameFagUsage)
	flag.StringVar(&ProgramName, "n", "", nameFagUsage+" (shorthand)")

	flag.BoolVar(&OverrideFile, "override", false, overrideFlagUsage)
	flag.BoolVar(&OverrideFile, "o", false, overrideFlagUsage+" (shorthand)")
}
