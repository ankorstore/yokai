package healthcheck_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/healthcheck/testdata/probes"
	"github.com/stretchr/testify/assert"
)

func TestNewCheckerProbeRegistration(t *testing.T) {
	t.Parallel()

	successProbe := probes.NewSuccessProbe()

	registration := healthcheck.NewCheckerProbeRegistration(successProbe, healthcheck.Startup, healthcheck.Liveness)

	assert.IsType(t, &healthcheck.CheckerProbeRegistration{}, registration)

	assert.Equal(t, successProbe, registration.Probe())
	assert.Equal(t, []healthcheck.ProbeKind{healthcheck.Startup, healthcheck.Liveness}, registration.Kinds())
}

func TestNewChecker(t *testing.T) {
	t.Parallel()

	checker := healthcheck.NewChecker()

	assert.IsType(t, &healthcheck.Checker{}, checker)
	assert.Len(t, checker.Probes(), 0)
}

func TestCheckerProbeRegistration(t *testing.T) {
	t.Parallel()

	checker := healthcheck.NewChecker()

	assert.Len(t, checker.Probes(), 0)

	successProbe := probes.NewSuccessProbe()
	failureProbe := probes.NewFailureProbe()

	checker.RegisterProbe(successProbe, healthcheck.Startup, healthcheck.Liveness)
	checker.RegisterProbe(failureProbe, healthcheck.Startup, healthcheck.Readiness)

	assert.Len(t, checker.Probes(), 2)
	assert.Len(t, checker.Probes(healthcheck.Startup, healthcheck.Liveness, healthcheck.Readiness), 2)
	assert.Len(t, checker.Probes(healthcheck.Startup), 2)
	assert.Len(t, checker.Probes(healthcheck.Liveness), 1)
	assert.Equal(t, successProbe.Name(), checker.Probes(healthcheck.Liveness)[0].Name())
	assert.Len(t, checker.Probes(healthcheck.Readiness), 1)
	assert.Equal(t, failureProbe.Name(), checker.Probes(healthcheck.Readiness)[0].Name())

	checker.RegisterProbe(successProbe)
	assert.Len(t, checker.Probes(), 2)
	assert.Len(t, checker.Probes(healthcheck.Startup, healthcheck.Liveness, healthcheck.Readiness), 2)
	assert.Len(t, checker.Probes(healthcheck.Startup), 2)
	assert.Len(t, checker.Probes(healthcheck.Liveness), 1)
	assert.Len(t, checker.Probes(healthcheck.Readiness), 2)
}

func TestCheckerCheck(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	checker := healthcheck.NewChecker()

	assert.Len(t, checker.Probes(), 0)
	assert.True(t, checker.Check(ctx, healthcheck.Startup).Success)
	assert.True(t, checker.Check(ctx, healthcheck.Liveness).Success)
	assert.True(t, checker.Check(ctx, healthcheck.Readiness).Success)

	successProbe := probes.NewSuccessProbe()
	failureProbe := probes.NewFailureProbe()

	checker.RegisterProbe(successProbe, healthcheck.Startup, healthcheck.Liveness)
	checker.RegisterProbe(failureProbe, healthcheck.Startup, healthcheck.Readiness)

	result := checker.Check(ctx, healthcheck.Startup)
	assert.False(t, result.Success)

	data, err := json.Marshal(result)
	assert.Nil(t, err)
	assert.Equal(t,
		`{"success":false,"probes":{"failureProbe":{"success":false,"message":"some failure"},"successProbe":{"success":true,"message":"some success"}}}`,
		string(data),
	)

	result = checker.Check(ctx, healthcheck.Liveness)
	assert.True(t, result.Success)

	data, err = json.Marshal(result)
	assert.Nil(t, err)
	assert.Equal(t,
		`{"success":true,"probes":{"successProbe":{"success":true,"message":"some success"}}}`,
		string(data),
	)

	result = checker.Check(ctx, healthcheck.Readiness)
	assert.False(t, result.Success)

	data, err = json.Marshal(result)
	assert.Nil(t, err)
	assert.Equal(t,
		`{"success":false,"probes":{"failureProbe":{"success":false,"message":"some failure"}}}`,
		string(data),
	)

	checker.RegisterProbe(successProbe)

	result = checker.Check(ctx, healthcheck.Readiness)
	assert.False(t, result.Success)

	data, err = json.Marshal(result)
	assert.Nil(t, err)
	assert.Equal(t,
		`{"success":false,"probes":{"failureProbe":{"success":false,"message":"some failure"},"successProbe":{"success":true,"message":"some success"}}}`,
		string(data),
	)
}
