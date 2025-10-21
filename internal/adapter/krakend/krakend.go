package krakend

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lapotkin/file-storage/internal/domain/json"
)

// Client represents a KrakenD API Gateway client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// BaseURL returns the base URL of the KrakenD gateway
func (c *Client) BaseURL() string {
	return c.baseURL
}

type Request struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
}

type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

func (c *Client) Do(req *Request) (*Response, error) {
	var bodyReader io.Reader
	if req.Body != nil {
		if rawBytes, ok := req.Body.([]byte); ok {
			bodyReader = bytes.NewReader(rawBytes)
		} else {
			jsonData, err := json.Marshal(req.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewReader(jsonData)
		}
	}

	httpReq, err := http.NewRequest(req.Method, c.baseURL+req.Path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if _, ok := req.Headers["Content-Type"]; !ok {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	if _, ok := req.Headers["Accept"]; !ok {
		httpReq.Header.Set("Accept", "application/json")
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: httpResp.StatusCode,
		Body:       body,
		Headers:    httpResp.Header,
	}, nil
}

func (r *Response) DecodeJSON(v interface{}) error {
	if err := json.Unmarshal(r.Body, v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

func (r *Response) ContentType() string {
	return r.Headers.Get("Content-Type")
}
