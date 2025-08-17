package enex

import (
	"bytes"
	"errors"
	"io"
	"regexp"
	"strings"
)

var reCDATA = regexp.MustCompile(`<!\[CDATA\[(.*?)\]\]>`)

// detectNestedCDATA reads the first 8KB to check for nested or malformed CDATA.
// Returns whether fixing is needed, and a reader that includes all data.
func detectNestedCDATA(r io.Reader) (bool, io.Reader, error) {
	const scanSize = 8192
	var buf bytes.Buffer
	limited := io.LimitReader(r, scanSize)
	tee := io.TeeReader(limited, &buf)

	chunk, err := io.ReadAll(tee)
	if err != nil && !errors.Is(err, io.EOF) {
		return false, nil, err
	}
	needsFix := hasNestedCDATA(string(chunk))
	return needsFix, io.MultiReader(&buf, r), nil
}

// hasNestedCDATA returns true if the input has unbalanced CDATA tags
// or any CDATA section contains another CDATA opening tag.
func hasNestedCDATA(input string) bool {
	openingTags := strings.Count(input, "<![CDATA[")
	closingTags := strings.Count(input, "]]>")
	if openingTags > 1 && openingTags != closingTags {
		return true
	}
	cdataStart := strings.Index(input, "<![CDATA[")
	if cdataStart != -1 {
		remaining := input[cdataStart+9:]
		cdataEnd := strings.Index(remaining, "]]>")
		if cdataEnd != -1 && strings.Contains(remaining[:cdataEnd], "<![CDATA[") {
			return true
		}
	}
	return false
}

// removeNestedCDATA removes nested CDATA tags recursively.
// Evernote sometimes produces nested CDATA, which is invalid XML.
func removeNestedCDATA(input string) string {
	output := reCDATA.ReplaceAllStringFunc(input, func(match string) string {
		submatch := reCDATA.FindStringSubmatch(match)
		if len(submatch) > 1 {
			return submatch[1]
		}
		return match
	})
	if output != input {
		return removeNestedCDATA(output)
	}
	return output
}
