package healthcheck

// ProbeKind is an enum for the supported kind of checks.
type ProbeKind int

const (
	Startup ProbeKind = iota
	Liveness
	Readiness
)

// String returns a string representation of the [ProbeKind].
//
//nolint:exhaustive
func (k ProbeKind) String() string {
	switch k {
	case Liveness:
		return "liveness"
	case Readiness:
		return "readiness"
	default:
		return "startup"
	}
}
