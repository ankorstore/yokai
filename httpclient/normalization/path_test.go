package normalization_test

import (
	"testing"

	"github.com/ankorstore/yokai/httpclient/normalization"
)

func TestMask(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		masks map[string]string
		path  string
		want  string
	}{
		"primary mask applied": {
			map[string]string{
				`/foo/(.+)/bar\?page=(.+)#baz`: "/foo/{fooId}/bar?page={pageId}#baz",
			},
			"/foo/1/bar?page=1#baz",
			"/foo/{fooId}/bar?page={pageId}#baz",
		},
		"secondary mask applied": {
			map[string]string{
				`/foo/(.+)/baz\?page=(.+)#baz`: "/foo/{fooId}/baz?page={pageId}#baz",
				`/foo/(.+)/bar\?page=(.+)#baz`: "/foo/{fooId}/bar?page={pageId}#baz",
			},
			"/foo/1/bar?page=1#baz",
			"/foo/{fooId}/bar?page={pageId}#baz",
		},
		"primary mask not applied": {
			map[string]string{
				`/foo/(.+)/bar\?page=(.+)#baz`: "/foo/{fooId}/bar?page={pageId}#baz",
			},
			"/foo/1/bar?pages=1#baz",
			"/foo/1/bar?pages=1#baz",
		},
		"primary mask applied on invalid regexp": {
			map[string]string{
				`(.`: "/foo/{fooId}/bar?page={pageId}#baz",
			},
			"/foo/1/bar?page=1#baz",
			"/foo/1/bar?page=1#baz",
		},
		"no mask applied on empty masks list": {
			map[string]string{},
			"/foo/1/bar?page=1#baz",
			"/foo/1/bar?page=1#baz",
		},
	}

	for name, tt := range tests {
		got := normalization.NormalizePath(tt.masks, tt.path)
		if got != tt.want {
			t.Errorf("%s: expected %s, got %s", name, tt.want, got)
		}
	}
}
