package handlers

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// DialogHandler handles native dialog operations
type DialogHandler struct {
	ctx context.Context
}

// NewDialogHandler creates a new DialogHandler
func NewDialogHandler() *DialogHandler {
	return &DialogHandler{}
}

// SetContext sets the Wails context
func (h *DialogHandler) SetContext(ctx context.Context) {
	h.ctx = ctx
}

// OpenFileDialog opens a native file selection dialog
func (h *DialogHandler) OpenFileDialog(title string) (string, error) {
	return runtime.OpenFileDialog(h.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PostMe Files (*.postme)",
				Pattern:     "*.postme",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
}

// SaveFileDialog opens a native file save dialog
func (h *DialogHandler) SaveFileDialog(title string, defaultFilename string) (string, error) {
	return runtime.SaveFileDialog(h.ctx, runtime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultFilename,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PostMe Files (*.postme)",
				Pattern:     "*.postme",
			},
		},
	})
}
