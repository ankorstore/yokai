package worker_test

import "github.com/ankorstore/yokai/worker"

type TestMiddleware struct {
	Func worker.MiddlewareFunc
}

func (m *TestMiddleware) Name() string {
	return "TestMiddleware"
}

func (m *TestMiddleware) Handle() worker.MiddlewareFunc {
	return m.Func
}
