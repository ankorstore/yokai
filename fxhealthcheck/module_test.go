package fxhealthcheck_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxhealthcheck/testdata/factory"
	"github.com/ankorstore/yokai/fxhealthcheck/testdata/probes"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModule(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var checker *healthcheck.Checker

	fxtest.New(
		t,
		fx.NopLogger,
		fxhealthcheck.FxHealthcheckModule,
		fx.Options(
			fxhealthcheck.AsCheckerProbe(probes.NewSuccessProbe),
			fxhealthcheck.AsCheckerProbe(probes.NewFailureProbe, healthcheck.Liveness, healthcheck.Readiness),
		),
		fx.Populate(&checker),
	).RequireStart().RequireStop()

	// startup probes checks
	result := checker.Check(ctx, healthcheck.Startup)
	assert.True(t, result.Success)

	data, err := json.Marshal(result)
	assert.Nil(t, err)
	assert.Equal(t,
		`{"success":true,"probes":{"successProbe":{"success":true,"message":"some success"}}}`,
		string(data),
	)

	// liveness probes checks
	result = checker.Check(ctx, healthcheck.Liveness)
	assert.False(t, result.Success)

	data, err = json.Marshal(result)
	assert.Nil(t, err)
	assert.Equal(t,
		`{"success":false,"probes":{"failureProbe":{"success":false,"message":"some failure"},"successProbe":{"success":true,"message":"some success"}}}`,
		string(data),
	)

	// readiness probes checks
	result = checker.Check(ctx, healthcheck.Readiness)
	assert.False(t, result.Success)

	data, err = json.Marshal(result)
	assert.Nil(t, err)
	assert.Equal(t,
		`{"success":false,"probes":{"failureProbe":{"success":false,"message":"some failure"},"successProbe":{"success":true,"message":"some success"}}}`,
		string(data),
	)
}

func TestModuleDecoration(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var checker *healthcheck.Checker

	fxtest.New(
		t,
		fx.NopLogger,
		fxhealthcheck.FxHealthcheckModule,
		fx.Decorate(factory.NewTestCheckerFactory),
		fx.Populate(&checker),
	).RequireStart().RequireStop()

	// NewTestCheckerFactory registers failureProbe only for readiness
	assert.True(t, checker.Check(ctx, healthcheck.Startup).Success)
	assert.True(t, checker.Check(ctx, healthcheck.Liveness).Success)
	assert.False(t, checker.Check(ctx, healthcheck.Readiness).Success)
}
