package vc

type ProgramOptions struct {
    ProgramId int
	AppFile   string
	Name      string
	Notes     string
    MobilityFile string
    WebxPanelFile string
    ProjectFile   string 
    CwsFile string
    StartNow string
}

type ProgramOptsFunc func(*ProgramOptions)

func WithFile(opt *ProgramOptions) {

}
