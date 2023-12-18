package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

type NewRoomForm struct {
	form     *huh.Form
	result   *vc.ProgramUploadResult
	progress progress.Model
	running  bool
	err      error
	edit     bool
}

var roomOptions *vc.RoomOptions

// func validateProgramFile(file string) error {
// 	if strings.HasSuffix(file, ".cpz") || strings.HasSuffix(file, ".zip") || strings.HasSuffix(file, ".lpz") {
// 		return nil
// 	}
// 	return fmt.Errorf("INVALID FILE EXTENSION %s", file)
// }

// func validateProgramName(name string) error {
// 	if len(name) < 5 {
// 		return fmt.Errorf("NAME %s MUST HAVE AT LEAST 5 CHARATERS", name)
// 	}
// 	return nil
// }

func NewRoomFormModel() NewRoomForm {
	roomOptions = &vc.RoomOptions{}

	p := progress.New(progress.WithDefaultGradient())
	p.Width = program.width

	return NewRoomForm{
		progress: p,
		running:  false,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("NAME").
					Title("Enter Friendly Name").
					Prompt("📛  ").
					Placeholder("My friendly program name").
					//Validate(validateProgramName).
					Value(&roomOptions.Name),

				huh.NewInput().
					Key("NOTES").
					Title("Enter Notes").
					Prompt("📝  ").
					Placeholder("My seemingly pointless notes").
					Value(&roomOptions.Notes),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func EditRoomFormModel(programEntry *vc.ProgramEntry) NewRoomForm {
	roomOptions = &vc.RoomOptions{
		ProgramLibraryId: int(programEntry.ProgramID),
		Name:             programEntry.FriendlyName,
		Notes:            programEntry.Notes,
	}

	p := progress.New(progress.WithDefaultGradient())
	p.Width = program.width

	return NewRoomForm{
		edit:     true,
		running:  false,
		progress: p,
		form: huh.NewForm(
			huh.NewGroup(

				huh.NewInput().
					Key("NAME").
					Title("Enter Friendly Name").
					Prompt("🖊  ").
					Placeholder("My friendly program name").
					Validate(validateProgramName).
					Value(&roomOptions.Name),

				huh.NewInput().
					Key("NOTES").
					Title("Enter Notes").
					Prompt("📝  ").
					Placeholder("My seemingly pointless notes").
					Value(&roomOptions.Notes),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func (m NewRoomForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m NewRoomForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case error:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 20
		return m, nil

	case vc.ProgramUploadResult:
		m.result = &msg
		return m, m.progress.IncrPercent(1.0)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case progressTick:
		if m.progress.Percent() == 1.0 {
			return m, nil
		}
		return m, tea.Batch(roomCreatedTickCmd(), m.progress.IncrPercent(0.20))

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+q":
			return ReturnToPrograms(), tea.Batch(tick, ProgramsQuery)

		case "ctrl+n":
			form := NewRoomFormModel()
			return form, form.Init()
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted && !m.running {
			m.running = true
			return m, tea.Batch(SumbitNewRoomForm(&m), roomCreatedTickCmd())
		}
	}
	return m, cmd
}

func (m NewRoomForm) View() string {
	s := ""
	if m.edit {
		s += GreyedOutText.Render("\n🖊 Edit Running Program Instance\n")
	} else {
		s += GreyedOutText.Render("\n🆕 Create New Program Instance\n")
	}

	s += "\n" + m.form.View()

	if m.progress.Percent() != 0.0 {
		s += "\n" + m.progress.View() + "\n\n"
	}
	if m.err != nil {
		s += RenderErrorBox("error uploading new program file", m.err)
		s += GreyedOutText.Render("\n\n esc return * ctrl+n reset form")
		return s
	}

	if m.result != nil {
		var resultMessage string
		resultMessage += "\n RESULT           " + m.result.Result
		resultMessage += "\n\n PROGRAM ID:      " + fmt.Sprintf("%d", m.result.ProgramID)
		resultMessage += "\n\n PROGRAM NAME:    " + m.result.FriendlyName
		resultMessage += "\n\n STATUS CODE:     " + fmt.Sprintf("%d", m.result.Code)
		resultMessage += "\n"

		s += RenderMessageBox(1000).Render(resultMessage)
		s += GreyedOutText.Render("\n\n esc return * ctrl+n reset form")
	}

	return s
}

func SumbitNewRoomForm(m *NewRoomForm) tea.Cmd {
	if m.form.State != huh.StateCompleted {
		return nil
	}
	return func() tea.Msg {
		if m.edit {
			return EditRoom(roomOptions)
		}
		return CreateRoom(roomOptions)
	}
}

func roomCreatedTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return progressTick(t)
	})
}
