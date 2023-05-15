package akahu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

const (
	unenrichedTransactionJson = "{ \"_id\": \"trans_1111111111111111111111111\", \"_account\": \"acc_1111111111111111111111111\", \"_connection\": \"conn_1111111111111111111111111\", \"created_at\": \"2020-01-01T01:00:00.000Z\", \"updated_at\": \"2020-01-01T02:00:00.000Z\", \"date\": \"2020-01-01T00:00:00.000Z\", \"description\": \"{RAW TRANSACTION DESCRIPTION}\", \"amount\": -5.5, \"balance\": 100, \"type\": \"EFTPOS\" }"
	enrichedTransactionJson   = "{ \"_id\": \"trans_1111111111111111111111111\", \"_account\": \"acc_1111111111111111111111111\", \"_connection\": \"conn_1111111111111111111111111\", \"created_at\": \"2020-01-01T01:00:00.000Z\", \"updated_at\": \"2020-01-01T02:00:00.000Z\", \"date\": \"2020-01-01T00:00:00.000Z\", \"description\": \"{RAW TRANSACTION DESCRIPTION}\", \"amount\": -5.5, \"balance\": 100, \"type\": \"EFTPOS\", \"merchant\": { \"_id\": \"merchant_1111111111111111111111111\", \"name\": \"Bob's Pizza\" }, \"category\": { \"_id\": \"nzfcc_1111111111111111111111111\", \"name\": \"Cafes and restaurants\", \"groups\": { \"personal_finance\": { \"_id\": \"group_clasr0ysw0011hk4m6hlk9fq0\", \"name\": \"Lifestyle\" } } } }"
)

var (
	createdAt, _ = time.Parse(time.RFC3339, "2020-01-01T01:00:00.000Z")
	updatedAt, _ = time.Parse(time.RFC3339, "2020-01-01T02:00:00.000Z")
	date, _      = time.Parse(time.RFC3339, "2020-01-01T00:00:00.000Z")
)

func TestTransactionsService_Get(t *testing.T) {
	tests := []struct {
		name                string
		jsonResponse        string
		statusCode          int
		expected            *TransactionResponse
		expectedAPIResponse *APIResponse
	}{
		{
			name:         "with unenriched response",
			jsonResponse: fmt.Sprintf(itemResponseJson, unenrichedTransactionJson),
			statusCode:   http.StatusOK,
			expected: &TransactionResponse{
				Id:          "trans_1111111111111111111111111",
				Account:     "acc_1111111111111111111111111",
				Connection:  "conn_1111111111111111111111111",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				Date:        date,
				Description: "{RAW TRANSACTION DESCRIPTION}",
				Amount:      decimal.NewFromFloat(-5.5),
				Balance:     decimal.NewFromInt(100),
				Type:        "EFTPOS",
			},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:         "with enriched response",
			jsonResponse: fmt.Sprintf(itemResponseJson, enrichedTransactionJson),
			statusCode:   http.StatusOK,
			expected: &TransactionResponse{
				Id:          "trans_1111111111111111111111111",
				Account:     "acc_1111111111111111111111111",
				Connection:  "conn_1111111111111111111111111",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				Date:        date,
				Description: "{RAW TRANSACTION DESCRIPTION}",
				Amount:      decimal.NewFromFloat(-5.5),
				Balance:     decimal.NewFromInt(100),
				Type:        "EFTPOS",
				Merchant: &Merchant{
					Id:   "merchant_1111111111111111111111111",
					Name: "Bob's Pizza",
				},
				Category: &Category{
					Id:   "nzfcc_1111111111111111111111111",
					Name: "Cafes and restaurants",
					Groups: &Groups{
						PersonalFinance: &PersonalFinance{
							Id:   "group_clasr0ysw0011hk4m6hlk9fq0",
							Name: "Lifestyle",
						},
					},
				},
			},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:                "with error response",
			jsonResponse:        errorResponseJson,
			statusCode:          http.StatusBadRequest,
			expected:            nil,
			expectedAPIResponse: expectedErrorAPIResponse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := setupClient(t, test.jsonResponse, http.MethodGet, test.statusCode, func(r *http.Request) {
				testTokenRequestHeaders(t, r, "app_token_123", "user_token_1")
			})

			actual, res, err := client.Transactions.Get(context.TODO(), "user_token_1", "id_1")
			testClientResponse(t, test.expected, actual, err)
			testClientAPIResponse(t, test.expectedAPIResponse, res, err)
		})
	}
}

