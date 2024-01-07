package vc

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	TOKENREQUEST = "Token"
)

type VcApiToken interface {
	GetTokens() ([]ApiToken, VirtualControlError)
	CreateToken(readonly bool, description string) (ApiToken, VirtualControlError)
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
	request := ApiTokenRequest{
		Status: 2,
	}

	if readonly {
		request.Status = 1
	}

	jsonValue, err := json.Marshal(request)
	if err != nil {
		return ApiToken{}, NewServerError(500, err)
	}

	resp, err := v.client.Post(TOKENREQUEST, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return ApiToken{}, NewServerError(resp.StatusCode, err)
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

type ApiTokenRequest struct {
	Description string      `json:"Description"`
	Status      TokenStatus `json:"Status"`
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
