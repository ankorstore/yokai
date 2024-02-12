package fxcron

import "strings"

// Contains returns true if a provided string is found in a list of strings.
func Contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}

// Sanitize transforms a given string to not contain spaces or dashes, and to be in lower case.
func Sanitize(str string) string {
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ReplaceAll(str, "-", "_")

	return strings.ToLower(str)
}
