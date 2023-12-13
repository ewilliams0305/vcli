package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ewilliams0305/VC4-CLI/vc"
)

type NewProgramForm struct {
	form *huh.Form // huh.Form is just a tea.Model
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

func validateProgramName(name string) error{
 if len(name) < 5 {
   fmt.Errorf("NAME %s MUST HAVE AT LEAST 5 CHARATERS", name)
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
					Prompt("📂").
					Placeholder("/home/user/my_progam.cpz").
					Validate(validateProgramFile).
					Value(&filePath),

				huh.NewInput().
					Key("NAME").
					Title("Enter Friendly Name").
					Prompt("🍌").
					Placeholder("My friendly program name").
					Validate(validateProgramName).
     Value(&fileName),

				huh.NewInput().
					Key("NOTES").
					Title("Enter Notes").
					Prompt("📝").
					Placeholder("My seemingly pointless notes").
					Value(&notes),
			),
		),
	}
}

func (m NewProgramForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m NewProgramForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
 switch msg := msg.(type) {
	 case vc.ProgramUploadResult:
   // GOT A NEW PROGRAM LOADED RESULT; RENDER AND RETURN

 }
 
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted {

			return ReturnToPrograms(), SubmitNewProgram(vc.ProgramOptions{
				AppFile: m.form.GetString("FILE"),
				Name:    m.form.GetString("NAME"),
				Notes:   m.form.GetString("NOTES"),
			})
		}
	}

	return m, cmd
}

func (m NewProgramForm) View() string {
	return m.form.View()
}

func SubmitNewProgram(options vc.ProgramOptions) tea.Cmd {

	return func() tea.Msg {
		return CreateNewProgram(options)
	}
}
