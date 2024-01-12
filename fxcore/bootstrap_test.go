package fxcore_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxcore/testdata/probes"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

type testCtxKey struct{}

func TestBootstrapApp(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var cfg *config.Config

	app := fxcore.NewBootstrapper().BootstrapApp(fx.Populate(&cfg))

	ctx := context.Background()

	err := app.Start(ctx)
	assert.NoError(t, err)

	err = app.Stop(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "core-app", cfg.AppName())
	assert.Equal(t, config.AppEnvDev, cfg.AppEnv())
	assert.False(t, cfg.AppDebug())
	assert.Equal(t, "0.1.0", cfg.AppVersion())
}

func TestBootstrapAppWithContext(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	ctx := context.WithValue(context.Background(), testCtxKey{}, "test-value")

	var popCtx context.Context

	app := fxcore.NewBootstrapper().WithContext(ctx).BootstrapApp(fx.Populate(&popCtx))

	err := app.Start(ctx)
	assert.NoError(t, err)

	err = app.Stop(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "test-value", popCtx.Value(testCtxKey{}))
}

func TestBootstrapAppWithOptions(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var checker *healthcheck.Checker

	app := fxcore.
		NewBootstrapper().
		WithOptions(
			fxhealthcheck.AsCheckerProbe(probes.NewSuccessProbe),
			fxhealthcheck.AsCheckerProbe(probes.NewFailureProbe, healthcheck.Readiness),
		).
		BootstrapApp(fx.Populate(&checker))

	ctx := context.Background()

	err := app.Start(ctx)
	assert.NoError(t, err)

	err = app.Stop(ctx)
	assert.NoError(t, err)

	result := checker.Check(context.Background(), healthcheck.Readiness)
	assert.False(t, result.Success)

	for probeName, probeResult := range result.ProbesResults {
		if probeName == "success" {
			assert.True(t, probeResult.Success)
			assert.Equal(t, "success", probeResult.Message)
		}
		if probeName == "failure" {
			assert.False(t, probeResult.Success)
			assert.Equal(t, "failure", probeResult.Message)
		}
	}
}

func TestRunAppWithContext(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	shutdown := func(sd fx.Shutdowner, lc fx.Lifecycle) {
		lc.Append(fx.Hook{
			OnStart: func(context.Context) error {
				return sd.Shutdown()
			},
		})
	}

	ctx := context.WithValue(context.Background(), testCtxKey{}, "test-value")

	var popCtx context.Context

	fxcore.
		NewBootstrapper().
		WithContext(ctx).
		WithOptions(fx.Invoke(shutdown)).
		RunApp(fx.Populate(&popCtx))

	assert.Equal(t, "test-value", popCtx.Value(testCtxKey{}))
}

func TestRunAppWithOptions(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	shutdown := func(sd fx.Shutdowner, lc fx.Lifecycle) {
		lc.Append(fx.Hook{
			OnStart: func(context.Context) error {
				return sd.Shutdown()
			},
		})
	}

	var checker *healthcheck.Checker

	fxcore.
		NewBootstrapper().
		WithOptions(
			fxhealthcheck.AsCheckerProbe(probes.NewSuccessProbe),
			fxhealthcheck.AsCheckerProbe(probes.NewFailureProbe, healthcheck.Readiness),
		).
		WithOptions(fx.Invoke(shutdown)).
		RunApp(fx.Populate(&checker))

	result := checker.Check(context.Background(), healthcheck.Readiness)
	assert.False(t, result.Success)

	for probeName, probeResult := range result.ProbesResults {
		if probeName == "success" {
			assert.True(t, probeResult.Success)
			assert.Equal(t, "success", probeResult.Message)
		}
		if probeName == "failure" {
			assert.False(t, probeResult.Success)
			assert.Equal(t, "failure", probeResult.Message)
		}
	}
}

func TestBootstrapTestApp(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var cfg *config.Config

	fxcore.
		NewBootstrapper().
		BootstrapTestApp(t, fx.Populate(&cfg)).
		RequireStart().
		RequireStop()

	assert.Equal(t, "core-app", cfg.AppName())
	assert.Equal(t, config.AppEnvTest, cfg.AppEnv())
	assert.True(t, cfg.AppDebug())
	assert.Equal(t, "0.1.0", cfg.AppVersion())
}

func TestBootstrapTestAppWithContext(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	ctx := context.WithValue(context.Background(), testCtxKey{}, "test-value")

	var popCtx context.Context

	fxcore.
		NewBootstrapper().
		WithContext(ctx).
		BootstrapTestApp(t, fx.Populate(&popCtx)).
		RequireStart().
		RequireStop()

	assert.Equal(t, "test-value", popCtx.Value(testCtxKey{}))
}

func TestBootstrapTestAppWithOptions(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var checker *healthcheck.Checker

	fxcore.
		NewBootstrapper().
		WithOptions(
			fxhealthcheck.AsCheckerProbe(probes.NewSuccessProbe),
			fxhealthcheck.AsCheckerProbe(probes.NewFailureProbe, healthcheck.Readiness),
		).
		BootstrapTestApp(t, fx.Populate(&checker)).
		RequireStart().
		RequireStop()

	result := checker.Check(context.Background(), healthcheck.Readiness)
	assert.False(t, result.Success)

	for probeName, probeResult := range result.ProbesResults {
		if probeName == "success" {
			assert.True(t, probeResult.Success)
			assert.Equal(t, "success", probeResult.Message)
		}
		if probeName == "failure" {
			assert.False(t, probeResult.Success)
			assert.Equal(t, "failure", probeResult.Message)
		}
	}
}

func TestRunTestApp(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var cfg *config.Config

	fxcore.
		NewBootstrapper().
		RunTestApp(t, fx.Populate(&cfg))

	assert.Equal(t, "core-app", cfg.AppName())
	assert.Equal(t, config.AppEnvTest, cfg.AppEnv())
	assert.True(t, cfg.AppDebug())
	assert.Equal(t, "0.1.0", cfg.AppVersion())
}

func TestRunTestAppWithContext(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	ctx := context.WithValue(context.Background(), testCtxKey{}, "test-value")

	var popCtx context.Context

	fxcore.
		NewBootstrapper().
		WithContext(ctx).
		RunTestApp(t, fx.Populate(&popCtx))

	assert.Equal(t, "test-value", popCtx.Value(testCtxKey{}))
}

func TestRunTestAppWithOptions(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var checker *healthcheck.Checker

	fxcore.
		NewBootstrapper().
		WithOptions(
			fxhealthcheck.AsCheckerProbe(probes.NewSuccessProbe),
			fxhealthcheck.AsCheckerProbe(probes.NewFailureProbe, healthcheck.Readiness),
		).
		RunTestApp(t, fx.Populate(&checker))

	result := checker.Check(context.Background(), healthcheck.Readiness)
	assert.False(t, result.Success)

	for probeName, probeResult := range result.ProbesResults {
		if probeName == "success" {
			assert.True(t, probeResult.Success)
			assert.Equal(t, "success", probeResult.Message)
		}
		if probeName == "failure" {
			assert.False(t, probeResult.Success)
			assert.Equal(t, "failure", probeResult.Message)
		}
	}
}
