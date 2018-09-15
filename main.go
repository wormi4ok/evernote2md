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
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/integrii/flaggy"
	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/file"
	"github.com/wormi4ok/evernote2md/internal"
	pb "gopkg.in/cheggaaa/pb.v1"
)

var version = "dev"

func main() {
	var input string
	var outputDir = "./notes"

	flaggy.SetName("evernote2md")
	flaggy.SetDescription(" Convert Evernote notes exported in *.enex format to markdown files")
	flaggy.SetVersion(version)

	flaggy.AddPositionalValue(&input, "input", 1, true, "Evernote export file")
	flaggy.String(&outputDir, "o", "outputDir", "Directory where markdown files will be created")

	flaggy.DefaultParser.ShowHelpOnUnexpected = false
	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/wormi4ok/evernote2md"

	flaggy.Parse()

	run(input, outputDir)
}

func run(input, output string) {
	var assetsDir = output + "/img"

	f, err := os.Open(input)
	failWhen(err)
	defer f.Close()

	export, err := enex.Decode(f)
	failWhen(err)

	err = os.MkdirAll(output, os.ModePerm)
	failWhen(err)

	progress := pb.StartNew(len(export.Notes))
	progress.Prefix("Notes:")
	for _, note := range export.Notes {
		md, err := internal.Converter{AssetsDir: assetsDir}.Convert(&note)
		failWhen(err)
		mdFile := filepath.FromSlash(output + "/" + file.BaseName(note.Title) + ".md")
		output, err := os.Create(mdFile)
		failWhen(err)
		_, err = io.Copy(output, bytes.NewReader(md.Content))
		failWhen(err)
		for _, res := range md.Media {
			err = file.Save(assetsDir, res.Name, bytes.NewReader(res.Content))
			if err != nil {
				log.Fatal(err)
			}
		}
		failWhen(err)
		progress.Increment()
	}
	progress.FinishPrint("Done!")
}

func failWhen(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
