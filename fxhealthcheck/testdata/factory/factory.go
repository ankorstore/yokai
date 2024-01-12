package factory

import (
	"github.com/ankorstore/yokai/fxhealthcheck/testdata/probes"
	"github.com/ankorstore/yokai/healthcheck"
)

type TestCheckerFactory struct{}

func NewTestCheckerFactory() healthcheck.CheckerFactory {
	return &TestCheckerFactory{}
}

func (f *TestCheckerFactory) Create(options ...healthcheck.CheckerOption) (*healthcheck.Checker, error) {
	checker := healthcheck.NewChecker()
	checker.RegisterProbe(probes.NewFailureProbe(), healthcheck.Readiness)

	return checker, nil
}
