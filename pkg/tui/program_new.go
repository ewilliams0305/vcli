package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

type NewProgramForm struct {
	form   *huh.Form
	result *vc.ProgramUploadResult
	err    error
}

var (
	filePath string
	fileName string
	notes    string
)

func validateProgramFile(file string) error {
	if strings.HasSuffix(file, ".cpz") || strings.HasSuffix(file, ".zip") || strings.HasSuffix(file, ".lpz") {
		return nil
	}
	return fmt.Errorf("INVALID FILE EXTENSION %s", file)
}

func validateProgramName(name string) error {
	if len(name) < 5 {
		return fmt.Errorf("NAME %s MUST HAVE AT LEAST 5 CHARATERS", name)
	}
	return nil
}

func NewProgramFormModel() NewProgramForm {
	return NewProgramForm{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("FILE").
					Title("Enter local file path").
					Prompt("📂  ").
					Placeholder("/home/user/my_progam.cpz").
					Validate(validateProgramFile).
					Value(&filePath),

				huh.NewInput().
					Key("NAME").
					Title("Enter Friendly Name").
					Prompt("🍌  ").
					Placeholder("My friendly program name").
					Validate(validateProgramName).
					Value(&fileName),

				huh.NewInput().
					Key("NOTES").
					Title("Enter Notes").
					Prompt("📝  ").
					Placeholder("My seemingly pointless notes").
					Value(&notes),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func (m NewProgramForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m NewProgramForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case error:
		m.err = msg
		return m, nil

	case vc.ProgramUploadResult:
		m.result = &msg
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return ReturnToPrograms(), tea.Batch(tick, ProgramsQuery)

		case "ctrl+n":

			fileName = ""
			filePath = ""
			notes = ""
			form := NewProgramFormModel()
			return form, form.Init()

		case "ctrl+q":

			form := NewProgramFormModel()
			return form, form.Init()

		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted {

			return m, SumbitNewProgramForm(m.form)
		}
	}

	return m, cmd
}

func (m NewProgramForm) View() string {
	s := m.form.View()

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

func SumbitNewProgramForm(f *huh.Form) tea.Cmd {
	if f.State != huh.StateCompleted {
		return nil
	}

	return func() tea.Msg {
		options := vc.ProgramOptions{
			AppFile: f.GetString("FILE"),
			Name:    f.GetString("NAME"),
			Notes:   f.GetString("NOTES"),
		}
		return CreateNewProgram(options)
	}
}
