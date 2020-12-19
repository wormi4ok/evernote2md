package file

import (
	"os"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
)

// Max path length is 255 - 9 bytes for extension (.md) in multibyte encoding
const maxNameBytes int = 246

// Additional rule for Windows
var illegalChars = regexp.MustCompile(`[\s\\|"'<>&_=+:?*]`)

// ChangeFileTimes uses SetFileTime syscall in Windows implementation
// which supports updating both creation and modification dates
func ChangeFileTimes(dir, name string, ctime, mtime time.Time) error {
	path := filepath.FromSlash(dir + "/" + name)
	ctimeSpec := syscall.NsecToFiletime(ctime.UnixNano())
	mtimeSpec := syscall.NsecToFiletime(mtime.UnixNano())

	fd, err := syscall.Open(path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	return syscall.SetFileTime(fd, &ctimeSpec, &mtimeSpec, &mtimeSpec)
}
