package main

import "testing"

func TestSavedWindowPositionUsableForSingleDisplay(t *testing.T) {
	displays := []displaySize{{width: 1512, height: 982}}

	tests := []struct {
		name string
		x    int
		y    int
		want bool
	}{
		{name: "normal position", x: 120, y: 80, want: true},
		{name: "partially off left but visible", x: -1000, y: 80, want: true},
		{name: "fully off left", x: -1300, y: 80, want: false},
		{name: "top edge too far above screen", x: 120, y: -700, want: false},
		{name: "off right", x: 2000, y: 80, want: false},
		{name: "off bottom", x: 120, y: 1200, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := savedWindowPositionUsableForDisplays(tt.x, tt.y, 1200, 800, displays)
			if got != tt.want {
				t.Fatalf("savedWindowPositionUsableForDisplays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSavedWindowPositionUsableForMultipleDisplays(t *testing.T) {
	displays := []displaySize{
		{width: 1512, height: 982},
		{width: 1920, height: 1080},
	}

	tests := []struct {
		name string
		x    int
		y    int
		want bool
	}{
		{name: "external display on left", x: -1800, y: 100, want: true},
		{name: "external display on right", x: 3000, y: 100, want: true},
		{name: "far outside horizontal envelope", x: 9000, y: 100, want: false},
		{name: "far outside vertical envelope", x: 100, y: 7000, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := savedWindowPositionUsableForDisplays(tt.x, tt.y, 1200, 800, displays)
			if got != tt.want {
				t.Fatalf("savedWindowPositionUsableForDisplays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClampWindowSizeForDisplays(t *testing.T) {
	displays := []displaySize{{width: 1440, height: 900}}

	tests := []struct {
		name       string
		width      int
		height     int
		wantWidth  int
		wantHeight int
	}{
		{name: "clamps oversized window", width: 3000, height: 2000, wantWidth: 1440, wantHeight: 900},
		{name: "keeps minimum window size", width: 400, height: 300, wantWidth: 800, wantHeight: 600},
		{name: "keeps usable saved size", width: 1200, height: 800, wantWidth: 1200, wantHeight: 800},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, gotHeight := clampWindowSizeForDisplays(tt.width, tt.height, displays)
			if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
				t.Fatalf("clampWindowSizeForDisplays() = (%d, %d), want (%d, %d)",
					gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
			}
		})
	}
}
