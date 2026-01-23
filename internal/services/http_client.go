package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/SoulTraitor/postme/internal/models"
)

// HTTPClient handles HTTP request execution
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new HTTPClient
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{},
	}
}

// ExecuteRequest represents the request to execute
type ExecuteRequest struct {
	Method   string            `json:"method"`
	URL      string            `json:"url"`
	Headers  []models.KeyValue `json:"headers"`
	Body     string            `json:"body"`
	BodyType string            `json:"bodyType"`
	Timeout  float64           `json:"timeout"`
}

// Execute executes an HTTP request
func (c *HTTPClient) Execute(ctx context.Context, req ExecuteRequest) (*models.Response, error) {
	// Create context with timeout if specified
	if req.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(req.Timeout*float64(time.Second)))
		defer cancel()
	}

	// Create body reader
	var bodyReader io.Reader
	if req.Body != "" && req.BodyType != "none" {
		bodyReader = strings.NewReader(req.Body)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set headers
	for _, h := range req.Headers {
		if h.Enabled && h.Key != "" {
			httpReq.Header.Set(h.Key, h.Value)
		}
	}

	// Set content-type based on body type
	if req.BodyType == "json" && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	} else if req.BodyType == "xml" && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/xml")
	} else if req.BodyType == "text" && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "text/plain")
	}

	// Execute request
	startTime := time.Now()
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	duration := time.Since(startTime).Milliseconds()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract headers
	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &models.Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    headers,
		Body:       string(body),
		Size:       int64(len(body)),
		Duration:   duration,
	}, nil
}

// BuildRequestHeadersJSON builds JSON string from headers
func BuildRequestHeadersJSON(headers []models.KeyValue) string {
	data, _ := json.Marshal(headers)
	return string(data)
}

// BuildResponseHeadersJSON builds JSON string from response headers
func BuildResponseHeadersJSON(headers map[string]string) string {
	data, _ := json.Marshal(headers)
	return string(data)
}
