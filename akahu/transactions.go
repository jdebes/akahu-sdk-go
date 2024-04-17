package akahu

import (
	"context"
	"net/http"
	"path"
	"time"

	"github.com/shopspring/decimal"
)

const transactionsPath = "transactions"
const pendingPath = "pending"

type TransactionsService service

type Merchant struct {
	Id      string  `json:"_id"`
	Name    string  `json:"name"`
	Website *string `json:"website"`
}

type PersonalFinance struct {
	Id   string `json:"_id"`
	Name string `json:"name"`
}

type Groups struct {
	PersonalFinance *PersonalFinance `json:"personal_finance"`
}

type Category struct {
	Id     string  `json:"_id"`
	Name   string  `json:"name"`
	Groups *Groups `json:"groups"`
}

type Conversion struct {
	Amount   *decimal.Decimal `json:"amount"`
	Currency *string          `json:"currency"`
	Rate     *decimal.Decimal `json:"rate"`
}

type Meta struct {
	Particulars  *string     `json:"particulars"`
	Code         *string     `json:"code"`
	Reference    *string     `json:"reference"`
	OtherAccount *string     `json:"other_account"`
	Conversion   *Conversion `json:"conversion"`
}

type TransactionResponse struct {
	Id          string          `json:"_id"`
	Account     string          `json:"_account"`
	Connection  string          `json:"_connection"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Date        time.Time       `json:"date"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	Balance     decimal.Decimal `json:"balance"`
	Type        string          `json:"type"`
	Merchant    *Merchant       `json:"merchant"`
	Category    *Category       `json:"category"`
	Meta        *Meta           `json:"meta"`
}

// List gets a list of settled transactions within the 'start' and 'end' time range.
//
// Akahu docs: https://developers.akahu.nz/reference/get_transactions
func (s *TransactionsService) List(ctx context.Context, userAccessToken string, startTime, endTime time.Time) ([]TransactionResponse, *APIResponse, error) {
	return s.list(ctx, transactionsPath, userAccessToken, startTime, endTime)
}

// ListPending gets a list of pending transactions within the 'start' and 'end' time range.
//
// Akahu docs: https://developers.akahu.nz/reference/get_transactions-pending
func (s *TransactionsService) ListPending(ctx context.Context, userAccessToken string, startTime, endTime time.Time) ([]TransactionResponse, *APIResponse, error) {
	return s.list(ctx, path.Join(transactionsPath, pendingPath), userAccessToken, startTime, endTime)
}

// Get fetches an individual transaction from one of the user's connected accounts.
// All returned dates are in UTC.
//
// Akahu docs: https://developers.akahu.nz/reference/get_transactions-id
func (s *TransactionsService) Get(ctx context.Context, userAccessToken, id string) (*TransactionResponse, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodGet, path.Join(transactionsPath, id), nil, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return nil, nil, err
	}

	var accounts itemResponse[TransactionResponse]
	res, err := s.client.do(ctx, r, &accounts)
	if err != nil {
		return nil, nil, err
	}

	return accounts.Item, res, nil
}

// GetByIds fetches a list of transactions from one of the user's connected accounts by Transaction ids.
// All returned dates are in UTC.
//
// Akahu docs: https://developers.akahu.nz/reference/post_transactions-ids
func (s *TransactionsService) GetByIds(ctx context.Context, userAccessToken string, ids ...string) ([]TransactionResponse, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodPost, path.Join(transactionsPath, "ids"), ids, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return nil, nil, err
	}

	var accounts collectionResponse[TransactionResponse]
	res, err := s.client.do(ctx, r, &accounts)
	if err != nil {
		return nil, nil, err
	}

	return accounts.Items, res, nil
}

func (s *TransactionsService) list(ctx context.Context, urlPath, userAccessToken string, startTime, endTime time.Time) ([]TransactionResponse, *APIResponse, error) {
	params := paramsWithDateRange(startTime, endTime)
	encodedPath := pathWithParams(urlPath, params)

	r, err := s.client.newRequest(http.MethodGet, encodedPath, nil, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return nil, nil, err
	}

	var accounts collectionResponse[TransactionResponse]
	res, err := s.client.do(ctx, r, &accounts)
	if err != nil {
		return nil, nil, err
	}

	return accounts.Items, res, nil
}
