package file

import "regexp"

// Max path length is 255 - 9 bytes for extension (.md) in multibyte encoding
const maxPathLength int = 246

// Additional rule for
var illegalChars = regexp.MustCompile(`[\s\\|"'<>&_=+:?*]`)
