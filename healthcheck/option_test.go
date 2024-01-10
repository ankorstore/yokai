package healthcheck_test

import (
	"testing"

	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/healthcheck/testdata/probes"
	"github.com/stretchr/testify/assert"
)

func TestWithProbe(t *testing.T) {
	t.Parallel()

	probe := probes.NewSuccessProbe()

	opt := healthcheck.DefaultCheckerOptions()
	healthcheck.WithProbe(probe)(&opt)

	assert.Equal(t, probe, opt.Registrations[probe.Name()].Probe())
	assert.Equal(
		t,
		[]healthcheck.ProbeKind{
			healthcheck.Startup,
			healthcheck.Liveness,
			healthcheck.Readiness,
		},
		opt.Registrations[probe.Name()].Kinds(),
	)
}

func TestWithAlreadyRegisteredProbe(t *testing.T) {
	t.Parallel()

	probe := probes.NewSuccessProbe()

	opt := healthcheck.DefaultCheckerOptions()

	healthcheck.WithProbe(probe)(&opt)

	assert.Equal(t, probe, opt.Registrations[probe.Name()].Probe())
	assert.Equal(
		t,
		[]healthcheck.ProbeKind{
			healthcheck.Startup,
			healthcheck.Liveness,
			healthcheck.Readiness,
		},
		opt.Registrations[probe.Name()].Kinds(),
	)

	healthcheck.WithProbe(probe, healthcheck.Startup)(&opt)

	assert.Equal(t, probe, opt.Registrations[probe.Name()].Probe())
	assert.Equal(
		t,
		[]healthcheck.ProbeKind{
			healthcheck.Startup,
		},
		opt.Registrations[probe.Name()].Kinds(),
	)
}
