package tui

import "github.com/ewilliams0305/VC4-CLI/pkg/vc"

var (
	server       vc.VirtualControl
	app          *MainModel
	programsView *ProgramsModel
	roomsModel   *RoomsTableModel
	iptable      *IpTableModel
)
