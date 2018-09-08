// DOcumentation for my program
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
)

var version = "dev"

func main() {
	var input = "./data/Evernote.enex"
	var outputDir = "./data/notes"

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
	}
}

func failWhen(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
