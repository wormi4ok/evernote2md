# Evernote to Markdown converter

[![Build Status](https://github.com/wormi4ok/evernote2md/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/wormi4ok/evernote2md/actions/workflows/ci.yml)
[![Docker Image Size](https://img.shields.io/docker/image-size/wormi4ok/evernote2md)](https://hub.docker.com/r/wormi4ok/evernote2md/)
[![Homebrew](https://repology.org/badge/version-for-repo/homebrew/evernote2md.svg?header=Homebrew)](https://repology.org/project/evernote2md/versions)
[![Code Coverage](https://qlty.sh/gh/wormi4ok/projects/evernote2md/coverage.svg)](https://qlty.sh/gh/wormi4ok/projects/evernote2md)
[![Go Report Card](https://goreportcard.com/badge/github.com/wormi4ok/evernote2md)](https://goreportcard.com/report/github.com/wormi4ok/evernote2md)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/wormi4ok/evernote2md)](https://pkg.go.dev/github.com/wormi4ok/evernote2md)

Evernote2md is a CLI tool to convert Evernote notes exported in *.enex format to a directory with markdown files.

Key features:

* Zero dependencies - download and run
* Creates one markdown file per note ( with optional frontmatter e.g. for [Jekyll](https://jekyllrb.com/docs/front-matter/) )
* Converts attachments to files ( two directories will be created: `image` for images and `file` for other attachments
  e.g. pdf files )
* Retains correct links to attachments
* Inserts Evernote tags in notes as text entries with customizable formatting
* Shows highlighted Evernote text
* Sets file created and modified date equal to the note attributes

### Installation

Using [Homebrew](https://brew.sh) package manager:

```
brew install evernote2md
```

Manually:

[Download the latest release](https://github.com/wormi4ok/evernote2md/releases/latest) for your OS.

> ##### Note for macOS users!
> Please, check this [wiki](https://github.com/wormi4ok/evernote2md/wiki/macOS-FAQ) page if you have problems running the tool.

### How to use

```
evernote2md (flags) [input] [outputDir]
```

`input` can be a file, a directory with exported files, or a glob pattern (like `exports/My*.enex`, `exports/**/*.enex` for example).

If `outputDir` is not specified, `./notes` is used.

An option `--tagTemplate` allows to change the way tags are formatted.
See [wiki article](https://github.com/wormi4ok/evernote2md/wiki/Custom-tag-template) for more information.

Flag `--help` shows all available options.

To put exported notes in folders or structure in another custom way I recommend trying [mdmv](https://github.com/wormi4ok/mdmv) - Move Markdown files tool.

#### With Docker

```
docker run -t --rm -v "$PWD":/tmp -w /tmp wormi4ok/evernote2md:latest (flags) [input] [outputDir]
```

### How to export notes from Evernote

Here is a link to an article in Evernote Help Center:

[How to back up (export) and restore (import) notes and notebooks](https://help.evernote.com/hc/en-us/articles/209005557-Export-notes)

Newer versions of the Evernote App do not allow selecting more than 50 notes at a time.
Consider [exporting entire Notebook](https://github.com/wormi4ok/evernote2md/wiki/Export-a-notebook) instead.
