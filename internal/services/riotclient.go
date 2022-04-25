package services

import (
	"context"
	"io"
	"net/http"
	"riot-developer-proxy/internal/valueobjects"
)

type RiotClient struct {
	client    http.Client
	baseURI   string
	riotToken *string
}

func NewRiotClient(client http.Client, baseURI string) *RiotClient {
	return &RiotClient{
		client:  client,
		baseURI: baseURI,
	}
}

func (rc *RiotClient) WithToken(token string) {
	rc.riotToken = &token
}

func (rc *RiotClient) DoReq(ctx context.Context, method string, path string) (*valueobjects.ProxyResponse, error) {
	req, err := http.NewRequestWithContext(ctx, method, rc.baseURI+path, nil)

	if err != nil {
		return nil, err
	}

	if rc.riotToken != nil {
		req.Header.Add("X-Riot-Token", *rc.riotToken)
	}

	resp, err := rc.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	response := &valueobjects.ProxyResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}

	return response, nil
}
