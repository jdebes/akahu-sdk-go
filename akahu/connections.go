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

// List Gets a list of all connected financial institutions that users can connect to your Akahu application.
//
// Akahu docs: https://developers.akahu.nz/reference/get_connections
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

// Get fetches an individual financial institution connection.
//
// Akahu docs: https://developers.akahu.nz/reference/get_connections-id
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
