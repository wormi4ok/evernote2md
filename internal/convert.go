package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	TagFormat        string
	EnableHighlights bool
}

// Convert Evernote file to markdown
func (c *Converter) Convert(note *enex.Note) (*markdown.Note, error) {
	var md markdown.Note
	md.Media = map[string]markdown.Resource{}

	if err := mapResources(note, md); err != nil {
		return nil, err
	}

	html, err := normalizeHTML(note.Content, NewReplacerMedia(md.Media), &Code{}, &ExtraDiv{}, &TextFormatter{})
	if err != nil {
		return nil, err
	}

	content := prependTags(note.Tags, string(html))
	content = prependTitle(note.Title, content)

	var b bytes.Buffer
	err = markdown.Convert(&b, strings.NewReader(content), c.EnableHighlights)
	if err != nil {
		return nil, err
	}

	md.Content = regexp.MustCompile(`\n{3,}`).ReplaceAllLiteral(b.Bytes(), []byte("\n\n"))
	md.Content = append(bytes.TrimRight(md.Content, "\n"), '\n')

	md.CTime = convertEvernoteDate(note.Created)
	md.MTime = convertEvernoteDate(note.Updated)
	return &md, nil
}

func mapResources(note *enex.Note, md markdown.Note) error {
	names := map[string]int{}
	r := note.Resources
	for i := range r {
		p, err := ioutil.ReadAll(decoder(r[i].Data))
		if err != nil {
			return err
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
	return nil
}

func prependTags(tags []string, content string) string {
	var tt []string
	for _, t := range tags {
		tt = append(tt, fmt.Sprintf("<code>%s</code>", t))
	}
	return strings.Join(tt, " ") + "<br>" + content
}

func prependTitle(title, content string) string {
	return fmt.Sprintf("<h1>%s</h1>", title) + content
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
