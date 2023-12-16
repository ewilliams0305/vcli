package vc

const (
	DEVICEINFO = "DeviceInfo"
)

type VcInfoApi interface {
	DeviceInfo() (DeviceInfo, VirtualControlError)
}

func (v *VC) DeviceInfo() (DeviceInfo, VirtualControlError) {
	return getDeviceInfo(v)
}

func getDeviceInfo(vc *VC) (DeviceInfo, VirtualControlError) {

	var deviceData DeviceInformationResponse
	err := vc.getBody(DEVICEINFO, &deviceData)

	if err != nil {
		return emptyDeviceInfo(), NewServerError(500, err)
	}

	return deviceData.Device.DeviceInfo, nil
}

type DeviceInformationResponse struct {
	Device DeviceContext `json:"Device"`
}

type DeviceContext struct {
	DeviceInfo DeviceInfo `json:"DeviceInfo"`
}

type DeviceInfo struct {
	ID                 string `json:"ID"`
	Model              string `json:"Model"`
	Category           string `json:"Category"`
	Manufacturer       string `json:"Manufacturer"`
	DeviceID           string `json:"DeviceId"`
	Name               string `json:"Name"`
	ApplicationVersion string `json:"ApplicationVersion"`
	BuildDate          string `json:"BuildDate"`
	DeviceKey          string `json:"DeviceKey"`
	MACAddress         string `json:"MacAddress"`
	Version            string `json:"Version"`
	PythonVersion      string `json:"PythonVersion"`
	MonoVersion        string `json:"MonoVersion"`
}

func emptyDeviceInfo() DeviceInfo {
	return DeviceInfo{}
}
