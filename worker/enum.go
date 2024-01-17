package worker

// WorkerStatus is an enum for the possible statuses of a workers.
type WorkerStatus int

const (
	Unknown WorkerStatus = iota
	Deferred
	Running
	Success
	Error
)

// String returns a string representation of the [WorkerStatus].
//
//nolint:exhaustive
func (s WorkerStatus) String() string {
	switch s {
	case Deferred:
		return "deferred"
	case Running:
		return "running"
	case Success:
		return "success"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}
