package httpserver

import (
	"strings"
)

// MatchPrefix returns true if a given prefix matches an item of a given prefixes list.
func MatchPrefix(prefixes []string, str string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}

	return false
}
