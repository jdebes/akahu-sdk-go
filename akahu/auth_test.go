package akahu

import (
	"context"
	"net/http"
	"testing"
)

const exchangeJson = "{ \"success\": true, \"access_token\": \"user_token_1111111111111111111111111\", \"token_type\": \"bearer\", \"scope\": \"IDENTITY_BASIC ACCOUNTS TRANSACTIONS\" }"

func TestAuthService_BuildAuthorizationURL(t *testing.T) {
	client := NewClient(nil, "app_token_123", "appsecret123", "https://example.com/auth/akahu")
	email := "test_user@gmail.com"
	connection := "conn_1234"
	responseType := "codetest"
	state := "1234567890"
	scope := "ENDURING_CONSENT_TEST"

	tests := []struct {
		name     string
		opts     AuthorizationURLOptions
		expected string
	}{
		{
			name:     "with all defaults configurations",
			opts:     AuthorizationURLOptions{},
			expected: "https://oauth.akahu.io/?client_id=app_token_123&redirect_uri=https%3A%2F%2Fexample.com%2Fauth%2Fakahu&response_type=code&scope=ENDURING_CONSENT",
		},
		{
			name: "with email and connection configured",
			opts: AuthorizationURLOptions{
				Email:      &email,
				Connection: &connection,
			},
			expected: "https://oauth.akahu.io/?client_id=app_token_123&connection=conn_1234&email=test_user%40gmail.com&redirect_uri=https%3A%2F%2Fexample.com%2Fauth%2Fakahu&response_type=code&scope=ENDURING_CONSENT",
		},
		{
			name: "with all options configured",
			opts: AuthorizationURLOptions{
				ResponseType: &responseType,
				Email:        &email,
				Connection:   &connection,
				Scope:        &scope,
				State:        &state,
			},
			expected: "https://oauth.akahu.io/?client_id=app_token_123&connection=conn_1234&email=test_user%40gmail.com&redirect_uri=https%3A%2F%2Fexample.com%2Fauth%2Fakahu&response_type=codetest&scope=ENDURING_CONSENT_TEST&state=1234567890",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := client.Auth.BuildAuthorizationURL(test.opts); actual != test.expected {
				t.Errorf("expected %v, actual %v", test.expected, actual)
			}
		})
	}
}

func TestAuthService_Exchange(t *testing.T) {
	tests := []struct {
		name                string
		jsonResponse        string
		statusCode          int
		expected            *ExchangeResponse
		expectedAPIResponse *APIResponse
	}{
		{
			name:         "with success response",
			jsonResponse: exchangeJson,
			statusCode:   http.StatusOK,
			expected: &ExchangeResponse{
				AccessToken: "user_token_1111111111111111111111111",
				TokenType:   "bearer",
				Scope:       "IDENTITY_BASIC ACCOUNTS TRANSACTIONS",
			},
			expectedAPIResponse: expectedSuccessAPIResponse,
		},
		{
			name:                "with error response",
			jsonResponse:        errorResponseJsonWithError,
			statusCode:          http.StatusBadRequest,
			expected:            nil,
			expectedAPIResponse: expectedErrorAPIResponse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := setupClient(t, test.jsonResponse, http.MethodPost, test.statusCode)

			actual, res, err := client.Auth.Exchange(context.TODO(), "code")
			testClientResponse(t, test.expected, actual, err)
			testClientAPIResponse(t, test.expectedAPIResponse, res, err)
		})
	}
}
