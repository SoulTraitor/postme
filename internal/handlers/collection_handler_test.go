package handlers

import "testing"

func TestSanitizeExportFilename(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "Team/API:Prod?", want: "Team_API_Prod_.postme"},
		{name: "  collection.postme  ", want: "collection.postme"},
		{name: "CON", want: "_CON.postme"},
		{name: "...", want: "collection.postme"},
		{name: "", want: "collection.postme"},
	}

	for _, tt := range tests {
		if got := sanitizeExportFilename(tt.name); got != tt.want {
			t.Fatalf("sanitizeExportFilename(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}
