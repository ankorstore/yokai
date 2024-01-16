package worker

import "strings"

func Sanitize(str string) string {
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ReplaceAll(str, "-", "_")

	return strings.ToLower(str)
}
