package vc

type DeviceInfo struct {
	DeviceKey          string `json:"DeviceKey"`
	DeviceId           string `json:"DeviceId"`
	ApplicationVersion string `json:"ApplicationVersion"`
	Manufacturer       string `json:"Manufacturer"`
	ID                 string `json:"ID"`
	MacAddress         string `json:"MacAddress"`
	Model1             string `json:"Model1"`
	Category           string `json:"Category"`
	BuildDate          string `json:"BuildDate"`
	Version            string `json:"Version"`
	Name               string `json:"Name"`
	PythonVersion      string `json:"PythonVersion"`
	MonoVersion        string `json:"MonoVersion"`
}

type Device struct {
	DeviceInfo DeviceInfo `json:"DeviceInfo"`
}
