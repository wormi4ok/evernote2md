package internal

import (
	"testing"

	"github.com/wormi4ok/evernote2md/encoding/enex"
)

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

func Test_guessName(t *testing.T) {
	tests := []struct {
		name string
		res  enex.Resource
		want string
	}{
		{"filename", enex.Resource{Attributes: enex.Attributes{Filename: "A.png"}}, "A.png"},
		{"sourceUrl", enex.Resource{Attributes: enex.Attributes{SourceUrl: "http://petrashov.ru/C.jpg"}}, "C.jpg"},
		{"ID", enex.Resource{ID: "A"}, "A"},
		{"type", enex.Resource{Type: "C"}, "C"},
		{"order of the fields", enex.Resource{ID: "A", Attributes: enex.Attributes{
			Filename:  "!",
			SourceUrl: "?",
		}}, "!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := guessName(tt.res); got != tt.want {
				t.Errorf("guessName for %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
