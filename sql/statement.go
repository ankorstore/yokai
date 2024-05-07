package sql

import (
	"context"
	"database/sql/driver"
)

var (
	_ driver.Stmt             = (*Statement)(nil)
	_ driver.StmtExecContext  = (*Statement)(nil)
	_ driver.StmtQueryContext = (*Statement)(nil)
)

// Statement is a SQL driver statement wrapping a driver.Stmt.
//
//nolint:containedctx
type Statement struct {
	base          driver.Stmt
	context       context.Context
	query         string
	configuration *Configuration
}

// NewStatement returns a new Statement.
//
//nolint:contextcheck
func NewStatement(base driver.Stmt, ctx context.Context, query string, configuration *Configuration) *Statement {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Statement{
		base:          base,
		context:       ctx,
		query:         query,
		configuration: configuration,
	}
}

// Close closes the Statement.
func (s *Statement) Close() error {
	return s.base.Close()
}

// NumInput returns the number of inputs of the Statement.
func (s *Statement) NumInput() int {
	return s.base.NumInput()
}

// Exec executes a statement and returns a driver.Result.
func (s *Statement) Exec(args []driver.Value) (driver.Result, error) {
	event := NewHookEvent(s.configuration.System(), StatementExecOperation, s.query, args)

	s.applyBeforeHooks(event)

	event.Start()
	//nolint:staticcheck
	res, err := s.base.Exec(args)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	if res != nil {
		lastInsertId, lastInsertIdErr := res.LastInsertId()
		if lastInsertIdErr == nil {
			event.SetLastInsertId(lastInsertId)
		}

		rowsAffected, rowsAffectedErr := res.RowsAffected()
		if rowsAffectedErr == nil {
			event.SetRowsAffected(rowsAffected)
		}
	}

	s.applyAfterHooks(event)

	return res, err
}

// ExecContext executes a statement for a context and returns a driver.Result.
func (s *Statement) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	s.context = ctx

	engine, ok := s.base.(driver.StmtExecContext)
	if !ok {
		return s.Exec(ConvertNamedValuesToValues(args))
	}

	event := NewHookEvent(s.configuration.System(), StatementExecContextOperation, s.query, args)

	s.applyBeforeHooks(event)

	event.Start()
	//nolint:contextcheck
	res, err := engine.ExecContext(s.context, args)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	if res != nil {
		lastInsertId, lastInsertIdErr := res.LastInsertId()
		if lastInsertIdErr == nil {
			event.SetLastInsertId(lastInsertId)
		}

		rowsAffected, rowsAffectedErr := res.RowsAffected()
		if rowsAffectedErr == nil {
			event.SetRowsAffected(rowsAffected)
		}
	}

	s.applyAfterHooks(event)

	return res, err
}

// Query executes a statement and returns a driver.Rows.
func (s *Statement) Query(args []driver.Value) (driver.Rows, error) {
	event := NewHookEvent(s.configuration.System(), StatementQueryOperation, s.query, args)

	s.applyBeforeHooks(event)

	event.Start()
	//nolint:staticcheck
	rows, err := s.base.Query(args)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	s.applyAfterHooks(event)

	return rows, err
}

// QueryContext executes a statement for a context and returns a driver.Rows.
func (s *Statement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	s.context = ctx

	engine, ok := s.base.(driver.StmtQueryContext)
	if !ok {
		return s.Query(ConvertNamedValuesToValues(args))
	}

	event := NewHookEvent(s.configuration.System(), StatementQueryContextOperation, s.query, args)

	s.applyBeforeHooks(event)

	event.Start()
	//nolint:contextcheck
	rows, err := engine.QueryContext(s.context, args)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	s.applyAfterHooks(event)

	return rows, err
}

func (s *Statement) applyBeforeHooks(event *HookEvent) {
	for _, h := range s.configuration.Hooks() {
		s.context = h.Before(s.context, event)
	}
}

func (s *Statement) applyAfterHooks(event *HookEvent) {
	for _, h := range s.configuration.Hooks() {
		h.After(s.context, event)
	}
}
