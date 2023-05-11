package akahu

import (
	"context"
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
	ItemId string `json:"item_id"`
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

func (s *WebhooksService) List(ctx context.Context, userAccessToken string) ([]WebhookResponse, *http.Response, error) {
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

func (s *WebhooksService) GetPublicKey(ctx context.Context) (*string, *http.Response, error) {
	r, err := s.client.newRequest(http.MethodGet, publicKeyPath, nil, withBasicAuthRequestConfig())
	if err != nil {
		return nil, nil, err
	}

	var publicKey itemResponse[string]
	res, err := s.client.do(ctx, r, &publicKey)
	if err != nil {
		return nil, nil, err
	}

	return &publicKey.Item, res, nil
}

func (s *WebhooksService) ListEvents(ctx context.Context, userAccessToken, status string, startTime, endTime time.Time) ([]WebHookEventResponse, *http.Response, error) {
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

func (s *WebhooksService) Subscribe(ctx context.Context, userAccessToken string, body WebhookSubscribeRequest) (*string, *http.Response, error) {
	r, err := s.client.newRequest(http.MethodPost, webhooksPath, body, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return nil, nil, err
	}

	var webhookSubscribe WebhookSubscribeResponse
	res, err := s.client.do(ctx, r, &webhookSubscribe)
	if err != nil {
		return nil, nil, err
	}

	return &webhookSubscribe.ItemId, res, nil
}

func (s *WebhooksService) Unsubscribe(ctx context.Context, userAccessToken, id string) (bool, *http.Response, error) {
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