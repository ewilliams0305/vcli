package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

type NewTokenForm struct {
	form     *huh.Form
	result   *vc.ApiToken
	progress progress.Model
	running  bool
	err      error
	edit     bool
}

var tokenFormDescription string
var tokenFormValue string
var tokenFormReadonly bool

func NewTokenFormModel() NewTokenForm {

	p := progress.New(progress.WithDefaultGradient())
	p.Width = app.width

	return NewTokenForm{
		progress: p,
		running:  false,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("DESCRIPTION").
					Title("Enter a description for the new API token").
					Prompt("ℹ  ").
					Placeholder("external vcli control").
					Validate(validateString).
					Value(&tokenFormDescription),

				huh.NewConfirm().
					Key("READONLY").
					Title("Reaonly API Key").
					Value(&tokenFormReadonly),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func EditTokenFormModel(token *vc.ApiToken) NewTokenForm {

	tokenFormDescription = token.Description
	tokenFormReadonly = tokenStatusToBool(token.Status)
	tokenFormValue = token.Token

	p := progress.New(progress.WithDefaultGradient())
	p.Width = app.width

	return NewTokenForm{
		progress: p,
		edit:     true,
		running:  false,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("DESCRIPTION").
					Title("Enter a description for the new API token").
					Prompt("ℹ  ").
					Placeholder("external vcli control").
					Validate(validateString).
					Value(&tokenFormDescription),

				huh.NewConfirm().
					Key("READONLY").
					Title("Reaonly API Key").
					Value(&tokenFormReadonly),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func (m NewTokenForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m NewTokenForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case error:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 20
		return m, nil

	case vc.ApiToken:
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
		return m, tea.Batch(tokenCreatedTickCmd(), m.progress.IncrPercent(0.20))

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+q":
			return InitialTokensModel(app.width, app.height), tea.Batch(tick, QueryTokens)

		case "ctrl+n":
			form := NewTokenFormModel()
			return form, tea.Batch(form.Init(), nil)
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted && !m.running {
			m.running = true

			if !m.edit {
				return m, tea.Batch(CreateToken(tokenFormDescription, tokenFormReadonly), tokenCreatedTickCmd())
			}
			return m, tea.Batch(EditToken(tokenFormDescription, tokenFormReadonly, tokenFormValue), tokenCreatedTickCmd())
		}
	}
	return m, cmd
}

func (m NewTokenForm) View() string {
	s := ""
	if m.edit {
		s += GreyedOutText.Render("\n🖊 Edit API Token\n")
	} else {
		s += GreyedOutText.Render("\n🆕 Create New Api Token\n")
	}

	s += "\n" + m.form.View()

	if m.progress.Percent() != 0.0 {
		s += "\n" + m.progress.View() + "\n\n"
	}
	if m.err != nil {
		s += RenderErrorBox("error creating new api token", m.err)
		s += GreyedOutText.Render("\n\n esc return * ctrl+n reset form")
		return s
	}

	if m.result != nil {

		s += RenderMessageBox(app.width).Render(fmt.Sprintf("API TOKEN: %s\n\n", m.result.Token))
		s += GreyedOutText.Render("\n\n esc return * ctrl+n reset form")
	}

	return s
}

func tokenStatusToBool(status vc.TokenStatus) bool {
	return status == vc.ReadOnlyToken
}

// func tokenReadonlyToStatus(readonly bool) vc.TokenStatus {
// 	if readonly {
// 		return vc.ReadOnlyToken
// 	}

// 	return vc.ReadWriteToken
// }

func tokenCreatedTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return progressTick(t)
	})
}
