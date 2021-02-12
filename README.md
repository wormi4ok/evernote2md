# Evernote to Markdown converter

[![Build Status](https://github.com/wormi4ok/evernote2md/workflows/Test/badge.svg)](https://github.com/wormi4ok/evernote2md/actions)
[![Docker Build Status](https://img.shields.io/docker/cloud/build/wormi4ok/evernote2md?label=Docker)](https://hub.docker.com/r/wormi4ok/evernote2md/)
[![Homebrew](https://repology.org/badge/version-for-repo/homebrew/evernote2md.svg?header=Homebrew)](https://repology.org/project/evernote2md/versions)
[![codecov](https://codecov.io/gh/wormi4ok/evernote2md/branch/master/graph/badge.svg)](https://codecov.io/gh/wormi4ok/evernote2md)
[![Go Report Card](https://goreportcard.com/badge/github.com/wormi4ok/evernote2md)](https://goreportcard.com/report/github.com/wormi4ok/evernote2md)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/wormi4ok/evernote2md)](https://pkg.go.dev/github.com/wormi4ok/evernote2md)

Evernote2md is a CLI tool to convert Evernote notes exported in *.enex format to a directory with markdown files.

Key features:

* Zero dependencies - download and run
* Creates one markdown file per note
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

`input` can be a file, a directory with exported files, or a glob pattern (like `exports/My*.enex`, for example).

If `outputDir` is not specified, `./notes` is used. Add optional `--folders` flag to put every note in a separate folder.

An option `--tagTemplate` allows to change the way tags are formatted. 
See [wiki article](https://github.com/wormi4ok/evernote2md/wiki/Custom-tag-template) for more information.

Flag `--help` shows all available options.

#### With Docker

```
docker run -t --rm -v "$PWD":/tmp -w /tmp wormi4ok/evernote2md:latest (flags) [input] [outputDir]
```

### How to export notes from Evernote

Here is a link to an article in Evernote Help Center:

[How to back up (export) and restore (import) notes and notebooks](https://help.evernote.com/hc/en-us/articles/209005557-Export-notes)

Newer versions of the Evernote App do not allow selecting more than 50 notes at a time.
Consider [exporting entire Notebook](https://help.evernote.com/hc/articles/360053159414) instead.

-----
Made with ❤ using IDE from JetBrains.

[![JetBrains](.github/powered_by.svg)](https://www.jetbrains.com/?from=evernote2md)
