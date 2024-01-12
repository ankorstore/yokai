package fxhealthcheck_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxhealthcheck/testdata/probes"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/stretchr/testify/assert"
)

func TestAsCheckerProbe(t *testing.T) {
	t.Parallel()

	result := fxhealthcheck.AsCheckerProbe(probes.NewSuccessProbe, healthcheck.Startup)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result))
}
