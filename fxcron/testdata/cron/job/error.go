package job

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxcron/testdata/cron/tracker"
	"go.opentelemetry.io/otel/attribute"
)

type ErrorCron struct {
	config  *config.Config
	tracker *tracker.CronExecutionTracker
}

func NewErrorCron(cfg *config.Config, trk *tracker.CronExecutionTracker) *ErrorCron {
	return &ErrorCron{
		config:  cfg,
		tracker: trk,
	}
}

func (c *ErrorCron) Name() string {
	return "error"
}

func (c *ErrorCron) Run(ctx context.Context) error {
	jobName := fxcron.CtxCronJobName(ctx)

	c.tracker.TrackJobExecution(jobName)
	e := c.tracker.JobExecutions(jobName)

	ctx, span := fxcron.CtxTracer(ctx).Start(ctx, "error cron job span")
	defer span.End()

	span.SetAttributes(attribute.Int("Run", e))

	fxcron.CtxLogger(ctx).Info().Int("run", e).Msgf("error cron job log from %s", c.config.AppName())

	return fmt.Errorf("error cron job error")
}
