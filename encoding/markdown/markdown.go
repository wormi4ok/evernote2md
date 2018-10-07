package markdown

import (
	"io"

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
	// Note is a markdown representation of some valuable knowledge
	// which combines media resources and text represented in markdown format
	Note struct {
		Content []byte
		Media   map[string]Resource
	}

	// Resource is a media resource related to a mardown note
	Resource struct {
		Name    string
		Type    ResourceType
		Content []byte
	}
)

// Convert wraps a call to external dependency to provide
// stable interface for package users
func Convert(w io.Writer, r io.Reader) error {
	return godown.Convert(w, r, nil)
}
