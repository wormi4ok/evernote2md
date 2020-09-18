package internal

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"path"
	"regexp"
	"strings"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/file"
)

var reImg = regexp.MustCompile(`^image/[\w]+`)

var reFileAndExt = regexp.MustCompile(`(.*)(\.[\w\d]+)`)

func decoder(d enex.Data) io.Reader {
	if d.Encoding == "base64" {
		return base64.NewDecoder(base64.StdEncoding, bytes.NewReader(bytes.TrimSpace(d.Content)))
	}

	return bytes.NewReader(d.Content)
}

func isImage(mimeType string) bool {
	return reImg.MatchString(mimeType)
}

func name(r enex.Resource) (name string, extension string) {
	name = guessName(r)
	// Try to split a file into name and extension
	ff := reFileAndExt.FindStringSubmatch(name)
	if len(ff) < 2 {
		// Guess the extension by the mime type
		return file.BaseName(name), guessExt(r.Mime)
	}

	// Return sanitized filename
	return file.BaseName(ff[len(ff)-2]), ff[len(ff)-1]
}

// guessName of the res with the following priority:
// 1. Filename attribute
// 2. SourceUrl attribute
// 3. ID of the res
// 4. File type as name
func guessName(r enex.Resource) string {
	switch {
	case r.Attributes.Filename != "":
		return r.Attributes.Filename
	case r.Attributes.SourceUrl != "":
		return strings.TrimSpace(path.Base(r.Attributes.SourceUrl))
	case r.ID != "":
		return r.ID
	default:
		return r.Type
	}
}

func guessExt(mimeType string) string {
	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(ext) == 0 {
		return ""
	}
	return ext[0]
}
