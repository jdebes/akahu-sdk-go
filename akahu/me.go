package akahu

import (
	"context"
	"net/http"
	"time"
)

type MeService service

type MeResponse struct {
	Id            string     `json:"_id"`
	CreatedAt     *time.Time `json:"created_at"`
	Email         string     `json:"email"`
	Mobile        *string    `json:"mobile"`
	FirstName     *string    `json:"first_name"`
	LastName      *string    `json:"last_name"`
	PreferredName *string    `json:"preferred_name"`
}

func (s *MeService) Get(ctx context.Context, userAccessToken string) (*MeResponse, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodGet, "me", nil, withTokenRequestConfig(userAccessToken))
	if err != nil {
		return nil, nil, err
	}

	var meResponse itemResponse[MeResponse]
	res, err := s.client.do(ctx, r, &meResponse)
	if err != nil {
		return nil, nil, err
	}

	return meResponse.Item, res, nil
}
