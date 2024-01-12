package fxhealthcheck_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxhealthcheck/testdata/probes"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/stretchr/testify/assert"
)

func TestNewCheckerProbeRegistry(t *testing.T) {
	t.Parallel()

	param := fxhealthcheck.FxCheckerProbeRegistryParam{
		Probes:      []healthcheck.CheckerProbe{},
		Definitions: []fxhealthcheck.CheckerProbeDefinition{},
	}
	registry := fxhealthcheck.NewFxCheckerProbeRegistry(param)

	assert.IsType(t, &fxhealthcheck.CheckerProbeRegistry{}, registry)
}

func TestResolveCheckerProbesRegistrationsSuccess(t *testing.T) {
	t.Parallel()

	param := fxhealthcheck.FxCheckerProbeRegistryParam{
		Probes: []healthcheck.CheckerProbe{
			probes.NewSuccessProbe(),
			probes.NewFailureProbe(),
		},
		Definitions: []fxhealthcheck.CheckerProbeDefinition{
			fxhealthcheck.NewCheckerProbeDefinition("*probes.SuccessProbe", healthcheck.Liveness),
			fxhealthcheck.NewCheckerProbeDefinition("*probes.FailureProbe", healthcheck.Readiness),
		},
	}

	registry := fxhealthcheck.NewFxCheckerProbeRegistry(param)

	registrations, err := registry.ResolveCheckerProbesRegistrations()
	assert.NoError(t, err)

	assert.Len(t, registrations, 2)
	assert.IsType(t, &probes.SuccessProbe{}, registrations[0].Probe())
	assert.Equal(t, []healthcheck.ProbeKind{healthcheck.Liveness}, registrations[0].Kinds())
	assert.IsType(t, &probes.FailureProbe{}, registrations[1].Probe())
	assert.Equal(t, []healthcheck.ProbeKind{healthcheck.Readiness}, registrations[1].Kinds())
}

func TestResolveCheckerProbesRegistrationsFailure(t *testing.T) {
	t.Parallel()

	param := fxhealthcheck.FxCheckerProbeRegistryParam{
		Probes: []healthcheck.CheckerProbe{
			probes.NewSuccessProbe(),
		},
		Definitions: []fxhealthcheck.CheckerProbeDefinition{
			fxhealthcheck.NewCheckerProbeDefinition("invalid", healthcheck.Liveness),
		},
	}

	registry := fxhealthcheck.NewFxCheckerProbeRegistry(param)

	_, err := registry.ResolveCheckerProbesRegistrations()
	assert.Error(t, err)
	assert.Equal(t, "cannot find checker probe implementation for type invalid", err.Error())
}
