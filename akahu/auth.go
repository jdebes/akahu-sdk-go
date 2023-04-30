package akahu

import (
	"context"
	"net/http"
)

const authPath = "token"

type AuthService service

type exchangeRequest struct {
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ExchangeResponse struct {
	successResponse
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (s *AuthService) Exchange(ctx context.Context, code string) (*ExchangeResponse, *http.Response, error) {
	body := exchangeRequest{
		Code:         code,
		RedirectURI:  s.client.RedirectURI.String(),
		ClientID:     s.client.AppIDToken,
		ClientSecret: s.client.AppSecret,
	}

	r, err := s.client.newRequest(http.MethodPost, authPath, body)
	if err != nil {
		return nil, nil, err
	}

	var exchangeResponse ExchangeResponse
	res, err := s.client.do(ctx, r, &exchangeResponse)
	if err != nil {
		return nil, nil, err
	}

	return &exchangeResponse, res, nil
}

func (s *AuthService) RevokeToken(ctx context.Context, userAccessToken string) (bool, *http.Response, error) {
	r, err := s.client.newRequest(http.MethodDelete, authPath, nil, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return false, nil, err
	}

	var successResponse successResponse
	res, err := s.client.do(ctx, r, &successResponse)
	if err != nil {
		return false, nil, err
	}

	return successResponse.Success, res, nil
}
