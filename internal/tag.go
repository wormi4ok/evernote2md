package internal

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

const DefaultTagFormat = "`{{tag}}`"

const tagToken = "{{tag}}"

var spaces = regexp.MustCompile(`\s+`)

func (c *Converter) prependTags(note *enex.Note, md *markdown.Note) {
	if c.err != nil {
		return
	}

	var tt [][]byte
	for _, t := range note.Tags {
		// Default tag format allows spaces in tags, but for custom formats
		// we replace all spaces with underscores to prevent word splitting
		if c.TagFormat != DefaultTagFormat {
			t = spaces.ReplaceAllString(t, "_")
		}

		tt = append(tt, []byte(strings.Replace(c.TagFormat, tagToken, t, 1)))
	}

	md.Content = append([]byte("\n\n"), md.Content...)
	md.Content = append(bytes.Join(tt, []byte(" ")), md.Content...)
}
