package probes

import (
	"context"

	"github.com/ankorstore/yokai/healthcheck"
)

type TestProbe struct{}

func NewTestProbe() *TestProbe {
	return &TestProbe{}
}

func (p *TestProbe) Name() string {
	return "testProbe"
}

func (p *TestProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	return healthcheck.NewCheckerProbeResult(false, "test")
}
