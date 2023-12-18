package vc

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

const (
	Running          RoomStatus = "Running"
	Aborted          RoomStatus = "Aborted"
	Stopped          RoomStatus = "Stopped"
	Stopping         RoomStatus = "Stopping"
	Starting         RoomStatus = "Starting"
	PROGRAMINSTANCES            = "ProgramInstance"
)

type VcRoomApi interface {
	GetRooms() (Rooms, VirtualControlError)
	StartRoom(id string) (bool, VirtualControlError)
	StopRoom(id string) (bool, VirtualControlError)
	DebugRoom(id string, enable bool) (bool, VirtualControlError)
	RestartRoom(id string) (bool, VirtualControlError)
	CreateRoom(options RoomOptions) (RoomCreatedResult, VirtualControlError)
	EditRoom(options RoomOptions) (RoomCreatedResult, VirtualControlError)
	DeleteRoom(id string) VirtualControlError
}

func (v *VC) GetRooms() (Rooms, VirtualControlError) {

	rooms, err := getProgramInstances(v)
	if err != nil {
		return make(Rooms, 0), err
	}

	programs, err := getProgramLibrary(v)
	if err != nil {
		return make(Rooms, 0), err
	}

	return mapRoomsToPrograms(rooms, programs)
}

func mapRoomsToPrograms(instances ProgramInstanceLibrary, library ProgramsLibrary) (rooms Rooms, err error) {

	var roomsModel = make(Rooms, 0)

	for _, r := range instances {
		id := fmt.Sprint(r.ProgramLibraryID)
		prog, hasMatch := library[id]

		if hasMatch {
			roomsModel = append(roomsModel, NewRoom(r, prog))
			continue
		}
		return roomsModel, fmt.Errorf("Room %s has no matching code %s, WTF", r.ProgramInstanceID, r.ProgramInstanceID)
	}

	comparById := func(a, b Room) int {
		return cmp.Compare(a.ID, b.ID)
	}

	slices.SortFunc(roomsModel, comparById)

	return roomsModel, nil
}

func (v *VC) StartRoom(id string) (bool, VirtualControlError) {
	return putRoomAction(v, id, "Start", true)
}
func (v *VC) StopRoom(id string) (bool, VirtualControlError) {
	return putRoomAction(v, id, "Stop", true)
}
func (v *VC) RestartRoom(id string) (bool, VirtualControlError) {
	return putRoomAction(v, id, "Restart", true)
}
func (v *VC) DebugRoom(id string, enable bool) (bool, VirtualControlError) {
	return putRoomAction(v, id, "DebuggingEnabled", enable)
}

func (v *VC) CreateRoom(options RoomOptions) (RoomCreatedResult, VirtualControlError) {
	return postRoom(v, options)
}

func (v *VC) EditRoom(options RoomOptions) (RoomCreatedResult, VirtualControlError) {
	return putRoom(v, options)
}

func (v *VC) DeleteRoom(id string) VirtualControlError {
	return deleteRoom(v, id)
}

func getProgramInstances(server *VC) (ProgramInstanceLibrary, VirtualControlError) {

	var results ProgramInstanceResponse
	err := server.getBody(PROGRAMINSTANCES, &results)

	if err != nil {
		return ProgramInstanceLibrary{}, NewServerError(500, err)
	}
	return results.Device.Programs.ProgramInstanceLibrary, nil
}

