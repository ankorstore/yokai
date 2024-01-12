package fxhttpserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareKindAsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    fxhttpserver.MiddlewareKind
		expected string
	}{
		{fxhttpserver.GlobalUse, "global-use"},
		{fxhttpserver.GlobalPre, "global-pre"},
		{fxhttpserver.Attached, "attached"},
		{fxhttpserver.MiddlewareKind(1000), "global-use"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			actual := tt.input.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
