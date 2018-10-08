package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"regexp"
	"strings"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

// Convert Evernote file to markdown
func Convert(note *enex.Note) (*markdown.Note, error) {
	var md markdown.Note
	md.Media = map[string]markdown.Resource{}

	r := note.Resources
	for i := range r {
		p, err := ioutil.ReadAll(decoder(r[i].Data))
		if err != nil {
			return nil, err
		}

		rType := markdown.File
		if isImage(r[i].Mime) {
			rType = markdown.Image
		}

		mdr := markdown.Resource{
			Name:    r[i].Attributes.Filename,
			Type:    rType,
			Content: p,
		}
		if mdr.Name == "" {
			mdr.Name = r[i].ID + guessExt(r[i].Mime)
		}

		md.Media[r[i].ID] = mdr
	}

	html, err := convertEnMediaToHTML(note.Content, md.Media)
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

func decoder(d enex.Data) io.Reader {
	if d.Encoding == "base64" {
		return base64.NewDecoder(base64.StdEncoding, bytes.NewReader(d.Content))
	}

	return bytes.NewReader(d.Content)
}

func guessExt(mimeType string) string {
	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(ext) == 0 {
		return ""
	}
	return ext[0]
}

var reImg = regexp.MustCompile(`^image/[\w]+`)

func isImage(mimeType string) bool {
	return reImg.MatchString(mimeType)
}
