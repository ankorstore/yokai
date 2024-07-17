package healthcheck

import (
	"context"
	"fmt"
	"strings"

	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/worker"
)

// DefaultProbeName is the name of the worker probe.
const DefaultProbeName = "worker"

// WorkerProbe is a probe compatible with the [healthcheck] module.
//
// [healthcheck]: https://github.com/ankorstore/yokai/tree/main/healthcheck
type WorkerProbe struct {
	name string
	pool *worker.WorkerPool
}

// NewWorkerProbe returns a new [WorkerProbe].
func NewWorkerProbe(pool *worker.WorkerPool) *WorkerProbe {
	return &WorkerProbe{
		name: DefaultProbeName,
		pool: pool,
	}
}

// Name returns the name of the [WorkerProbe].
func (p *WorkerProbe) Name() string {
	return p.name
}

// SetName sets the name of the [WorkerProbe].
func (p *WorkerProbe) SetName(name string) *WorkerProbe {
	p.name = name

	return p
}

// Check returns a successful [healthcheck.CheckerProbeResult] if the worker pool executions are all in healthy status.
func (p *WorkerProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	success := true
	messages := []string{}

	for name, execution := range p.pool.Executions() {
		if execution.Status() == worker.Unknown || execution.Status() == worker.Error {
			success = false
		}

		messages = append(messages, fmt.Sprintf("%s: %s", name, execution.Status()))
	}

	return healthcheck.NewCheckerProbeResult(success, strings.Join(messages, ", "))
}
