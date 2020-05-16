package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

// Convert Evernote file to markdown
func Convert(note *enex.Note) (*markdown.Note, error) {
	var md markdown.Note
	md.Media = map[string]markdown.Resource{}

	if err := mapResources(note, md); err != nil {
		return nil, err
	}

	html, err := normalizeHTML(note.Content, NewReplacerMedia(md.Media), &Code{})
	if err != nil {
		return nil, err
	}

	content := prependTags(note.Tags, string(html))
	content = prependTitle(note.Title, content)

	var b bytes.Buffer
	err = markdown.Convert(&b, strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	md.Content = regexp.MustCompile(`\n{3,}`).ReplaceAllLiteral(b.Bytes(), []byte("\n\n"))
	md.Content = append(bytes.TrimRight(md.Content, "\n"), '\n')

	return &md, nil
}

func mapResources(note *enex.Note, md markdown.Note) error {
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

		mdr := markdown.Resource{
			Name:    sanitize(r[i].Attributes.Filename),
			Type:    rType,
			Content: p,
		}
		if mdr.Name == "" {
			mdr.Name = r[i].ID + guessExt(r[i].Mime)
		}

		md.Media[r[i].ID] = mdr
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
