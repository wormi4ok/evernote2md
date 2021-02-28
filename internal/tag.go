package internal

import (
	"regexp"
	"strings"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

// DefaultTagTemplate format if none specified
const DefaultTagTemplate = "`{{tag}}`"

const tagToken = "{{tag}}"

var spaces = regexp.MustCompile(`\s+`)

func (c *Converter) tagList(note *enex.Note, tagTemplate string, joinString string, spacesToUnderscores bool) string {
	var tt []string

	for _, t := range note.Tags {
		// Default tag template allows spaces in tags, but for custom templates
		// we replace all spaces with underscores to prevent word splitting
		if spacesToUnderscores {
			t = spaces.ReplaceAllString(t, "_")
		}
		tt = append(tt, strings.Replace(tagTemplate, tagToken, t, 1))
	}
	return strings.Join(tt, joinString)
}
func (c *Converter) prependTags(note *enex.Note, md *markdown.Note) {
	if c.err != nil {
		return
	}
	md.Content = append([]byte("\n\n"), md.Content...)
	md.Content = append([]byte(c.tagList(note, c.TagTemplate, " ", c.TagTemplate != DefaultTagTemplate)), md.Content...)
}
