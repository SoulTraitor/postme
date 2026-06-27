//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>
#include <dispatch/dispatch.h>

typedef struct PMWindowBounds {
	int ok;
	int x;
	int y;
	int width;
	int height;
} PMWindowBounds;

typedef struct PMDisplayBounds {
	int x;
	int y;
	int width;
	int height;
} PMDisplayBounds;

typedef struct PMDisplayList {
	int ok;
	int count;
	PMDisplayBounds displays[16];
} PMDisplayList;

typedef struct PMSetWindowBoundsArgs {
	int ok;
	int x;
	int y;
	int width;
	int height;
} PMSetWindowBoundsArgs;

static NSWindow* PMPostmeWindow() {
	NSWindow *window = [NSApp mainWindow];
	if (window == nil) {
		window = [NSApp keyWindow];
	}
	if (window == nil) {
		for (NSWindow *candidate in [NSApp windows]) {
			if ([candidate isVisible]) {
				window = candidate;
				break;
			}
		}
	}
	return window;
}

static void PMGetWindowBoundsOnMain(void *ptr) {
	PMWindowBounds *result = (PMWindowBounds*)ptr;
	NSWindow *window = PMPostmeWindow();
	if (window == nil) {
		result->ok = 0;
		return;
	}

	NSRect frame = [window frame];
	result->ok = 1;
	result->x = (int)frame.origin.x;
	result->y = (int)frame.origin.y;
	result->width = (int)frame.size.width;
	result->height = (int)frame.size.height;
}

static PMWindowBounds PMGetWindowBounds() {
	PMWindowBounds result = {0, 0, 0, 0, 0};
	if ([NSThread isMainThread]) {
		PMGetWindowBoundsOnMain(&result);
	} else {
		dispatch_sync_f(dispatch_get_main_queue(), &result, PMGetWindowBoundsOnMain);
	}
	return result;
}

static void PMSetWindowBoundsOnMain(void *ptr) {
	PMSetWindowBoundsArgs *args = (PMSetWindowBoundsArgs*)ptr;
	NSWindow *window = PMPostmeWindow();
	if (window == nil || args->width <= 0 || args->height <= 0) {
		args->ok = 0;
		return;
	}

	NSRect frame = NSMakeRect((CGFloat)args->x, (CGFloat)args->y, (CGFloat)args->width, (CGFloat)args->height);
	[window setFrame:frame display:YES animate:NO];
	args->ok = 1;
}

static int PMSetWindowBounds(int x, int y, int width, int height) {
	PMSetWindowBoundsArgs args = {0, x, y, width, height};
	if ([NSThread isMainThread]) {
		PMSetWindowBoundsOnMain(&args);
	} else {
		dispatch_sync_f(dispatch_get_main_queue(), &args, PMSetWindowBoundsOnMain);
	}
	return args.ok;
}

static void PMGetDisplayBoundsOnMain(void *ptr) {
	PMDisplayList *result = (PMDisplayList*)ptr;
	NSArray<NSScreen *> *screens = [NSScreen screens];
	int count = (int)[screens count];
	if (count > 16) {
		count = 16;
	}

	result->ok = 1;
	result->count = count;
	for (int i = 0; i < count; i++) {
		NSScreen *screen = [screens objectAtIndex:i];
		NSRect frame = [screen visibleFrame];
		result->displays[i].x = (int)frame.origin.x;
		result->displays[i].y = (int)frame.origin.y;
		result->displays[i].width = (int)frame.size.width;
		result->displays[i].height = (int)frame.size.height;
	}
}

static PMDisplayList PMGetDisplayBounds() {
	PMDisplayList result = {0, 0, {{0, 0, 0, 0}}};
	if ([NSThread isMainThread]) {
		PMGetDisplayBoundsOnMain(&result);
	} else {
		dispatch_sync_f(dispatch_get_main_queue(), &result, PMGetDisplayBoundsOnMain);
	}
	return result;
}
*/
import "C"

func nativeWindowBoundsSupported() bool {
	return true
}

func getNativeWindowBounds() (windowBounds, bool) {
	bounds := C.PMGetWindowBounds()
	if bounds.ok != 1 {
		return windowBounds{}, false
	}

	return windowBounds{
		x:            int(bounds.x),
		y:            int(bounds.y),
		width:        int(bounds.width),
		height:       int(bounds.height),
		positionMode: windowPositionModeNativeMac,
	}, true
}

func setNativeWindowBounds(bounds windowBounds) bool {
	return C.PMSetWindowBounds(C.int(bounds.x), C.int(bounds.y), C.int(bounds.width), C.int(bounds.height)) == 1
}

func getNativeDisplayBounds() ([]displayBounds, bool) {
	list := C.PMGetDisplayBounds()
	if list.ok != 1 {
		return nil, false
	}

	displays := make([]displayBounds, 0, int(list.count))
	for i := 0; i < int(list.count); i++ {
		display := list.displays[i]
		if display.width <= 0 || display.height <= 0 {
			continue
		}

		displays = append(displays, displayBounds{
			x:      int(display.x),
			y:      int(display.y),
			width:  int(display.width),
			height: int(display.height),
		})
	}

	return displays, true
}
