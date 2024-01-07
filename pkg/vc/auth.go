package vc

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
