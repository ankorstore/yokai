package sql

import (
	"context"
	"database/sql/driver"
)

//nolint:containedctx
type Transaction struct {
	base          driver.Tx
	context       context.Context
	configuration *Configuration
}

//nolint:contextcheck
func NewTransaction(base driver.Tx, ctx context.Context, configuration *Configuration) *Transaction {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Transaction{
		base:          base,
		context:       ctx,
		configuration: configuration,
	}
}

func (t *Transaction) Commit() error {
	event := NewHookEvent(t.configuration.System(), TransactionCommitOperation, "", nil)

	t.applyBeforeHooks(event)

	event.Start()
	err := t.base.Commit()
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	t.applyAfterHooks(event)

	return err
}

func (t *Transaction) Rollback() error {
	event := NewHookEvent(t.configuration.System(), TransactionRollbackOperation, "", nil)

	t.applyBeforeHooks(event)

	event.Start()
	err := t.base.Rollback()
	event.Stop()
	if err != nil {
		event.SetError(err)
	}

	t.applyAfterHooks(event)

	return err
}

func (t *Transaction) applyBeforeHooks(event *HookEvent) {
	for _, h := range t.configuration.Hooks() {
		t.context = h.Before(t.context, event)
	}
}

func (t *Transaction) applyAfterHooks(event *HookEvent) {
	for _, h := range t.configuration.Hooks() {
		h.After(t.context, event)
	}
}
