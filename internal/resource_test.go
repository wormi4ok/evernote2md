package internal

import "testing"

func Test_guessExt(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		want     string
	}{
		{"image", "image/png", ".png"},
		{"unknown mime type", "unknown", ""},
		{"empty input", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := guessExt(tt.mimeType); got != tt.want {
				t.Errorf("guessExt for %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
