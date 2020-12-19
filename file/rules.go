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
	filePathToModify := filepath.FromSlash(dir + "/" + name)
	if _, err := os.Stat(filePathToModify); os.IsNotExist(err) {
		return fmt.Errorf("change file timestamps %s: %w", filePathToModify, err)
	}

	touchCmd, err := exec.LookPath("touch")
	if err != nil {
		return os.Chtimes(filePathToModify, mtime, mtime)
	}
	changeMtime := exec.Command(touchCmd, "-mt", mtime.Format(touchTimeFormat), filePathToModify)
	if err := changeMtime.Run(); err != nil {
		return err
	}
	changeCTime := exec.Command(touchCmd, "-t", ctime.Format(touchTimeFormat), filePathToModify)
	if err := changeCTime.Run(); err != nil {
		return err
	}
	return nil
}
