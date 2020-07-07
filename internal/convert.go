package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
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
