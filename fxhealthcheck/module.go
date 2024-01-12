package fxhealthcheck

import (
	"github.com/ankorstore/yokai/healthcheck"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "healthcheck"

// FxHealthcheckModule is the [Fx] healthcheck module.
//
// [Fx]: https://github.com/uber-go/fx
var FxHealthcheckModule = fx.Module(
	ModuleName,
	fx.Provide(
		healthcheck.NewDefaultCheckerFactory,
		NewFxCheckerProbeRegistry,
		NewFxChecker,
	),
)

// FxCheckerParam allows injection of the required dependencies in [NewFxChecker].
type FxCheckerParam struct {
	fx.In
	Factory  healthcheck.CheckerFactory
	Registry *CheckerProbeRegistry
}

// NewFxChecker returns a new [healthcheck.Checker].
func NewFxChecker(p FxCheckerParam) (*healthcheck.Checker, error) {
	registrations, err := p.Registry.ResolveCheckerProbesRegistrations()
	if err != nil {
		return nil, err
	}

	options := []healthcheck.CheckerOption{}
	for _, registration := range registrations {
		options = append(options, healthcheck.WithProbe(registration.Probe(), registration.Kinds()...))
	}

	return p.Factory.Create(options...)
}
