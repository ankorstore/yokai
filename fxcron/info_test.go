package fxcron_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxcron/testdata/cron/job"
	"github.com/go-co-op/gocron/v2"
	"github.com/stretchr/testify/assert"
)

func TestFxCronModuleInfo(t *testing.T) {
	t.Parallel()

	cronJob := job.NewDummyCron()
	cronJobExpression := `*/1 * * * * *`
	cronJobOptions := []gocron.JobOption(nil)

	scheduler, err := gocron.NewScheduler()
	assert.NoError(t, err)

	_, err = scheduler.NewJob(
		gocron.CronJob(cronJobExpression, true),
		gocron.NewTask(
			func() {
				err := cronJob.Run(context.Background())
				assert.NoError(t, err)
			},
		),
	)
	assert.NoError(t, err)

	param := fxcron.FxCronJobRegistryParam{
		CronJobs: []fxcron.CronJob{cronJob},
		CronJobsDefinitions: []fxcron.CronJobDefinition{
			fxcron.NewCronJobDefinition(fxcron.GetType(cronJob), cronJobExpression, cronJobOptions...),
		},
	}

	registry := fxcron.NewFxCronJobRegistry(param)

	info := fxcron.NewFxCronModuleInfo(scheduler, registry)

	assert.IsType(t, &fxcron.FxCronModuleInfo{}, info)
	assert.Equal(t, fxcron.ModuleName, info.Name())

	assert.Equal(
		t,
		map[string]interface{}{
			"jobs": map[string]interface{}{
				"scheduled": map[string]interface{}{},
				"unscheduled": map[string]interface{}{
					"dummy": map[string]interface{}{
						"expression": cronJobExpression,
						"type":       fxcron.GetType(cronJob),
					},
				},
			},
		},
		info.Data(),
	)
}

func TestFxCronModuleInfoError(t *testing.T) {
	t.Parallel()

	cronJobExpression := `*/1 * * * * *`

	scheduler, err := gocron.NewScheduler()
	assert.NoError(t, err)

	param := fxcron.FxCronJobRegistryParam{
		CronJobs: []fxcron.CronJob{},
		CronJobsDefinitions: []fxcron.CronJobDefinition{
			fxcron.NewCronJobDefinition("invalid", cronJobExpression),
		},
	}

	registry := fxcron.NewFxCronJobRegistry(param)

	info := fxcron.NewFxCronModuleInfo(scheduler, registry)

	assert.Equal(
		t,
		map[string]interface{}{
			"jobs": map[string]interface{}{
				"scheduled":   fxcron.NON_AVAILABLE,
				"unscheduled": fxcron.NON_AVAILABLE,
			},
		},
		info.Data(),
	)
}
