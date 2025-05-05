package server

import (
	"reflect"
	"runtime"
	"strings"
)

// FuncName returns a readable func name for code browsing purposes
func FuncName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// Sanitize transforms a given string to not contain spaces or dashes, and to be in lower case.
func Sanitize(str string) string {
	san := strings.ReplaceAll(str, " ", "_")
	san = strings.ReplaceAll(san, "-", "_")

	return strings.ToLower(san)
}

// Split trims and splits a provided string by comma.
func Split(str string) []string {
	return strings.Split(strings.ReplaceAll(str, " ", ""), ",")
}

// Contain returns true if a given string can be found in a given slice of strings.
func Contain(list []string, item string) bool {
	for _, i := range list {
		if strings.ToLower(i) == strings.ToLower(item) {
			return true
		}
	}

	return false
}
