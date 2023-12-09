package vc

import (
	"encoding/json"
	"errors"
	"io"
)

const (
	DEVICEINFO = "DeviceInfo"
)

func getDeviceInfo(vc *vc) (DeviceInfo, VirtualControlError) {

	resp, err := vc.client.Get(vc.url + DEVICEINFO)
	if err != nil {
		return emptyDeviceInfo(), newServerError(500, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return emptyDeviceInfo(), newServerError(resp.StatusCode, errors.New("FAILED TO GET DEVICE INFO"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return emptyDeviceInfo(), newServerError(resp.StatusCode, err)
	}

	var deviceData DeviceInformationResponse
	err = json.Unmarshal(body, &deviceData)
	if err != nil {
		return emptyDeviceInfo(), newServerError(resp.StatusCode, err)
	}

	return deviceData.Device.DeviceInfo, nil
}

func emptyDeviceInfo() DeviceInfo {
	return DeviceInfo{
		DeviceKey:  "",
		MACAddress: "00:00:00:00:00:00",
	}
}
