package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

const DefaultFrontMatterTemplate = `---
date: '{{.CTime}}'
updated_at: '{{.MTime}}'
title: {{.Title}}
tags: [ {{.TagList}} ]
{{- if .Attributes.SourceUrl}}
url: {{.Attributes.SourceUrl -}}
{{end -}}
{{- if .Attributes.Latitude}}
latitude: {{.Attributes.Latitude -}}
{{end -}}
{{- if .Attributes.Longitude}}
longitude: {{.Attributes.Longitude -}}
{{end -}}
{{- if .Attributes.Altitude}}
altitude: {{.Attributes.Altitude -}}
{{end -}}
{{- if .Attributes.Source}}
source: {{.Attributes.Source -}}
{{end}}

---

`

// Converter holds configuration options to control conversion
type Converter struct {
	TagTemplate         string
	EnableHighlights    bool
	EnableFrontMatter   bool
	FrontMatterTemplate string

	// err holds an error during conversion
	// Every conversion step should check this field and skip execution if it is not empty
	err error
}

// NewConverter creates a Converter with valid tagTemplate
func NewConverter(tagTemplate string, enableFrontMatter bool, frontMatterTemplate string, enableHighlights bool) (*Converter, error) {
	if tagTemplate == "" {
		tagTemplate = DefaultTagTemplate
	}
	if frontMatterTemplate == "" {
		frontMatterTemplate = DefaultFrontMatterTemplate
	}

	if strings.Count(tagTemplate, tagToken) != 1 {
		return nil, errors.New("tag format should contain exactly one {{tag}} template variable")
	}

	return &Converter{TagTemplate: tagTemplate, EnableFrontMatter: enableFrontMatter, FrontMatterTemplate: frontMatterTemplate, EnableHighlights: enableHighlights}, nil
}

// Convert an Evernote file to markdown
func (c *Converter) Convert(note *enex.Note) (*markdown.Note, error) {
	md := new(markdown.Note)
	md.Media = map[string]markdown.Resource{}

	c.mapResources(note, md)
	c.normalizeHTML(note, md, NewReplacerMedia(md.Media), &Code{}, &ExtraDiv{}, &TextFormatter{}, &EmptyAnchor{})
	c.toMarkdown(note, md)
	c.prependTags(note, md)
	c.prependTitle(note, md)
	c.trimSpaces(note, md)
	c.addDates(note, md)
	if c.EnableFrontMatter {
		c.addFrontMatter(note, md)
	}

	return md, c.err
}

func (c *Converter) mapResources(note *enex.Note, md *markdown.Note) {
	names := map[string]int{}
	r := note.Resources
	for i := range r {
		p, err := io.ReadAll(decoder(r[i].Data))
		if c.err = err; err != nil {
			return
		}

		rType := markdown.File
		if isImage(r[i].Mime) {
			rType = markdown.Image
		}
		name, ext := name(r[i])

		// Ensure the name is unique
		if cnt, exist := names[name+ext]; exist {
			names[name+ext] = cnt + 1
			name = fmt.Sprintf("%s-%d", name, cnt)
		} else {
			names[name+ext] = 1
		}

		mdr := markdown.Resource{
			Name:    name + ext,
			Type:    rType,
			Content: p,
		}

		if r[i].ID != "" {
			md.Media[r[i].ID] = mdr
		} else {
			md.Media[strconv.Itoa(i)] = mdr
		}
	}
}

func (c *Converter) prependTitle(note *enex.Note, md *markdown.Note) {
	if c.err != nil {
		return
	}

	md.Content = append([]byte(fmt.Sprintf("# %s\n\n", note.Title)), md.Content...)
}

func (c *Converter) toMarkdown(note *enex.Note, md *markdown.Note) {
	if c.err != nil {
		return
	}
	var b bytes.Buffer
	err := markdown.Convert(&b, bytes.NewReader(note.Content), c.EnableHighlights)
	if c.err = err; err != nil {
		return
	}

	md.Content = b.Bytes()
}

func (c *Converter) trimSpaces(_ *enex.Note, md *markdown.Note) {
	if c.err != nil {
		return
	}

	md.Content = regexp.MustCompile(`\n{3,}`).ReplaceAllLiteral(md.Content, []byte("\n\n"))
	md.Content = append(bytes.TrimRight(md.Content, "\n"), '\n')
}

func (c *Converter) addDates(note *enex.Note, md *markdown.Note) {
	if c.err != nil {
		return
	}

	md.CTime = convertEvernoteDate(note.Created)
	md.MTime = convertEvernoteDate(note.Updated)
}

const dateFrontMatterFormat = "2006-01-02 15:04:05 -0700"

func (c *Converter) addFrontMatter(note *enex.Note, md *markdown.Note) {
	data := struct {
		CTime      string
		MTime      string
		Title      string
		Attributes enex.NoteAttributes
		TagList    string
	}{
		md.CTime.Format(dateFrontMatterFormat),
		md.MTime.Format(dateFrontMatterFormat),
		note.Title,
		note.Attributes,
		c.tagList(note, "'{{tag}}'", ", ", false),
	}
	tmpl, err := template.New("frontMatter").Parse(c.FrontMatterTemplate)
	if err != nil {
		panic(err)
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		panic(err)
	}
	md.Content = append(b.Bytes(), md.Content...)
}

const evernoteDateFormat = "20060102T150405Z"

// 20180109T173725Z -> 2018-01-09T17:37:25Z
func convertEvernoteDate(evernoteDate string) time.Time {
	converted, err := time.Parse(evernoteDateFormat, evernoteDate)
	if err != nil {
		log.Printf("[DEBUG] Could not convert time /%s: %s, using today instead", evernoteDate, err.Error())
		converted = time.Now()
	}

	return converted
}
