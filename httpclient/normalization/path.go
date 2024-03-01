package normalization

import (
	"regexp"
)

// NormalizePath normalizes a path if matching one of the provided masks, or returns the original path instead.
//
// For example: NormalizePath(map[string]string{"/foo/(.+)", "/foo/{id}"}, "/foo/1") will return "/foo/{id}".
func NormalizePath(masks map[string]string, path string) string {
	for pattern, mask := range masks {
		re, err := regexp.Compile(pattern)
		if err == nil {
			matched := re.MatchString(path)
			if matched {
				return mask
			}
		}
	}

	return path
}
