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
	"github.com/hashicorp/logutils"
	"github.com/integrii/flaggy"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/internal"
)

var version = "dev"

func init() {
	flaggy.SetName("evernote2md")
	flaggy.SetDescription(" Convert Evernote notes exported in *.enex format to markdown files")
	flaggy.SetVersion(version)

	flaggy.DefaultParser.ShowHelpOnUnexpected = false
	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/wormi4ok/evernote2md"
}

func main() {
	var input string
	var outputDir = filepath.FromSlash("./notes")
	var outputOverride string
	var folders, noHighlights, resetTimestamps, debug bool

	flaggy.AddPositionalValue(&input, "input", 1, true, "Evernote export file")
	flaggy.AddPositionalValue(&outputDir, "output", 2, false, "Output directory")
	flaggy.String(&outputOverride, "o", "outputDir", "Directory where markdown files will be created")

	flaggy.Bool(&folders, "", "folders", "Put every note in a separate folder")
	flaggy.Bool(&noHighlights, "", "noHighlights", "Disable converting Evernote highlights to inline HTML tags")
	flaggy.Bool(&resetTimestamps, "", "resetTimestamps", "Create files ignoring timestamps in the note attributes")
	flaggy.Bool(&debug, "v", "debug", "Show debug output")

	flaggy.Parse()

	if len(outputOverride) > 0 {
		outputDir = outputOverride
	}

	output := newNoteFilesDir(outputDir, folders, !resetTimestamps)
	converter := internal.Converter{EnableHighlights: !noHighlights}

	setLogLevel(debug)

	run(input, output, newProgressBar(debug), converter)
}

func run(input string, output *noteFilesDir, progress *pb.ProgressBar, c internal.Converter) {
	i, err := os.Open(input)
	failWhen(err)

	export, err := enex.Decode(i)
	failWhen(err)

	err = i.Close()
	failWhen(err)

	err = os.MkdirAll(output.Path(), os.ModePerm)
	failWhen(err)

	progress.SetTotal(int64(len(export.Notes)))
	progress.Start()
	n := export.Notes
	for i := range n {
		md, err := c.Convert(&n[i])
		failWhen(err)
		err = output.SaveNote(n[i].Title, md)
		failWhen(err)

		progress.Increment()
	}
	progress.Finish()
	fmt.Println("Done!")
}

const progressBarTmpl = `Notes: {{counters .}} {{bar . "[" "=" ">" "_" "]" }} {{percent .}} {{etime .}}`

func newProgressBar(debug bool) *pb.ProgressBar {
	progress := new(pb.ProgressBar)
	progress.SetTemplateString(progressBarTmpl)
	if debug {
		progress.SetWriter(new(bytes.Buffer))
	}
	return progress
}

func setLogLevel(debug bool) {
	var logLevel logutils.LogLevel = "WARN"

	if debug {
		logLevel = "DEBUG"
	}

	log.SetOutput(&logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logLevel,
		Writer:   os.Stderr,
	})
}

func failWhen(err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("[ERROR] %w", err))
	}
}
