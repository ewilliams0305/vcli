package vc

import (
	"fmt"
	"net/http"
)

const (
	virtualBaseUrl = "/VirtualControl/config/api/"
)

// The virtual control server will need to be controlled one of two ways:
//
// - https with a authorization header
// - http local host without auth.
//
// Controlling local host :
// -- requires the /VirtualControl/config/api/ be removed from the route
// -- Local host port 5000
// -- use http NOT https
//
// Controlling external:
// -- full path urls https://[ServerURL]/VirtualControl/config/api/
// -- use of https, I'm not writing this for unsecured servers.
// -- accept self signed certs, this si way too common not too!
// -- header  "Authorization: [Token]"
type VirtualControl interface {
	Config() *VirtualConfig

	VcProgramApi
	VcInfoApi
	VcRoomApi
}

type VC struct {
	client   *http.Client
	url      string
	http     bool
	port     int
	hostname string
	token    string
}

type VirtualConfig struct {
	http     bool
	port     *int
	hostname *string
	token    *string
}

// Create VC Clients
func NewLocalVC() VirtualControl {
	return &VC{
		client:   createLocalClient(),
		url:      LOCALHOSTURL,
		http:     true,
		port:     5000,
		hostname: "127.0.0.1",
		token:    "",
	}
}

func NewRemoteVC(host string, token string) VirtualControl {
	return &VC{
		client:   createRemoteClient(token),
		url:      baseUrl(host),
		http:     false,
		port:     5000,
		hostname: host,
		token:    token,
	}
}

func baseUrl(host string) string {
	return fmt.Sprintf("https://%s%s", host, virtualBaseUrl)
}

// Implement the VC Interface
func (v *VC) Config() *VirtualConfig {
	return &VirtualConfig{
		http:     v.http,
		port:     &v.port,
		hostname: &v.hostname,
		token:    &v.token,
	}
}

type ActionResponse struct {
	Actions []ActionData `json:"Actions"`
}

type ActionData struct {
	Operation    string                 `json:"Operation"`
	Results      []ActionResponseResult `json:"Results"`
	TargetObject string                 `json:"TargetObject"`
	Version      string                 `json:"Version"`
}

type ActionResponseResult struct {
	Path       string         `json:"path"`
	Object     ActionRoomData `json:"object"`
	StatusInfo string         `json:"StatusInfo"`
	StatusID   int64          `json:"StatusId"`
}

type ActionRoomData struct {
	ID                  int64  `json:"id"`
	ProgramInstanceID   string `json:"programInstanceId"`
	ProgramLibraryID    int64  `json:"ProgramLibraryId"`
	UserFile            string `json:"UserFile"`
	Status              string `json:"Status"`
	Name                string `json:"Name"`
	Level               string `json:"level"`
	AddressSetsLocation bool   `json:"AddressSetsLocation"`
	Location            string `json:"Location"`
	Longitude           string `json:"Longitude"`
	Latitude            string `json:"Latitude"`
	TimeZone            string `json:"TimeZone"`
	ConfigurationLink   string `json:"ConfigurationLink"`
	ProcessID           string `json:"ProcessId"`
	XpanelURL           string `json:"XpanelUrl"`
	ProgramSlotID       int64  `json:"ProgramSlotId"`
	Notes               string `json:"Notes"`
	DebuggingEnabled    bool   `json:"DebuggingEnabled"`
}
