package akahu

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"net/http"
	"path"
	"time"
)

type WebhooksService service

const (
	webhooksPath      = "webhooks"
	webhookEventsPath = "webhook-events"
	publicKeyPath     = "keys"
)

// WebhookType represents the types of webhooks that Akahu Supports.
type WebhookType string

const (
	Token       WebhookType = "TOKEN"
	Identity                = "IDENTITY"
	Account                 = "ACCOUNT"
	Transaction             = "TRANSACTION"
	Payment                 = "PAYMENT"
	Transfer                = "TRANSFER"
	Income                  = "INCOME"
)

type WebhookEventStatus string

const (
	Sent   WebhookEventStatus = "SENT"
	Failed                    = "FAILED"
	Retry                     = "RETRY"
)

type WebhookResponse struct {
	Id           string    `json:"_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastCalledAt time.Time `json:"last_called_at"`
	State        string    `json:"state"`
	Url          string    `json:"url"`
}

type WebhookSubscribeRequest struct {
	WebhookType WebhookType `json:"webhook_type"`
	State       string      `json:"state"`
}

type WebhookSubscribeResponse struct {
	successResponse
	ItemId *string `json:"item_id"`
}

type WebHookEventPayload struct {
	successResponse
	WebhookType `json:"webhook_type"`
	WebhookCode string `json:"webhook_code"`
}

type WebHookEventResponse struct {
	Id           string              `json:"_id"`
	Hook         string              `json:"hook"`
	Status       WebhookEventStatus  `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	LastFailedAt time.Time           `json:"last_failed_at"`
	Payload      WebHookEventPayload `json:"payload"`
}

// List gets the active webhook subscriptions that your app has created for the user.
//
// Akahu docs: https://developers.akahu.nz/reference/get_webhooks
func (s *WebhooksService) List(ctx context.Context, userAccessToken string) ([]WebhookResponse, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodGet, webhooksPath, nil, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return nil, nil, err
	}

	var webhooks collectionResponse[WebhookResponse]
	res, err := s.client.do(ctx, r, &webhooks)
	if err != nil {
		return nil, nil, err
	}

	return webhooks.Items, res, nil
}

// GetPublicKey gets one of the public keys that Akahu uses to sign webhooks.
//
// Akahu docs: https://developers.akahu.nz/reference/get_keys-id
func (s *WebhooksService) GetPublicKey(ctx context.Context, id string) (*string, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodGet, path.Join(publicKeyPath, id), nil, withBasicAuthRequestConfig())
	if err != nil {
		return nil, nil, err
	}

	var publicKey itemResponse[string]
	res, err := s.client.do(ctx, r, &publicKey)
	if err != nil {
		return nil, nil, err
	}

	return publicKey.Item, res, nil
}

// ListEvents gets a list of webhook events that have been published to your application by Akahu within the 'start' and 'end' time range.
//
// Akahu docs: https://developers.akahu.nz/reference/get_webhook-events
func (s *WebhooksService) ListEvents(ctx context.Context, userAccessToken, status string, startTime, endTime time.Time) ([]WebHookEventResponse, *APIResponse, error) {
	params := paramsWithDateRange(startTime, endTime)
	params.Add("status", status)
	encodedPath := pathWithParams(webhookEventsPath, params)

	r, err := s.client.newRequest(http.MethodGet, encodedPath, nil, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return nil, nil, err
	}

	var events collectionResponse[WebHookEventResponse]
	res, err := s.client.do(ctx, r, &events)
	if err != nil {
		return nil, nil, err
	}

	return events.Items, res, nil
}

// Subscribe creates a new webhook subscription for the user.
//
// Akahu docs: https://developers.akahu.nz/reference/post_webhooks
func (s *WebhooksService) Subscribe(ctx context.Context, userAccessToken string, body WebhookSubscribeRequest) (*string, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodPost, webhooksPath, body, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return nil, nil, err
	}

	var webhookSubscribe WebhookSubscribeResponse
	res, err := s.client.do(ctx, r, &webhookSubscribe)
	if err != nil {
		return nil, nil, err
	}

	return webhookSubscribe.ItemId, res, nil
}

// Unsubscribe deletes a webhook subscription that your application has previously created for the user.
//
// Akahu docs: https://developers.akahu.nz/reference/delete_webhooks-id
func (s *WebhooksService) Unsubscribe(ctx context.Context, userAccessToken, id string) (bool, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodDelete, path.Join(webhooksPath, id), nil, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return false, nil, err
	}

	var webhookDelete successResponse
	res, err := s.client.do(ctx, r, &webhookDelete)
	if err != nil {
		return false, nil, err
	}

	return webhookDelete.Success, res, nil
}

// ValidateWebhookSignature validates the payload of a Akahu webhook request against private RSA key held by Akahu.
// The public key can be obtained via the WebhooksService.GetPublicKey method.
// The signature is a base64 encoded string that is included in the 'X-Akahu-Signature' header of the request.
//
// See the Webhooks reference for more information: https://developers.akahu.nz/docs/reference-webhooks.
func ValidateWebhookSignature(publicKey, signature string, body []byte) (bool, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return false, errors.New("ssh: no key found")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(body)

	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], decodedSignature) == nil, nil
}
