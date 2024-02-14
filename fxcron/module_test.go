package fxcron_test

import (
	"strings"
	"testing"
	"time"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxcron/testdata/cron/job"
	"github.com/ankorstore/yokai/fxcron/testdata/cron/tracker"
	"github.com/ankorstore/yokai/fxcron/testdata/factory"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxgenerate/fxgeneratetest/uuid"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/go-co-op/gocron/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

//nolint:maintidx
func TestModule(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_ENV", "test")
	t.Setenv("CRON_START_IMMEDIATELY", "true")
	t.Setenv("CRON_METRICS_BUCKETS", "1,10,100")
	t.Setenv("CRON_METRICS_NAMESPACE", "foo")
	t.Setenv("CRON_METRICS_SUBSYSTEM", "bar")

	var cronTracker *tracker.CronExecutionTracker
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	// limits to avoid flaky tests
	expectedSuccessRuns := 3
	expectedErrorRuns := 2
	expectedPanicRuns := 1

	// deterministic test cron job execution id
	cronJobExecutionId := "testCronJobExecutionID"

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxcron.FxCronModule,
		// deterministic test cron job execution id
		fx.Provide(
			fx.Annotate(
				func() string {
					return cronJobExecutionId
				},
				fx.ResultTags(`name:"generate-test-uuid-value"`),
			),
		),
		fx.Options(
			// cron jobs execution tracker
			fx.Provide(tracker.NewCronExecutionTracker),
			// cron jobs registration
			fxcron.AsCronJob(job.NewSuccessCron, `*/1 * * * * *`, gocron.WithLimitedRuns(uint(expectedSuccessRuns))),
			fxcron.AsCronJob(job.NewErrorCron, `*/1 * * * * *`, gocron.WithLimitedRuns(uint(expectedErrorRuns))),
			fxcron.AsCronJob(job.NewPanicCron, `*/1 * * * * *`, gocron.WithLimitedRuns(uint(expectedPanicRuns))),
		),
		// deterministic generator for cron job execution id
		fx.Decorate(uuid.NewFxTestUuidGeneratorFactory),
		// extraction
		fx.Populate(&cronTracker, &logBuffer, &traceExporter, &metricsRegistry),
		// invoke scheduler
		fx.Invoke(func(scheduler gocron.Scheduler) {}),
	).RequireStart()

	// 5 seconds for cron jobs to run
	time.Sleep(4 * time.Second)

	app.RequireStop()

	// success cron assertions
	assert.Equal(t, expectedSuccessRuns, cronTracker.JobExecutions("success"))

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":              "info",
		"service":            "test",
		"system":             "cron",
		"cronJob":            "success",
		"cronJobExecutionID": cronJobExecutionId,
		"message":            "job execution start",
	})

	for i := 1; i <= expectedSuccessRuns; i++ {
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":              "info",
			"service":            "test",
			"system":             "cron",
			"cronJob":            "success",
			"cronJobExecutionID": cronJobExecutionId,
			"message":            "success cron job log from test",
			"run":                i,
		})
	}

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":              "info",
		"service":            "test",
		"system":             "cron",
		"cronJob":            "success",
		"cronJobExecutionID": cronJobExecutionId,
		"message":            "success cron job log from test",
		"run":                expectedSuccessRuns + 1,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":              "info",
		"service":            "test",
		"system":             "cron",
		"cronJob":            "success",
		"cronJobExecutionID": cronJobExecutionId,
		"message":            "job execution success",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"cron success",
		attribute.String("CronJob", "success"),
		attribute.String("CronJobExecutionID", cronJobExecutionId),
	)

	for i := 1; i <= expectedSuccessRuns; i++ {
		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"success cron job span",
			attribute.String("CronJob", "success"),
			attribute.String("CronJobExecutionID", cronJobExecutionId),
			attribute.Int("Run", i),
		)
	}

	tracetest.AssertHasNotTraceSpan(
		t,
		traceExporter,
		"cron success",
		attribute.String("CronJob", "success"),
		attribute.String("CronJobExecutionID", cronJobExecutionId),
		attribute.Int("Run", expectedSuccessRuns+1),
	)

	// error cron assertions
	assert.Equal(t, expectedErrorRuns, cronTracker.JobExecutions("error"))

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":              "info",
		"service":            "test",
		"system":             "cron",
		"cronJob":            "error",
		"cronJobExecutionID": cronJobExecutionId,
		"message":            "job execution start",
	})

	for i := 1; i <= expectedErrorRuns; i++ {
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":              "info",
			"service":            "test",
			"system":             "cron",
			"cronJob":            "error",
			"cronJobExecutionID": cronJobExecutionId,
			"message":            "error cron job log from test",
			"run":                i,
		})
	}

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":              "info",
		"service":            "test",
		"system":             "cron",
		"cronJob":            "error",
		"cronJobExecutionID": cronJobExecutionId,
		"message":            "error cron job log from test",
		"run":                expectedErrorRuns + 1,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":              "error",
		"service":            "test",
		"system":             "cron",
		"cronJob":            "error",
		"cronJobExecutionID": cronJobExecutionId,
		"message":            "job execution error",
		"error":              "error cron job error",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"cron error",
		attribute.String("CronJob", "error"),
		attribute.String("CronJobExecutionID", cronJobExecutionId),
	)

	for i := 1; i <= expectedErrorRuns; i++ {
		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"error cron job span",
			attribute.String("CronJob", "error"),
			attribute.String("CronJobExecutionID", cronJobExecutionId),
			attribute.Int("Run", i),
		)
	}

	tracetest.AssertHasNotTraceSpan(
		t,
		traceExporter,
		"cron error",
		attribute.String("CronJob", "error"),
		attribute.String("CronJobExecutionID", cronJobExecutionId),
		attribute.Int("Run", expectedErrorRuns+1),
	)

	// panic cron assertions
	assert.Equal(t, expectedPanicRuns, cronTracker.JobExecutions("panic"))

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":              "info",
		"service":            "test",
		"system":             "cron",
		"cronJob":            "panic",
		"cronJobExecutionID": cronJobExecutionId,
		"message":            "job execution start",
	})

	for i := 1; i <= expectedPanicRuns; i++ {
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":              "info",
			"service":            "test",
			"system":             "cron",
			"cronJob":            "panic",
			"cronJobExecutionID": cronJobExecutionId,
			"message":            "panic cron job log from test",
			"run":                i,
		})
	}

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":              "info",
		"service":            "test",
		"system":             "cron",
		"cronJob":            "panic",
		"cronJobExecutionID": cronJobExecutionId,
		"message":            "panic cron job log from test",
		"run":                expectedPanicRuns + 1,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":              "error",
		"service":            "test",
		"system":             "cron",
		"cronJob":            "panic",
		"cronJobExecutionID": cronJobExecutionId,
		"message":            "job execution panic",
		"panic":              "panic cron job panic",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"cron panic",
		attribute.String("CronJob", "panic"),
		attribute.String("CronJobExecutionID", cronJobExecutionId),
	)

	tracetest.AssertHasNotTraceSpan(
		t,
		traceExporter,
		"cron panic",
		attribute.String("CronJob", "panic"),
		attribute.String("CronJobExecutionID", cronJobExecutionId),
		attribute.Int("Run", expectedPanicRuns),
	)

	// cron metrics assertions
	expected := `
		# HELP foo_bar_job_execution_total Total number of cron job executions
		# TYPE foo_bar_job_execution_total counter
        foo_bar_job_execution_total{job="success",status="success"} 3
        foo_bar_job_execution_total{job="error",status="error"} 2
        foo_bar_job_execution_total{job="panic",status="error"} 1
	`

	err := testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expected),
		"foo_bar_job_execution_total",
	)
	assert.NoError(t, err)
}

