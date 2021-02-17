package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

func TestNoteFilesDir_SaveNote(t *testing.T) {
	log.SetOutput(io.Discard)
	tmpDir := t.TempDir()
	wantDate := time.Unix(1608463260, 0)

	d := newNoteFilesDir(tmpDir, false, true)
	md := fakeNote(wantDate)
	err := d.SaveNote("test_note", md)
	if err != nil {
		t.Errorf("SaveNote returned error: %s", err.Error())
	}

	stat := shouldExist(t, tmpDir+"/test_note.md")
	if stat.ModTime() != wantDate {
		t.Errorf("Timestamp doesn't match, got =  %s, want = %s", stat.ModTime().String(), wantDate.String())
	}
	shouldExist(t, tmpDir+"/image/test.jpg")
}

// Test non-default flag states
func TestNoteFilesDir_Flags(t *testing.T) {
	tmpDir := t.TempDir()
	fixedDate := time.Unix(1608463260, 0)
	d := newNoteFilesDir(tmpDir, true, false)

	md := fakeNote(fixedDate)
	err := d.SaveNote("test_note", md)
	if err != nil {
		t.Errorf("SaveNote returned error: %s", err.Error())
	}

	stat := shouldExist(t, tmpDir+"/test_note/README.md")
	if stat.ModTime() == fixedDate {
		t.Errorf("Timestamp matches the fixed date =  %s, want = %s", stat.ModTime().String(), time.Now().String())
	}
}

// Test that notes don't overwrite each other
func TestNoteFilesDir_UniqueNames(t *testing.T) {
	tmpDir := t.TempDir()
	d := newNoteFilesDir(tmpDir, false, false)

	md := fakeNote(time.Now())
	err := d.SaveNote("test_note", md)
	if err != nil {
		t.Errorf("SaveNote returned error: %s", err.Error())
	}

	err = d.SaveNote("test_note", md)
	if err != nil {
		t.Errorf("SaveNote returned error: %s", err.Error())
	}

	shouldExist(t, tmpDir+"/test_note-1.md")
}

// Test that notes with identical names but different casing don't override each other
func TestNoteFilesDir_UniqueNames_CaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()
	d := newNoteFilesDir(tmpDir, false, false)

	md := fakeNote(time.Now())
	err := d.SaveNote("TEST_note", md)
	if err != nil {
		t.Errorf("SaveNote returned error: %s", err.Error())
	}

	err = d.SaveNote("test_note", md)
	if err != nil {
		t.Errorf("SaveNote returned error: %s", err.Error())
	}

	shouldExist(t, tmpDir+"/test_note-1.md")
}

func fakeNote(wantDate time.Time) *markdown.Note {
	return &markdown.Note{
		Content: []byte(`12345`),
		Media: map[string]markdown.Resource{
			"123": {
				Name:    "test.jpg",
				Type:    "image",
				Content: []byte(`fakeContent`),
			},
		},
		CTime: wantDate,
		MTime: wantDate,
	}
}

func shouldExist(t *testing.T, path string) os.FileInfo {
	wantPath := filepath.FromSlash(path)
	stat, err := os.Stat(wantPath)
	if err != nil && os.IsNotExist(err) {
		t.Errorf("%s was not created", wantPath)
	}

	return stat
}
