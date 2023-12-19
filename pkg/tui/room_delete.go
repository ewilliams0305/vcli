package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

type DeleteRoomForm struct {
	room *vc.Room
	form *huh.Form
}

var (
	roomDeleteConfirm bool
)

func DeleteRoomFormModel(room *vc.Room) DeleteRoomForm {
	return DeleteRoomForm{
		room: room,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to delete room %s?", room.ID)).
					Value(&roomDeleteConfirm).
					Affirmative("Delete").
					Negative("Cancel"),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func (m DeleteRoomForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m DeleteRoomForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "shift+tab", "ctrl+q":
			return ReturnRoomsModel(), tea.Batch(RoomsQuery, tick)
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted {

			if roomDeleteConfirm {
				return ReturnRoomsModel(), tea.Batch(RoomsQuery, tick, DeleteRoom(m.room.ID))
			}
			return ReturnRoomsModel(), tea.Batch(RoomsQuery, tick)
		}
	}
	return m, cmd
}

func (m DeleteRoomForm) View() string {
	s := m.form.View()
	return s
}
