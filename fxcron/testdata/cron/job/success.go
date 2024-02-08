package job

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxcron/testdata/cron/tracker"
	"go.opentelemetry.io/otel/attribute"
)

type SuccessCron struct {
	config  *config.Config
	tracker *tracker.CronExecutionTracker
}

func NewSuccessCron(cfg *config.Config, trk *tracker.CronExecutionTracker) *SuccessCron {
	return &SuccessCron{
		config:  cfg,
		tracker: trk,
	}
}

func (c *SuccessCron) Name() string {
	return "success"
}

func (c *SuccessCron) Run(ctx context.Context) error {
	jobName := fxcron.CtxCronJobName(ctx)

	c.tracker.TrackJobExecution(jobName)
	e := c.tracker.JobExecutions(jobName)

	ctx, span := fxcron.CtxTracer(ctx).Start(ctx, "success cron job span")
	defer span.End()

	span.SetAttributes(attribute.Int("Run", e))

	fxcron.CtxLogger(ctx).Info().Int("run", e).Msgf("success cron job log from %s", c.config.AppName())

	return nil
}
