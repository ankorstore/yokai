package server

import (
	"reflect"
	"runtime"
	"strings"
)

// FuncName returns a readable func name for code browsing purposes.
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
