package vc

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	PROGRAMLIBRARY = "ProgramLibrary"
)

type VcProgramApi interface {
	GetPrograms() (Programs, VirtualControlError)
	CreateProgram(options ProgramOptions) (status int, err VirtualControlError)
}

func (v *VC) GetPrograms() (Programs, VirtualControlError) {
	progs, err := getProgramLibrary(v)
	if err != nil {
		return make(Programs, 0), err
	}

	p := make([]ProgramEntry, 0, len(progs))
	for _, value := range progs {
		p = append(p, value)
	}

	comparById := func(a, b ProgramEntry) int {
		return cmp.Compare(a.ProgramID, b.ProgramID)
	}

	slices.SortFunc(p, comparById)

	return p, nil
}

func (v *VC) CreateProgram(options ProgramOptions) (status int, err VirtualControlError) {
	return postProgram(v, options)
}

func getProgramLibrary(vc *VC) (ProgramsLibrary, VirtualControlError) {

	var results ProgramLibraryResponse
	err := vc.getBody(PROGRAMLIBRARY, &results)

	if err != nil {
		return ProgramsLibrary{}, NewServerError(500, err)
	}
	return results.Device.Programs.ProgramLibrary, nil
}

// UPLOADS A NEW PROGRAM TO THE APPLIANCE
func postProgram(vc *VC, options ProgramOptions) (status int, err error) {

	if !programIsValid(options.AppFile) {
		return 0, errors.New("INVALID FILE EXTENSION")
	}

	file, err := os.Open(options.AppFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	form := &bytes.Buffer{}
	writer := multipart.NewWriter(form)

	part, err := writer.CreateFormFile("AppFile", filepath.Base(options.AppFile))
	if err != nil {
		return 0, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return 0, err
	}

	file.Close()

	addFormField(writer, "filetype", "AppFile")
	addFormField(writer, "FriendlyName", options.Name)
	addFormField(writer, "Notes", options.Notes)

	err = writer.Close()
	if err != nil {
		return 0, err
	}

	request, err := http.NewRequest("POST", vc.url+PROGRAMLIBRARY, form)
	if err != nil {
		return 0, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := vc.client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return response.StatusCode, NewServerError(response.StatusCode, errors.New("FAILED TO UPLOAD FILE"))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, NewServerError(response.StatusCode, err)
	}

	actions := ActionResponse[ProgramEntry]{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return 500, NewServerError(response.StatusCode, err)
	}

	return response.StatusCode, NewServerError(response.StatusCode, errors.New("FILE FAILED TO UPLOAD"))
}

func programIsValid(file string) bool {

	return strings.HasSuffix(file, ".cpz") ||
		strings.HasSuffix(file, ".zip") ||
		strings.HasSuffix(file, ".lpz")
}

func addFormField(writer *multipart.Writer, key string, value string) {

	fieldWriter, err := writer.CreateFormField(key)
	if err != nil {
		//updater.Logger.Warn("Error creating form field:", err)
		return
	}
	_, err = fieldWriter.Write([]byte(value))
	if err != nil {
		//updater.Logger.Warn("Error writing form field value:", err)
		return
	}
}

type ProgramLibraryResponse struct {
	Device LibraryContext `json:"Device"`
}

type LibraryContext struct {
	Programs ProgramsContext `json:"Programs"`
}

type ProgramsContext struct {
	ProgramLibrary ProgramsLibrary `json:"ProgramLibrary"`
}

type ProgramsLibrary map[string]ProgramEntry

type Programs []ProgramEntry

type ProgramEntry struct {
	ProgramID         int16  `json:"ProgramId"`
	FriendlyName      string `json:"FriendlyName"`
	Notes             string `json:"Notes"`
	AppFile           string `json:"AppFile"`
	AppFileTS         string `json:"AppFileTS"`
	MobilityFile      string `json:"MobilityFile"`
	MobilityFileTS    string `json:"MobilityFileTS"`
	WebxPanelFile     string `json:"WebxPanelFile"`
	WebxPanelFileTS   string `json:"WebxPanelFileTS"`
	ProjectFile       string `json:"ProjectFile"`
	ProjectFileTS     string `json:"ProjectFileTS"`
	CwsFile           string `json:"CwsFile"`
	CwsFileTS         string `json:"CwsFileTS"`
	ProgramType       string `json:"ProgramType"`
	ProgramName       string `json:"ProgramName"`
	CompileDateTime   string `json:"CompileDateTime"`
	CresDBVersion     string `json:"CresDBVersion"`
	DeviceDBVersion   string `json:"DeviceDBVersion"`
	IncludeDATVersion string `json:"IncludeDatVersion"`
}
