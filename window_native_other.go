//go:build !darwin

package main

func nativeWindowBoundsSupported() bool {
	return false
}

func getNativeWindowBounds() (windowBounds, bool) {
	return windowBounds{}, false
}

func setNativeWindowBounds(bounds windowBounds) bool {
	return false
}

func getNativeDisplayBounds() ([]displayBounds, bool) {
	return nil, false
}
