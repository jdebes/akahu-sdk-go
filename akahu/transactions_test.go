package akahu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

const (
	unenrichedTransactionJson = "{ \"_id\": \"trans_1111111111111111111111111\", \"_account\": \"acc_1111111111111111111111111\", \"_connection\": \"conn_1111111111111111111111111\", \"created_at\": \"2020-01-01T01:00:00.000Z\", \"updated_at\": \"2020-01-01T02:00:00.000Z\", \"date\": \"2020-01-01T00:00:00.000Z\", \"description\": \"{RAW TRANSACTION DESCRIPTION}\", \"amount\": -5.5, \"balance\": 100, \"type\": \"EFTPOS\" }"
	enrichedTransactionJson   = "{ \"_id\": \"trans_1111111111111111111111111\", \"_account\": \"acc_1111111111111111111111111\", \"_connection\": \"conn_1111111111111111111111111\", \"created_at\": \"2020-01-01T01:00:00.000Z\", \"updated_at\": \"2020-01-01T02:00:00.000Z\", \"date\": \"2020-01-01T00:00:00.000Z\", \"description\": \"{RAW TRANSACTION DESCRIPTION}\", \"amount\": -5.5, \"balance\": 100, \"type\": \"EFTPOS\", \"merchant\": { \"_id\": \"merchant_1111111111111111111111111\", \"name\": \"Bob's Pizza\" }, \"category\": { \"_id\": \"nzfcc_1111111111111111111111111\", \"name\": \"Cafes and restaurants\", \"groups\": { \"personal_finance\": { \"_id\": \"group_clasr0ysw0011hk4m6hlk9fq0\", \"name\": \"Lifestyle\" } } } }"
	itemResponseJson          = "{ \"success\": true, \"item\": %s }"
	collectionResponseJson    = "{ \"success\": true, \"items\": [%s] }"
)

var (
	createdAt, _ = time.Parse(time.RFC3339, "2020-01-01T01:00:00.000Z")
	updatedAt, _ = time.Parse(time.RFC3339, "2020-01-01T02:00:00.000Z")
	date, _      = time.Parse(time.RFC3339, "2020-01-01T00:00:00.000Z")
)

type clientTest func(r *http.Request)

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func testTokenRequestHeaders(t *testing.T, r *http.Request, appToken, userAccessToken string) {
	akahuIdHeader := r.Header.Get("X-Akahu-ID")
	authorizationHeader := r.Header.Get("Authorization")

	if akahuIdHeader != appToken {
		t.Fatalf("expected header X-Akahu-ID, actual %s", akahuIdHeader)
	}

	if authorizationHeader != fmt.Sprintf("Bearer %s", userAccessToken) {
		t.Fatalf("expected header Authorization, actual %s", authorizationHeader)
	}
}

func testClientResponse(t *testing.T, expected, actual interface{}, err error) {
	if err != nil {
		t.Fatalf("client request returned err %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %+v, actual %+v", expected, actual)
	}
}

func setupClient(t *testing.T, mockedResponse, expectedHttpMethod string, clientTests ...clientTest) *Client {
	mockHttpClient := http.Client{Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != expectedHttpMethod {
			t.Fatalf("expected method %s, actual %s", expectedHttpMethod, req.Method)
		}

		testTokenRequestHeaders(t, req, "app_token_123", "user_token_1")

		for _, ct := range clientTests {
			ct(req)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(mockedResponse)),
		}, nil
	})}

	return NewClient(&mockHttpClient, "app_token_123", "", "")
}

func TestTransactionsService_Get(t *testing.T) {
	tests := []struct {
		name         string
		jsonResponse string
		expected     *TransactionResponse
	}{
		{
			name:         "with unenriched response",
			jsonResponse: fmt.Sprintf(itemResponseJson, unenrichedTransactionJson),
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
		},
		{
			name:         "with enriched response",
			jsonResponse: fmt.Sprintf(itemResponseJson, enrichedTransactionJson),
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
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := setupClient(t, test.jsonResponse, http.MethodGet)

			actual, _, err := client.Transactions.Get(context.TODO(), "user_token_1", "id_1")
			testClientResponse(t, test.expected, actual, err)
		})
	}
}

func TestTransactionsService_GetByIds(t *testing.T) {
	tests := []struct {
		name         string
		jsonResponse string
		ids          []string
		expected     []TransactionResponse
	}{
		{
			name:         "with empty response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, ""),
			ids:          []string{"id_1"},
			expected:     []TransactionResponse{},
		},
		{
			name:         "with multiple ids",
			jsonResponse: fmt.Sprintf(collectionResponseJson, ""),
			ids:          []string{"id_1", "id_2", "id_3"},
			expected:     []TransactionResponse{},
		},
		{
			name:         "with single unenriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, unenrichedTransactionJson),
			ids:          []string{"id_1"},
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
		},
		{
			name:         "with single enriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, enrichedTransactionJson),
			ids:          []string{"id_1"},
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
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := setupClient(t, test.jsonResponse, http.MethodPost, func(r *http.Request) {
				var ids []string
				_ = json.NewDecoder(r.Body).Decode(&ids)

				if !reflect.DeepEqual(ids, test.ids) {
					t.Fatalf("Expected request body %+v, actual %+v", test.ids, ids)
				}
			})

			actual, _, err := client.Transactions.GetByIds(context.TODO(), "user_token_1", test.ids...)
			testClientResponse(t, test.expected, actual, err)
		})
	}
}

func TestTransactionsService_List(t *testing.T) {
	tests := []struct {
		name         string
		jsonResponse string
		expected     []TransactionResponse
	}{
		{
			name:         "with empty response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, ""),
			expected:     []TransactionResponse{},
		},
		{
			name:         "with single unenriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, unenrichedTransactionJson),
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
		},
		{
			name:         "with single enriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, enrichedTransactionJson),
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
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := setupClient(t, test.jsonResponse, http.MethodGet)

			actual, _, err := client.Transactions.List(context.TODO(), "user_token_1", time.Now(), time.Now())
			testClientResponse(t, test.expected, actual, err)
		})
	}
}

func TestTransactionsService_ListPending(t *testing.T) {
	tests := []struct {
		name         string
		jsonResponse string
		expected     []TransactionResponse
	}{
		{
			name:         "with empty response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, ""),
			expected:     []TransactionResponse{},
		},
		{
			name:         "with single unenriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, unenrichedTransactionJson),
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
		},
		{
			name:         "with single enriched response",
			jsonResponse: fmt.Sprintf(collectionResponseJson, enrichedTransactionJson),
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
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := setupClient(t, test.jsonResponse, http.MethodGet)

			actual, _, err := client.Transactions.ListPending(context.TODO(), "user_token_1", time.Now(), time.Now())
			testClientResponse(t, test.expected, actual, err)
		})
	}
}
