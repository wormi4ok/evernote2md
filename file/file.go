package file

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/exp/utf8string"
)

// Mon Jan 2 15:04:05 -0700 MST 2006 represented as yyyyMMddhhmm
const TOUCH_TIME_FORMAT = "200601021504"
const maxPathLength = 200

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

// Match the file times with the evernote metadata
// uses touch if available to change both cdate and mdate
// uses os.Chtimes to change only the mdate
func ChangeFileTimes(dir, name string, ctime, mtime time.Time) {
	var err error
	filePathToModify := filepath.FromSlash((dir + "/" + name))
	_, err = os.Stat(filePathToModify)
	if os.IsNotExist(err) {
		// file doesn't exist
		log.Printf("Tried to change file creation times, file not found %s", filePathToModify)
		return
	}
	_, err = exec.LookPath("touch")
	if err != nil {
		os.Chtimes(filePathToModify, mtime, mtime)
		return
	}
	changeMtime := exec.Command("touch", "-mt", mtime.Format(TOUCH_TIME_FORMAT), filePathToModify)
	if err := changeMtime.Run(); err != nil {
		log.Panic(err)
	}
	changeCTime := exec.Command("touch", "-t", ctime.Format(TOUCH_TIME_FORMAT), filePathToModify)
	if err := changeCTime.Run(); err != nil {
		log.Panic(err)
	}
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

	// Trim filename to max allowed number of characters
	utf8 := utf8string.NewString(s)
	if utf8.RuneCount() > maxPathLength {
		s = utf8.Slice(0, maxPathLength)
	}

	return s
}
