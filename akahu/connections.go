package akahu

import (
	"context"
	"net/http"
	"path"
)

const connectionsPath = "connections"

type ConnectionsService service

type ConnectionResponse struct {
	Id   string  `json:"_id"`
	Name string  `json:"name"`
	Url  *string `json:"url"`
	Logo string  `json:"logo"`
}

func (s *ConnectionsService) List(ctx context.Context) ([]ConnectionResponse, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodGet, connectionsPath, nil, withBasicAuthRequestConfig())
	if err != nil {
		return nil, nil, err
	}

	var connections collectionResponse[ConnectionResponse]
	res, err := s.client.do(ctx, r, &connections)
	if err != nil {
		return nil, nil, err
	}

	return connections.Items, res, nil
}

func (s *ConnectionsService) Get(ctx context.Context, connectionId string) (*ConnectionResponse, *APIResponse, error) {
	r, err := s.client.newRequest(http.MethodGet, path.Join(connectionsPath, connectionId), nil, withBasicAuthRequestConfig())
	if err != nil {
		return nil, nil, err
	}

	var connections itemResponse[ConnectionResponse]
	res, err := s.client.do(ctx, r, &connections)
	if err != nil {
		return nil, nil, err
	}

	return connections.Item, res, nil
}
