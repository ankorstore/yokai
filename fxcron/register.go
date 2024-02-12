package fxcron

import (
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/fx"
)

// AsCronJob registers a cron job into Fx, with an optional list of [gocron.JobOption].
func AsCronJob(j any, expression string, options ...gocron.JobOption) fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				j,
				fx.As(new(CronJob)),
				fx.ResultTags(`group:"cron-jobs"`),
			),
		),
		fx.Supply(
			fx.Annotate(
				NewCronJobDefinition(GetReturnType(j), expression, options...),
				fx.As(new(CronJobDefinition)),
				fx.ResultTags(`group:"cron-jobs-definitions"`),
			),
		),
	)
}
