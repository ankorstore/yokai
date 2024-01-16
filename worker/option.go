package worker

import "github.com/ankorstore/yokai/generate/uuid"

const (
	DefaultDeferredStartThreshold = 0
	DefaultMaxExecutionsAttempts  = 1
	DefaultMetricsNamespace       = ""
	DefaultMetricsSubsystem       = ""
)

// PoolOptions are options for the [WorkerPoolFactory] implementations.
type PoolOptions struct {
	GlobalDeferredStartThreshold float64
	GlobalMaxExecutionsAttempts  int
	Metrics                      *WorkerMetrics
	Generator                    uuid.UuidGenerator
	Registrations                map[string]*WorkerRegistration
}

// DefaultWorkerPoolOptions are the default options used in the [DefaultWorkerPoolFactory].
func DefaultWorkerPoolOptions() PoolOptions {
	return PoolOptions{
		GlobalDeferredStartThreshold: DefaultDeferredStartThreshold,
		GlobalMaxExecutionsAttempts:  DefaultMaxExecutionsAttempts,
		Metrics:                      NewWorkerMetrics(DefaultMetricsNamespace, DefaultMetricsSubsystem),
		Generator:                    uuid.NewDefaultUuidGenerator(),
		Registrations:                make(map[string]*WorkerRegistration),
	}
}

// WorkerPoolOption are functional options for the [WorkerPoolFactory] implementations.
type WorkerPoolOption func(o *PoolOptions)

// WithGlobalDeferredStartThreshold is used to specify the global workers deferred start threshold, in seconds.
func WithGlobalDeferredStartThreshold(threshold float64) WorkerPoolOption {
	return func(o *PoolOptions) {
		o.GlobalDeferredStartThreshold = threshold
	}
}

// WithGlobalMaxExecutionsAttempts is used to specify the global workers max execution attempts.
func WithGlobalMaxExecutionsAttempts(max int) WorkerPoolOption {
	return func(o *PoolOptions) {
		o.GlobalMaxExecutionsAttempts = max
	}
}

// WithMetrics is used to specify the [WorkerMetrics] to use by the [WorkerPool].
func WithMetrics(metrics *WorkerMetrics) WorkerPoolOption {
	return func(o *PoolOptions) {
		o.Metrics = metrics
	}
}

// WithGenerator is used to specify the [uuid.UuidGenerator] to use by the [WorkerPool].
func WithGenerator(generator uuid.UuidGenerator) WorkerPoolOption {
	return func(o *PoolOptions) {
		o.Generator = generator
	}
}

// WithWorker is used to register a [Worker] in the [WorkerPool], with an optional list of [WorkerPoolOption].
func WithWorker(worker Worker, options ...WorkerExecutionOption) WorkerPoolOption {
	return func(o *PoolOptions) {
		o.Registrations[worker.Name()] = NewWorkerRegistration(worker, options...)
	}
}

// ExecutionOptions are options for the [Worker] executions.
type ExecutionOptions struct {
	DeferredStartThreshold float64
	MaxExecutionsAttempts  int
}

// DefaultWorkerExecutionOptions are the default options for the [Worker] executions.
func DefaultWorkerExecutionOptions() ExecutionOptions {
	return ExecutionOptions{
		DeferredStartThreshold: DefaultDeferredStartThreshold,
		MaxExecutionsAttempts:  DefaultMaxExecutionsAttempts,
	}
}

// WorkerExecutionOption are functional options for the [Worker] executions.
type WorkerExecutionOption func(o *ExecutionOptions)

// WithDeferredStartThreshold is used to specify the worker deferred start threshold, in seconds.
func WithDeferredStartThreshold(t float64) WorkerExecutionOption {
	return func(o *ExecutionOptions) {
		o.DeferredStartThreshold = t
	}
}

// WithMaxExecutionsAttempts is used to specify the worker max execution attempts.
func WithMaxExecutionsAttempts(l int) WorkerExecutionOption {
	return func(o *ExecutionOptions) {
		o.MaxExecutionsAttempts = l
	}
}
