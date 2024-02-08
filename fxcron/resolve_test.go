package fxcron_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxcron/testdata/cron/job"
	"github.com/go-co-op/gocron/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewResolvedCronJob(t *testing.T) {
	t.Parallel()

	cronJob := job.NewDummyCron()
	expression := `* * * * *`
	options := []gocron.JobOption(nil)

	resolvedCronJob := fxcron.NewResolvedCronJob(cronJob, expression, options...)

	assert.IsType(t, &fxcron.ResolvedCronJob{}, resolvedCronJob)
	assert.Equal(t, cronJob, resolvedCronJob.Implementation())
	assert.Equal(t, expression, resolvedCronJob.Expression())
	assert.Equal(t, options, resolvedCronJob.Options())
}
