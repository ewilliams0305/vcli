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
 result *vc.UploadProgramResult
 err error
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
   m.result = *msg
   return m, nil
  
  case error:
   m.err = msg
   return m, nil

  case tea.KeyMsg:
		 switch msg.String() {

		 case "esc":
			 return ReturnToPrograms(), tea.Batch(tick, QueryPrograms)
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
 s:= m.form.View()

 if m.err != nil {
  s += "\n\n" m.err.Error()
 }
 
 if m.result != nil {
  s += "\n\n" + m.result.Code
  // DISPLAY RESULT HERE AND MAYBE NO FORM
 } 
 return s
}

func SumbitNewProgramForm(f *huh.Form) tea.Cmd{
   if f == huh.StateCompleted {

			return CreateNewProgram(options)(vc.ProgramOptions{
				AppFile: f.GetString("FILE"),
				Name:    f.GetString("NAME"),
				Notes:   f.GetString("NOTES"),
			})
  }
}

func SubmitNewProgram(options vc.ProgramOptions) tea.Cmd {

	return func() tea.Msg {
		return CreateNewProgram(options)
	}
}
