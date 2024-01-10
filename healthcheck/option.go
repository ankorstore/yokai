package healthcheck

// Options are options for the [CheckerFactory] implementations.
type Options struct {
	Registrations map[string]*CheckerProbeRegistration
}

// DefaultCheckerOptions are the default options used in the [DefaultCheckerFactory].
func DefaultCheckerOptions() Options {
	return Options{
		Registrations: map[string]*CheckerProbeRegistration{},
	}
}

// CheckerOption are functional options for the [CheckerFactory] implementations.
type CheckerOption func(o *Options)

// WithProbe is used to register a [CheckerProbe] for an optional list of [ProbeKind].
// If no [ProbeKind] was provided, the [CheckerProbe] will be registered for all kinds.
func WithProbe(probe CheckerProbe, kinds ...ProbeKind) CheckerOption {
	return func(o *Options) {
		if len(kinds) == 0 {
			kinds = []ProbeKind{Startup, Liveness, Readiness}
		}

		if _, ok := o.Registrations[probe.Name()]; ok {
			o.Registrations[probe.Name()].kinds = kinds
		} else {
			o.Registrations[probe.Name()] = NewCheckerProbeRegistration(probe, kinds...)
		}
	}
}
