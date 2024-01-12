package fxhttpserver

// MiddlewareKind is an enum for the middleware kinds (global, pre, post).
type MiddlewareKind int

const (
	GlobalUse MiddlewareKind = iota
	GlobalPre
	Attached
)

// String returns a string representation of a [MiddlewareKind].
func (k MiddlewareKind) String() string {
	switch k {
	case GlobalUse:
		return "global-use"
	case GlobalPre:
		return "global-pre"
	case Attached:
		return "attached"
	default:
		return "global-use"
	}
}
