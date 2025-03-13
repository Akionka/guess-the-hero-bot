package stratz

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const apiURL = "https://api.stratz.com/graphql"

type Client struct {
	httpClient *http.Client
	token      string
}

func NewClient(httpClient *http.Client, token string) *Client {
	return &Client{
		httpClient: httpClient,
		token:      token,
	}
}

type queryRequest struct {
	Query string `json:"query"`
}

type queryResponse[T any] struct {
	Error error `json:"error"`
	Data  *T    `json:"data"`
}

func (c *Client) query(ctx context.Context, query string) (io.ReadCloser, error) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(queryRequest{
		Query: query,
	}); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "STRATZ_API")
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}

type QueryError struct {
	Message string `json:"message"`
}

func (e QueryError) Error() string {
	return e.Message
}
