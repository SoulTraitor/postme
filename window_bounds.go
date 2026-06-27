package main

import (
	"context"

	"github.com/SoulTraitor/postme/internal/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	minWindowWidth              = 800
	minWindowHeight             = 600
	minVisibleWindowPixels      = 80
	fallbackDesktopSpanLimit    = 10000
	windowPositionModeNativeMac = "native-mac"
	windowPositionModeWails     = "wails"
)

type displayBounds struct {
	x      int
	y      int
	width  int
	height int
}

type windowBounds struct {
	x            int
	y            int
	width        int
	height       int
	positionMode string
}

func restoreSavedWindowBounds(ctx context.Context, savedState *models.AppState, windowWidth, windowHeight int) {
	if savedState == nil {
		return
	}
	if savedState.WindowMaximized && !nativeWindowBoundsSupported() {
		return
	}

	displays := getDisplayBounds(ctx)
	width, height := clampWindowSizeForDisplays(windowWidth, windowHeight, displays)

	if nativeWindowBoundsSupported() {
		restoreSavedNativeWindowBounds(ctx, savedState, width, height, displays)
		return
	}

	if width != windowWidth || height != windowHeight {
		runtime.WindowSetSize(ctx, width, height)
	}

	if savedState.WindowX == nil || savedState.WindowY == nil {
		runtime.WindowCenter(ctx)
		return
	}

	x, y := *savedState.WindowX, *savedState.WindowY
	if runtimeWindowPositionUsableForDisplays(x, y, width, height, displays) {
		runtime.WindowSetPosition(ctx, x, y)
		return
	}

	runtime.WindowCenter(ctx)
}

func restoreSavedNativeWindowBounds(ctx context.Context, savedState *models.AppState, width, height int, displays []displayBounds) {
	if savedState.WindowX == nil ||
		savedState.WindowY == nil ||
		savedState.WindowPositionMode != windowPositionModeNativeMac {
		runtime.WindowCenter(ctx)
		return
	}

	bounds := windowBounds{
		x:            *savedState.WindowX,
		y:            *savedState.WindowY,
		width:        width,
		height:       height,
		positionMode: windowPositionModeNativeMac,
	}
	if nativeWindowPositionUsableForDisplays(bounds.x, bounds.y, bounds.width, bounds.height, displays) &&
		setNativeWindowBounds(bounds) {
		return
	}

	runtime.WindowCenter(ctx)
}

func getCurrentWindowBounds(ctx context.Context) (windowBounds, bool) {
	if bounds, ok := getNativeWindowBounds(); ok {
		return bounds, true
	}

	width, height := runtime.WindowGetSize(ctx)
	x, y := runtime.WindowGetPosition(ctx)
	return windowBounds{
		x:            x,
		y:            y,
		width:        width,
		height:       height,
		positionMode: windowPositionModeWails,
	}, true
}

func shouldSaveWindowBounds(ctx context.Context, bounds windowBounds) bool {
	if bounds.width < minWindowWidth || bounds.height < minWindowHeight {
		return false
	}

	displays := getDisplayBounds(ctx)
	if bounds.positionMode == windowPositionModeNativeMac {
		return nativeWindowPositionUsableForDisplays(bounds.x, bounds.y, bounds.width, bounds.height, displays)
	}

	return runtimeWindowPositionUsableForDisplays(bounds.x, bounds.y, bounds.width, bounds.height, displays)
}

func getDisplayBounds(ctx context.Context) []displayBounds {
	if displays, ok := getNativeDisplayBounds(); ok {
		return displays
	}

	screens, err := runtime.ScreenGetAll(ctx)
	if err != nil {
		return nil
	}

	displays := make([]displayBounds, 0, len(screens))
	for _, screen := range screens {
		width, height := screenLogicalSize(screen)
		if width > 0 && height > 0 {
			displays = append(displays, displayBounds{width: width, height: height})
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

func clampWindowSizeForDisplays(width, height int, displays []displayBounds) (int, int) {
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

func runtimeWindowPositionUsableForDisplays(x, y, width, height int, displays []displayBounds) bool {
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

func nativeWindowPositionUsableForDisplays(x, y, width, height int, displays []displayBounds) bool {
	if width <= 0 || height <= 0 {
		return false
	}

	if len(displays) == 0 {
		return x > -fallbackDesktopSpanLimit &&
			x < fallbackDesktopSpanLimit &&
			y > -fallbackDesktopSpanLimit &&
			y < fallbackDesktopSpanLimit
	}

	minVisibleWidth := minInt(width, minVisibleWindowPixels)
	minVisibleHeight := minInt(height, minVisibleWindowPixels)
	windowTop := y + height

	for _, display := range displays {
		if display.width <= 0 || display.height <= 0 {
			continue
		}

		displayRight := display.x + display.width
		displayTop := display.y + display.height
		if spanOverlapsBy(x, width, display.x, displayRight, minVisibleWidth) &&
			spanOverlapsBy(y, height, display.y, displayTop, minVisibleHeight) &&
			topEdgeFallsNearDisplay(windowTop, display.y, displayTop, minVisibleWindowPixels) {
			return true
		}
	}

	return false
}

func displayEnvelope(displays []displayBounds) (left, right, top, bottom int, ok bool) {
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

func maxDisplaySize(displays []displayBounds) (int, int) {
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

func topEdgeFallsNearDisplay(top, boundStart, boundEnd, tolerance int) bool {
	return top >= boundStart+tolerance && top <= boundEnd+tolerance
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
