package worker

import "strings"

func Sanitize(str string) string {
	san := strings.ReplaceAll(str, " ", "_")
	san = strings.ReplaceAll(san, "-", "_")

	return strings.ToLower(san)
}
