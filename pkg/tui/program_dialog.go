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
	edit   bool
}

var programOptions *vc.ProgramOptions

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
	programOptions = &vc.ProgramOptions{}

	return NewProgramForm{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("FILE").
					Title("Enter local program file path").
					Prompt("📂  ").
					Placeholder("/home/user/my_progam.cpz").
					Validate(validateProgramFile).
					Value(&programOptions.AppFile),

				huh.NewInput().
					Key("NAME").
					Title("Enter Friendly Name").
					Prompt("📛  ").
					Placeholder("My friendly program name").
					Validate(validateProgramName).
					Value(&programOptions.Name),

				huh.NewInput().
					Key("NOTES").
					Title("Enter Notes").
					Prompt("📝  ").
					Placeholder("My seemingly pointless notes").
					Value(&programOptions.Notes),

				huh.NewInput().
					Key("MOBILITY").
					Title("Enter local mobile project path").
					Prompt("📱  ").
					Placeholder("/home/user/mobile.Core3z").
					Value(&programOptions.MobilityFile),

				huh.NewInput().
					Key("XPANEL").
					Title("Enter local xpanel path").
					Prompt("❌  ").
					Placeholder("/home/user/xpanel.ch5z").
					Value(&programOptions.WebxPanelFile),

				huh.NewInput().
					Key("TOUCHPANEL").
					Title("Enter local touch panel project path").
					Prompt("📲  ").
					Placeholder("/home/user/mytp.vtz").
					Value(&programOptions.ProjectFile),

				huh.NewInput().
					Key("CONFIGURATION").
					Title("Enter local configuration webpage path").
					Prompt("⚙  ").
					Placeholder("/home/user/my_dist.zip").
					Value(&programOptions.CwsFile),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func EditProgramFormModel(program *vc.ProgramEntry) NewProgramForm {
	programOptions = &vc.ProgramOptions{
		ProgramId:     int(program.ProgramID),
		AppFile:       program.AppFile,
		Name:          program.FriendlyName,
		Notes:         program.Notes,
		MobilityFile:  program.MobilityFile,
		WebxPanelFile: program.WebxPanelFile,
		CwsFile:       program.CwsFile,
		StartNow:      false,
	}

	return NewProgramForm{
		edit: true,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("FILE").
					Title("Enter local file path").
					Prompt("📂  ").
					Placeholder("/home/user/my_progam.cpz").
					Validate(validateProgramFile).
					Value(&programOptions.AppFile),

				huh.NewInput().
					Key("NAME").
					Title("Enter Friendly Name").
					Prompt("🖊  ").
					Placeholder("My friendly program name").
					Validate(validateProgramName).
					Value(&programOptions.Name),

				huh.NewInput().
					Key("NOTES").
					Title("Enter Notes").
					Prompt("📝  ").
					Placeholder("My seemingly pointless notes").
					Value(&programOptions.Notes),

				huh.NewInput().
					Key("MOBILITY").
					Title("Enter local mobile project path").
					Prompt("📱  ").
					Placeholder("/home/user/mobile.Core3z").
					Value(&programOptions.MobilityFile),

				huh.NewInput().
					Key("XPANEL").
					Title("Enter local xpanel path").
					Prompt("❌  ").
					Placeholder("/home/user/xpanel.ch5z").
					Value(&programOptions.WebxPanelFile),

				huh.NewInput().
					Key("TOUCHPANEL").
					Title("Enter local touch panel project path").
					Prompt("📲  ").
					Placeholder("/home/user/mytp.vtz").
					Value(&programOptions.ProjectFile),

				huh.NewInput().
					Key("CONFIGURATION").
					Title("Enter local configuration webpage path").
					Prompt("⚙  ").
					Placeholder("/home/user/my_dist.zip").
					Value(&programOptions.CwsFile),

				huh.NewConfirm().
					Title("Would you like to restart effected rooms?").
					Value(&programOptions.StartNow),
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
		case "esc", "ctrl+q":
			return ReturnToPrograms(), tea.Batch(tick, ProgramsQuery)

		case "ctrl+n":
			form := NewProgramFormModel()
			return form, form.Init()

		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted {
			return m, SumbitNewProgramForm(&m)
		}
	}
	return m, cmd
}

func (m NewProgramForm) View() string {
	s := ""
	if m.edit {
		s += GreyedOutText.Render("\n🖊 Edit Program Entry\n")
	} else {
		s += GreyedOutText.Render("\n🆕 Create New Program Entry\n")
	}
	s += "\n" + m.form.View()

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

func SumbitNewProgramForm(m *NewProgramForm) tea.Cmd {
	if m.form.State != huh.StateCompleted {
		return nil
	}
	return func() tea.Msg {
		if m.edit {
			return EditProgram(*programOptions)
		}
		return CreateNewProgram(*programOptions)
	}
}
