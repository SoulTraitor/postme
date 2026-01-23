package models

import "time"

// Folder represents a folder within a collection (no nesting allowed)
type Folder struct {
	ID           int64     `json:"id" db:"id"`
	CollectionID int64     `json:"collectionId" db:"collection_id"`
	Name         string    `json:"name" db:"name"`
	SortOrder    int       `json:"sortOrder" db:"sort_order"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}
