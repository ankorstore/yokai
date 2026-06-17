package fxtrace

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestBoundedContext_NoParentDeadline_UsesCap(t *testing.T) {
	ctx, cancel := boundedContext(context.Background(), 100*time.Millisecond)
	defer cancel()

	deadline, ok := ctx.Deadline()
	assert.True(t, ok, "child must have a deadline derived from the cap")
	assert.WithinDuration(t, time.Now().Add(100*time.Millisecond), deadline, 50*time.Millisecond)
}

func TestBoundedContext_ParentDeadlineSmallerThanCap_UsesHalfRemaining(t *testing.T) {
	// Parent deadline is 200ms away; cap is 5s. half-remaining (~100ms) wins.
	parent, parentCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer parentCancel()

	ctx, cancel := boundedContext(parent, 5*time.Second)
	defer cancel()

	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	remaining := time.Until(deadline)
	assert.Less(t, remaining, 150*time.Millisecond, "child deadline must be tighter than parent (half of remaining)")
	assert.Greater(t, remaining, 0*time.Millisecond)
}

func TestBoundedContext_ParentDeadlineLargerThanCap_UsesCap(t *testing.T) {
	parent, parentCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer parentCancel()

	ctx, cancel := boundedContext(parent, 100*time.Millisecond)
	defer cancel()

	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.WithinDuration(t, time.Now().Add(100*time.Millisecond), deadline, 50*time.Millisecond)
}

func TestBoundedContext_ExpiredParent_FallsBackToCancel(t *testing.T) {
	// Parent is already past its deadline => half-remaining <= 0 path: we keep
	// the cap (the inner `half > 0` guard rejects negative). Make cap also
	// non-positive to exercise the `timeout <= 0 -> WithCancel` branch.
	parent, parentCancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer parentCancel()

	ctx, cancel := boundedContext(parent, 0)
	defer cancel()

	// With cap=0 and an already-expired parent, we expect no fresh deadline
	// from boundedContext itself — the parent's expired deadline still shows
	// through (Go contract: child inherits parent's deadline when it's tighter).
	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, deadline.Before(time.Now()))
}

func TestBestEffortStop_SwallowsErrorAndLogs(t *testing.T) {
	logger := log.FromZerolog(zerolog.Nop())

	var called bool
	bestEffortStop(context.Background(), "boom", func(context.Context) error {
		called = true

		return errors.New("simulated failure")
	}, logger)

	assert.True(t, called, "fn must be invoked")
	// No panic, no propagation — the swallow is the test.
}

func TestBestEffortStop_RespectsBoundedTimeout(t *testing.T) {
	logger := log.FromZerolog(zerolog.Nop())

	var observedDeadline time.Time
	bestEffortStop(context.Background(), "deadline-probe", func(ctx context.Context) error {
		d, ok := ctx.Deadline()
		assert.True(t, ok, "bestEffortStop must pass a bounded context")
		observedDeadline = d

		return nil
	}, logger)

	assert.WithinDuration(t, time.Now().Add(shutdownCap), observedDeadline, shutdownCap)
}
