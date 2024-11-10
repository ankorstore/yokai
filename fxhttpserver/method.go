package fxhttpserver

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	// MethodPropfind can be used on collection and property resources.
	MethodPropfind = "PROPFIND"
	// MethodReport Method can be used to get information about a resource, see rfc 3253
	MethodReport = "REPORT"
	// AllMethods is a shortcut to specify all valid methods.
	AllMethods = "*"
)

var validMethods = []string{
	http.MethodConnect,
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
	http.MethodTrace,
	MethodPropfind,
	MethodReport,
}

func ExtractMethods(methods string) ([]string, error) {
	if methods == AllMethods {
		return validMethods, nil
	}

	var extractedMethods []string

	for _, method := range Split(methods) {
		method = strings.ToUpper(method)

		if Contains(validMethods, method) {
			extractedMethods = append(extractedMethods, method)
		} else {
			return nil, fmt.Errorf("invalid HTTP method %q", method)
		}
	}

	return extractedMethods, nil
}
