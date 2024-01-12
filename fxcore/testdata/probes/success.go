package probes

import (
	"context"

	"github.com/ankorstore/yokai/healthcheck"
)

type SuccessProbe struct{}

func NewSuccessProbe() *SuccessProbe {
	return &SuccessProbe{}
}

func (p *SuccessProbe) Name() string {
	return "successProbe"
}

func (p *SuccessProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	return healthcheck.NewCheckerProbeResult(true, "success")
}