func TestModuleInfo(t *testing.T) {
	startAt := time.Now().Add(5 * time.Second)

	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_ENV", "test")
	t.Setenv("CRON_START_IMMEDIATELY", "false")
	t.Setenv("CRON_START_AT", startAt.Format(time.RFC3339))
	t.Setenv("CRON_CONCURRENCY_LIMIT_MODE", "reschedule")
	t.Setenv("CRON_SINGLETON_ENABLED", "true")
	t.Setenv("CRON_SINGLETON_MODE", "reschedule")

	var modulesInfo []any

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxcron.FxCronModule,
		fx.Options(
			// cron jobs registration
			fxcron.AsCronJob(job.NewDummyCron, `*/1 * * * * *`),
		),
		// extraction
		fx.Populate(
			fx.Annotate(
				&modulesInfo,
				fx.ParamTags(`group:"core-module-infos"`),
			),
		),
		// invoke scheduler
		fx.Invoke(func(scheduler gocron.Scheduler) {}),
	).RequireStart()

	// scheduling assertions
	info, ok := modulesInfo[0].(*fxcron.FxCronModuleInfo)
	if !ok {
		t.Error("expected type *fxcron.FxCronModuleInfo")
	}

	assert.Equal(
		t,
		map[string]interface{}{
			"jobs": map[string]interface{}{
				"scheduled": map[string]interface{}{
					"dummy": map[string]interface{}{
						"expression": `*/1 * * * * *`,
						"last_run":   time.Time{}.Format(time.RFC3339),
						"next_run":   startAt.Format(time.RFC3339),
						"type":       "*job.DummyCron",
					},
				},
				"unscheduled": map[string]interface{}{},
			},
		},
		info.Data(),
	)

	app.RequireStop()
}

func TestModuleDecoration(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_ENV", "test")
	t.Setenv("MODULES_CRON_METRICS_COLLECT_ENABLED", "false")

	var scheduler gocron.Scheduler

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxcron.FxCronModule,
		fx.Decorate(factory.NewTestCronSchedulerFactory),
		fx.Populate(&scheduler),
	).RequireStart().RequireStop()

	assert.Len(t, scheduler.Jobs(), 0)
}
