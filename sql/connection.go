package sql

import (
	"context"
	"database/sql/driver"
)

type Connection struct {
	base          driver.Conn
	configuration *Configuration
}

func NewConnection(base driver.Conn, configuration *Configuration) *Connection {
	return &Connection{
		base:          base,
		configuration: configuration,
	}
}

func (c *Connection) Exec(query string, args []driver.Value) (driver.Result, error) {
	engine, ok := c.base.(driver.Execer)
	if !ok {
		return nil, driver.ErrSkip
	}

	event := NewHookEvent(c.configuration.System(), ConnectionExecOperation, query, args)

	ctx := c.applyBeforeHooks(context.Background(), event)

	event.Start()
	res, err := engine.Exec(query, args)
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

	c.applyAfterHooks(ctx, event)

	return res, err
}

func (c *Connection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	engine, ok := c.base.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	event := NewHookEvent(c.configuration.System(), ConnectionExecContextOperation, query, args)

	ctx = c.applyBeforeHooks(ctx, event)

	event.Start()
	res, err := engine.ExecContext(ctx, query, args)
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

	c.applyAfterHooks(ctx, event)

	return res, err
}

func (c *Connection) Query(query string, args []driver.Value) (driver.Rows, error) {
	engine, ok := c.base.(driver.Queryer)
	if !ok {
		return nil, driver.ErrSkip
	}

	event := NewHookEvent(c.configuration.System(), ConnectionQueryOperation, query, args)

	ctx := c.applyBeforeHooks(context.Background(), event)

	event.Start()
	rows, err := engine.Query(query, args)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return rows, err
}

func (c *Connection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	engine, ok := c.base.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}

	event := NewHookEvent(c.configuration.System(), ConnectionQueryContextOperation, query, args)

	ctx = c.applyBeforeHooks(ctx, event)

	event.Start()
	rows, err := engine.QueryContext(ctx, query, args)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return rows, err
}

func (c *Connection) Prepare(query string) (driver.Stmt, error) {
	event := NewHookEvent(c.configuration.System(), ConnectionPrepareOperation, query, nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	event.Start()
	stmt, err := c.base.Prepare(query)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return NewStatement(stmt, nil, query, c.configuration), err
}

func (c *Connection) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	event := NewHookEvent(c.configuration.System(), ConnectionPrepareContextOperation, query, nil)

	ctx = c.applyBeforeHooks(ctx, event)

	if engine, ok := c.base.(driver.ConnPrepareContext); ok {
		event.Start()
		stmt, err := engine.PrepareContext(ctx, query)
		event.Stop()
		if err != nil {
			event.SetError(err)
		}

		c.applyAfterHooks(ctx, event)

		return NewStatement(stmt, ctx, query, c.configuration), err
	} else {
		event.Start()
		stmt, err := c.base.Prepare(query)
		event.Stop()
		if err != nil {
			event.SetError(err)
		}

		c.applyAfterHooks(ctx, event)

		return NewStatement(stmt, ctx, query, c.configuration), err
	}
}

func (c *Connection) Begin() (driver.Tx, error) {
	event := NewHookEvent(c.configuration.System(), ConnectionBeginOperation, "", nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	event.Start()
	tx, err := c.base.Begin()
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return NewTransaction(tx, ctx, c.configuration), err
}

func (c *Connection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	event := NewHookEvent(c.configuration.System(), ConnectionBeginTxOperation, "", nil)

	ctx = c.applyBeforeHooks(ctx, event)

	if engine, ok := c.base.(driver.ConnBeginTx); ok {
		event.Start()
		tx, err := engine.BeginTx(ctx, opts)
		event.Stop()
		if err != nil {
			event.SetError(err)
		}

		c.applyAfterHooks(ctx, event)

		return NewTransaction(tx, ctx, c.configuration), err
	} else {
		event.Start()
		tx, err := c.base.Begin()
		event.Stop()
		if err != nil {
			event.SetError(err)
		}

		c.applyAfterHooks(ctx, event)

		return NewTransaction(tx, ctx, c.configuration), err
	}
}

func (c *Connection) Ping(ctx context.Context) error {
	engine, ok := c.base.(driver.Pinger)
	if !ok {
		return driver.ErrSkip
	}

	event := NewHookEvent(c.configuration.System(), ConnectionPingOperation, "", nil)

	ctx = c.applyBeforeHooks(ctx, event)

	event.Start()
	err := engine.Ping(ctx)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return err
}

func (c *Connection) ResetSession(ctx context.Context) error {
	engine, ok := c.base.(driver.SessionResetter)
	if !ok {
		return driver.ErrSkip
	}

	event := NewHookEvent(c.configuration.System(), ConnectionResetSessionOperation, "", nil)

	ctx = c.applyBeforeHooks(context.Background(), event)

	event.Start()
	err := engine.ResetSession(ctx)
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return err

}

func (c *Connection) Close() error {
	event := NewHookEvent(c.configuration.System(), ConnectionCloseOperation, "", nil)

	ctx := c.applyBeforeHooks(context.Background(), event)

	event.Start()
	err := c.base.Close()
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	c.applyAfterHooks(ctx, event)

	return err
}

func (c *Connection) applyBeforeHooks(ctx context.Context, event *HookEvent) context.Context {
	for _, h := range c.configuration.Hooks() {
		ctx = h.Before(ctx, event)
	}

	return ctx
}

func (c *Connection) applyAfterHooks(ctx context.Context, event *HookEvent) {
	for _, h := range c.configuration.Hooks() {
		h.After(ctx, event)
	}
}
