package probes

import (
	"context"

	"github.com/ankorstore/yokai/healthcheck"
)

type FailureProbe struct{}

func NewFailureProbe() *FailureProbe {
	return &FailureProbe{}
}

func (p *FailureProbe) Name() string {
	return "failureProbe"
}

func (p *FailureProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	return healthcheck.NewCheckerProbeResult(false, "some failure")
}
