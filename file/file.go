package file

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// Mon Jan 2 15:04:05 -0700 MST 2006 represented as yyyyMMddhhmm
	touchTimeFormat = "200601021504"

	// OS allow 255 character for filenames = 252 + 3 (.md)
	maxNameChars = 252
)

var (
	baseNameSeparators = regexp.MustCompile(`[./]`)

	dashes = regexp.MustCompile(`[\-_]{2,}`)
)

// Save a new file in a given dir with the following content.
// Creates a directory if necessary.
func Save(dir, name string, content io.Reader) error {
	if len(name) == 0 {
		return nil
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	output, err := os.Create(filepath.FromSlash(dir + "/" + name))
	if err != nil {
		return err
	}

	if _, err = io.Copy(output, content); err != nil {
		_ = output.Close()
		return err
	}

	return output.Close()
}

// BaseName normalizes a given string to use it as a safe filename
func BaseName(s string) string {
	// Replace separator characters with a dash
	s = baseNameSeparators.ReplaceAllString(s, "-")

	// Remove any trailing space to avoid ending on -
	s = strings.Trim(s, " ")

	// Replace inappropriate characters with an underscore
	s = illegalChars.ReplaceAllString(s, "_")

	// Remove any multiple dashes caused by replacements above
	s = dashes.ReplaceAllString(s, "-")

	// Check file name length in bytes
	if len(s) < maxNameChars {
		return s
	}

	// Trim filename to the max allowed number of bytes
	var sb strings.Builder
	var i = 0
	for index, c := range s {
		if index >= maxNameBytes || i >= maxNameChars {
			return sb.String()
		}
		sb.WriteRune(c)
		i++
	}

	return s
}
