package healthcheck

import (
	"context"
)

// CheckerProbe is the interface for the probes executed by the [Checker].
type CheckerProbe interface {
	Name() string
	Check(ctx context.Context) *CheckerProbeResult
}

// CheckerProbeResult is the result of a [CheckerProbe] execution.
type CheckerProbeResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// NewCheckerProbeResult returns a [CheckerProbeResult], with a probe execution status and feedback message.
func NewCheckerProbeResult(success bool, message string) *CheckerProbeResult {
	return &CheckerProbeResult{
		Success: success,
		Message: message,
	}
}
