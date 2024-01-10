package healthcheck_test

import (
	"testing"

	"github.com/ankorstore/yokai/healthcheck"
	"github.com/stretchr/testify/assert"
)

func TestProbeKindAsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		kind     healthcheck.ProbeKind
		expected string
	}{
		{healthcheck.Startup, "startup"},
		{healthcheck.Readiness, "readiness"},
		{healthcheck.Liveness, "liveness"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.kind.String())
	}
}
