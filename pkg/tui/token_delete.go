package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

type DeleteTokenForm struct {
	form  *huh.Form
	token *vc.ApiToken
}

var (
	tokenDeleteConfirm bool
)

func DeleteTokenFormModel(token *vc.ApiToken) DeleteTokenForm {
	return DeleteTokenForm{
		token: token,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to delete %s api?", token.Description)).
					Value(&tokenDeleteConfirm).
					Affirmative("Yes").
					Negative("Cancel"),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func (m DeleteTokenForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m DeleteTokenForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case bool:
		if msg {
			return InitialTokensModel(app.width, app.height), tea.Batch(DeleteToken(m.token.Token), QueryTokens, tick)
		}
		return InitialTokensModel(app.width, app.height), tea.Batch(QueryTokens, tick)
	case tea.KeyMsg:
		switch msg.String() {
		case "shift+tab":
			return InitialTokensModel(app.width, app.height), tea.Batch(QueryTokens, tick)
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted {

			return m, SumbitDeleteTokenForm(&m)
		}
	}
	return m, cmd
}

func (m DeleteTokenForm) View() string {
	s := m.form.View()
	return s
}

func SumbitDeleteTokenForm(m *DeleteTokenForm) tea.Cmd {
	if m.form.State != huh.StateCompleted {
		return nil
	}
	return deleteTokenConfirmation
}

func deleteTokenConfirmation() tea.Msg {
	return tokenDeleteConfirm
}
