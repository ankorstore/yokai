package fxhealthcheck

import "github.com/ankorstore/yokai/healthcheck"

// CheckerProbeDefinition is the interface for probes definitions.
type CheckerProbeDefinition interface {
	ReturnType() string
	Kinds() []healthcheck.ProbeKind
}

type checkerProbeDefinition struct {
	returnType string
	kinds      []healthcheck.ProbeKind
}

// NewCheckerProbeDefinition returns a new [CheckerProbeDefinition].
func NewCheckerProbeDefinition(returnType string, kinds ...healthcheck.ProbeKind) CheckerProbeDefinition {
	return &checkerProbeDefinition{
		returnType: returnType,
		kinds:      kinds,
	}
}

// ReturnType returns the probe return type.
func (c *checkerProbeDefinition) ReturnType() string {
	return c.returnType
}

// Kinds returns the probe registration kinds.
func (c *checkerProbeDefinition) Kinds() []healthcheck.ProbeKind {
	return c.kinds
}
