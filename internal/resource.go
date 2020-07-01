package internal

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"regexp"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/file"
)

var reImg = regexp.MustCompile(`^image/[\w]+`)

var reFileAndExt = regexp.MustCompile(`(.*)(\.[^.]+)`)

func decoder(d enex.Data) io.Reader {
	if d.Encoding == "base64" {
		return base64.NewDecoder(base64.StdEncoding, bytes.NewReader(d.Content))
	}

	return bytes.NewReader(d.Content)
}

func isImage(mimeType string) bool {
	return reImg.MatchString(mimeType)
}

func guessName(r enex.Resource) (name string, extension string) {
	// Use ID as name if there is no resource name
	if r.Attributes.Filename == "" {
		return r.ID, guessExt(r.Mime)
	}

	// Try to split a file into name and extension
	ff := reFileAndExt.FindStringSubmatch(r.Attributes.Filename)
	if len(ff) < 2 {
		// Use only filename if there is no extension
		return file.BaseName(r.Attributes.Filename), ""
	}

	// Return sanitized filename
	return file.BaseName(ff[len(ff)-2]), ff[len(ff)-1]
}

func guessExt(mimeType string) string {
	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(ext) == 0 {
		return ""
	}
	return ext[0]
}
