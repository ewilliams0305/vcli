package vc4

import (
	"encoding/json"
	"io"
)

func GetDeviceInfo() (*DeviceInfo, error) {

	resp, err := client.Get(makeUrl(DEVICEINFO))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, newResponseError(resp.StatusCode)
	}

	var deviceData Device
	err = json.Unmarshal(body, &deviceData)
	if err != nil {
		return nil, err
	}

	return &deviceData.DeviceInfo, nil
}
