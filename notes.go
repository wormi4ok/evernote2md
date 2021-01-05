package main

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/wormi4ok/evernote2md/encoding/markdown"
	"github.com/wormi4ok/evernote2md/file"
)

// noteFilesDir saves markdown notes in a directory on the filesystem
type noteFilesDir struct {
	path string

	// flags modifying the logic for saving notes
	flagFolders    bool
	flagTimestamps bool

	// A map to keep track of what notes are already created
	names map[string]int
}

func newNoteFilesDir(output string, folders, timestamps bool) *noteFilesDir {
	return &noteFilesDir{
		path:           output,
		flagFolders:    folders,
		flagTimestamps: timestamps,
		names:          map[string]int{},
	}
}

// SaveNote along with media resources
func (d *noteFilesDir) SaveNote(title string, md *markdown.Note) error {
	path := d.path
	if d.flagFolders {
		path = filepath.FromSlash(d.path + "/" + d.uniqueName(title))
		title = "README.md"
	} else {
		title = d.uniqueName(title) + ".md"
	}

	log.Printf("[DEBUG] Saving file %s/%s", path, title)
	if err := file.Save(path, title, bytes.NewReader(md.Content)); err != nil {
		return fmt.Errorf("save file %s: %w", path+"/"+title, err)
	}

	if d.flagTimestamps {
		if err := file.ChangeFileTimes(path, title, md.CTime, md.MTime); err != nil {
			// Continue processing on error
			log.Printf("[WARN] Error updating file times for a file: %s", title)
		}
	}

	for _, res := range md.Media {
		mediaPath := filepath.FromSlash(path + "/" + string(res.Type))
		log.Printf("[DEBUG] Saving attachment %s/%s", mediaPath, res.Name)
		if err := file.Save(mediaPath, res.Name, bytes.NewReader(res.Content)); err != nil {
			return fmt.Errorf("save resource %s: %w", mediaPath+"/"+res.Name, err)
		}
	}

	return nil
}

func (d *noteFilesDir) Path() string {
	return d.path
}

// uniqueName returns a unique note name
func (d *noteFilesDir) uniqueName(title string) string {
	name := file.BaseName(title)
	index := strings.ToLower(name)

	if k, exist := d.names[index]; exist {
		d.names[index] = k + 1
		name = fmt.Sprintf("%s-%d", name, k)
	} else {
		d.names[index] = 1
	}

	return name
}
