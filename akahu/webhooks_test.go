package akahu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestWebhooksService_Unsubscribe(t *testing.T) {
	tests := []struct {
		name                string
		jsonResponse        string
		statusCode          int
		expected            bool
		expectedAPIResponse *APIResponse
	}{
		{
			name:                "with success response",
			jsonResponse:        "{ \"success\": true }",
			statusCode:          http.StatusOK,
			expected:            true,
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:                "with error response",
			jsonResponse:        errorResponseJson,
			statusCode:          http.StatusBadRequest,
			expected:            false,
			expectedAPIResponse: expectedErrorAPIResponse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := setupClient(t, test.jsonResponse, http.MethodDelete, test.statusCode, func(r *http.Request) {
				testTokenRequestHeaders(t, r, "app_token_123", "user_token_1")
			})

			actual, res, err := client.Webhooks.Unsubscribe(context.TODO(), "user_token_1", "id_1")
			testClientResponse(t, test.expected, actual, err)
			testClientAPIResponse(t, test.expectedAPIResponse, res, err)
		})
	}
}

func TestWebhooksService_Subscribe(t *testing.T) {
	expectedId := "hook_1111111111111111111111111"

	tests := []struct {
		name                string
		body                WebhookSubscribeRequest
		jsonResponse        string
		statusCode          int
		expected            *string
		expectedAPIResponse *APIResponse
	}{
		{
			name: "with success response",
			body: WebhookSubscribeRequest{
				WebhookType: Token,
				State:       "state123",
			},
			jsonResponse:        "{\"success\": true, \"item_id\": \"hook_1111111111111111111111111\" }",
			statusCode:          http.StatusOK,
			expected:            &expectedId,
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name: "with error response",
			body: WebhookSubscribeRequest{
				WebhookType: Token,
				State:       "state123",
			},
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

				var body WebhookSubscribeRequest
				_ = json.NewDecoder(r.Body).Decode(&body)

				if !reflect.DeepEqual(test.body, body) {
					t.Fatalf("expected request body %+v, actual %+v", test.body, body)
				}

			})

			actual, res, err := client.Webhooks.Subscribe(context.TODO(), "user_token_1", test.body)
			testClientResponse(t, test.expected, actual, err)
			testClientAPIResponse(t, test.expectedAPIResponse, res, err)
		})
	}
}

func TestAccountsService_List(t *testing.T) {
	webhookJson := "{ \"_id\": \"hook_1111111111111111111111111\", \"created_at\": \"2020-04-08T23:15:39.917Z\", \"updated_at\": \"2020-04-09T23:15:39.917Z\", \"last_called_at\": \"2020-04-10T23:15:39.917Z\", \"state\": \"foobarbaz\", \"url\": \"https://webhooks.myapp.com/akahu\" }"

	createdAt, _ := time.Parse(time.RFC3339, "2020-04-08T23:15:39.917Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2020-04-09T23:15:39.917Z")
	lastFailedAt, _ := time.Parse(time.RFC3339, "2020-04-10T23:15:39.917Z")

	tests := []struct {
		name                string
		jsonResponse        string
		statusCode          int
		expected            []WebhookResponse
		expectedAPIResponse *APIResponse
	}{
		{
			name:                "with empty response",
			jsonResponse:        fmt.Sprintf(collectionResponseJson, ""),
			statusCode:          http.StatusOK,
			expected:            []WebhookResponse{},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:         "with single response item",
			jsonResponse: fmt.Sprintf(collectionResponseJson, webhookJson),
			statusCode:   http.StatusOK,
			expected: []WebhookResponse{
				{
					Id:           "hook_1111111111111111111111111",
					CreatedAt:    createdAt,
					UpdatedAt:    updatedAt,
					LastCalledAt: lastFailedAt,
					State:        "foobarbaz",
					Url:          "https://webhooks.myapp.com/akahu",
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

			actual, res, err := client.Webhooks.List(context.TODO(), "user_token_1")
			testClientResponse(t, test.expected, actual, err)
			testClientAPIResponse(t, test.expectedAPIResponse, res, err)
		})
	}
}

func TestWebhooksService_GetPublicKey(t *testing.T) {
	publicKey := "-----BEGIN RSA PUBLIC KEY----- { PEM ENCODED PUBLIC KEY } -----END RSA PUBLIC KEY-----"

	tests := []struct {
		name                string
		jsonResponse        string
		statusCode          int
		expected            *string
		expectedAPIResponse *APIResponse
	}{
		{
			name:                "with valid public key response",
			jsonResponse:        fmt.Sprintf(itemResponseJson, "\""+publicKey+"\""),
			statusCode:          http.StatusOK,
			expected:            &publicKey,
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
				testBasicRequestHeaders(t, r)
			})

			actual, res, err := client.Webhooks.GetPublicKey(context.TODO())
			testClientResponse(t, test.expected, actual, err)
			testClientAPIResponse(t, test.expectedAPIResponse, res, err)
		})
	}
}

