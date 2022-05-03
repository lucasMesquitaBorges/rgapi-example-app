package rgapis

import (
	"context"
	"fmt"
	"net/http"
)

type rgapiClient struct {
	client     *http.Client
	rgapiToken string
}

const RGAPI_TOKEN_HEADER_NAME = "X-Riot-Token"
const RGAPI_BASE_URI = "https://%s.api.riotgames.com%s"

func NewRGAPIClient(rgapiToken string) *rgapiClient {
	return &rgapiClient{
		client:     http.DefaultClient,
		rgapiToken: rgapiToken,
	}
}

func (rc *rgapiClient) DoReq(ctx context.Context, method string, region string, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf(RGAPI_BASE_URI, region, path),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add(RGAPI_TOKEN_HEADER_NAME, rc.rgapiToken)

	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
