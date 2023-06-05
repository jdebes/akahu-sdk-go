package akahu

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

const (
	itemResponseJson             = "{ \"success\": true, \"item\": %s }"
	collectionResponseJson       = "{ \"success\": true, \"items\": [%s] }"
	errorResponseJsonWithMessage = "{ \"success\": false, \"message\": \"Error\" }"
	errorResponseJsonWithError   = "{ \"success\": false, \"error\": \"Error\" }"
)

var (
	expectedErrorAPIResponse = &APIResponse{
		Success: false,
		Message: "Error",
	}
	expectedSuccessAPIResponse = &APIResponse{
		Success: true,
	}
)

type clientTest func(r *http.Request)

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func setupClient(t *testing.T, mockedResponse, expectedHttpMethod string, respStatus int, clientTests ...clientTest) *Client {
	mockHttpClient := http.Client{Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if method := req.Method; method != expectedHttpMethod {
			t.Fatalf("expected method %s, actual %s", expectedHttpMethod, method)
		}

		for _, ct := range clientTests {
			ct(req)
		}

		return &http.Response{
			StatusCode: respStatus,
			Body:       io.NopCloser(strings.NewReader(mockedResponse)),
		}, nil
	})}

	return NewClient(&mockHttpClient, "app_token_123", "appSecret123", "")
}

func testTokenRequestHeaders(t *testing.T, r *http.Request, appToken, userAccessToken string) {
	if akahuIdHeader := r.Header.Get("X-Akahu-ID"); akahuIdHeader != appToken {
		t.Fatalf("expected header X-Akahu-ID %s, actual %s", appToken, akahuIdHeader)
	}

	expectedAuthHeader := fmt.Sprintf("Bearer %s", userAccessToken)
	if authorizationHeader := r.Header.Get("Authorization"); authorizationHeader != expectedAuthHeader {
		t.Fatalf("expected header Authorization %s, actual %s", expectedAuthHeader, authorizationHeader)
	}
}

func testBasicRequestHeaders(t *testing.T, r *http.Request) {
	// Ensure basic auth header is base64 encoded "app_token_123:appSecret123"
	expectedAuthHeader := fmt.Sprintf("Basic %s", "YXBwX3Rva2VuXzEyMzphcHBTZWNyZXQxMjM=")

	if authorizationHeader := r.Header.Get("Authorization"); authorizationHeader != expectedAuthHeader {
		t.Fatalf("expected header Authorization %s, actual %s", expectedAuthHeader, authorizationHeader)
	}
}

func testClientResponse(t *testing.T, expected, actual interface{}, err error) {
	if err != nil {
		t.Fatalf("client request returned err %v", err)
	}

	expectedCmp := expected
	if expectedRaw := reflect.ValueOf(expected); expectedRaw.Kind() == reflect.Ptr {
		expectedCmp = expectedRaw.Elem()
	}

	actualCmp := actual
	if actualRaw := reflect.ValueOf(actual); actualRaw.Kind() == reflect.Ptr {
		actualCmp = actualRaw.Elem()
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %+v, actual %+v", expectedCmp, actualCmp)
	}
}

func testClientAPIResponse(t *testing.T, expected, actual *APIResponse, err error) {
	if err != nil {
		t.Fatalf("client request returned err %v", err)
	}

	if expected.Success != actual.Success {
		t.Fatalf("expected APIResponse Success %t, actual %t", expected.Success, actual.Success)
	}

	if expected.Message != actual.Message {
		t.Fatalf("expected APIResponse Message %s, actual %s", expected.Message, actual.Message)
	}
}
