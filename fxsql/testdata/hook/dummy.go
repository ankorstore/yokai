package hook

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/sql"
)

type DummyHook struct{}

func NewDummyHook() *DummyHook {
	return &DummyHook{}
}

func (h *DummyHook) Before(ctx context.Context, event *sql.HookEvent) context.Context {
	log.CtxLogger(ctx).Info().Msgf("DummyHook: before %s", event.Operation())

	return ctx
}

func (h *DummyHook) After(ctx context.Context, event *sql.HookEvent) {
	log.CtxLogger(ctx).Info().Msgf("DummyHook: after %s", event.Operation())
}
