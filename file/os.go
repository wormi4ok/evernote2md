// +build !windows

package file

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

// Max path length in bytes, determined empirically.
// P.S. Don't trust Apple documentation
const maxNameBytes int = 704

// Semicolon is not allowed in MacOS and spaces is just my personal preference
var illegalChars = regexp.MustCompile(`[\s:]`)

// ChangeFileTimes matches the file times with the Evernote metadata
//
// Uses touch if available to change both creation  and modification date
// Otherwise it falls back to os.Chtimes to change only the modification date
func ChangeFileTimes(dir, name string, ctime, mtime time.Time) error {
	path := filepath.FromSlash(dir + "/" + name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("change file timestamps %s: %w", path, err)
	}

	touchCmd, err := exec.LookPath("touch")
	if err != nil {
		return os.Chtimes(path, mtime, mtime)
	}

	// On macOS, first touch set both creation date and modification date
	// On Linux, this touch will be ignored by second touch. There is no easy way to setting creation date
	changeCTime := exec.Command(touchCmd, "-mt", ctime.Format(touchTimeFormat), path)
	if err := changeCTime.Run(); err != nil {
		return err
	}
	// On macOS, second touch updates the modification date and the creation date is preserved
	changeMtime := exec.Command(touchCmd, "-mt", mtime.Format(touchTimeFormat), path)
	if err := changeMtime.Run(); err != nil {
		return err
	}

	return nil
}
