// +build !windows

package file

import "regexp"

// Max path length is 1000 - 9 bytes for extension (.md) in multibyte encoding
const maxPathLength int = 991

// Semicolon is not allowed in MacOS and spaces is just my personal preference
var illegalChars = regexp.MustCompile(`[\s:]`)
