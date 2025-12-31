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
		{"image", "image/png", ".png"},
		{"image", "image/jpeg", ".jpg"},
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
		{"empty filename prefers ID over SourceUrl", enex.Resource{
			ID: "1xxx590685x61x4xxx1x24xxxxx0097x",
			Attributes: enex.Attributes{
				Filename:  "",
				SourceUrl: "en-cache://tokenKey%3D%22AuthToken%3AUser%3A00000000%22+0x00xx00-0000-000x-0000-xx00x0xxx0x0+1xxx590685x61x4xxx1x24xxxxx0097x+https%3A%2F%2Fpublic.www.evernote.com%2Fresources%2Fx000%2F000000x0-0x0x-000x-xx0x-x0000000000x",
			},
		}, "1xxx590685x61x4xxx1x24xxxxx0097x"},
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
