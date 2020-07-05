# Evernote to Markdown converter

[![Build Status](https://github.com/wormi4ok/evernote2md/workflows/Test/badge.svg)](https://github.com/wormi4ok/evernote2md/actions)
[![Docker Build Status](https://img.shields.io/docker/build/wormi4ok/evernote2md.svg)](https://hub.docker.com/r/wormi4ok/evernote2md/)
[![codecov](https://codecov.io/gh/wormi4ok/evernote2md/branch/master/graph/badge.svg)](https://codecov.io/gh/wormi4ok/evernote2md)
[![GoDoc](https://godoc.org/github.com/wormi4ok/evernote2md?status.svg)](http://godoc.org/github.com/wormi4ok/evernote2md)
[![Go Report Card](https://goreportcard.com/badge/github.com/wormi4ok/evernote2md)](https://goreportcard.com/report/github.com/wormi4ok/evernote2md)

Evernote2md is a CLI tool to convert Evernote notes exported in *.enex format to a directory with markdown files.

Key features:

* Zero dependencies - download and run 
* Creates one markdown file per note
* Converts attachments to files ( two directories will be created: `image` for images and `file` for other attachments e.g. pdf files )
* Retains correct links to attachments
* Inserts evernote tags in notes as text entries

### Installation

[Download a release](https://github.com/wormi4ok/evernote2md/releases/latest) for your OS.

### How to use

#### With binary

```
evernote2md [input] [outputDir]
```

If outputDir is not specified, `./notes` is used.

#### With docker

```
docker run -t --rm -v "$PWD":/tmp -w /tmp wormi4ok/evernote2md:latest [input] [outputDir]
```

### How to export notes from Evernote

Here is a link to an article in Evernote Help Center:

[How to back up (export) and restore (import) notes and notebooks](https://help.evernote.com/hc/en-us/articles/209005557-How-to-back-up-export-and-restore-import-notes-and-notebooks)

