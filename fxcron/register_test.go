package fxcron_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxcron"
	"github.com/ankorstore/yokai/fxcron/testdata/cron/job"
	"github.com/go-co-op/gocron/v2"
	"github.com/stretchr/testify/assert"
)

func TestAsCronJob(t *testing.T) {
	t.Parallel()

	result := fxcron.AsCronJob(job.NewDummyCron, `* * * * *`, gocron.WithLimitedRuns(1))

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result))
}
