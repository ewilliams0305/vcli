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
	DeviceInfo() (DeviceInfo, VirtualControlError)
	ProgramLibrary() (ProgramsLibrary, VirtualControlError)
 LicenseInfo() (LicenseInfo, VirtualControlError)
 
 VcRoomApi
}

type VcRoomApi interface {
 ProgramInstances() (ProgramInstanceLibrary, VirtualControlError)
 StartRoom(id string) (bool, VirtualControlError)
	StopRoom(id string) (bool, VirtualControlError)
	RestartRoom(id string) (bool, VirtualControlError)
 EditRoom(id string, name string, notes string) (ActionResult, VirtualControlError)
 AddRoom(id string) (ActionResult, VirtualControlError)
 DeleteRoom(id string) (ActionResult, VirtualControlError)  
}

type vc struct {
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

type ActionResponse struct {
 Status string
 Result string
}

// Create VC Clients

func NewLocalVC() VirtualControl {
	return &vc{
		client:   createLocalClient(),
		url:      LOCALHOSTURL,
		http:     true,
		port:     5000,
		hostname: "127.0.0.1",
		token:    "",
	}
}

func NewRemoteVC(host string, token string) VirtualControl {
	return &vc{
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
func (v *vc) Config() *VirtualConfig {
	return &VirtualConfig{
		http:     v.http,
		port:     &v.port,
		hostname: &v.hostname,
		token:    &v.token,
	}
}

func (v *vc) DeviceInfo() (DeviceInfo, VirtualControlError) {
	return getDeviceInfo(v)
}
func (v *vc) ProgramInstances() (ProgramInstanceLibrary, VirtualControlError) {
	return getProgramInstances(v)
}
func (v *vc) ProgramLibrary() (ProgramsLibrary, VirtualControlError) {
	return getProgramLibrary(v)
}
func (v *vc) StartRoom(id string) (bool, VirtualControlError) {
	return putRoomAction(v, id, "Start")
}
func (v *vc) StopRoom(id string) (bool, VirtualControlError) {
	return putRoomAction(v, id, "Stop")
}
func (v *vc) RestartRoom(id string) (bool, VirtualControlError) {
	return putRoomAction(v, id, "Restart")
}
