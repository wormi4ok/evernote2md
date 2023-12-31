package markdown

import (
	"io"
	"time"

	"github.com/mattn/godown"
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
func Convert(w io.Writer, r io.Reader, highlights, escapeSpecialChars bool) error {
	rules := []godown.CustomRule{
		&TodoItem{}, // Handling checkboxes is always enabled
	}

	if highlights {
		rules = append(rules, &HighlightedText{})
	}

	return godown.Convert(w, r, &godown.Option{
		CustomRules: rules,
		DoNotEscape: !escapeSpecialChars,
	})
}
