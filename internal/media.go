package internal

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"

	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

var htmlFormat = map[markdown.ResourceType]string{
	markdown.Image: `<img src="%s/%s" alt="%s" />`,
	markdown.File:  `<a href="./%s/%s">%s</a>`,
}

func convertEnMediaToHTML(b []byte, rr map[string]markdown.Resource) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if isCode(n) {
			n.Data = "pre"
		}
		if isMedia(n) {
			if res, ok := rr[hashAttr(n)]; ok {
				replaceNode(n, res)
			} else {
				if _, ok := rr[""]; ok && len(rr) == 1 {
					replaceNode(n, rr[""])
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	var out bytes.Buffer
	html.Render(&out, doc)

	return out.Bytes(), nil
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

func hashAttr(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "hash" {
			return a.Val
		}
	}

	return ""
}

func isMedia(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "en-media"
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

func resourceReference(res markdown.Resource) string {
	return fmt.Sprintf(htmlFormat[res.Type], res.Type, res.Name, res.Name)
}