func TestTransactionsService_GetByIds(t *testing.T) {
	tests := []struct {
		name                string
		jsonResponse        string
		ids                 []string
		statusCode          int
		expected            []TransactionResponse
		expectedAPIResponse *APIResponse
	}{
		{
			name:                "with empty response",
			jsonResponse:        fmt.Sprintf(collectionResponseJson, ""),
			ids:                 []string{"id_1"},
			statusCode:          http.StatusOK,
			expected:            []TransactionResponse{},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:                "with multiple ids",
			jsonResponse:        fmt.Sprintf(collectionResponseJson, ""),
			ids:                 []string{"id_1", "id_2", "id_3"},
			statusCode:          http.StatusOK,
			expected:            []TransactionResponse{},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:         "with single unenriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, unenrichedTransactionJson),
			ids:          []string{"id_1"},
			statusCode:   http.StatusOK,
			expected: []TransactionResponse{
				{
					Id:          "trans_1111111111111111111111111",
					Account:     "acc_1111111111111111111111111",
					Connection:  "conn_1111111111111111111111111",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
					Date:        date,
					Description: "{RAW TRANSACTION DESCRIPTION}",
					Amount:      decimal.NewFromFloat(-5.5),
					Balance:     decimal.NewFromInt(100),
					Type:        "EFTPOS",
				},
			},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:         "with single enriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, enrichedTransactionJson),
			ids:          []string{"id_1"},
			statusCode:   http.StatusOK,
			expected: []TransactionResponse{
				{
					Id:          "trans_1111111111111111111111111",
					Account:     "acc_1111111111111111111111111",
					Connection:  "conn_1111111111111111111111111",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
					Date:        date,
					Description: "{RAW TRANSACTION DESCRIPTION}",
					Amount:      decimal.NewFromFloat(-5.5),
					Balance:     decimal.NewFromInt(100),
					Type:        "EFTPOS",
					Merchant: &Merchant{
						Id:   "merchant_1111111111111111111111111",
						Name: "Bob's Pizza",
					},
					Category: &Category{
						Id:   "nzfcc_1111111111111111111111111",
						Name: "Cafes and restaurants",
						Groups: &Groups{
							PersonalFinance: &PersonalFinance{
								Id:   "group_clasr0ysw0011hk4m6hlk9fq0",
								Name: "Lifestyle",
							},
						},
					},
				},
			},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:                "with error response",
			jsonResponse:        errorResponseJson,
			statusCode:          http.StatusBadRequest,
			expected:            nil,
			expectedAPIResponse: expectedErrorAPIResponse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := setupClient(t, test.jsonResponse, http.MethodPost, test.statusCode, func(r *http.Request) {
				testTokenRequestHeaders(t, r, "app_token_123", "user_token_1")

				var ids []string
				_ = json.NewDecoder(r.Body).Decode(&ids)

				if !reflect.DeepEqual(ids, test.ids) {
					t.Fatalf("Expected request body %+v, actual %+v", test.ids, ids)
				}
			})

			actual, res, err := client.Transactions.GetByIds(context.TODO(), "user_token_1", test.ids...)
			testClientResponse(t, test.expected, actual, err)
			testClientAPIResponse(t, test.expectedAPIResponse, res, err)
		})
	}
}

