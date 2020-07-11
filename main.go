// Evernote2md is a cli tool to convert Evernote notes exported in *.enex format
// to a directory with markdown files.
//
// Usage:
//   evernote2md <file> [-o <outputDir>]
//
// If outputDir is not specified, current workdir is used.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"github.com/integrii/flaggy"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
	"github.com/wormi4ok/evernote2md/file"
	"github.com/wormi4ok/evernote2md/internal"
)

var version = "dev"

func main() {
	var input string
	var outputDir = filepath.FromSlash("./notes")
	var outputOverride string
	var folders bool

	flaggy.SetName("evernote2md")
	flaggy.SetDescription(" Convert Evernote notes exported in *.enex format to markdown files")
	flaggy.SetVersion(version)

	flaggy.AddPositionalValue(&input, "input", 1, true, "Evernote export file")
	flaggy.AddPositionalValue(&outputDir, "output", 2, false, "Output directory")
	flaggy.String(&outputOverride, "o", "outputDir", "Directory where markdown files will be created")

	flaggy.Bool(&folders, "", "folders", "Put every note in a separate folder")

	flaggy.DefaultParser.ShowHelpOnUnexpected = false
	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/wormi4ok/evernote2md"

	flaggy.Parse()

	if len(outputOverride) > 0 {
		outputDir = outputOverride
	}

	run(input, outputDir, folders)
}

const progressBarTmpl = `Notes: {{counters .}} {{bar . "[" "=" ">" "_" "]" }} {{percent .}} {{etime .}}`

// A map to keep track of what notes are already created
var notes = map[string]int{}

func run(input, output string, folders bool) {
	i, err := os.Open(input)
	failWhen(err)

	export, err := enex.Decode(i)
	failWhen(err)

	err = i.Close()
	failWhen(err)

	err = os.MkdirAll(output, os.ModePerm)
	failWhen(err)

	progress := pb.StartNew(len(export.Notes))
	progress.SetTemplateString(progressBarTmpl)

	n := export.Notes
	for i := range n {
		md, err := internal.Convert(&n[i])
		failWhen(err)
		if folders {
			path := output + "/" + uniqueName(n[i].Title)
			err = saveNote(path, "README.md", md)
		} else {
			err = saveNote(output, uniqueName(n[i].Title)+".md", md)
		}
		failWhen(err)

		progress.Increment()
	}
	progress.Finish()
	fmt.Println("Done!")
}

// saveNote along with media resources
func saveNote(path string, title string, md *markdown.Note) error {
	if err := file.Save(path, title, bytes.NewReader(md.Content)); err != nil {
		return fmt.Errorf("save file %s: %w", path+"/"+title, err)
	}
	for _, res := range md.Media {
		mediaPath := path + "/" + string(res.Type)
		if err := file.Save(mediaPath, res.Name, bytes.NewReader(res.Content)); err != nil {
			return fmt.Errorf("save resource %s: %w", mediaPath+"/"+res.Name, err)
		}
	}

	return nil
}

// uniqueName returns a unique note name
func uniqueName(title string) string {
	name := file.BaseName(title)
	if k, exist := notes[name]; exist {
		notes[name] = k + 1
		name = fmt.Sprintf("%s-%d", file.BaseName(title), k)
	} else {
		notes[name] = 1
	}

	return name
}

func failWhen(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
