package factory

import (
	"github.com/ankorstore/yokai/worker"
)

type TestWorkerPoolFactory struct{}

func NewTestWorkerPoolFactory() worker.WorkerPoolFactory {
	return &TestWorkerPoolFactory{}
}

func (f *TestWorkerPoolFactory) Create(options ...worker.WorkerPoolOption) (*worker.WorkerPool, error) {
	return worker.NewWorkerPool(worker.WithGlobalMaxExecutionsAttempts(99)), nil
}
