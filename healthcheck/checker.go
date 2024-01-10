package healthcheck

import "context"

// CheckerResult is the result of a [Checker] check.
// It contains a global status, and a list of [CheckerProbeResult] corresponding to each probe execution.
type CheckerResult struct {
	Success       bool                           `json:"success"`
	ProbesResults map[string]*CheckerProbeResult `json:"probes"`
}

// CheckerProbeRegistration represents a registration of a [CheckerProbe] in the [Checker].
type CheckerProbeRegistration struct {
	probe CheckerProbe
	kinds []ProbeKind
}

// NewCheckerProbeRegistration returns a [CheckerProbeRegistration], and accepts a [CheckerProbe] and an optional list of [ProbeKind].
// If no [ProbeKind] is provided, the [CheckerProbe] will be registered to be executed on all kinds of checks.
func NewCheckerProbeRegistration(probe CheckerProbe, kinds ...ProbeKind) *CheckerProbeRegistration {
	return &CheckerProbeRegistration{
		probe: probe,
		kinds: kinds,
	}
}

// Probe returns the [CheckerProbe] of the [CheckerProbeRegistration].
func (r *CheckerProbeRegistration) Probe() CheckerProbe {
	return r.probe
}

// Kinds returns the list of [ProbeKind] of the [CheckerProbeRegistration].
func (r *CheckerProbeRegistration) Kinds() []ProbeKind {
	return r.kinds
}

// Match returns true if the [CheckerProbeRegistration] match any of the provided [ProbeKind] list.
func (r *CheckerProbeRegistration) Match(kinds ...ProbeKind) bool {
	for _, kind := range kinds {
		for _, registrationKind := range r.kinds {
			if registrationKind == kind {
				return true
			}
		}
	}

	return false
}

// Checker provides the possibility to register several [CheckerProbe] and execute them.
type Checker struct {
	registrations map[string]*CheckerProbeRegistration
}

// NewChecker returns a [Checker] instance.
func NewChecker() *Checker {
	return &Checker{
		registrations: map[string]*CheckerProbeRegistration{},
	}
}

// Probes returns the list of [CheckerProbe] registered for the provided list of [ProbeKind].
// If no [ProbeKind] is provided, probes matching all kinds will be returned.
func (c *Checker) Probes(kinds ...ProbeKind) []CheckerProbe {
	var probes []CheckerProbe

	if len(kinds) == 0 {
		kinds = []ProbeKind{Startup, Liveness, Readiness}
	}

	for _, registration := range c.registrations {
		if registration.Match(kinds...) {
			probes = append(probes, registration.probe)
		}
	}

	return probes
}

// RegisterProbe registers a [CheckerProbe] for an optional list of [ProbeKind].
// If no [ProbeKind] is provided, the [CheckerProbe] will be registered for all kinds.
func (c *Checker) RegisterProbe(probe CheckerProbe, kinds ...ProbeKind) *Checker {
	if len(kinds) == 0 {
		kinds = []ProbeKind{Startup, Liveness, Readiness}
	}

	if _, ok := c.registrations[probe.Name()]; ok {
		c.registrations[probe.Name()].kinds = kinds
	} else {
		c.registrations[probe.Name()] = NewCheckerProbeRegistration(probe, kinds...)
	}

	return c
}

// Check executes all the registered probes for a [ProbeKind], passes a [context.Context] to each of them, and returns a [CheckerResult].
// The [CheckerResult] is successful if all probes executed with success.
func (c *Checker) Check(ctx context.Context, kind ProbeKind) *CheckerResult {
	probeResults := map[string]*CheckerProbeResult{}

	success := true
	for name, registration := range c.registrations {
		if registration.Match(kind) {
			pr := registration.probe.Check(ctx)

			success = success && pr.Success
			probeResults[name] = pr
		}
	}

	return &CheckerResult{
		Success:       success,
		ProbesResults: probeResults,
	}
}
