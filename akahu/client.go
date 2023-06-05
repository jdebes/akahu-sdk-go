package akahu

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	jsonContentType = "application/json"
	defaultBaseURL  = "https://api.akahu.io/v1/"
	akahuIDHeader   = "X-Akahu-ID"
)

type Client struct {
	client *http.Client

	BaseURL     *url.URL
	RedirectURI *url.URL
	AppSecret   string
	AppIDToken  string

	Accounts     *AccountsService
	Auth         *AuthService
	Me           *MeService
	Connections  *ConnectionsService
	Transactions *TransactionsService
	Webhooks     *WebhooksService
}

type successResponse struct {
	Success bool `json:"success"`
}

type itemResponse[T any] struct {
	successResponse
	Item *T `json:"item"`
}

type collectionResponse[T any] struct {
	successResponse
	Items []T `json:"items"`
}

type errorResponse struct {
	successResponse
	Message *string `json:"message"`
	Error   *string `json:"error"`
}

type APIResponse struct {
	Success bool
	Message string
	*http.Response
}

type service struct {
	client *Client
}

type requestConfig func(req *http.Request, client *Client)

func withTokenRequestConfig(userAccessToken string) requestConfig {
	return func(req *http.Request, c *Client) {
		req.Header.Set(akahuIDHeader, c.AppIDToken)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userAccessToken))
	}
}

func withBasicAuthRequestConfig() requestConfig {
	return func(req *http.Request, c *Client) {
		credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.AppIDToken, c.AppSecret)))

		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", credentials))
	}
}

func NewClient(httpClient *http.Client, appIDToken, appSecret, redirectUri string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	// TODO handle error
	parsedRedirectUri, _ := url.Parse(redirectUri)

	c := &Client{
		client:      httpClient,
		BaseURL:     baseURL,
		RedirectURI: parsedRedirectUri,
		AppIDToken:  appIDToken,
		AppSecret:   appSecret,
	}
	c.Accounts = &AccountsService{client: c}
	c.Auth = &AuthService{client: c}
	c.Me = &MeService{client: c}
	c.Connections = &ConnectionsService{client: c}
	c.Transactions = &TransactionsService{client: c}
	c.Webhooks = &WebhooksService{client: c}

	return c
}

func (c *Client) newRequest(method, urlPath string, body interface{}, requestConfigs ...requestConfig) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlPath)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", jsonContentType)
	}
	req.Header.Set("Accept", jsonContentType)

	for _, rc := range requestConfigs {
		rc(req, c)
	}

	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*APIResponse, error) {
	res, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		err = decoder.Decode(&v)
		if err != nil {
			return nil, err
		}

		return &APIResponse{
			Success:  true,
			Response: res,
		}, nil
	}

	var errResp errorResponse
	err = decoder.Decode(&errResp)
	if err != nil {
		return nil, err
	}

	var message string
	if errResp.Message != nil {
		message = *errResp.Message
	} else if errResp.Error != nil {
		message = *errResp.Error
	}

	return &APIResponse{
		Success:  errResp.Success,
		Message:  message,
		Response: res,
	}, nil
}
