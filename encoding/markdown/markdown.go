package markdown

type (
	// Note is a markdown representation of some text
	Note struct {
		Content []byte
		Media   map[string]Resource
	}

	// Resource is a resource
	Resource struct {
		Name    string
		Content []byte
	}
)
