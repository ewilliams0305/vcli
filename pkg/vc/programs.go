package vc

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
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
	CreateProgram(options ProgramOptions) (result ProgramUploadResult, err VirtualControlError)
	EditProgram(options ProgramOptions) (result ProgramUploadResult, err VirtualControlError)
	DeleteProgram(id int) (result ProgramDeleteResult, err VirtualControlError)
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

func (v *VC) CreateProgram(options ProgramOptions) (result ProgramUploadResult, err VirtualControlError) {
	return postProgram(v, options)
}

func (v *VC) EditProgram(options ProgramOptions) (result ProgramUploadResult, err VirtualControlError) {
	return editProgram(v, options)
}

func (v *VC) DeleteProgram(id int) (result ProgramDeleteResult, err VirtualControlError) {
	return deleteProgram(v, id)
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
func postProgram(vc *VC, options ProgramOptions) (result ProgramUploadResult, err error) {

	if !programIsValid(options.AppFile) {
		return ProgramUploadResult{}, errors.New("INVALID FILE EXTENSION")
	}

	file, err := os.Open(options.AppFile)
	if err != nil {
		return ProgramUploadResult{}, err
	}
	defer file.Close()

	form := &bytes.Buffer{}
	writer := multipart.NewWriter(form)

	part, err := writer.CreateFormFile("AppFile", filepath.Base(options.AppFile))
	if err != nil {
		return ProgramUploadResult{}, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return ProgramUploadResult{}, err
	}

	file.Close()

	addFormField(writer, "filetype", "AppFile")
	addFormField(writer, "FriendlyName", options.Name)
	addFormField(writer, "Notes", options.Notes)

	err = writer.Close()
	if err != nil {
		return ProgramUploadResult{}, err
	}

	request, err := http.NewRequest("POST", vc.url+PROGRAMLIBRARY, form)
	if err != nil {
		return ProgramUploadResult{}, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := vc.client.Do(request)
	if err != nil {
		return ProgramUploadResult{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return ProgramUploadResult{}, NewServerError(response.StatusCode, errors.New("FAILED TO UPLOAD FILE"))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ProgramUploadResult{}, NewServerError(response.StatusCode, err)
	}

	actions := ActionResponse[any]{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return ProgramUploadResult{}, NewServerError(response.StatusCode, err)
	}

	actionResult := actions.Actions[0].Results[0]

	// UHM THIS IS A WAY WAY BROKEN REST API
	if actionResult.StatusID != 0 {
		return ProgramUploadResult{}, fmt.Errorf("FAILED UPLOADING NEW PROGRAM \n\nREASON: %s", actionResult.StatusInfo)
	}

	actionProgram := ActionResponse[ProgramEntry]{}
	err = json.Unmarshal(body, &actionProgram)
	if err != nil {
		return ProgramUploadResult{}, NewServerError(response.StatusCode, err)
	}

	return NewProgramUploadResult(&actionProgram), nil
}

func editProgram(vc *VC, options ProgramOptions) (result ProgramUploadResult, err error) {

	form := &bytes.Buffer{}
	writer := multipart.NewWriter(form)
	defer writer.Close()

	_, err = validateAndCreateFileHeader(writer, options.AppFile, "AppFile", []string{".cpz", ".lpz", ".zip"})
	_, err = validateAndCreateFileHeader(writer, options.MobilityFile, "MobilityFile", []string{".zip", ".Core3z"})
	_, err = validateAndCreateFileHeader(writer, options.ProjectFile, "ProjectFile", []string{".zip", "ch5z", ".vtz", ".Core3z"})
	_, err = validateAndCreateFileHeader(writer, options.WebxPanelFile, "WebxPanelFile", []string{".zip", "ch5z", ".vtz", ".Core3z"})
	_, err = validateAndCreateFileHeader(writer, options.CwsFile, "CwsFile", []string{".zip"})

	if err != nil {
		return ProgramUploadResult{}, err
	}

	if len(options.MobilityFile) > 0 {
		addFormFile(options.MobilityFile, "MobilityFile", writer)
	}
	addFormField(writer, "ProgramId", fmt.Sprintf("%d", options.ProgramId))
	addFormField(writer, "FriendlyName", options.Name)
	addFormField(writer, "Notes", options.Notes)

	if options.StartNow {
		addFormField(writer, "StartNow", "true")
	}

	err = writer.Close()
	if err != nil {
		return ProgramUploadResult{}, err
	}

	request, err := http.NewRequest("PUT", vc.url+PROGRAMLIBRARY, form)
	if err != nil {
		return ProgramUploadResult{}, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := vc.client.Do(request)
	if err != nil {
		return ProgramUploadResult{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return ProgramUploadResult{}, NewServerError(response.StatusCode, errors.New("FAILED TO UPLOAD FILE"))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ProgramUploadResult{}, NewServerError(response.StatusCode, err)
	}

	actions := ActionResponse[any]{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return ProgramUploadResult{}, NewServerError(response.StatusCode, err)
	}

	actionResult := actions.Actions[0].Results[0]

	// UHM THIS IS A WAY WAY BROKEN REST API
	if actionResult.StatusID != 0 {
		return ProgramUploadResult{}, fmt.Errorf("FAILED UPLOADING NEW PROGRAM \n\nREASON: %s", actionResult.StatusInfo)
	}

	actionProgram := ActionResponse[ProgramEntry]{}
	err = json.Unmarshal(body, &actionProgram)
	if err != nil {
		return ProgramUploadResult{}, NewServerError(response.StatusCode, err)
	}

	return NewProgramUploadResult(&actionProgram), nil
}

func deleteProgram(vc *VC, id int) (result ProgramDeleteResult, err error) {

	request, err := http.NewRequest("DELETE", vc.url+PROGRAMLIBRARY+fmt.Sprintf("/%d", id), nil)
	if err != nil {
		return ProgramDeleteResult{}, err
	}

	response, err := vc.client.Do(request)
	if err != nil {
		return ProgramDeleteResult{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return ProgramDeleteResult{}, NewServerError(response.StatusCode, errors.New("FAILED TO DELETE PROGRAM"))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ProgramDeleteResult{}, NewServerError(response.StatusCode, err)
	}

	actions := ActionResponse[any]{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return ProgramDeleteResult{}, NewServerError(response.StatusCode, err)
	}

	return NewProgramDeleteResult(&actions), nil
}

func programIsValid(file string) bool {

	return strings.HasSuffix(file, ".cpz") ||
		strings.HasSuffix(file, ".zip") ||
		strings.HasSuffix(file, ".lpz")
}

func programFileIsFullPath(file string) bool {
	return strings.ContainsAny(file, "/") || strings.ContainsAny(file, "\\")
}

func validateAndCreateFileHeader(writer *multipart.Writer, file string, key string, extensions []string) (bool, error) {
	if !programFileIsFullPath(file) {
		return false, nil
	}

	if !validateProgramExtensions(file, extensions) {
		return false, fmt.Errorf("FILE %s HAS INVALID EXTENSION", file)
	}

	err := addFormFile(file, key, writer)
	if err != nil {
		return false, err
	}
	return true, nil
}

func validateProgramExtensions(file string, extensions []string) bool {
	for _, e := range extensions {
		if strings.HasSuffix(file, e) {
			return true
		}
	}
	return false
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

type ProgramUploadResult struct {
	ProgramID    int16
	FriendlyName string
	Result       string
	Code         int16
	Success      bool
}

func NewProgramUploadResult(actions *ActionResponse[ProgramEntry]) ProgramUploadResult {

	p := actions.Actions[0].Results[0].Object
	s := actions.Actions[0].Results[0].StatusInfo
	c := actions.Actions[0].Results[0].StatusID

	return ProgramUploadResult{
		ProgramID:    p.ProgramID,
		FriendlyName: p.FriendlyName,
		Result:       s,
		Code:         c,
		Success:      c == 0,
	}
}

type ProgramDeleteResult struct {
	Result  string
	Code    int16
	Success bool
}

func NewProgramDeleteResult(actions *ActionResponse[any]) ProgramDeleteResult {

	s := actions.Actions[0].Results[0].StatusInfo
	c := actions.Actions[0].Results[0].StatusID

	return ProgramDeleteResult{

		Result:  s,
		Code:    c,
		Success: c == 0,
	}
}
