package models

import "time"

// AppState represents the application state (single row)
type AppState struct {
	ID                int64     `json:"id" db:"id"`
	WindowWidth       int       `json:"windowWidth" db:"window_width"`
	WindowHeight      int       `json:"windowHeight" db:"window_height"`
	WindowX           *int      `json:"windowX" db:"window_x"`
	WindowY           *int      `json:"windowY" db:"window_y"`
	WindowMaximized   bool      `json:"windowMaximized" db:"window_maximized"`
	SidebarOpen       bool      `json:"sidebarOpen" db:"sidebar_open"`
	SidebarWidth      int       `json:"sidebarWidth" db:"sidebar_width"`
	LayoutDirection   string    `json:"layoutDirection" db:"layout_direction"`
	SplitRatio        int       `json:"splitRatio" db:"split_ratio"`
	Theme             string    `json:"theme" db:"theme"`
	ActiveEnvID       *int64    `json:"activeEnvId" db:"active_env_id"`
	RequestTimeout    float64   `json:"requestTimeout" db:"request_timeout"`
	AutoLocateSidebar bool      `json:"autoLocateSidebar" db:"auto_locate_sidebar"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
}

// SidebarState represents the expanded/collapsed state of sidebar items
type SidebarState struct {
	ID       int64  `json:"id" db:"id"`
	ItemType string `json:"itemType" db:"item_type"`
	ItemID   int64  `json:"itemId" db:"item_id"`
	Expanded bool   `json:"expanded" db:"expanded"`
}

// TabSession represents a tab session
type TabSession struct {
	ID           int64      `json:"id" db:"id"`
	TabID        string     `json:"tabId" db:"tab_id"`
	RequestID    *int64     `json:"requestId" db:"request_id"`
	Title        string     `json:"title" db:"title"`
	SortOrder    int        `json:"sortOrder" db:"sort_order"`
	IsActive     bool       `json:"isActive" db:"is_active"`
	IsDirty      bool       `json:"isDirty" db:"is_dirty"`
	Method       string     `json:"method" db:"method"`
	URL          string     `json:"url" db:"url"`
	Headers      []KeyValue `json:"headers" db:"-"`
	HeadersJSON  string     `json:"-" db:"headers"`
	Params       []KeyValue `json:"params" db:"-"`
	ParamsJSON   string     `json:"-" db:"params"`
	Body         string     `json:"body" db:"body"`
	BodyType     string     `json:"bodyType" db:"body_type"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
}
