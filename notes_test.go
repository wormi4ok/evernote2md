package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

func TestNoteFilesDir_SaveNote(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	tmpDir := t.TempDir()
	now := time.Now()

	d := newNoteFilesDir(tmpDir, false, true)
	md := &markdown.Note{
		Content: []byte(`12345`),
		Media:   nil,
		CTime:   now,
		MTime:   now,
	}
	err := d.SaveNote("test_note", md)
	if err != nil {
		t.Errorf("SaveNote returned error: %s", err.Error())
	}

	wantPath := filepath.FromSlash(tmpDir + "/test_note.md")
	if _, err := os.Stat(wantPath); os.IsNotExist(err) {
		t.Errorf("%s was not created", wantPath)
	}
}
