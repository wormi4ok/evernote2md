package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/mattn/godown"
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
	var b bytes.Buffer
	err := godown.Convert(&b, strings.NewReader(content), nil)
	if err != nil {
		return nil, err
	}
	title := []byte(fmt.Sprintf("# %s\n\n", note.Title))
	md.Content = append(title, bytes.TrimRight(b.Bytes(), "\n")...)
	md.Content = append(md.Content, '\n')

	return &md, nil
}

func decoder(d enex.Data) io.Reader {
	if d.Encoding == "base64" {
		return base64.NewDecoder(base64.StdEncoding, bytes.NewReader(d.Content))
	}

	return bytes.NewReader(d.Content)
}
