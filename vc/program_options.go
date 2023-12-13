package vc

type ProgramOptions struct {
	AppFile string
	Name    string
	Notes   string
}

type ProgramOptsFunc func(*ProgramOptions)

func WithFile(opt *ProgramOptions) {

}

// func DefaultProgramOptions() ProgramOptions {
// 	return ProgramOptions{
// 		appFile: "/home/program.zip",
// 		name:    "New Program",
// 		notes:   "Default Program Info",
// 	}
// }

// func NewProgramOptions(opts ...ProgramOptsFunc) *ProgramOptions {
// 	defaultOpts := DefaultProgramOptions()

// 	for _, optFn := range opts {
// 		optFn(&defaultOpts)
// 	}

// 	return &ProgramOptions{}
// }
