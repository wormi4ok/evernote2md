package file

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
)

// Max path length according to fixLongPath function is 248 - 3 bytes for extension (.md)
const maxNameBytes int = 245

// Additional rule for Windows
var illegalChars = regexp.MustCompile(`[\s\\|"'<>&_=+:?*]`)

// ChangeFileTimes uses SetFileTime syscall in Windows implementation
// which supports updating both creation and modification dates
func ChangeFileTimes(dir, name string, ctime, mtime time.Time) error {
	path := filepath.FromSlash(dir + "/" + name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("change file timestamps %s: %w", path, err)
	}
	ctimeSpec := syscall.NsecToFiletime(ctime.UnixNano())
	mtimeSpec := syscall.NsecToFiletime(mtime.UnixNano())

	fd, err := syscall.Open(path, os.O_RDWR, 644)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	return syscall.SetFileTime(fd, &ctimeSpec, &mtimeSpec, &mtimeSpec)
}
