package vc

import (
	"cmp"
	"errors"
	"fmt"
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
	RestartRoom(id string) (bool, VirtualControlError)
	// EditRoom(id string, name string, notes string) (ActionResult, VirtualControlError)
	// AddRoom(id string) (ActionResult, VirtualControlError)
	// DeleteRoom(id string) (ActionResult, VirtualControlError)
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

	var roomsModel = make([]Room, 0)

	for _, r := range instances {
		id := fmt.Sprint(r.ProgramLibraryID)
		prog, hasMatch := library[id]

		if hasMatch {
			roomsModel = append(roomsModel, NewRoom(r, prog))
			continue
		}
		return roomsModel, fmt.Errorf("Room %s has no matching code %s, WTF!", r.ProgramInstanceID, r.ProgramInstanceID)
	}

	comparById := func(a, b Room) int {
		return cmp.Compare(a.ID, b.ID)
	}

	slices.SortFunc(roomsModel, comparById)

	return roomsModel, nil
}

func (v *VC) StartRoom(id string) (bool, VirtualControlError) {
	return putRoomAction(v, id, "Start")
}
func (v *VC) StopRoom(id string) (bool, VirtualControlError) {
	return putRoomAction(v, id, "Stop")
}
func (v *VC) RestartRoom(id string) (bool, VirtualControlError) {
	return putRoomAction(v, id, "Restart")
}

func getProgramInstances(server *VC) (ProgramInstanceLibrary, VirtualControlError) {

	var results ProgramInstanceResponse
	err := server.getBody(PROGRAMINSTANCES, &results)

	if err != nil {
		return ProgramInstanceLibrary{}, NewServerError(500, err)
	}
	return results.Device.Programs.ProgramInstanceLibrary, nil
}

func putRoomAction(server *VC, id string, action string) (bool, VirtualControlError) {

	apiUrl := server.url
	resource := PROGRAMINSTANCES
	data := url.Values{}
	data.Set("ProgramInstanceId", id)
	data.Set(action, "true")

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
