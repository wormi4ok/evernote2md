package markdown

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mattn/godown"
	"golang.org/x/net/html"
)

// ResourceType gives a hint on the way to represent Resource
type ResourceType string

const (
	// Image can be displayed using common ![]() syntax
	Image ResourceType = "image"
	// File should be referenced as an external resource []()
	File ResourceType = "file"
)

type (
	// Note is a markdown representation of valuable knowledge
	// that combines media resources and text represented in markdown format
	Note struct {
		Content []byte
		Media   map[string]Resource
		CTime   time.Time
		MTime   time.Time
	}

	// Resource is a media resource related to a markdown note
	Resource struct {
		Name    string
		Type    ResourceType
		Content []byte
	}
)

// Convert wraps a call to external dependency to provide
// stable interface for package users
func Convert(w io.Writer, r io.Reader, highlights bool, bold bool) error {
	var rules []godown.CustomRule
	var formatter = &Formatter{}
	if highlights {
		formatter.highlights = true
	}
	if bold {
		formatter.bold = true
	}

	rules = append(rules, formatter)

	return godown.Convert(w, r, &godown.Option{CustomRules: rules})
}

// Formatter is a parsing rule to convert Evernote highlights to HTML spans with a background color
type Formatter struct {
	highlights bool
	bold       bool
}

// Rule implements godown.CustomRule interface to extend basic conversion rules and
// convert text highlighted in Evernote to an inline HTML `span` tag with a custom background color
func (r *Formatter) Rule(next godown.WalkFunc) (string, godown.WalkFunc) {
	return "span", func(node *html.Node, w io.Writer, nest int, option *godown.Option) {
		if node.Attr == nil {
			next(node, w, nest, option)
			return
		}

		for _, attr := range node.Attr {
			if r.highlights && attr.Key == "style" && strings.Contains(attr.Val, "-evernote-highlight:true") {
				_, _ = fmt.Fprint(w, `<span style="background-color: #ffaaaa">`)
				next(node, w, nest, option)
				_, _ = fmt.Fprint(w, "</span>")
			} else if r.bold && attr.Key == "style" && strings.Contains(attr.Val, "font-weight: bold") {
				_, _ = fmt.Fprint(w, `<span style="font-weight:bold">`)
				next(node, w, nest, option)
				_, _ = fmt.Fprint(w, "</span>")
			} else {
				next(node, w, nest, option)
			}
		}
	}
}
