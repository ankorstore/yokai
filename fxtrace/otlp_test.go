package fxtrace_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestCreateOtlpGrpcDialRetryPolicyWithDefaults(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "noop")
	t.Setenv("SAMPLER_TYPE", "always-on")

	var cfg *config.Config

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fx.Populate(&cfg),
	).RequireStart().RequireStop()

	retryPolicy := fxtrace.BuildOtlpGrpcDialRetryPolicy(cfg)

	assert.Equal(
		t,
		`{
            "methodConfig": [{
                "waitForReady": true,
                "retryPolicy": {
                    "MaxAttempts": 4,
                    "InitialBackoff": "0.1s",
                    "MaxBackoff": "1s",
                    "BackoffMultiplier": 2,
                    "RetryableStatusCodes": [ "UNAVAILABLE" ]
                }
            }]
        }`,
		retryPolicy,
	)
}

func TestCreateOtlpGrpcDialRetryPolicyWithConfig(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "noop")
	t.Setenv("SAMPLER_TYPE", "always-on")
	t.Setenv("APP_ENV", "test")

	var cfg *config.Config

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fx.Populate(&cfg),
	).RequireStart().RequireStop()

	retryPolicy := fxtrace.BuildOtlpGrpcDialRetryPolicy(cfg)

	assert.Equal(
		t,
		`{
            "methodConfig": [{
                "waitForReady": true,
                "retryPolicy": {
                    "MaxAttempts": 8,
                    "InitialBackoff": "0.2s",
                    "MaxBackoff": "2s",
                    "BackoffMultiplier": 2,
                    "RetryableStatusCodes": [ "UNAVAILABLE", "INTERNAL" ]
                }
            }]
        }`,
		retryPolicy,
	)
}
