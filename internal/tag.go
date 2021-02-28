package internal

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

// DefaultTagTemplate format if none specified
const DefaultTagTemplate = "`{{tag}}`"

const tagToken = "{{tag}}"

var spaces = regexp.MustCompile(`\s+`)

func (c *Converter) tagList(note *enex.Note, md *markdown.Note) []byte {
	var tt [][]byte

	for _, t := range note.Tags {
		// Default tag template allows spaces in tags, but for custom templates
		// we replace all spaces with underscores to prevent word splitting
		if c.TagTemplate != DefaultTagTemplate {
			t = spaces.ReplaceAllString(t, "_")
		}

		tt = append(tt, []byte(strings.Replace(c.TagTemplate, tagToken, t, 1)))
	}
	return bytes.Join(tt, []byte(" "))
}
func (c *Converter) prependTags(note *enex.Note, md *markdown.Note) {
	if c.err != nil {
		return
	}
	md.Content = append([]byte("\n\n"), md.Content...)
	md.Content = append(c.tagList(note, md), md.Content...)
}

func (c *Converter) frontMatterTags(note *enex.Note, md *markdown.Note) {
	if c.err != nil {
		return
	}

	var tt [][]byte
	for _, t := range note.Tags {
		// Default tag template allows spaces in tags, but for custom templates
		// we replace all spaces with underscores to prevent word splitting
		if c.TagTemplate != DefaultTagTemplate {
			t = spaces.ReplaceAllString(t, "_")
		}

		tt = append(tt, []byte(strings.Replace(c.TagTemplate, tagToken, t, 1)))
	}

	md.Content = append([]byte("\n\n"), md.Content...)
	md.Content = append(bytes.Join(tt, []byte(" ")), md.Content...)
}
