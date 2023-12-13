package vc

type ProgramOptions struct {
	AppFile string
	Name    string
	Notes   string
}

type ProgramOptsFunc func(*ProgramOptions)

func WithFile(opt *ProgramOptions) {

}
