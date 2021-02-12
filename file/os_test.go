package file_test

import (
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"github.com/wormi4ok/evernote2md/file"
)

func TestChangeFileTimes(t *testing.T) {
	now := time.Now().Add(-1 * time.Minute)
	dir := t.TempDir()
	if _, err := os.Create(path.Join(dir, "test.md")); err != nil {
		t.Fatal(err)
	}

	if err := file.ChangeFileTimes(dir, "test.md", now, now); err != nil {
		t.Fatal(err)
	}

	stat, err := os.Stat(path.Join(dir, "test.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := now.Truncate(time.Minute)
	if got := stat.ModTime(); got != want {
		t.Errorf("Modification time mismatch. want = %v ,got = %v", want, got)
	}
}

func TestChangeFileTimes_NoFile(t *testing.T) {
	err := file.ChangeFileTimes(t.TempDir(), "not_exist.md", time.Now(), time.Now())
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}
