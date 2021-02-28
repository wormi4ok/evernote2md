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
	"time"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

// Converter holds configuration options to control conversion
type Converter struct {
	TagTemplate      string
	EnableHighlights bool

	// err holds an error during conversion
	// Every conversion step should check this field and skip execution if it is not empty
	err error
}

// NewConverter creates a Converter with valid tagTemplate
func NewConverter(tagTemplate string, enableHighlights bool) (*Converter, error) {
	if tagTemplate == "" {
		tagTemplate = DefaultTagTemplate
	}

	if strings.Count(tagTemplate, tagToken) != 1 {
		return nil, errors.New("tag format should contain exactly one {{tag}} template variable")
	}

	return &Converter{TagTemplate: tagTemplate, EnableHighlights: enableHighlights}, nil
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
	c.addFrontMatter(note, md)

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
	var frontMatter = fmt.Sprintf("date: %s\nupdated_at: %s\ntitle: %s\ntags: %s\n",
		md.CTime.Format(dateFrontMatterFormat),
		md.MTime.Format(dateFrontMatterFormat),
		note.Title, c.tagList(note, md))
	if len(note.Attributes.SourceUrl) > 0 {
		frontMatter += fmt.Sprintf("url: %s\n", note.Attributes.SourceUrl)
	}
	if len(note.Attributes.Latitude) > 0 {
		frontMatter += fmt.Sprintf("lattitude: %s\n", note.Attributes.Latitude)
	}
	if len(note.Attributes.Longitude) > 0 {
		frontMatter += fmt.Sprintf("longitude: %s\n", note.Attributes.Longitude)
	}
	if len(note.Attributes.Altitude) > 0 {
		frontMatter += fmt.Sprintf("altitude: %s\n", note.Attributes.Altitude)
	}
	if len(note.Attributes.Source) > 0 {
		frontMatter += fmt.Sprintf("source: %s\n", note.Attributes.Source)
	}
	var frontMatterBytes = append([]byte(frontMatter), []byte("\n---\n\n")...)
	md.Content = append(frontMatterBytes, md.Content...)
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
