package internal

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"

	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

// TagReplacer allows to manipulate HTML nodes in order
// to make custom tags
type TagReplacer interface {
	ReplaceTag(node *html.Node)
}

func normalizeHTML(b []byte, rr ...TagReplacer) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		for _, replacer := range rr {
			replacer.ReplaceTag(n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	var out bytes.Buffer
	if err := html.Render(&out, doc); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// Replaces custom <en-media> with corresponding tag
// <img> tag if it is an image and <a> for everything else
// to be able to download it as a file
type Media struct {
	resources map[string]markdown.Resource
}

var htmlFormat = map[markdown.ResourceType]string{
	markdown.Image: `<img src="%s/%s" alt="%s" />`,
	markdown.File:  `<a href="./%s/%s">%s</a>`,
}

func NewReplacerMedia(resources map[string]markdown.Resource) *Media {
	return &Media{resources: resources}
}

func (r *Media) ReplaceTag(n *html.Node) {
	if isMedia(n) {
		if res, ok := r.resources[hashAttr(n)]; ok {
			replaceNode(n, res)
		} else {
			if _, ok := r.resources[""]; ok && len(r.resources) == 1 {
				replaceNode(n, r.resources[""])
			}
		}
	}
}

func isMedia(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "en-media"
}

func hashAttr(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "hash" {
			return a.Val
		}
	}

	return ""
}

func replaceNode(n *html.Node, res markdown.Resource) {
	appendMedia(n, parseOne(resourceReference(res), n))
}

func appendMedia(note, media *html.Node) {
	p := note.Parent
	for isMedia(p) {
		p = p.Parent
	}
	p.AppendChild(media)
	p.AppendChild(parseOne(`<br/>`, note)) // newline
}

// Since we control input, this wrapper gives a simple
// interface which will panic in case of bad strings
func parseOne(h string, context *html.Node) *html.Node {
	nodes, err := html.ParseFragment(strings.NewReader(h), context)
	if err != nil {
		panic("parseHtml: " + err.Error())
	}
	return nodes[0]
}

func resourceReference(res markdown.Resource) string {
	return fmt.Sprintf(htmlFormat[res.Type], res.Type, res.Name, res.Name)
}

// Code replaces div tag stylized to look like code blocks with actual <pre> tag
type Code struct{}

func (r *Code) ReplaceTag(n *html.Node) {
	if isCode(n) {
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "div" {
				p := n.Parent
				p.InsertBefore(parseOne("\n", p), n.NextSibling)
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(n)
		n.Data = "pre"
	}
}

func isCode(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "style" {
				return strings.Contains(a.Val, "-en-codeblock:true")
			}
		}
	}

	return false
}