func TestWebhooksService_GetEvents(t *testing.T) {
	jsonResponse := "{ \"_id\": \"hook_1111111111111111111111111\", \"hook\": \"hook_1111111111111111111111112\", \"status\": \"FAILED\", \"created_at\": \"2020-04-08T23:15:39.917Z\", \"updated_at\": \"2020-04-09T23:15:39.917Z\", \"last_failed_at\": \"2020-04-10T23:15:39.917Z\", \"payload\": { \"success\": true, \"webhook_type\": \"TOKEN\", \"webhook_code\": \"test_1234\" } }"

	createdAt, _ := time.Parse(time.RFC3339, "2020-04-08T23:15:39.917Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2020-04-09T23:15:39.917Z")
	lastFailedAt, _ := time.Parse(time.RFC3339, "2020-04-10T23:15:39.917Z")

	tests := []struct {
		name                string
		status              string
		startTime           string
		endTime             string
		jsonResponse        string
		statusCode          int
		expected            []WebHookEventResponse
		expectedAPIResponse *APIResponse
	}{
		{
			name:                "with empty response",
			status:              "test_1234",
			startTime:           "2020-10-01T00:00:00Z",
			endTime:             "2020-10-05T00:00:00Z",
			jsonResponse:        fmt.Sprintf(collectionResponseJson, ""),
			statusCode:          http.StatusOK,
			expected:            []WebHookEventResponse{},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:         "with single response item",
			status:       "test_1234",
			startTime:    "2020-10-01T00:00:00Z",
			endTime:      "2020-10-05T00:00:00Z",
			jsonResponse: fmt.Sprintf(collectionResponseJson, jsonResponse),
			statusCode:   http.StatusOK,
			expected: []WebHookEventResponse{
				{
					Id:           "hook_1111111111111111111111111",
					Hook:         "hook_1111111111111111111111112",
					Status:       Failed,
					CreatedAt:    createdAt,
					UpdatedAt:    updatedAt,
					LastFailedAt: lastFailedAt,
					Payload: WebHookEventPayload{
						successResponse: successResponse{
							Success: true,
						},
						WebhookType: Token,
						WebhookCode: "test_1234",
					},
				},
			},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:                "with error response",
			status:              "test_1234",
			startTime:           "2020-10-01T00:00:00Z",
			endTime:             "2020-10-05T00:00:00Z",
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

				params := r.URL.Query()

				if start := params.Get("start"); start != test.startTime {
					t.Fatalf("Expected start param %s, actual %s", test.startTime, start)
				}

				if end := params.Get("end"); end != test.endTime {
					t.Fatalf("Expected end param %s, actual %s", test.endTime, end)
				}

				if status := params.Get("status"); status != test.status {
					t.Fatalf("Expected status param %s, actual %s", test.status, status)
				}

			})

			start, _ := time.Parse(time.RFC3339, test.startTime)
			end, _ := time.Parse(time.RFC3339, test.endTime)
			actual, res, err := client.Webhooks.ListEvents(context.TODO(), "user_token_1", test.status, start, end)
			testClientResponse(t, test.expected, actual, err)
			testClientAPIResponse(t, test.expectedAPIResponse, res, err)
		})
	}
}
