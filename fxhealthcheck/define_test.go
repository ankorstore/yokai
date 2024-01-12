package fxhealthcheck_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/stretchr/testify/assert"
)

func TestNewCheckerProbeDefinition(t *testing.T) {
	t.Parallel()

	definition := fxhealthcheck.NewCheckerProbeDefinition("test", healthcheck.Liveness, healthcheck.Readiness)

	assert.Implements(t, (*fxhealthcheck.CheckerProbeDefinition)(nil), definition)
	assert.Equal(t, "test", definition.ReturnType())
	assert.Equal(t, []healthcheck.ProbeKind{healthcheck.Liveness, healthcheck.Readiness}, definition.Kinds())
}
