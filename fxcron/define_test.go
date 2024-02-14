package fxcron_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxcron"
	"github.com/go-co-op/gocron/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewCronJobDefinition(t *testing.T) {
	t.Parallel()

	definition := fxcron.NewCronJobDefinition("*TestCron", `* * * * *`)

	assert.Implements(t, (*fxcron.CronJobDefinition)(nil), definition)
	assert.Equal(t, "*TestCron", definition.ReturnType())
	assert.Equal(t, `* * * * *`, definition.Expression())
	assert.Equal(t, []gocron.JobOption(nil), definition.Options())
}
