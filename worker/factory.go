package worker

type WorkerPoolFactory interface {
	Create(options ...WorkerPoolOption) (*WorkerPool, error)
}

type DefaultWorkerPoolFactory struct{}

func NewDefaultWorkerPoolFactory() WorkerPoolFactory {
	return &DefaultWorkerPoolFactory{}
}

func (f *DefaultWorkerPoolFactory) Create(options ...WorkerPoolOption) (*WorkerPool, error) {
	return NewWorkerPool(options...), nil
}
