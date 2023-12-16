package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

type DeleteProgramForm struct {
	form    *huh.Form
	program *vc.ProgramEntry
}

var (
	progDeleteConfirm bool
)

func DeleteProgramFormModel(program *vc.ProgramEntry) DeleteProgramForm {
	return DeleteProgramForm{
		program: program,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to delete %s and any rooms using this program?", program.ProgramName)).
					Value(&progDeleteConfirm).
					Affirmative("Yes").
					Negative("Cancel"),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func (m DeleteProgramForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m DeleteProgramForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case deleteProgramConfirmationResult:
		if msg == 0 {
			return ReturnToPrograms(), tea.Batch(ProgramsQuery, tick)
		}

	case vc.ProgramDeleteResult:
		if msg.Success {
			return ReturnToPrograms(), tea.Batch(ProgramsQuery, tick)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "shift+tab":
			return ReturnToPrograms(), tea.Batch(ProgramsQuery, tick)
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted {

			return m, SumbitDeleteProgramForm(&m)
		}
	}

	return m, cmd
}

func (m DeleteProgramForm) View() string {
	s := m.form.View()
	return s
}

func SumbitDeleteProgramForm(m *DeleteProgramForm) tea.Cmd {
	if m.form.State != huh.StateCompleted {
		return nil
	}

	if progDeleteConfirm {
		return DeleteProgram(int(m.program.ProgramID))
	}

	return ReturnFromDelete
}

func ReturnFromDelete() tea.Msg {
	return deleteProgramConfirmationResult(0)
}

type deleteProgramConfirmationResult int
