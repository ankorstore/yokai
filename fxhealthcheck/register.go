package fxhealthcheck

import (
	"github.com/ankorstore/yokai/healthcheck"
	"go.uber.org/fx"
)

// AsCheckerProbe registers a [healthcheck.CheckerProbe] into Fx.
func AsCheckerProbe(p any, kinds ...healthcheck.ProbeKind) fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				p,
				fx.As(new(healthcheck.CheckerProbe)),
				fx.ResultTags(`group:"healthcheck-probes"`),
			),
		),
		fx.Supply(
			fx.Annotate(
				NewCheckerProbeDefinition(GetReturnType(p), kinds...),
				fx.As(new(CheckerProbeDefinition)),
				fx.ResultTags(`group:"healthcheck-probes-definitions"`),
			),
		),
	)
}
