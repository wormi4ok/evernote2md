package file

import (
	"io"
	"os"
	"regexp"
	"strings"
)

const maxPathLength = 200

var (
	baseNameSeparators = regexp.MustCompile(`[./]`)

	blacklist = regexp.MustCompile(`[ |"'<>&_=+:?]`)

	dashes = regexp.MustCompile(`[\-]+`)
)

// Save a new file in a given dir with the following content
// If directory doesn't exist it will create it
func Save(dir, name string, content io.Reader) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}
	if len(name) == 0 {
		return nil
	}

	output, err := os.Create(dir + "/" + name)
	if err != nil {
		return err
	}

	_, err = io.Copy(output, content)
	return err
}

// BaseName normalizes a given string to use it as a safe filename
func BaseName(s string) string {
	// Replace certain joining characters with a dash
	s = baseNameSeparators.ReplaceAllString(s, "-")

	// Remove any trailing space to avoid ending on -
	s = strings.Trim(s, " ")

	// Replace certain joining characters with a dash
	s = blacklist.ReplaceAllString(s, "_")

	// Remove any multiple dashes caused by replacements above
	s = dashes.ReplaceAllString(s, "-")

	if len(s) > maxPathLength {
		s = s[:maxPathLength]
	}

	return s
}
