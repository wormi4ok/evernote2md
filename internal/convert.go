package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

// Converter holds configuration options to control conversion
type Converter struct {
	TagFormat        string
	EnableHighlights bool

	// err holds an error during conversion
	// Every conversion step should check this field and skip execution if it is not empty
	err error
}

// Convert an Evernote file to markdown
func (c *Converter) Convert(note *enex.Note) (*markdown.Note, error) {
	md := new(markdown.Note)
	md.Media = map[string]markdown.Resource{}

	c.mapResources(note, md)
	c.normalizeHTML(note, md, NewReplacerMedia(md.Media), &Code{}, &ExtraDiv{}, &TextFormatter{})
	c.prependTags(note, md)
	c.prependTitle(note, md)
	c.toMarkdown(note, md)
	c.trimSpaces(note, md)
	c.addDates(note, md)

	return md, c.err
}

func (c *Converter) mapResources(note *enex.Note, md *markdown.Note) {
	names := map[string]int{}
	r := note.Resources
	for i := range r {
		p, err := ioutil.ReadAll(decoder(r[i].Data))
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

func (c *Converter) prependTags(note *enex.Note, _ *markdown.Note) {
	if c.err != nil {
		return
	}

	var tt [][]byte
	for _, t := range note.Tags {
		tt = append(tt, []byte(fmt.Sprintf("<code>%s</code>", t)))
	}
	note.Content = append([]byte("<br>"), note.Content...)
	note.Content = append(bytes.Join(tt, []byte(" ")), note.Content...)
}

func (c *Converter) prependTitle(note *enex.Note, _ *markdown.Note) {
	if c.err != nil {
		return
	}

	note.Content = append([]byte(fmt.Sprintf("<h1>%s</h1>", note.Title)), note.Content...)
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
