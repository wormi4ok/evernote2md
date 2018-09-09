# Evernote to Markdown converter

[![Build Status](https://travis-ci.org/wormi4ok/evernote2md.svg?branch=master)](https://travis-ci.org/wormi4ok/evernote2md)
[![codecov](https://codecov.io/gh/wormi4ok/evernote2md/branch/master/graph/badge.svg)](https://codecov.io/gh/wormi4ok/evernote2md)
[![GoDoc](https://godoc.org/github.com/wormi4ok/evernote2md?status.svg)](http://godoc.org/github.com/wormi4ok/evernote2md)
[![Go Report Card](https://goreportcard.com/badge/github.com/wormi4ok/evernote2md)](https://goreportcard.com/report/github.com/wormi4ok/evernote2md)

Evernote2md is a CLI tool to convert Evernote notes exported in *.enex format to a directory with markdown files.

### Installation

Download a release for your OS and architecture. 

### How to use

```
evernote2md export.enex -o ./notes
```

If outputDir is not specified, current `workdir` is used.

### How to export notes from Evernote

Here is a link to an article in Evernote Help Center:

[How to back up (export) and restore (import) notes and notebooks](https://help.evernote.com/hc/en-us/articles/209005557-How-to-back-up-export-and-restore-import-notes-and-notebooks)

