package akahu

import (
	"context"
	"net/http"
	"net/url"
)

const (
	authPath        = "token"
	authURLBasePath = "https://oauth.akahu.io/"
)

type AuthService service

type grantType string

const (
	authCode grantType = "authorization_code"
)

type exchangeRequest struct {
	GrantType    grantType `json:"grant_type"`
	Code         string    `json:"code"`
	RedirectURI  string    `json:"redirect_uri"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
}

type ExchangeResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type AuthorizationURLOptions struct {
	ResponseType *string
	Email        *string
	Connection   *string
	Scope        *string
	State        *string
}

func (s *AuthService) Exchange(ctx context.Context, code string) (*ExchangeResponse, *APIResponse, error) {
	body := exchangeRequest{
		GrantType:    authCode,
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

	if !res.Success {
		return nil, res, nil
	}

	return &exchangeResponse, res, nil
}

func (s *AuthService) RevokeToken(ctx context.Context, userAccessToken string) (bool, *APIResponse, error) {
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

func (s *AuthService) BuildAuthorizationURL(options AuthorizationURLOptions) string {
	var responseType string
	if options.ResponseType == nil {
		responseType = "code"
	} else {
		responseType = *options.ResponseType
	}
	var scope string
	if options.Scope == nil {
		scope = "ENDURING_CONSENT"
	} else {
		scope = *options.Scope
	}

	params := url.Values{}
	params.Add("response_type", responseType)
	params.Add("scope", scope)

	params.Add("client_id", s.client.AppIDToken)
	params.Add("redirect_uri", s.client.RedirectURI.String())

	if options.Email != nil {
		params.Add("email", *options.Email)
	}
	if options.Connection != nil {
		params.Add("connection", *options.Connection)
	}
	if options.State != nil {
		params.Add("state", *options.State)
	}

	authURL, _ := url.Parse(authURLBasePath)
	authURL.RawQuery = params.Encode()

	return authURL.String()
}
