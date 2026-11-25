package markdown

import (
	"fmt"
	"io"
	"strings"

	"github.com/mattn/godown"
	"golang.org/x/net/html"
)

// HighlightedText is a parsing rule to convert Evernote highlights to HTML spans with a background color
type HighlightedText struct{}

// Rule implements godown.CustomRule interface to extend basic conversion rules and
// convert text highlighted in Evernote to an inline HTML `span` tag with a custom background color
func (r *HighlightedText) Rule(next godown.WalkFunc) (string, godown.WalkFunc) {
	return "span", func(node *html.Node, w io.Writer, nest int, option *godown.Option) {
		highlighted := false
		for _, attr := range node.Attr {
			if attr.Key == "style" && strings.Contains(attr.Val, "-evernote-highlight:true") {
				highlighted = true
				break
			}
		}

		if highlighted {
			_, _ = fmt.Fprint(w, `<span style="background-color: #ffaaaa">`)
		}
		next(node, w, nest, option)
		if highlighted {
			_, _ = fmt.Fprint(w, "</span>")
		}
	}
}

// TodoItem is a parsing rule to convert Evernote checkboxes to corresponding GitHub Flavoured Markdown items
type TodoItem struct{}

// Rule implements godown.CustomRule interface to handle Evernote-specific "en-todo" tag
// It converts the tag to a Markdown format with correct "checked" state
func (r TodoItem) Rule(next godown.WalkFunc) (string, godown.WalkFunc) {
	return "en-todo", func(node *html.Node, w io.Writer, nest int, option *godown.Option) {
		for _, attr := range node.Attr {
			if attr.Key == "checked" && attr.Val == "true" {
				_, _ = fmt.Fprint(w, "[x] ")
				next(node, w, nest, option)
				return
			}
		}
		_, _ = fmt.Fprint(w, "[ ] ")
		next(node, w, nest, option)
	}
}