func TestTransactionsService_List(t *testing.T) {
	tests := []struct {
		name                string
		jsonResponse        string
		statusCode          int
		expected            []TransactionResponse
		expectedAPIResponse *APIResponse
	}{
		{
			name:                "with empty response",
			jsonResponse:        fmt.Sprintf(collectionResponseJson, ""),
			statusCode:          http.StatusOK,
			expected:            []TransactionResponse{},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:         "with single unenriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, unenrichedTransactionJson),
			statusCode:   http.StatusOK,
			expected: []TransactionResponse{
				{
					Id:          "trans_1111111111111111111111111",
					Account:     "acc_1111111111111111111111111",
					Connection:  "conn_1111111111111111111111111",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
					Date:        date,
					Description: "{RAW TRANSACTION DESCRIPTION}",
					Amount:      decimal.NewFromFloat(-5.5),
					Balance:     decimal.NewFromInt(100),
					Type:        "EFTPOS",
				},
			},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:         "with single enriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, enrichedTransactionJson),
			statusCode:   http.StatusOK,
			expected: []TransactionResponse{
				{
					Id:          "trans_1111111111111111111111111",
					Account:     "acc_1111111111111111111111111",
					Connection:  "conn_1111111111111111111111111",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
					Date:        date,
					Description: "{RAW TRANSACTION DESCRIPTION}",
					Amount:      decimal.NewFromFloat(-5.5),
					Balance:     decimal.NewFromInt(100),
					Type:        "EFTPOS",
					Merchant: &Merchant{
						Id:   "merchant_1111111111111111111111111",
						Name: "Bob's Pizza",
					},
					Category: &Category{
						Id:   "nzfcc_1111111111111111111111111",
						Name: "Cafes and restaurants",
						Groups: &Groups{
							PersonalFinance: &PersonalFinance{
								Id:   "group_clasr0ysw0011hk4m6hlk9fq0",
								Name: "Lifestyle",
							},
						},
					},
				},
			},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:                "with error response",
			jsonResponse:        errorResponseJson,
			statusCode:          http.StatusBadRequest,
			expected:            nil,
			expectedAPIResponse: expectedErrorAPIResponse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := setupClient(t, test.jsonResponse, http.MethodGet, test.statusCode, func(r *http.Request) {
				testTokenRequestHeaders(t, r, "app_token_123", "user_token_1")
			})

			actual, res, err := client.Transactions.List(context.TODO(), "user_token_1", time.Now(), time.Now())
			testClientResponse(t, test.expected, actual, err)
			testClientAPIResponse(t, test.expectedAPIResponse, res, err)
		})
	}
}

func TestTransactionsService_ListPending(t *testing.T) {
	tests := []struct {
		name                string
		jsonResponse        string
		statusCode          int
		expected            []TransactionResponse
		expectedAPIResponse *APIResponse
	}{
		{
			name:                "with empty response",
			jsonResponse:        fmt.Sprintf(collectionResponseJson, ""),
			statusCode:          http.StatusOK,
			expected:            []TransactionResponse{},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:         "with single unenriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, unenrichedTransactionJson),
			statusCode:   http.StatusOK,
			expected: []TransactionResponse{
				{
					Id:          "trans_1111111111111111111111111",
					Account:     "acc_1111111111111111111111111",
					Connection:  "conn_1111111111111111111111111",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
					Date:        date,
					Description: "{RAW TRANSACTION DESCRIPTION}",
					Amount:      decimal.NewFromFloat(-5.5),
					Balance:     decimal.NewFromInt(100),
					Type:        "EFTPOS",
				},
			},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:         "with single enriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, enrichedTransactionJson),
			statusCode:   http.StatusOK,
			expected: []TransactionResponse{
				{
					Id:          "trans_1111111111111111111111111",
					Account:     "acc_1111111111111111111111111",
					Connection:  "conn_1111111111111111111111111",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
					Date:        date,
					Description: "{RAW TRANSACTION DESCRIPTION}",
					Amount:      decimal.NewFromFloat(-5.5),
					Balance:     decimal.NewFromInt(100),
					Type:        "EFTPOS",
					Merchant: &Merchant{
						Id:   "merchant_1111111111111111111111111",
						Name: "Bob's Pizza",
					},
					Category: &Category{
						Id:   "nzfcc_1111111111111111111111111",
						Name: "Cafes and restaurants",
						Groups: &Groups{
							PersonalFinance: &PersonalFinance{
								Id:   "group_clasr0ysw0011hk4m6hlk9fq0",
								Name: "Lifestyle",
							},
						},
					},
				},
			},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:                "with error response",
			jsonResponse:        errorResponseJson,
			statusCode:          http.StatusBadRequest,
			expected:            nil,
			expectedAPIResponse: expectedErrorAPIResponse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := setupClient(t, test.jsonResponse, http.MethodGet, test.statusCode, func(r *http.Request) {
				testTokenRequestHeaders(t, r, "app_token_123", "user_token_1")
			})

			actual, res, err := client.Transactions.ListPending(context.TODO(), "user_token_1", time.Now(), time.Now())
			testClientResponse(t, test.expected, actual, err)
			testClientAPIResponse(t, test.expectedAPIResponse, res, err)
		})
	}
}
