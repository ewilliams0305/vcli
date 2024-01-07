package vc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	TOKENREQUEST = "Token"
)

type VcApiToken interface {
	GetTokens() ([]ApiToken, VirtualControlError)
	CreateToken(readonly bool, description string) (ApiToken, VirtualControlError)
	EditToken(readonly bool, description string, token string) (ApiToken, VirtualControlError)
	DeleteToken(token string) (bool, VirtualControlError)
}

func (v *VC) GetTokens() ([]ApiToken, VirtualControlError) {
	var results ApitTokenResponse
	err := v.getBody(TOKENREQUEST, &results)

	if err != nil {
		return make([]ApiToken, 0), NewServerError(500, err)
	}

	tokens := results.Device.Programs.TokenList
	return tokens, nil
}

func (v *VC) CreateToken(readonly bool, description string) (ApiToken, VirtualControlError) {
	request := CreateApiTokenRequest{
		Status:      2,
		Description: description,
	}

	if readonly {
		request.Status = 1
	}

	jsonValue, err := json.Marshal(request)
	if err != nil {
		return ApiToken{}, NewServerError(500, err)
	}

	resp, err := v.client.Post(v.url+TOKENREQUEST, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return ApiToken{}, NewServerError(500, err)
	}

	if resp.StatusCode != 200 {
		return ApiToken{}, NewServerError(resp.StatusCode, fmt.Errorf("FAILED CREATING API KEY"))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ApiToken{}, NewServerError(resp.StatusCode, err)
	}

	actions := ActionResponse[ApiToken]{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return ApiToken{}, NewServerError(resp.StatusCode, err)
	}

	return actions.Actions[0].Results[0].Object, nil
}

func (v *VC) EditToken(readonly bool, description string, token string) (ApiToken, VirtualControlError) {
	request := EditApiTokenRequest{
		Status:      2,
		Description: description,
		Token:       token,
	}

	if readonly {
		request.Status = 1
	}

	jsonValue, err := json.Marshal(request)
	if err != nil {
		return ApiToken{}, NewServerError(500, err)
	}

	req, reqErr := http.NewRequest("PUT", v.url+TOKENREQUEST, bytes.NewBuffer(jsonValue))
	if reqErr != nil {
		return ApiToken{}, NewServerError(500, reqErr)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := v.client.Do(req)
	if err != nil {
		return ApiToken{}, NewServerError(500, err)
	}

	if resp.StatusCode != 200 {
		return ApiToken{}, NewServerError(resp.StatusCode, fmt.Errorf("FAILED CREATING API KEY"))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ApiToken{}, NewServerError(resp.StatusCode, err)
	}

	actions := ActionResponse[ApiToken]{}
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return ApiToken{}, NewServerError(resp.StatusCode, err)
	}

	return actions.Actions[0].Results[0].Object, nil
}

func (v *VC) DeleteToken(token string) (bool, VirtualControlError) {

	req, reqErr := http.NewRequest("DELETE", v.url+TOKENREQUEST+"/"+token, nil)
	if reqErr != nil {
		return false, NewServerError(500, reqErr)
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return false, NewServerError(resp.StatusCode, err)
	}

	return true, nil
}

type CreateApiTokenRequest struct {
	Description string      `json:"Description"`
	Status      TokenStatus `json:"Status"`
}

type EditApiTokenRequest struct {
	Description string      `json:"Description"`
	Status      TokenStatus `json:"Status"`
	Token       string      `json:"Token"`
}

type ApitTokenResponse struct {
	Device struct {
		Programs struct {
			TokenList []ApiToken `json:"TokenList"`
		} `json:"Programs"`
	} `json:"Device"`
}

type ApiToken struct {
	Token       string      `json:"Token"`
	Status      TokenStatus `json:"Status"`
	Description string      `json:"Description"`
	Level       string      `json:"Level"`
}

type TokenStatus int

const (
	ReadOnlyToken  TokenStatus = 1
	ReadWriteToken TokenStatus = 2
)
