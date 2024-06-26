package akahu

import (
	"context"
	"net/http"
	"path"
	"time"

	"github.com/shopspring/decimal"
)

const accountsPath = "accounts"

type AccountsService service

type AccountResponse struct {
	ID          string `json:"_id"`
	Credentials string `json:"_credentials"`
	Connection  struct {
		Name string `json:"name"`
		Logo string `json:"logo"`
		Id   string `json:"_id"`
	} `json:"connection"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Meta   struct {
		Holder string `json:"holder"`
	} `json:"meta"`
	Refreshed struct {
		Balance      time.Time `json:"balance"`
		Meta         time.Time `json:"meta"`
		Transactions time.Time `json:"transactions"`
	} `json:"refreshed"`
	FormattedAccount string `json:"formatted_account"`
	Balance          struct {
		Available decimal.Decimal `json:"available"`
		Currency  string          `json:"currency"`
		Current   decimal.Decimal `json:"current"`
		Limit     decimal.Decimal `json:"limit"`
		Overdrawn bool            `json:"overdrawn"`
	} `json:"balance"`
	Attributes []string `json:"attributes"`
	Branch     struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Phone       string `json:"phone"`
	} `json:"branch"`
	Type string `json:"type"`
}

// List gets a list of all accounts that the user has connected to your application.
//
// Akahu docs: https://developers.akahu.nz/reference/get_accounts
func (s *AccountsService) List(ctx context.Context, userAccessToken string) ([]AccountResponse, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodGet, accountsPath, nil, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return nil, nil, err
	}

	var accounts collectionResponse[AccountResponse]
	res, err := s.client.do(ctx, r, &accounts)
	if err != nil {
		return nil, nil, err
	}

	return accounts.Items, res, nil
}

// Get an individual account that the user has connected to your application.
//
// Akahu docs: https://developers.akahu.nz/reference/get_accounts-id
func (s *AccountsService) Get(ctx context.Context, userAccessToken, ID string) (*AccountResponse, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodGet, path.Join(accountsPath, ID), nil, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return nil, nil, err
	}

	var accounts itemResponse[AccountResponse]
	res, err := s.client.do(ctx, r, &accounts)
	if err != nil {
		return nil, nil, err
	}

	return accounts.Item, res, nil
}

// Revoke your application's access to one of the user's connected accounts and its associated data, including transactions.
//
// Akahu docs: https://developers.akahu.nz/reference/delete_accounts-id
func (s *AccountsService) Revoke(ctx context.Context, userAccessToken, ID string) (bool, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodDelete, path.Join(accountsPath, ID), nil, withTokenRequestConfig(userAccessToken))
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
