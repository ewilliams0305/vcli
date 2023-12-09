package vc

type DeviceInformationResponse struct {
	Device Device `json:"Device"`
}

type Device struct {
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
