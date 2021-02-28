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
	var input, outputOverride string
	var outputDir = filepath.FromSlash("./notes")
	var tagTemplate = internal.DefaultTagTemplate
	var folders, noHighlights, resetTimestamps, addFrontMatter, debug bool

	flaggy.AddPositionalValue(&input, "input", 1, true, "Evernote export file, directory or a glob pattern")
	flaggy.AddPositionalValue(&outputDir, "output", 2, false, "Output directory")

	flaggy.String(&tagTemplate, "t", "tagTemplate", "Define how Evernote tags are formatted")
	flaggy.String(&outputOverride, "o", "outputDir", "Override the directory where markdown files will be created")

	flaggy.Bool(&folders, "", "folders", "Put every note in a separate folder")
	flaggy.Bool(&noHighlights, "", "noHighlights", "Disable converting Evernote highlights to inline HTML tags")
	flaggy.Bool(&resetTimestamps, "", "resetTimestamps", "Create files ignoring timestamps in the note attributes")
	flaggy.Bool(&addFrontMatter, "", "addFrontMatter", "Prepend FrontMatter to markdown files")
	flaggy.Bool(&debug, "v", "debug", "Show debug output")

	flaggy.Parse()

	if len(outputOverride) > 0 {
		outputDir = outputOverride
	}

	files, err := matchInput(input)
	failWhen(err)
	output := newNoteFilesDir(outputDir, folders, !resetTimestamps)
	converter, err := internal.NewConverter(tagTemplate, addFrontMatter, "", !noHighlights)
	failWhen(err)

	setLogLevel(debug)

	run(files, output, newProgressBar(debug), converter)
}

func run(files []string, output *noteFilesDir, progress *pb.ProgressBar, c *internal.Converter) {
	export := decodeFiles(files)

	log.Printf("[DEBUG] Creating a directory: %s", output.Path())
	err := os.MkdirAll(output.Path(), os.ModePerm)
	failWhen(err)

	progress.SetTotal(int64(len(export.Notes)))
	progress.Start()
	n := export.Notes
	for i := range n {
		log.Printf("[DEBUG] Converting a note: %s", n[i].Title)
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

// matchInput finds all files matching input pattern
// If input is a path to a directory, it will search for *.enex files inside the directory
func matchInput(input string) ([]string, error) {
	var (
		files []string
		err   error
	)
	if input == "" {
		input, err = os.Getwd()
	} else {
		input, err = filepath.Abs(input)
	}
	if err != nil {
		return nil, err
	}

	// If input is a directory, find all *.enex files and return
	if info, err := os.Stat(input); err == nil && info.IsDir() {
		files, err := filepath.Glob(filepath.FromSlash(input + "/*.enex"))
		if files != nil {
			return files, err
		}
	}

	// User glob patterns may include directories that we filter out
	matches, err := filepath.Glob(input)
	for _, match := range matches {
		if info, err := os.Stat(match); err == nil && !info.IsDir() {
			files = append(files, match)
		}
	}
	if files == nil {
		err = fmt.Errorf("[ERROR] No enex files found in the path: %s", input)
	}

	return files, err
}

// decodeFiles creates a single Evernote export from multiple input files
func decodeFiles(files []string) *enex.Export {
	export := new(enex.Export)
	for _, file := range files {
		fd, err := os.Open(file)
		failWhen(err)

		log.Printf("[DEBUG] Decoding a file: %s", file)
		e, err := enex.Decode(fd)
		failWhen(err)

		err = fd.Close()
		failWhen(err)
		export.Notes = append(export.Notes, e.Notes...)
	}

	return export
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
