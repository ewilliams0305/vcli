package vc

const (
	Running RoomStatus = "Running"
	Aborted RoomStatus = "Aborted"
	Stopped RoomStatus = "Stopped"
)

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
	ProgramLibraryID    int64  `json:"ProgramLibraryId"`
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
