package fxworker_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxworker"
	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/stretchr/testify/assert"
)

// TestMiddleware is a simple middleware implementation for testing.
type TestMiddleware struct{}

func (m *TestMiddleware) Name() string {
	return "TestMiddleware"
}

func (m *TestMiddleware) Handle() worker.MiddlewareFunc {
	return func(next worker.HandlerFunc) worker.HandlerFunc {
		return func(ctx context.Context) error {
			// Simple pass-through middleware
			return next(ctx)
		}
	}
}

// NewTestMiddleware returns a new TestMiddleware.
func NewTestMiddleware() worker.Middleware {
	return &TestMiddleware{}
}

// simpleMiddlewareFunc is a concrete middleware function for testing.
func simpleMiddlewareFunc(next worker.HandlerFunc) worker.HandlerFunc {
	return func(ctx context.Context) error {
		// Simple pass-through middleware
		return next(ctx)
	}
}

func TestAsWorker(t *testing.T) {
	t.Parallel()

	// Test with execution options
	result1 := fxworker.AsWorker(workers.NewClassicWorker, worker.WithMaxExecutionsAttempts(2))
	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result1))

	// Test with concrete middleware
	result2 := fxworker.AsWorker(workers.NewClassicWorker, worker.MiddlewareFunc(simpleMiddlewareFunc))
	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result2))

	// Test with middleware constructor
	result3 := fxworker.AsWorker(workers.NewClassicWorker, NewTestMiddleware)
	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result3))

	// Test with mixed options
	result4 := fxworker.AsWorker(
		workers.NewClassicWorker,
		worker.WithMaxExecutionsAttempts(2),
		worker.MiddlewareFunc(simpleMiddlewareFunc),
		NewTestMiddleware,
	)
	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result4))
}
