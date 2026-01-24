package models

import "time"

// KeyValue represents a key-value pair with enabled state
type KeyValue struct {
	Key     string `json:"key" db:"key"`
	Value   string `json:"value" db:"value"`
	Enabled bool   `json:"enabled" db:"enabled"`
	Type    string `json:"type,omitempty" db:"type"` // "text" or "file" for form-data
}

// Request represents an HTTP request
type Request struct {
	ID           int64      `json:"id" db:"id"`
	CollectionID int64      `json:"collectionId" db:"collection_id"`
	FolderID     *int64     `json:"folderId" db:"folder_id"`
	Name         string     `json:"name" db:"name"`
	Method       string     `json:"method" db:"method"`
	URL          string     `json:"url" db:"url"`
	Headers      []KeyValue `json:"headers" db:"-"`
	HeadersJSON  string     `json:"-" db:"headers"`
	Params       []KeyValue `json:"params" db:"-"`
	ParamsJSON   string     `json:"-" db:"params"`
	Body         string     `json:"body" db:"body"`
	BodyType     string     `json:"bodyType" db:"body_type"`
	SortOrder    int        `json:"sortOrder" db:"sort_order"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
}
