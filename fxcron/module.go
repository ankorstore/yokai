package fxcron

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/go-co-op/gocron/v2"
	"github.com/prometheus/client_golang/prometheus"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const (
	ModuleName                           = "cron"
	LogRecordFieldCronJobName            = "cronJob"
	LogRecordFieldCronJobExecutionId     = "cronJobExecutionID"
	TraceSpanAttributeCronJobName        = "CronJob"
	TraceSpanAttributeCronJobExecutionId = "CronJobExecutionID"
)

// FxCronModule is the [Fx] cron module.
//
// [Fx]: https://github.com/uber-go/fx
var FxCronModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewDefaultCronSchedulerFactory,
		NewFxCronJobRegistry,
		NewFxCron,
		fx.Annotate(
			NewFxCronModuleInfo,
			fx.As(new(interface{})),
			fx.ResultTags(`group:"core-module-infos"`),
		),
	),
)

// FxCronParam allows injection of the required dependencies in [NewFxCron].
type FxCronParam struct {
	fx.In
	LifeCycle       fx.Lifecycle
	Generator       uuid.UuidGenerator
	TracerProvider  oteltrace.TracerProvider
	Factory         CronSchedulerFactory
	Config          *config.Config
	Registry        *CronJobRegistry
	Logger          *log.Logger
	MetricsRegistry *prometheus.Registry
}

// NewFxCron returns a new [gocron.Scheduler].
//
//nolint:cyclop,gocognit
func NewFxCron(p FxCronParam) (gocron.Scheduler, error) {
	appDebug := p.Config.AppDebug()

	// logger
	cronLogger := log.FromZerolog(p.Logger.ToZerolog().With().Str("system", ModuleName).Logger())

	// tracer provider
	tracerProvider := AnnotateTracerProvider(p.TracerProvider)

	// scheduler
	cronSchedulerOptions, err := buildSchedulerOptions(p.Config)
	if err != nil {
		p.Logger.Error().Err(err).Msg("cron scheduler options creation error")

		return nil, err
	}

	cronScheduler, err := p.Factory.Create(cronSchedulerOptions...)
	if err != nil {
		p.Logger.Error().Err(err).Msg("cron scheduler creation error")

		return nil, err
	}

	// jobs logs
	cronJobLogExecution := p.Config.GetBool("modules.cron.log.enabled") || appDebug
	cronJobLogExclusions := p.Config.GetStringSlice("modules.cron.log.exclude")

	// jobs traces
	cronJobTraceExecution := p.Config.GetBool("modules.cron.trace.enabled")
	cronJobTraceExclusions := p.Config.GetStringSlice("modules.cron.trace.exclude")

	// jobs metrics
	cronJobMetricsNamespace := p.Config.GetString("modules.cron.metrics.collect.namespace")
	if cronJobMetricsNamespace == "" {
		cronJobMetricsNamespace = p.Config.AppName()
	}

	cronJobMetricsSubsystem := p.Config.GetString("modules.cron.metrics.collect.subsystem")
	if cronJobMetricsSubsystem == "" {
		cronJobMetricsSubsystem = ModuleName
	}

	var cronJobMetrics *CronJobMetrics
	if cronJobMetricsBuckets := p.Config.GetString("modules.cron.metrics.buckets"); cronJobMetricsBuckets != "" {
		var buckets []float64

		for _, s := range strings.Split(strings.ReplaceAll(cronJobMetricsBuckets, " ", ""), ",") {
			f, err := strconv.ParseFloat(s, 64)
			if err == nil {
				buckets = append(buckets, f)
			}
		}

		cronJobMetrics = NewCronJobMetricsWithBuckets(cronJobMetricsNamespace, cronJobMetricsSubsystem, buckets)
	} else {
		cronJobMetrics = NewCronJobMetrics(cronJobMetricsNamespace, cronJobMetricsSubsystem)
	}

	if p.Config.GetBool("modules.cron.metrics.collect.enabled") {
		err = cronJobMetrics.Register(p.MetricsRegistry)
		if err != nil {
			p.Logger.Error().Err(err).Msg("cron scheduler metrics registration error")

			return nil, err
		}
	}

	// jobs registration
	cronJobs, err := p.Registry.ResolveCronJobs()
	if err != nil {
		p.Logger.Error().Err(err).Msg("cron jobs resolution error")

		return nil, err
	}

	for _, cronJob := range cronJobs {
		// var scoping
		currentCronJob := cronJob

		currentCronJobName := currentCronJob.Implementation().Name()
		currentJobOptions := append(currentCronJob.Options(), gocron.WithName(currentCronJobName))
		currentCronJobLogExecution := !Contains(cronJobLogExclusions, currentCronJobName)
		currentCronJobTraceExecution := !Contains(cronJobTraceExclusions, currentCronJobName)

		_, err = cronScheduler.NewJob(
			gocron.CronJob(
				currentCronJob.Expression(),
				p.Config.GetBool("modules.cron.scheduler.seconds"),
			),
			gocron.NewTask(
				func() {
					currentCronJobExecutionId := p.Generator.Generate()

					currentCronJobCtx := context.WithValue(context.Background(), CtxCronJobNameKey{}, currentCronJobName)
					currentCronJobCtx = context.WithValue(currentCronJobCtx, CtxCronJobExecutionIdKey{}, currentCronJobExecutionId)
					currentCronJobCtx = context.WithValue(currentCronJobCtx, trace.CtxKey{}, tracerProvider)

					var currentCronJobExecutionTraceSpan oteltrace.Span
					if cronJobTraceExecution && currentCronJobTraceExecution {
						currentCronJobCtx, currentCronJobExecutionTraceSpan = tracerProvider.
							Tracer(ModuleName).
							Start(currentCronJobCtx, fmt.Sprintf("%s %s", ModuleName, currentCronJobName))
					}

					currentCronJobLogger := log.FromZerolog(
						cronLogger.
							ToZerolog().
							With().
							Str(LogRecordFieldCronJobName, currentCronJobName).
							Str(LogRecordFieldCronJobExecutionId, currentCronJobExecutionId).
							Logger(),
					)

					currentCronJobCtx = currentCronJobLogger.WithContext(currentCronJobCtx)

					defer func(s oteltrace.Span, t time.Time) {
						if cronJobTraceExecution && currentCronJobTraceExecution && s != nil {
							s.End()
						}

						cronJobMetrics.ObserveCronJobExecutionDuration(currentCronJobName, time.Since(t).Seconds())

						if r := recover(); r != nil {
							cronJobMetrics.IncrementCronJobExecutionError(currentCronJobName)
							currentCronJobLogger.Error().Str("panic", fmt.Sprintf("%v", r)).Msg("job execution panic")
						}
					}(currentCronJobExecutionTraceSpan, time.Now())

					if cronJobLogExecution && currentCronJobLogExecution {
						currentCronJobLogger.Info().Msg("job execution start")
					}

					runErr := currentCronJob.Implementation().Run(currentCronJobCtx)

					if runErr != nil {
						cronJobMetrics.IncrementCronJobExecutionError(currentCronJobName)
						currentCronJobLogger.Error().Err(runErr).Msg("job execution error")
					} else {
						cronJobMetrics.IncrementCronJobExecutionSuccess(currentCronJobName)

						if cronJobLogExecution && currentCronJobLogExecution {
							currentCronJobLogger.Info().Msg("job execution success")
						}
					}
				},
			),
			currentJobOptions...,
		)

		if err != nil {
			cronLogger.Error().Err(err).Msgf("job registration error for job %s with %s", currentCronJobName, currentCronJob.Expression())

			return nil, err
		} else {
			cronLogger.Debug().Msgf("job registration success for job %s with %s", currentCronJobName, currentCronJob.Expression())
		}
	}

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			cronLogger.Debug().Msg("starting cron scheduler")

			cronScheduler.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			cronLogger.Debug().Msg("stopping cron scheduler")

			return cronScheduler.Shutdown()
		},
	})

	return cronScheduler, nil
}

