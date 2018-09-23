package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

var re = regexp.MustCompile(`<en-media type="[\w\/]*" hash="([a-z0-9]+)"/>`)

// Converter is the main entity that can convert
// *.enex notes to markdown representation
type Converter struct {
	AssetsDir string
}

// Convert Evernote file to markdown
func (c Converter) Convert(note *enex.Note) (*markdown.Note, error) {
	var md markdown.Note
	md.Media = map[string]markdown.Resource{}

	content := re.ReplaceAllString(string(note.Content), `<img src="`+c.AssetsDir+`/$1"><br>`)
	for _, res := range note.Resources {
		content = strings.Replace(content, res.ID, res.Attributes.Filename, 1)

		p, err := ioutil.ReadAll(decoder(res.Data))
		if err != nil {
			return nil, err
		}
		mdr := markdown.Resource{
			Name:    res.Attributes.Filename,
			Content: p,
		}

		md.Media[res.ID] = mdr
	}

	content = prependTags(note.Tags, content)
	content = prependTitle(note.Title, content)

	var b bytes.Buffer
	err := markdown.Convert(&b, strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	md.Content = regexp.MustCompile(`\n{3,}`).ReplaceAllLiteral(b.Bytes(), []byte("\n\n"))
	md.Content = append(bytes.TrimRight(md.Content, "\n"), '\n')

	return &md, nil
}

func prependTags(tags []string, content string) string {
	var tt []string
	for _, t := range tags {
		tt = append(tags, fmt.Sprintf("<code>%s</code>", t))
	}
	return strings.Join(tt, "") + "<br>" + content
}

func prependTitle(title, content string) string {
	return fmt.Sprintf("<h1>%s</h1>", title) + content
}

func decoder(d enex.Data) io.Reader {
	if d.Encoding == "base64" {
		return base64.NewDecoder(base64.StdEncoding, bytes.NewReader(d.Content))
	}

	return bytes.NewReader(d.Content)
}
