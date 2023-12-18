package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ewilliams0305/VC4-CLI/pkg/vc"
)

type NewRoomForm struct {
	form     *huh.Form
	result   *vc.RoomCreatedResult
	progress progress.Model
	running  bool
	err      error
	edit     bool
}

var roomOptions *vc.RoomOptions
var selectProg vc.ProgramEntry

func validateRoomId(file string) error {
	if strings.ContainsAny(file, " !@#$%^&*()_+{}[]|\\<,>.?/") {
		return fmt.Errorf("ROOM ID CANNOT CONTAIN SPECIAL CHARACTERS OR SPACES %s", file)
	}
	return nil
}

func validateRoomName(name string) error {
	if len(name) < 5 {
		return fmt.Errorf("NAME %s MUST HAVE AT LEAST 5 CHARATERS", name)
	}
	return nil
}

var originalRoomId *string

func validateEditRoomId(id string) error {
	if id != *originalRoomId {
		return fmt.Errorf("CANNOT EDIT THE ROOM ID, VALUE MUST BE %s", *originalRoomId)
	}
	return nil
}

func NewRoomFormModel() NewRoomForm {
	roomOptions = &vc.RoomOptions{
		AddressSetsLocation: true,
	}

	p := progress.New(progress.WithDefaultGradient())
	p.Width = app.width

	return NewRoomForm{
		progress: p,
		running:  false,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("NAME").
					Title("Enter Friendly Name").
					Prompt("📛  ").
					Placeholder("My friendly program name").
					Validate(validateRoomName).
					Value(&roomOptions.Name),

				huh.NewInput().
					Key("ROOM ID").
					Title("Enter Room ID").
					Prompt("🆔  ").
					Placeholder("ROOM404").
					Validate(validateRoomId).
					Value(&roomOptions.ProgramInstanceId),

				huh.NewInput().
					Key("NOTES").
					Title("Enter Notes").
					Prompt("📝  ").
					Placeholder("My seemingly pointless notes").
					Value(&roomOptions.Notes),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func NewRoomFormModelWithPrograms(programs vc.Programs) NewRoomForm {
	roomOptions = &vc.RoomOptions{
		AddressSetsLocation: true,
	}

	p := progress.New(progress.WithDefaultGradient())
	p.Width = app.width
	progs := make([]vc.ProgramEntry, len(programs))
	copy(progs, programs)

	return NewRoomForm{
		progress: p,
		running:  false,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[vc.ProgramEntry]().
					Title("Select a program").
					Options(huh.NewOptions[vc.ProgramEntry](progs...)...).
					Value(&selectProg).
					Description("the selected program will be instantiated"),

				huh.NewInput().
					Key("NAME").
					Title("Enter Friendly Name").
					Prompt("📛  ").
					Placeholder("My friendly program name").
					Validate(validateRoomName).
					Value(&roomOptions.Name),

				huh.NewInput().
					Key("ROOM ID").
					Title("Enter Room ID").
					Prompt("🆔  ").
					Placeholder("ROOM404").
					Validate(validateRoomId).
					Value(&roomOptions.ProgramInstanceId),

				huh.NewInput().
					Key("NOTES").
					Title("Enter Notes").
					Prompt("📝  ").
					Placeholder("My seemingly pointless notes").
					Value(&roomOptions.Notes),

				huh.NewInput().
					Key("ADDRESS").
					Title("Location").
					Prompt("🏚  ").
					Placeholder("404 Bad Address Location").
					Value(&roomOptions.Location),

				huh.NewConfirm().
					Title("Address Sets Location").
					Value(&roomOptions.AddressSetsLocation),

				huh.NewInput().
					Key("TIMEZONE").
					Title("Time Zone").
					Prompt("⏲  ").
					Placeholder("+/- numeric value").
					Value(&roomOptions.TimeZone),

				huh.NewInput().
					Key("LAT").
					Title("Latitude").
					Prompt("🌐  ").
					Placeholder("39.352862").
					Value(&roomOptions.Latitude),

				huh.NewInput().
					Key("LONG").
					Title("Longitude").
					Prompt("🌐  ").
					Placeholder("-76.407341").
					Value(&roomOptions.Longitude),

				huh.NewInput().
					Key("USER_FILE").
					Title("Upload User File").
					Prompt("👤  ").
					Placeholder("/home/user/myconfig.json").
					Value(&roomOptions.UserFile),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func NewRoomFromProgramFormModel(programEntry *vc.ProgramEntry) NewRoomForm {
	roomOptions = &vc.RoomOptions{
		ProgramLibraryId:    int(programEntry.ProgramID),
		AddressSetsLocation: true,
	}

	p := progress.New(progress.WithDefaultGradient())
	p.Width = app.width

	return NewRoomForm{
		progress: p,
		running:  false,
		form: huh.NewForm(
			huh.NewGroup(

				huh.NewSelect[vc.ProgramEntry]().
					Title("Select a program").
					Options(huh.NewOption[vc.ProgramEntry](programEntry.ProgramName, *programEntry)).
					Value(&selectProg).
					Description(fmt.Sprintf("%s program will be instantiated", programEntry.FriendlyName)),

				huh.NewInput().
					Key("NAME").
					Title("Enter Friendly Name").
					Prompt("📛  ").
					Placeholder("My friendly program name").
					Validate(validateRoomName).
					Value(&roomOptions.Name),

				huh.NewInput().
					Key("ROOM ID").
					Title("Enter Room ID").
					Prompt("🆔  ").
					Placeholder("ROOM404").
					Validate(validateRoomId).
					Value(&roomOptions.ProgramInstanceId),

				huh.NewInput().
					Key("NOTES").
					Title("Enter Notes").
					Prompt("📝  ").
					Placeholder("My seemingly pointless notes").
					Value(&roomOptions.Notes),

				huh.NewInput().
					Key("ADDRESS").
					Title("Location").
					Prompt("🏚  ").
					Placeholder("404 Bad Address Location").
					Value(&roomOptions.Location),

				huh.NewConfirm().
					Title("Address Sets Location").
					Value(&roomOptions.AddressSetsLocation),

				huh.NewInput().
					Key("TIMEZONE").
					Title("Time Zone").
					Prompt("⏲  ").
					Placeholder("+/- numeric value").
					Value(&roomOptions.TimeZone),

				huh.NewInput().
					Key("LAT").
					Title("Latitude").
					Prompt("🌐  ").
					Placeholder("39.352862").
					Value(&roomOptions.Latitude),

				huh.NewInput().
					Key("LONG").
					Title("Longitude").
					Prompt("🌐  ").
					Placeholder("-76.407341").
					Value(&roomOptions.Longitude),

				huh.NewInput().
					Key("USER_FILE").
					Title("Upload User File").
					Prompt("👤  ").
					Placeholder("/home/user/myconfig.json").
					Value(&roomOptions.UserFile),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func EditRoomFormModel(room *vc.Room) NewRoomForm {
	originalRoomId = &room.ID
	roomOptions = &vc.RoomOptions{
		ProgramInstanceId:   room.ID,
		ProgramLibraryId:    int(room.ProgramID),
		Name:                room.Name,
		Notes:               room.Notes,
		Location:            room.Location,
		Latitude:            room.Latitude,
		Longitude:           room.Longitude,
		TimeZone:            room.TimeZone,
		AddressSetsLocation: true,
	}

	p := progress.New(progress.WithDefaultGradient())
	p.Width = app.width

	return NewRoomForm{
		edit:     true,
		running:  false,
		progress: p,
		form: huh.NewForm(
			huh.NewGroup(

				huh.NewInput().
					Key("NAME").
					Title("Enter Friendly Name").
					Prompt("📛  ").
					Placeholder("My friendly program name").
					Validate(validateRoomName).
					Value(&roomOptions.Name),

				huh.NewInput().
					Key("ROOM ID").
					Title("Enter Room ID").
					Prompt("🆔  ").
					Placeholder(*originalRoomId).
					Validate(validateEditRoomId).
					Value(&roomOptions.ProgramInstanceId),

				huh.NewInput().
					Key("NOTES").
					Title("Enter Notes").
					Prompt("📝  ").
					Placeholder("My seemingly pointless notes").
					Value(&roomOptions.Notes),

				huh.NewInput().
					Key("ADDRESS").
					Title("Location").
					Prompt("🏚  ").
					Placeholder("404 Bad Address Location").
					Value(&roomOptions.Location),

				huh.NewConfirm().
					Title("Address Sets Location").
					Value(&roomOptions.AddressSetsLocation),

				huh.NewInput().
					Key("TIMEZONE").
					Title("Time Zone").
					Prompt("⏲  ").
					Placeholder("+/- numeric value").
					Value(&roomOptions.TimeZone),

				huh.NewInput().
					Key("LAT").
					Title("Latitude").
					Prompt("🌐  ").
					Placeholder("39.352862").
					Value(&roomOptions.Latitude),

				huh.NewInput().
					Key("LONG").
					Title("Longitude").
					Prompt("🌐  ").
					Placeholder("-76.407341").
					Value(&roomOptions.Longitude),

				huh.NewInput().
					Key("USER_FILE").
					Title("Upload User File").
					Prompt("👤  ").
					Placeholder("/home/user/myconfig.json").
					Value(&roomOptions.UserFile),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func (m NewRoomForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m NewRoomForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case error:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 20
		return m, nil

	case vc.Programs:

		return NewRoomFormModelWithPrograms(msg), nil

	case vc.RoomCreatedResult:
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
		return m, tea.Batch(roomCreatedTickCmd(), m.progress.IncrPercent(0.20))

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+q":
			return ReturnRoomsModel(), tea.Batch(tick, RoomsQuery)

		case "ctrl+n":
			form := NewRoomFormModel()
			return form, tea.Batch(form.Init(), ProgramsQuery)
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f

		if m.form.State == huh.StateCompleted && !m.running {
			m.running = true

			if !m.edit {
				roomOptions.ProgramLibraryId = int(selectProg.ProgramID)
				return m, tea.Batch(CreateRoom(*roomOptions), roomCreatedTickCmd())
			}
			return m, tea.Batch(EditRoom(*roomOptions), roomCreatedTickCmd())
		}
	}
	return m, cmd
}

func (m NewRoomForm) View() string {
	s := ""
	if m.edit {
		s += GreyedOutText.Render("\n🖊 Edit Running Program Instance\n")
	} else {
		s += GreyedOutText.Render("\n🆕 Create New Program Instance\n")
	}

	s += "\n" + m.form.View()

	if m.progress.Percent() != 0.0 {
		s += "\n" + m.progress.View() + "\n\n"
	}
	if m.err != nil {
		s += RenderErrorBox("error uploading new program file", m.err)
		s += GreyedOutText.Render("\n\n esc return * ctrl+n reset form")
		return s
	}

	if m.result != nil {
		var resultMessage string
		resultMessage += "\n RESULT           " + m.result.Message
		resultMessage += "\n\n STATUS CODE:     " + fmt.Sprintf("%d", m.result.Code)
		resultMessage += "\n"

		s += RenderMessageBox(1000).Render(resultMessage)
		s += GreyedOutText.Render("\n\n esc return * ctrl+n reset form")
	}

	return s
}

func roomCreatedTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return progressTick(t)
	})
}
