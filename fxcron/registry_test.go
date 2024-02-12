package fxcron_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxcron/testdata/cron/job"
	"github.com/go-co-op/gocron/v2"
	"github.com/stretchr/testify/assert"
)

const cronJobExpression = `* * * * *`

func TestNewFxCronJobRegistry(t *testing.T) {
	t.Parallel()

	param := fxcron.FxCronJobRegistryParam{
		CronJobs:            []fxcron.CronJob{},
		CronJobsDefinitions: []fxcron.CronJobDefinition{},
	}

	registry := fxcron.NewFxCronJobRegistry(param)

	assert.IsType(t, &fxcron.CronJobRegistry{}, registry)
}

func TestResolveCronJobsSuccess(t *testing.T) {
	t.Parallel()

	cronJob := job.NewDummyCron()
	cronJobOptions := []gocron.JobOption(nil)

	param := fxcron.FxCronJobRegistryParam{
		CronJobs: []fxcron.CronJob{cronJob},
		CronJobsDefinitions: []fxcron.CronJobDefinition{
			fxcron.NewCronJobDefinition(fxcron.GetType(cronJob), cronJobExpression, cronJobOptions...),
		},
	}

	registry := fxcron.NewFxCronJobRegistry(param)

	resolvedCronJobs, err := registry.ResolveCronJobs()
	assert.NoError(t, err)

	assert.Len(t, resolvedCronJobs, 1)
	assert.Equal(t, cronJob, resolvedCronJobs[0].Implementation())
	assert.Equal(t, cronJobExpression, resolvedCronJobs[0].Expression())
	assert.Equal(t, cronJobOptions, resolvedCronJobs[0].Options())
	assert.Equal(t, cronJob.Name(), resolvedCronJobs[0].Implementation().Name())
	assert.Nil(t, resolvedCronJobs[0].Implementation().Run(context.Background()))
}

func TestResolveCronJobsFailure(t *testing.T) {
	t.Parallel()

	cronJob := job.NewDummyCron()

	param := fxcron.FxCronJobRegistryParam{
		CronJobs: []fxcron.CronJob{cronJob},
		CronJobsDefinitions: []fxcron.CronJobDefinition{
			fxcron.NewCronJobDefinition("invalid", cronJobExpression),
		},
	}

	registry := fxcron.NewFxCronJobRegistry(param)

	_, err := registry.ResolveCronJobs()
	assert.Error(t, err)
	assert.Equal(t, "cannot find cron job implementation for type invalid", err.Error())
}
