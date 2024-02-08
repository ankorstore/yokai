package job

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxcron/testdata/cron/tracker"
)

type PanicCron struct {
	config  *config.Config
	tracker *tracker.CronExecutionTracker
}

func NewPanicCron(cfg *config.Config, trk *tracker.CronExecutionTracker) *PanicCron {
	return &PanicCron{
		config:  cfg,
		tracker: trk,
	}
}

func (c *PanicCron) Name() string {
	return "panic"
}

func (c *PanicCron) Run(ctx context.Context) error {
	jobName := fxcron.CtxCronJobName(ctx)

	c.tracker.TrackJobExecution(jobName)
	e := c.tracker.JobExecutions(jobName)

	fxcron.CtxLogger(ctx).Info().Int("run", e).Msgf("panic cron job log from %s", c.config.AppName())

	panic("panic cron job panic")
}
