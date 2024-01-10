package healthcheck

// CheckerFactory is the interface for [Checker] factories.
type CheckerFactory interface {
	Create(options ...CheckerOption) (*Checker, error)
}

// DefaultCheckerFactory is the default [CheckerFactory] implementation.
type DefaultCheckerFactory struct{}

// NewDefaultCheckerFactory returns a [DefaultCheckerFactory], implementing [CheckerFactory].
func NewDefaultCheckerFactory() CheckerFactory {
	return &DefaultCheckerFactory{}
}

// Create returns a new [Checker], and accepts a list of [CheckerOption].
// For example:
//
//	checker, _ := healthcheck.NewDefaultCheckerFactory().Create(
//		healthcheck.WithProbe(NewSomeProbe()),                        // registers for startup, readiness and liveness
//		healthcheck.WithProbe(NewOtherProbe(), healthcheck.Liveness), // registers for liveness  only
//	)
func (f *DefaultCheckerFactory) Create(options ...CheckerOption) (*Checker, error) {
	appliedOpts := DefaultCheckerOptions()
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	checker := NewChecker()

	for _, registration := range appliedOpts.Registrations {
		checker.RegisterProbe(registration.Probe(), registration.Kinds()...)
	}

	return checker, nil
}
