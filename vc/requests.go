package vc

import (
	"encoding/json"
	"errors"
	"io"
)

const (
	DEVICEINFO       = "DeviceInfo"
	PROGRAMINSTANCES = "ProgramInstance"
	PROGRAMLIBRARY   = "ProgramLibrary"
)

func (vc *vc) getBody(url string, result any) (err VirtualControlError) {

	resp, err := vc.client.Get(vc.url + url)
	if err != nil {
		return newServerError(500, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return newServerError(resp.StatusCode, errors.New("FAILED TO GET DEVICE INFO"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return newServerError(resp.StatusCode, err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return newServerError(resp.StatusCode, err)
	}

	return nil
}

func getDeviceInfo(vc *vc) (DeviceInfo, VirtualControlError) {

	var deviceData DeviceInformationResponse
	err := vc.getBody(DEVICEINFO, &deviceData)

	if err != nil {
		return emptyDeviceInfo(), newServerError(500, err)
	}

	return deviceData.Device.DeviceInfo, nil
}

func getProgramInstances(vc *vc) (ProgramInstanceLibrary, VirtualControlError) {

	var results ProgramInstanceResponse
	err := vc.getBody(PROGRAMINSTANCES, &results)

	if err != nil {
		return ProgramInstanceLibrary{}, newServerError(500, err)
	}
	return results.Device.Programs.ProgramInstanceLibrary, nil
}

func getProgramLibrary(vc *vc) (ProgramsLibrary, VirtualControlError) {

	var results ProgramLibraryResponse
	err := vc.getBody(PROGRAMLIBRARY, &results)

	if err != nil {
		return ProgramsLibrary{}, newServerError(500, err)
	}
	return results.Device.Programs.ProgramLibrary, nil
}