//nolint:cyclop
func buildSchedulerOptions(cfg *config.Config) ([]gocron.SchedulerOption, error) {
	var options []gocron.SchedulerOption

	// location, default local
	if cfgLocation := cfg.GetString("modules.cron.scheduler.location"); cfgLocation != "" {
		location, err := time.LoadLocation(cfgLocation)
		if err != nil {
			return nil, err
		}

		options = append(options, gocron.WithLocation(location))
	}

	// concurrency
	if cfg.GetBool("modules.cron.scheduler.concurrency.limit.enabled") {
		var mode gocron.LimitMode
		if cfg.GetString("modules.cron.scheduler.concurrency.limit.mode") == "reschedule" {
			mode = gocron.LimitModeReschedule
		} else {
			mode = gocron.LimitModeWait
		}

		options = append(options, gocron.WithLimitConcurrentJobs(cfg.GetUint("modules.cron.scheduler.concurrency.limit.max"), mode))
	}

	// stop timeout, default 10s
	if cfgStopTimeout := cfg.GetString("modules.cron.scheduler.stop.timeout"); cfgStopTimeout != "" {
		stopTimeout, err := time.ParseDuration(cfgStopTimeout)
		if err != nil {
			return nil, err
		}

		options = append(options, gocron.WithStopTimeout(stopTimeout))
	}

	// jobs global options
	var jobsOptions []gocron.JobOption

	// jobs execution start
	if cfg.GetBool("modules.cron.jobs.execution.start.immediately") {
		jobsOptions = append(jobsOptions, gocron.WithStartAt(gocron.WithStartImmediately()))
	} else if cfgJobsStartAt := cfg.GetString("modules.cron.jobs.execution.start.at"); cfgJobsStartAt != "" {
		jobsStartAt, err := time.Parse(time.RFC3339, cfgJobsStartAt)
		if err != nil {
			return nil, err
		}

		jobsOptions = append(jobsOptions, gocron.WithStartAt(gocron.WithStartDateTime(jobsStartAt)))
	}

	// jobs execution limit
	if cfg.GetBool("modules.cron.jobs.execution.limit.enabled") {
		jobsOptions = append(jobsOptions, gocron.WithLimitedRuns(cfg.GetUint("modules.cron.jobs.execution.limit.max")))
	}

	// jobs execution mode
	if cfg.GetBool("modules.cron.jobs.singleton.enabled") {
		var mode gocron.LimitMode
		if cfg.GetString("modules.cron.jobs.singleton.mode") == "reschedule" {
			mode = gocron.LimitModeReschedule
		} else {
			mode = gocron.LimitModeWait
		}
		jobsOptions = append(jobsOptions, gocron.WithSingletonMode(mode))
	}

	if len(jobsOptions) > 0 {
		options = append(options, gocron.WithGlobalJobOptions(jobsOptions...))
	}

	return options, nil
}
