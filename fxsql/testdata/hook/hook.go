package hook

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/sql"
)

type TestHook struct {
	config *config.Config
}

func NewTestHook(config *config.Config) *TestHook {
	return &TestHook{
		config: config,
	}
}

func (h *TestHook) Before(ctx context.Context, event *sql.HookEvent) context.Context {
	log.CtxLogger(ctx).Info().Msgf("%s before %s", h.config.GetString("config.hook_name"), event.Operation())

	return ctx
}

func (h *TestHook) After(ctx context.Context, event *sql.HookEvent) {
	log.CtxLogger(ctx).Info().Msgf("%s after %s", h.config.GetString("config.hook_name"), event.Operation())
}
