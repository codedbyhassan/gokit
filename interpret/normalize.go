package interpret

import "strings"

// Normalize trims surrounding whitespace and collapses repeated whitespace.
func Normalize(input string) string {
	return strings.Join(strings.Fields(input), " ")
}
