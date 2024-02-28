package status_test

import (
	"testing"

	"github.com/ankorstore/yokai/httpclient/status"
)

func TestNormalizeHTTPStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code int
		want string
	}{
		{"1xx normalization", 101, "1xx"},
		{"2xx normalization", 202, "2xx"},
		{"3xx normalization", 303, "3xx"},
		{"4xx normalization", 404, "4xx"},
		{"5xx normalization", 505, "5xx"},
	}

	for _, tt := range tests {
		got := status.NormalizeHTTPStatus(tt.code)

		if got != tt.want {
			t.Errorf("expected %s, got %s", tt.want, got)
		}
	}
}
