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
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
}
