package worker

// WorkerPoolFactory is the interface for [WorkerPool] factories.
type WorkerPoolFactory interface {
	Create(options ...WorkerPoolOption) (*WorkerPool, error)
}

// DefaultWorkerPoolFactory is the default [WorkerPoolFactory] implementation.
type DefaultWorkerPoolFactory struct{}

// NewDefaultWorkerPoolFactory returns a [DefaultWorkerPoolFactory], implementing [WorkerPoolFactory].
func NewDefaultWorkerPoolFactory() WorkerPoolFactory {
	return &DefaultWorkerPoolFactory{}
}

// Create returns a new [WorkerPool], and accepts a list of [WorkerPoolOption].
// For example:
//
//	var pool, _ = worker.NewDefaultWorkerPoolFactory().Create()
//
// is equivalent to:
//
//	var pool, _ = worker.NewDefaultWorkerPoolFactory().Create(
//		worker.WithGenerator(uuid.NewDefaultUuidGenerator()), // generator
//		worker.WithMetrics(worker.NewWorkerMetrics("", "")),  // metrics
//		worker.WithGlobalMaxExecutionsAttempts(1),            // no retries
//		worker.WithGlobalDeferredStartThreshold(0),           // no deferred start
//	)
func (f *DefaultWorkerPoolFactory) Create(options ...WorkerPoolOption) (*WorkerPool, error) {
	return NewWorkerPool(options...), nil
}
