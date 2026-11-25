package markdown_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wormi4ok/evernote2md/encoding/markdown"
)

func TestConvertHighlightedSpanWithMultipleAttributes(t *testing.T) {
	var out bytes.Buffer
	input := `<span class="x" style="-evernote-highlight:true">important</span>`

	err := markdown.Convert(&out, strings.NewReader(input), true, true)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	got := out.String()
	if count := strings.Count(got, "important"); count != 1 {
		t.Errorf("Convert() rendered text %d times, want 1: %q", count, got)
	}
	if count := strings.Count(got, `<span style="background-color: #ffaaaa">`); count != 1 {
		t.Errorf("Convert() rendered highlight wrapper %d times, want 1: %q", count, got)
	}
}
