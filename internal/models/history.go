package models

import "time"

// History represents a request history record
type History struct {
	ID              int64     `json:"id" db:"id"`
	RequestID       *int64    `json:"requestId" db:"request_id"`
	Method          string    `json:"method" db:"method"`
	URL             string    `json:"url" db:"url"`
	RequestHeaders  string    `json:"requestHeaders" db:"request_headers"`
	RequestBody     string    `json:"requestBody" db:"request_body"`
	StatusCode      *int      `json:"statusCode" db:"status_code"`
	ResponseHeaders string    `json:"responseHeaders" db:"response_headers"`
	ResponseBody    string    `json:"responseBody" db:"response_body"`
	DurationMs      *int64    `json:"durationMs" db:"duration_ms"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}
