package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc"
	"golang.org/x/term"
)

const (
	loadProgram   initialAction = 0
	createRoom    initialAction = 1
	loadAndCreate initialAction = 2
)

type initialAction int

type ActionsModel struct {
	state    appState
	err      error
	message  string
	action   initialAction
	progress progress.Model
	results  *actionResult
	banner   *BannerModel
}

type actionResult struct {
	message string
	status  int
}

func InitialActionModel(message string, action initialAction) *ActionsModel {

	w, _, _ := term.GetSize(int(os.Stdout.Fd()))
	prog := progress.New(progress.WithDefaultGradient())
	prog.Width = w
	var state appState

	if action == loadProgram {
		state = programs
	}

	if action == createRoom {
		state = rooms
	}

	if action == loadAndCreate {
		state = rooms
	}

	return &ActionsModel{
		state:    state,
		message:  message,
		action:   action,
		err:      nil,
		progress: prog,
		banner:   NewBanner("VCLI Quick Action", 0, w),
	}
}

func (m ActionsModel) Init() tea.Cmd {

	var cmd tea.Cmd

	if m.action == loadProgram {
		ops := &vc.ProgramOptions{
			AppFile: ProgramFile,
			Name:    ProgramName,
		}
		cmd = CreateProgramAction(ops)
		return tea.Batch(cmd, actionsTickCmd())
	}

	if m.action == loadAndCreate {

		pops := &vc.ProgramOptions{
			AppFile: ProgramFile,
			Name:    ProgramName,
		}
		rops := &vc.RoomOptions{
			Name:                RoomID,
			ProgramInstanceId:   RoomID,
			AddressSetsLocation: true,
		}
		cmd = CreateAndRunRoomAction(pops, rops)
		return tea.Batch(cmd, actionsTickCmd())
	}

	return tea.Batch(cmd, actionsTickCmd())
}

func (m ActionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		m.err = msg
		m.banner = NewBanner(m.message, BannerErrorState, m.progress.Width)

		return m, nil

	case vc.ProgramUploadResult:
		r := actionResult{
			message: msg.FriendlyName,
			status:  int(msg.Code),
		}
		m.results = &r
		return m, nil

	case vc.RoomCreatedResult:
		r := actionResult{
			message: msg.Message,
			status:  int(msg.Code),
		}
		m.results = &r
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case progressTick:
		if m.progress.Percent() == 1.0 {
			return m, nil
		}
		return m, tea.Batch(actionsTickCmd(), m.progress.IncrPercent(0.20))

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			main := InitialModel()
			return main, main.Init()

		case "ctrl+q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m ActionsModel) View() string {
	s := m.banner.View()
	s += "\n" + m.message + "\n"

	if m.progress.Percent() != 0.0 {
		s += "\n" + m.progress.View() + "\n\n"
	}

	if m.err != nil {
		s += RenderErrorBox("error performing intial actions", m.err)
		s += GreyedOutText.Render("\n\n press esc to return to main menu or ctrl+q to quit")
		return s
	}

	if m.results != nil {
		var resultMessage string
		resultMessage += "\n RESULT: " + m.results.message
		resultMessage += "\n\n STATUS CODE: " + fmt.Sprintf("%d", m.results.status)
		resultMessage += "\n"

		s += RenderMessageBox(1000).Render(resultMessage)
	}

	s += GreyedOutText.Render("press esc to return to main menu or ctrl+q to quit")
	return s
}

func actionsTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return progressTick(t)
	})
}

func CreateProgramAction(options *vc.ProgramOptions) tea.Cmd {

	if !OverrideFile {
		return func() tea.Msg {
			return CreateNewProgram(*options)
		}
	}

	return func() tea.Msg {
		programs, err := server.GetPrograms()
		if err != nil {
			return err
		}

		var prog *vc.ProgramEntry
		for _, p := range programs {
			if strings.HasSuffix(options.Name, p.FriendlyName) {
				prog = &p
				break
			}
		}
		if prog == nil {
			return nil
		}

		options.ProgramId = int(prog.ProgramID)
		options.StartNow = true
		return EditProgram(*options)
	}

}

func CreateAndRunRoomAction(progOps *vc.ProgramOptions, roomOps *vc.RoomOptions) tea.Cmd {

	return func() tea.Msg {
		return CreateAndRunProgram(progOps, roomOps)
	}
}

func CreateRoomAction(options *vc.RoomOptions) tea.Cmd {

	return func() tea.Msg {
		return CreateRoom(*options)
	}
}

func CreateErrorAction() tea.Msg {
	return fmt.Errorf("FAILED TO PROCESS THE APPLICATION FLAGS, VALID COMBINATIONS REQUIRE -F && -N || -F && -N && -R	")
}
