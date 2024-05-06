package sql

import (
	"context"
	"fmt"
	"time"
)

// Hook is the interface for database hooks.
type Hook interface {
	Before(context.Context, *HookEvent) context.Context
	After(context.Context, *HookEvent)
}

// HookEvent is representing an event provided to a database Hook.
type HookEvent struct {
	system       System
	operation    Operation
	startedAt    time.Time
	stoppedAt    time.Time
	query        string
	arguments    any
	lastInsertId int64
	rowsAffected int64
	err          error
}

// NewHookEvent returns a new HookEvent.
func NewHookEvent(system System, operation Operation, query string, arguments interface{}) *HookEvent {
	return &HookEvent{
		system:    system,
		operation: operation,
		query:     query,
		arguments: arguments,
	}
}

// System returns the HookEvent System.
func (e *HookEvent) System() System {
	return e.system
}

// Operation returns the HookEvent Operation.
func (e *HookEvent) Operation() Operation {
	return e.operation
}

// Query returns the HookEvent query.
func (e *HookEvent) Query() string {
	return e.query
}

// Arguments returns the HookEvent query arguments.
func (e *HookEvent) Arguments() any {
	return e.arguments
}

// LastInsertId returns the HookEvent database last inserted id.
func (e *HookEvent) LastInsertId() int64 {
	return e.lastInsertId
}

// RowsAffected returns the HookEvent database affected rows.
func (e *HookEvent) RowsAffected() int64 {
	return e.rowsAffected
}

// Error returns the HookEvent error.
func (e *HookEvent) Error() error {
	return e.err
}

// SetLastInsertId sets the HookEvent database last inserted id.
func (e *HookEvent) SetLastInsertId(lastInsertId int64) *HookEvent {
	e.lastInsertId = lastInsertId

	return e
}

// SetRowsAffected sets the HookEvent database affected rows.
func (e *HookEvent) SetRowsAffected(rowsAffected int64) *HookEvent {
	e.rowsAffected = rowsAffected

	return e
}

// SetError sets the HookEvent error.
func (e *HookEvent) SetError(err error) *HookEvent {
	e.err = err

	return e
}

// Start records the HookEvent start time.
func (e *HookEvent) Start() *HookEvent {
	e.startedAt = time.Now()

	return e
}

// Stop records the HookEvent stop time.
func (e *HookEvent) Stop() *HookEvent {
	e.stoppedAt = time.Now()

	return e
}

// Latency returns the HookEvent latency (duration between start and end times).
func (e *HookEvent) Latency() (time.Duration, error) {
	if e.startedAt.IsZero() {
		return time.Duration(0), fmt.Errorf("event was not started")
	}

	if e.stoppedAt.IsZero() {
		return time.Duration(0), fmt.Errorf("event was not stopped")
	}

	return e.stoppedAt.Sub(e.startedAt), nil
}
