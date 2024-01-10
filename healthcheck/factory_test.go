package healthcheck_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/healthcheck/testdata/probes"
	"github.com/stretchr/testify/assert"
)

func TestDefaultCheckerFactory(t *testing.T) {
	t.Parallel()

	factory := healthcheck.NewDefaultCheckerFactory()

	assert.IsType(t, &healthcheck.DefaultCheckerFactory{}, factory)
	assert.Implements(t, (*healthcheck.CheckerFactory)(nil), factory)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	successProbe := probes.NewSuccessProbe()
	failureProbe := probes.NewFailureProbe()

	factory := healthcheck.NewDefaultCheckerFactory()

	checker, err := factory.Create(
		healthcheck.WithProbe(successProbe),
	)

	assert.Nil(t, err)
	assert.IsType(t, &healthcheck.Checker{}, checker)
	assert.True(t, checker.Check(ctx, healthcheck.Startup).Success)
	assert.True(t, checker.Check(ctx, healthcheck.Liveness).Success)
	assert.True(t, checker.Check(ctx, healthcheck.Readiness).Success)

	checker, err = factory.Create(
		healthcheck.WithProbe(successProbe, healthcheck.Startup, healthcheck.Liveness),
		healthcheck.WithProbe(failureProbe, healthcheck.Startup, healthcheck.Readiness),
	)

	assert.Nil(t, err)
	assert.IsType(t, &healthcheck.Checker{}, checker)
	assert.False(t, checker.Check(ctx, healthcheck.Startup).Success)
	assert.True(t, checker.Check(ctx, healthcheck.Liveness).Success)
	assert.False(t, checker.Check(ctx, healthcheck.Readiness).Success)
}
