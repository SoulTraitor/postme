package models

// Configuration constants
const (
	// MaxHistoryRecords is the maximum number of history records to keep
	MaxHistoryRecords = 100

	// WindowResizeDebounce is the debounce time for window resize events in milliseconds
	WindowResizeDebounce = 500

	// ContentEditDebounce is the debounce time for content edit events in milliseconds
	ContentEditDebounce = 1000

	// DefaultWindowWidth is the default window width
	DefaultWindowWidth = 1200

	// DefaultWindowHeight is the default window height
	DefaultWindowHeight = 800

	// DefaultSidebarWidth is the default sidebar width
	DefaultSidebarWidth = 260

	// DefaultSplitRatio is the default split ratio percentage
	DefaultSplitRatio = 50

	// DefaultRequestTimeout is the default request timeout in seconds (0 means no limit)
	DefaultRequestTimeout = 30.0
)
