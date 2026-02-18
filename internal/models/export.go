package models

import "time"

// ExportFile is the top-level export file structure
type ExportFile struct {
	Version    int              `json:"version"`
	ExportedAt time.Time       `json:"exportedAt"`
	Collection ExportCollection `json:"collection"`
}

// ExportCollection represents a collection without IDs/timestamps
type ExportCollection struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Folders     []ExportFolder  `json:"folders"`
	Requests    []ExportRequest `json:"requests"`
}

// ExportFolder represents a folder without IDs/timestamps
type ExportFolder struct {
	Name      string          `json:"name"`
	SortOrder int             `json:"sortOrder"`
	Requests  []ExportRequest `json:"requests"`
}

// ExportRequest represents a request without IDs/timestamps
type ExportRequest struct {
	Name      string     `json:"name"`
	Method    string     `json:"method"`
	URL       string     `json:"url"`
	Headers   []KeyValue `json:"headers"`
	Params    []KeyValue `json:"params"`
	Body      string     `json:"body"`
	BodyType  string     `json:"bodyType"`
	SortOrder int        `json:"sortOrder"`
}
