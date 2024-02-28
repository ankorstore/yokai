package fxhttpclient

import "strings"

// Sanitize transforms a given string to not contain spaces or dashes, and to be in lower case.
func Sanitize(str string) string {
	san := strings.ReplaceAll(str, " ", "_")
	san = strings.ReplaceAll(san, "-", "_")

	return strings.ToLower(san)
}