func putRoomAction(server *VC, id string, action string, state bool) (bool, VirtualControlError) {

	var stateVal string
	if state {
		stateVal = "true"
	} else {
		stateVal = "false"
	}
	apiUrl := server.url
	resource := PROGRAMINSTANCES
	data := url.Values{}
	data.Set("ProgramInstanceId", id)
	data.Set(action, stateVal)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource

	req, err := http.NewRequest("PUT", server.url+PROGRAMINSTANCES, strings.NewReader(data.Encode()))
	if err != nil {
		return false, NewServerError(500, err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := server.client.Do(req)
	if err != nil {

		return false, NewServerError(resp.StatusCode, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, NewServerError(resp.StatusCode, errors.New("FAILED HTTP REQUEST"))
	}

	return true, nil
}

func postRoom(vc *VC, options RoomOptions) (result RoomCreatedResult, err error) {

	form := &bytes.Buffer{}
	writer := multipart.NewWriter(form)

	addFormField(writer, "Name", options.Name)
	addFormField(writer, "ProgramInstanceId", options.ProgramInstanceId)
	addFormField(writer, "ProgramLibraryId", fmt.Sprintf("%d", options.ProgramLibraryId))

	err = writer.Close()
	if err != nil {
		return RoomCreatedResult{}, err
	}

	request, err := http.NewRequest("POST", vc.url+PROGRAMINSTANCES, form)
	if err != nil {
		return RoomCreatedResult{}, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := vc.client.Do(request)
	if err != nil {
		return RoomCreatedResult{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return RoomCreatedResult{}, NewServerError(response.StatusCode, errors.New("FAILED TO UPLOAD FILE"))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return RoomCreatedResult{}, NewServerError(response.StatusCode, err)
	}

	actions := ActionResponse[any]{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return RoomCreatedResult{}, NewServerError(response.StatusCode, err)
	}

	actionResult := actions.Actions[0].Results[0]

	// UHM THIS IS A WAY WAY BROKEN REST API
	if actionResult.StatusID != 0 {
		return RoomCreatedResult{}, fmt.Errorf("FAILED CREATING NEW PROGRAM INSTANCE\n\nREASON: %s", actionResult.StatusInfo)
	}

	return RoomCreatedResult{
		Message: actionResult.StatusInfo,
		Code:    int(actionResult.StatusID),
		Success: actionResult.StatusID == 0,
	}, nil
}

func putRoom(vc *VC, options RoomOptions) (result RoomCreatedResult, err error) {

	form := &bytes.Buffer{}
	writer := multipart.NewWriter(form)

	addFormField(writer, "Name", options.Name)
	addFormField(writer, "ProgramInstanceId", options.ProgramInstanceId)
	addFormField(writer, "ProgramLibraryId", fmt.Sprintf("%d", options.ProgramLibraryId))

	err = writer.Close()
	if err != nil {
		return RoomCreatedResult{}, err
	}

	request, err := http.NewRequest("PUT", vc.url+PROGRAMINSTANCES, form)
	if err != nil {
		return RoomCreatedResult{}, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := vc.client.Do(request)
	if err != nil {
		return RoomCreatedResult{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return RoomCreatedResult{}, NewServerError(response.StatusCode, errors.New("FAILED TO UPLOAD FILE"))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return RoomCreatedResult{}, NewServerError(response.StatusCode, err)
	}

	actions := ActionResponse[any]{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return RoomCreatedResult{}, NewServerError(response.StatusCode, err)
	}

	actionResult := actions.Actions[0].Results[0]

	// UHM THIS IS A WAY WAY BROKEN REST API
	if actionResult.StatusID != 0 {
		return RoomCreatedResult{}, fmt.Errorf("FAILED CREATING NEW PROGRAM INSTANCE\n\nREASON: %s", actionResult.StatusInfo)
	}

	return RoomCreatedResult{
		Message: actionResult.StatusInfo,
		Code:    int(actionResult.StatusID),
		Success: actionResult.StatusID == 0,
	}, nil
}

func deleteRoom(vc *VC, id string) (err error) {

	form := &bytes.Buffer{}
	writer := multipart.NewWriter(form)

	err = writer.Close()
	if err != nil {
		return err
	}

	request, err := http.NewRequest("DEL", vc.url+PROGRAMINSTANCES, form)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := vc.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return NewServerError(response.StatusCode, errors.New("FAILED TO UPLOAD FILE"))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return NewServerError(response.StatusCode, err)
	}

	actions := ActionResponse[any]{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return NewServerError(response.StatusCode, err)
	}

	actionResult := actions.Actions[0].Results[0]

	// UHM THIS IS A WAY WAY BROKEN REST API
	if actionResult.StatusID != 0 {
		return fmt.Errorf("FAILED UPLOADING NEW PROGRAM \n\nREASON: %s", actionResult.StatusInfo)
	}

	return nil
}

type RoomStatus string

type ProgramInstanceResponse struct {
	Device DeviceProgramInstances `json:"Device"`
}

type DeviceProgramInstances struct {
	Programs ProgramInstances `json:"Programs"`
}

type ProgramInstances struct {
	ProgramInstanceLibrary ProgramInstanceLibrary `json:"ProgramInstanceLibrary"`
}

type ProgramInstanceLibrary map[string]ProgramInstance

type ProgramInstance struct {
	ID                  int64  `json:"id"`
	ProgramInstanceID   string `json:"ProgramInstanceId"`
	UserFile            string `json:"UserFile"`
	Status              string `json:"Status"`
	Name                string `json:"Name"`
	ProgramLibraryID    int    `json:"ProgramLibraryId"`
	Level               string `json:"Level"`
	AddressSetsLocation bool   `json:"AddressSetsLocation"`
	Location            string `json:"Location"`
	Longitude           string `json:"Longitude"`
	Latitude            string `json:"Latitude"`
	TimeZone            string `json:"TimeZone"`
	ConfigurationLink   string `json:"ConfigurationLink"`
	XpanelURL           string `json:"XpanelUrl"`
	Notes               string `json:"Notes"`
	DebuggingEnabled    bool   `json:"DebuggingEnabled"`
}

type Rooms []Room

// This is the actual model used by the TUI.
// This model includes all the information we need about a room and the program it is running.
type Room struct {
	ID                string
	Name              string
	Status            string
	Debugging         bool
	Location          string
	Longitude         string
	Latitude          string
	TimeZone          string
	ConfigurationLink string
	XpanelURL         string
	Notes             string

	ProgramID       int16  `json:"ProgramId"`
	ProgramName     string `json:"ProgramName"`
	ProgramFriendly string `json:"FriendlyName"`
	ProgramType     string `json:"ProgramType"`
	CompileDateTime string `json:"CompileDateTime"`
}

func NewRoom(i ProgramInstance, p ProgramEntry) Room {
	return Room{
		ID:                i.ProgramInstanceID,
		Name:              i.Name,
		Status:            i.Status,
		Debugging:         i.DebuggingEnabled,
		Location:          i.Location,
		Longitude:         i.Longitude,
		Latitude:          i.Latitude,
		TimeZone:          i.TimeZone,
		ConfigurationLink: i.ConfigurationLink,
		XpanelURL:         i.XpanelURL,
		Notes:             i.Notes,

		ProgramID:       p.ProgramID,
		ProgramFriendly: p.FriendlyName,
		ProgramType:     p.ProgramType,
		ProgramName:     p.ProgramName,
		CompileDateTime: p.CompileDateTime,
	}
}

type RoomCreatedResult struct {
	Success bool
	Message string
	Code    int
}
