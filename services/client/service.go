package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service struct {
	client HttpClient
}

func New(client HttpClient) *Service {
	return &Service{client: client}
}

// Get makes a request to the specified url
func (s *Service) Get(ctx context.Context, url string) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response status: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
