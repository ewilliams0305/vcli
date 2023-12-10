package vc

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
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

func putRoomAction(vc *vc, id string, action string) (bool, VirtualControlError) {

	apiUrl := vc.url
	resource := PROGRAMINSTANCES
	data := url.Values{}
	data.Set("ProgramInstanceId", id)
	data.Set(action, "true")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	//urlStr := u.String() // "https://api.com/user/"

	req, err := http.NewRequest("PUT", vc.url+PROGRAMINSTANCES, strings.NewReader(data.Encode()))
	if err != nil {
		return false, newServerError(500, err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := vc.client.Do(req)
	if err != nil {

		return false, newServerError(resp.StatusCode, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, newServerError(resp.StatusCode, errors.New("FAILED HTTP REQUEST"))
	}

	return true, nil
}
