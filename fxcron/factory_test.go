package fxcron_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxcron"
	"github.com/go-co-op/gocron/v2"
	"github.com/stretchr/testify/assert"
)

func TestDefaultCronSchedulerFactory(t *testing.T) {
	t.Parallel()

	factory := fxcron.NewDefaultCronSchedulerFactory()

	assert.IsType(t, &fxcron.DefaultCronSchedulerFactory{}, factory)
	assert.Implements(t, (*fxcron.CronSchedulerFactory)(nil), factory)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	factory := fxcron.NewDefaultCronSchedulerFactory()

	scheduler, err := factory.Create()
	assert.NoError(t, err)

	assert.Implements(t, (*gocron.Scheduler)(nil), scheduler)
}
