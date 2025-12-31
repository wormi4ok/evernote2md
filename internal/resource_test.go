package internal

import (
	"bytes"
	"encoding/base64"
	"io"
	"testing"

	"github.com/wormi4ok/evernote2md/encoding/enex"
)

func Test_guessExt(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		want     string
	}{
		{"PNG image", "image/png", ".png"},
		{"JPEG image", "image/jpeg", ".jpg"},
		{"GIF image", "image/gif", ".gif"},
		{"WebP image", "image/webp", ".webp"},
		{"SVG image", "image/svg+xml", ".svg"},
		{"TIFF image", "image/tiff", ".tiff"},
		{"BMP image", "image/bmp", ".bmp"},
		{"ICO image", "image/ico", ".ico"},
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

func Test_decoder(t *testing.T) {
	want := []byte("sample text")
	encoded := new(bytes.Buffer)
	b64encoder := base64.NewEncoder(base64.StdEncoding, encoded)

	if _, err := b64encoder.Write(want); err != nil {
		t.Error(err)
	}
	if err := b64encoder.Close(); err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
		data enex.Data
	}{
		{
			name: "not encoded",
			data: enex.Data{
				Encoding: "",
				Content:  want,
			},
		},
		{
			name: "base64 encoded",
			data: enex.Data{
				Encoding: "base64",
				Content:  encoded.Bytes(),
			},
		},
		{
			name: "base64 encoded - encoding value missing",
			data: enex.Data{
				Encoding: "",
				Content:  encoded.Bytes(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := io.ReadAll(decoder(tt.data))
			if err != nil {
				t.Error(err)
			}
			if !bytes.Equal(got, want) {
				t.Errorf("decoder() = %s, want %s", got, want)
			}
		})
	}
}
