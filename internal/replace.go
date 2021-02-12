package internal

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

// TagReplacer allows manipulating HTML nodes in order
// to present custom tags correctly in Markdown format after conversion
type TagReplacer interface {
	ReplaceTag(node *html.Node)
}

func (c *Converter) normalizeHTML(note *enex.Note, _ *markdown.Note, rr ...TagReplacer) {
	if c.err != nil {
		return
	}
	doc, err := html.Parse(bytes.NewReader(note.Content))
	if c.err = err; err != nil {
		return
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
	if c.err = html.Render(&out, doc); c.err != nil {
		return
	}
	note.Content = out.Bytes()
}

// Media tag replacer puts a standard HTML <img> tag
// instead of custom <en-media> tag if it is an image
// and <a> tag for everything else to be able to download it as a file
type Media struct {
	resources map[string]markdown.Resource

	// If identifiers are missing we use resources one by one
	cnt int
}

var htmlFormat = map[markdown.ResourceType]string{
	markdown.Image: `<img src="%s/%s" alt="%s" />`,
	markdown.File:  `<a href="./%s/%s">%s</a>`,
}

// NewReplacerMedia creates a Media TagReplacer using resources as a data source
func NewReplacerMedia(resources map[string]markdown.Resource) *Media {
	return &Media{resources: resources}
}

// ReplaceTag implements the TagReplacer interface
func (r *Media) ReplaceTag(n *html.Node) {
	if isMedia(n) {
		if res, ok := r.resources[hashAttr(n)]; ok {
			replaceNode(n, res)
			return
		}
		replaceNode(n, r.resources[strconv.Itoa(r.cnt)])
		r.cnt++
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

// Code replaces div tag stylized to look like code blocks with an actual <pre> tag
type Code struct{}

// ReplaceTag implements the TagReplacer interface
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

// ExtraDiv removes extra line break in tables and lists
type ExtraDiv struct{}

// ReplaceTag implements the TagReplacer interface
func (*ExtraDiv) ReplaceTag(n *html.Node) {
	if hasExtraDiv(n) {
		wrapper := n.FirstChild
		if wrapper != nil && wrapper.Data == "div" {
			content := wrapper.FirstChild
			if content == nil {
				return
			}
			wrapper.RemoveChild(content)
			if content.Data != "br" || content.FirstChild != nil {
				n.AppendChild(content)
			}
			n.RemoveChild(wrapper)
		}
	}
}

func hasExtraDiv(n *html.Node) bool {
	tagsToClean := []string{"li", "td", "th"}
	if n.Type == html.ElementNode {
		for _, tag := range tagsToClean {
			if tag == n.Data {
				return true
			}
		}
	}

	return false
}

// TextFormatter catches bold and italic, bold takes precedence
type TextFormatter struct{}

// ReplaceTag implements the TagReplacer interface
func (*TextFormatter) ReplaceTag(n *html.Node) {
	if isBold(n) {
		n.Data = "strong"
		n.Attr = []html.Attribute{}
	} else if isItalic(n) {
		n.Data = "i"
		n.Attr = []html.Attribute{}
	}
}

func isBold(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "span" {
		for _, a := range n.Attr {
			if a.Key == "style" {
				return strings.Contains(a.Val, "font-weight: bold")
			}
		}
	}

	return false
}

func isItalic(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "span" {
		for _, a := range n.Attr {
			if a.Key == "style" {
				return strings.Contains(a.Val, "font-style: italic")
			}
		}
	}

	return false
}

// EmptyAnchor removes anchor tags without text
type EmptyAnchor struct{}

// ReplaceTag implements the TagReplacer interface
func (*EmptyAnchor) ReplaceTag(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		if content := n.FirstChild; content == nil {
			n.Data = "span"
		}
	}
}
