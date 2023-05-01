// Evernote2md is a cli tool to convert Evernote notes exported in *.enex format
// to a directory with markdown files.
//
// Usage:
//
//	evernote2md <file> [-o <outputDir>]
//
// If outputDir is not specified, current workdir is used.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/hako/durafmt"
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
	converter, err := internal.NewConverter(tagTemplate, addFrontMatter, !noHighlights)
	failWhen(err)

	setLogLevel(debug)
	run(files, output, newSpinner(debug), converter)
}

func newSpinner(disabled bool) *spinner.Spinner {
	sp := spinner.New(spinner.CharSets[43], 200*time.Millisecond)
	if disabled {
		sp.Disable()
	}
	return sp
}

func run(files []string, output *noteFilesDir, sp *spinner.Spinner, c *internal.Converter) {
	log.Printf("[DEBUG] Creating a directory: %s", output.Path())
	err := os.MkdirAll(output.Path(), os.ModePerm)
	failWhen(err)

	cnt := 0
	start := time.Now()
	sp.Start()

	for _, file := range files {
		fd, err := os.Open(file)
		failWhen(err)

		log.Printf("[DEBUG] Decoding file: %s", file)
		d, err := enex.NewStreamDecoder(fd)
		if progressError(err, file, "Failed to decode file") {
			continue
		}

		for {
			note := enex.Note{}
			if err := d.Next(&note); err != nil {
				if err != io.EOF {
					log.Printf("Failed to decode the next note: %s", err)
				}
				break
			}
			md, innerErr := c.Convert(&note)
			if progressError(innerErr, note.Title, "Failed to convert note") {
				continue
			}
			innerErr = output.SaveNote(note.Title, md)
			if progressError(innerErr, note.Title, "Failed to save note") {
				continue
			}
			cnt++
		}
		err = fd.Close()
		failWhen(err)
	}
	sp.FinalMSG = fmt.Sprintf("Done!\nConverted %d notes in %s\n", cnt, durafmt.ParseShort(time.Since(start)))
	sp.Stop()
}

func progressError(err error, name string, text string) bool {
	if err != nil {
		fmt.Print("\r") // Erase current spinner
		log.Printf(`[ERROR] %s "%s": %s`, text, name, err)
		return true
	}
	return false
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
		err = fmt.Errorf("no enex files found in the path: %s", input)
	}

	return files, err
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
