package models

// Response represents an HTTP response
type Response struct {
	StatusCode int               `json:"statusCode"`
	Status     string            `json:"status"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Size       int64             `json:"size"`
	Duration   int64             `json:"duration"` // milliseconds
}
