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

func NewProgramFormModel() NewProgramForm {
	return NewProgramForm{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("FILE").
					Title("Enter local file path").
					//Prompt("📂").
					//Placeholder("/home/user/my_progam.cpz").
					Validate(validateProgramFile).
					Value(&filePath),

				huh.NewInput().
					Key("NAME").
					Title("Enter Friendly Name").
					//Prompt("🍌").
					//Placeholder("My friendly program name").
					Value(&fileName),

				huh.NewInput().
					Key("NOTES").
					Title("Enter Notes").
					//Prompt("📝").
					//Placeholder("My seemingly pointless notes").
					Value(&notes),
			),
		),
	}
}

func (m NewProgramForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m NewProgramForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted {

			path := m.form.GetString("FILE")
			name := m.form.GetString("NAME")
			notes := m.form.GetString("NOTES")

			return InitialProgramsModel(200, 200), SubmitNewProgram(vc.ProgramOptions{
				AppFile: path,
				Name:    name,
				Notes:   notes,
			})
		}
	}

	return m, cmd
}

func (m NewProgramForm) View() string {
	// if m.form.State == huh.StateCompleted {
	// 	path := m.form.GetString("FILE")
	// 	name := m.form.GetString("NAME")
	// 	notes := m.form.GetString("NOTES")
	// 	return fmt.Sprintf("You selected: %s, Lvl. %d", path, name, notes)
	// }
	return m.form.View()
}

func SubmitNewProgram(options vc.ProgramOptions) tea.Cmd {

	return func() tea.Msg {
		return CreateNewProgram(options)
	}
}
