package file

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Max path length is 255 - 9 bytes for extension (.md) in multibyte encoding
const MaxPathLength int = 246

var (
	baseNameSeparators = regexp.MustCompile(`[./]`)

	blacklist = regexp.MustCompile(`[\s|"'<>&_=+:?]`)

	dashes = regexp.MustCompile(`[\-_]{2,}`)
)

// Save a new file in a given dir with the following content.
// Creates a directory if necessary.
func Save(dir, name string, content io.Reader) error {
	if len(name) == 0 {
		return nil
	}

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	output, err := os.Create(filepath.FromSlash(dir + "/" + name))
	if err != nil {
		return err
	}

	_, err = io.Copy(output, content)

	err = output.Close()
	return err
}

// BaseName normalizes a given string to use it as a safe filename
func BaseName(s string) string {
	// Replace separator characters with a dash
	s = baseNameSeparators.ReplaceAllString(s, "-")

	// Remove any trailing space to avoid ending on -
	s = strings.Trim(s, " ")

	// Replace inappropriate characters with an underscore
	s = blacklist.ReplaceAllString(s, "_")

	// Remove any multiple dashes caused by replacements above
	s = dashes.ReplaceAllString(s, "-")

	// Check file name length in bytes
	if len(s) <= MaxPathLength {
		return s
	}

	// Trim filename to the max allowed number of bytes
	var sb strings.Builder
	for index, c := range s {
		if sb.Len()+index >= MaxPathLength {
			return sb.String()
		}
		sb.WriteRune(c)
	}

	return s
}
