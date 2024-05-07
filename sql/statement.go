package sql

import (
	"context"
	"database/sql/driver"
)

//nolint:containedctx
type Statement struct {
	base          driver.Stmt
	context       context.Context
	query         string
	configuration *Configuration
}

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

func (s *Statement) Close() error {
	return s.base.Close()
}

func (s *Statement) NumInput() int {
	return s.base.NumInput()
}

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
