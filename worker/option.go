package worker

import "github.com/ankorstore/yokai/generate/uuid"

const (
	DefaultDeferredStartThreshold = 0
	DefaultMaxExecutionsAttempts  = 1
	DefaultMetricsNamespace       = ""
	DefaultMetricsSubsystem       = ""
)

type PoolOptions struct {
	GlobalDeferredStartThreshold float64
	GlobalMaxExecutionsAttempts  int
	Metrics                      *WorkerMetrics
	Generator                    uuid.UuidGenerator
	Registrations                map[string]*WorkerRegistration
}

type WorkerPoolOption func(o *PoolOptions)

func DefaultWorkerPoolOptions() PoolOptions {
	return PoolOptions{
		GlobalDeferredStartThreshold: DefaultDeferredStartThreshold,
		GlobalMaxExecutionsAttempts:  DefaultMaxExecutionsAttempts,
		Metrics:                      NewWorkerMetrics(DefaultMetricsNamespace, DefaultMetricsSubsystem),
		Generator:                    uuid.NewDefaultUuidGenerator(),
		Registrations:                make(map[string]*WorkerRegistration),
	}
}

func WithGlobalDeferredStartThreshold(threshold float64) WorkerPoolOption {
	return func(o *PoolOptions) {
		o.GlobalDeferredStartThreshold = threshold
	}
}

func WithGlobalMaxExecutionsAttempts(max int) WorkerPoolOption {
	return func(o *PoolOptions) {
		o.GlobalMaxExecutionsAttempts = max
	}
}

func WithMetrics(metrics *WorkerMetrics) WorkerPoolOption {
	return func(o *PoolOptions) {
		o.Metrics = metrics
	}
}

func WithGenerator(generator uuid.UuidGenerator) WorkerPoolOption {
	return func(o *PoolOptions) {
		o.Generator = generator
	}
}

func WithWorker(worker Worker, options ...WorkerExecutionOption) WorkerPoolOption {
	return func(o *PoolOptions) {
		o.Registrations[worker.Name()] = NewWorkerRegistration(worker, options...)
	}
}

type ExecutionOptions struct {
	DeferredStartThreshold float64
	MaxExecutionsAttempts  int
}

type WorkerExecutionOption func(o *ExecutionOptions)

func DefaultWorkerExecutionOptions() ExecutionOptions {
	return ExecutionOptions{
		DeferredStartThreshold: DefaultDeferredStartThreshold,
		MaxExecutionsAttempts:  DefaultMaxExecutionsAttempts,
	}
}

func WithDeferredStartThreshold(t float64) WorkerExecutionOption {
	return func(o *ExecutionOptions) {
		o.DeferredStartThreshold = t
	}
}

func WithMaxExecutionsAttempts(l int) WorkerExecutionOption {
	return func(o *ExecutionOptions) {
		o.MaxExecutionsAttempts = l
	}
}
