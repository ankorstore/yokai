package fxcron

import "strings"

func Contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}

func Sanitize(str string) string {
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ReplaceAll(str, "-", "_")

	return strings.ToLower(str)
}
