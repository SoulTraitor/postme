package main

import (
	"context"

	"github.com/SoulTraitor/postme/internal/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	minWindowWidth           = 800
	minWindowHeight          = 600
	minVisibleWindowPixels   = 80
	fallbackDesktopSpanLimit = 10000
)

type displaySize struct {
	width  int
	height int
}

func restoreSavedWindowBounds(ctx context.Context, savedState *models.AppState, windowWidth, windowHeight int) {
	if savedState == nil || savedState.WindowMaximized {
		return
	}

	displays := getDisplaySizes(ctx)
	width, height := clampWindowSizeForDisplays(windowWidth, windowHeight, displays)
	if width != windowWidth || height != windowHeight {
		runtime.WindowSetSize(ctx, width, height)
	}

	if savedState.WindowX == nil || savedState.WindowY == nil {
		runtime.WindowCenter(ctx)
		return
	}

	x, y := *savedState.WindowX, *savedState.WindowY
	if savedWindowPositionUsableForDisplays(x, y, width, height, displays) {
		runtime.WindowSetPosition(ctx, x, y)
		return
	}

	runtime.WindowCenter(ctx)
}

func shouldSaveWindowBounds(ctx context.Context, x, y, width, height int) bool {
	if width < minWindowWidth || height < minWindowHeight {
		return false
	}

	return savedWindowPositionUsableForDisplays(x, y, width, height, getDisplaySizes(ctx))
}

func getDisplaySizes(ctx context.Context) []displaySize {
	screens, err := runtime.ScreenGetAll(ctx)
	if err != nil {
		return nil
	}

	displays := make([]displaySize, 0, len(screens))
	for _, screen := range screens {
		width, height := screenLogicalSize(screen)
		if width > 0 && height > 0 {
			displays = append(displays, displaySize{width: width, height: height})
		}
	}

	return displays
}

func screenLogicalSize(screen runtime.Screen) (int, int) {
	width := screen.Size.Width
	height := screen.Size.Height

	if width <= 0 {
		width = screen.Width
	}
	if height <= 0 {
		height = screen.Height
	}

	return width, height
}

func clampWindowSizeForDisplays(width, height int, displays []displaySize) (int, int) {
	width = maxInt(width, minWindowWidth)
	height = maxInt(height, minWindowHeight)

	maxWidth, maxHeight := maxDisplaySize(displays)
	if maxWidth <= 0 || maxHeight <= 0 {
		return width, height
	}

	maxWidth = maxInt(maxWidth, minWindowWidth)
	maxHeight = maxInt(maxHeight, minWindowHeight)

	return minInt(width, maxWidth), minInt(height, maxHeight)
}

func savedWindowPositionUsableForDisplays(x, y, width, height int, displays []displaySize) bool {
	if width <= 0 || height <= 0 {
		return false
	}

	left, right, top, bottom, ok := displayEnvelope(displays)
	if !ok {
		return x > -fallbackDesktopSpanLimit &&
			x < fallbackDesktopSpanLimit &&
			y > -fallbackDesktopSpanLimit &&
			y < fallbackDesktopSpanLimit
	}

	minVisibleWidth := minInt(width, minVisibleWindowPixels)

	return spanOverlapsBy(x, width, left, right, minVisibleWidth) &&
		startFallsWithin(y, top, bottom, minVisibleWindowPixels)
}

func displayEnvelope(displays []displaySize) (left, right, top, bottom int, ok bool) {
	validDisplays := 0
	totalWidth := 0
	totalHeight := 0
	maxWidth := 0
	maxHeight := 0

	for _, display := range displays {
		if display.width <= 0 || display.height <= 0 {
			continue
		}
		validDisplays++
		totalWidth += display.width
		totalHeight += display.height
		maxWidth = maxInt(maxWidth, display.width)
		maxHeight = maxInt(maxHeight, display.height)
	}

	if validDisplays == 0 {
		return 0, 0, 0, 0, false
	}

	if validDisplays == 1 {
		return 0, maxWidth, 0, maxHeight, true
	}

	return -totalWidth, totalWidth, -totalHeight, totalHeight, true
}

func maxDisplaySize(displays []displaySize) (int, int) {
	maxWidth := 0
	maxHeight := 0
	for _, display := range displays {
		maxWidth = maxInt(maxWidth, display.width)
		maxHeight = maxInt(maxHeight, display.height)
	}
	return maxWidth, maxHeight
}

func spanOverlapsBy(start, length, boundStart, boundEnd, minOverlap int) bool {
	overlap := minInt(start+length, boundEnd) - maxInt(start, boundStart)
	return overlap >= minOverlap
}

func startFallsWithin(start, boundStart, boundEnd, tolerance int) bool {
	return start >= boundStart-tolerance && start <= boundEnd-tolerance
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
