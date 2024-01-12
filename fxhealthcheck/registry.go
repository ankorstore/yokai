package fxhealthcheck

import (
	"fmt"

	"github.com/ankorstore/yokai/healthcheck"
	"go.uber.org/fx"
)

// CheckerProbeRegistry is the registry collecting probes and their definitions.
type CheckerProbeRegistry struct {
	probes      []healthcheck.CheckerProbe
	definitions []CheckerProbeDefinition
}

// FxCheckerProbeRegistryParam allows injection of the required dependencies in [NewFxCheckerProbeRegistry].
type FxCheckerProbeRegistryParam struct {
	fx.In
	Probes      []healthcheck.CheckerProbe `group:"healthcheck-probes"`
	Definitions []CheckerProbeDefinition   `group:"healthcheck-probes-definitions"`
}

// NewFxCheckerProbeRegistry returns as new [CheckerProbeRegistry].
func NewFxCheckerProbeRegistry(p FxCheckerProbeRegistryParam) *CheckerProbeRegistry {
	return &CheckerProbeRegistry{
		probes:      p.Probes,
		definitions: p.Definitions,
	}
}

// ResolveCheckerProbesRegistrations resolves [healthcheck.CheckerProbeRegistration] from their definitions.
func (r *CheckerProbeRegistry) ResolveCheckerProbesRegistrations() ([]*healthcheck.CheckerProbeRegistration, error) {
	registrations := []*healthcheck.CheckerProbeRegistration{}

	for _, definition := range r.definitions {
		implementation, err := r.lookupRegisteredCheckerProbe(definition.ReturnType())
		if err != nil {
			return nil, err
		}

		registrations = append(
			registrations,
			healthcheck.NewCheckerProbeRegistration(implementation, definition.Kinds()...),
		)
	}

	return registrations, nil
}

func (r *CheckerProbeRegistry) lookupRegisteredCheckerProbe(returnType string) (healthcheck.CheckerProbe, error) {
	for _, implementation := range r.probes {
		if GetType(implementation) == returnType {
			return implementation, nil
		}
	}

	return nil, fmt.Errorf("cannot find checker probe implementation for type %s", returnType)
}
