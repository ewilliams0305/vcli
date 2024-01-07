package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc"
	"golang.org/x/term"
)

type TokenModel struct {
	table         table.Model
	entries       []vc.ApiToken
	selected      vc.ApiToken
	err           error
	help          TokensHelpModel
	width, height int
	banner        *BannerModel
}

func InitialTokensModel(width, height int) *TokenModel {
	entries := make([]vc.ApiToken, 0)
	token := textarea.New()
	token.Placeholder = "selected token displayed here"
	token.Focus()
	tokens := &TokenModel{
		table:    newApiTokensTable(entries, 0, width),
		entries:  entries,
		selected: vc.ApiToken{},
		help:     NewtokensHelpModel(),
		width:    width,
		height:   height,
		banner:   NewBanner("MANAGE API TOKENS", BannerNormalState, width),
	}
	return tokens
}

func (m TokenModel) Init() tea.Cmd {
	return QueryTokens
}

func (m TokenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tickMsg:
		w, h, _ := term.GetSize(int(os.Stdout.Fd()))
		if w != m.width || h != m.height {
			m.width = w
			m.height = h
		}
		return m, tea.Batch(tick, QueryTokens)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table, cmd = m.table.Update(msg)
		return m, cmd

	case int:
		m.selected = m.entries[msg]
		return m, nil

	case []vc.ApiToken:
		m.err = nil
		m.entries = msg
		m.table.SetRows(getApiTokenRows(m.width, m.table.Cursor(), msg))

		if len(msg) != 0 && m.table.Cursor() >= 0 {
			m.selected = msg[m.table.Cursor()]
		}
		return m, nil

	case error:
		m.err = msg
		m.banner = NewBanner(msg.Error(), BannerErrorState, m.width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "up":
			m.table, cmd = m.table.Update(msg)
			m.selected = m.entries[m.table.Cursor()]
			return m, cmd

		case "down":
			m.table, cmd = m.table.Update(msg)
			m.selected = m.entries[m.table.Cursor()]
			return m, cmd

		case "ctrl+q", "q", "ctrl+c", "esc":
			return ReturnToHomeModel(auth), tea.Batch(tick, DeviceInfoCommand)
		case "ctrl+d", "delete":
			if len(m.selected.Token) > 10 {
				return DeleteTokenFormModel(&m.selected), nil
			}
		case "ctrl+n", "n":
			form := NewTokenFormModel()
			return form, form.Init()
		case "ctrl+e", "enter":
			if len(m.selected.Token) > 10 {
				form := EditTokenFormModel(&m.selected)
				return form, form.Init()
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TokenModel) View() string {
	s := m.banner.View() + "\n"
	s += BaseStyle.Render(m.table.View()) + "\n\n"
	s += RenderMessageBox(m.width).Render(fmt.Sprintf("API TOKEN: %s\n\n", m.selected.Token))
	if m.err != nil {
		s += RenderErrorBox("FAILED MANAGING API TOKENS", m.err)
	}

	s += m.help.renderHelpInfo()
	return s
}

func newApiTokensTable(entries []vc.ApiToken, cursor int, width int) table.Model {

	columns := getApiKeyColumns(width)
	rows := getApiTokenRows(width, cursor, entries)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(16),
		table.WithWidth(width),
	)
	t.Focus()

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Foreground(lipgloss.Color(AccentColor)).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(PrimaryLight)).
		Background(lipgloss.Color(PrimaryDark)).
		Bold(false)
	t.SetStyles(s)
	t.SetCursor(cursor)
	return t
}

func getApiKeyColumns(width int) []table.Column {
	return []table.Column{
		{Title: "", Width: 1},
		{Title: "DESCRIPTION", Width: 16},
		{Title: "READONLY", Width: 8},
		{Title: "TOKEN", Width: width - 35},
	}
}

func getApiTokenRows(width int, cursor int, entries []vc.ApiToken) []table.Row {
	rows := []table.Row{}

	for i, token := range entries {
		marker := ""
		if cursor == i {
			marker = "\u2192"
		}
		rows = append(rows, table.Row{marker, token.Description, GetReadonlyIcon(token.Status), token.Token})
	}
	return rows
}

func QueryTokens() tea.Msg {
	tokens, err := server.GetTokens()
	if err != nil {
		return err
	}
	return tokens
}

func CreateToken(description string, readonly bool) tea.Cmd {
	return func() tea.Msg {
		r, e := server.CreateToken(readonly, description)
		if e != nil {
			return e
		}
		return r
	}
}

func EditToken(description string, readonly bool, token string) tea.Cmd {
	return func() tea.Msg {
		r, e := server.EditToken(readonly, description, token)
		if e != nil {
			return e
		}
		return r
	}
}

func DeleteToken(token string) tea.Cmd {
	return func() tea.Msg {
		r, e := server.DeleteToken(token)
		if e != nil {
			return e
		}
		return r
	}
}
