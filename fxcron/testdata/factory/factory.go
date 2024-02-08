package factory

import (
	"github.com/ankorstore/yokai/fxcron"
	"github.com/go-co-op/gocron/v2"
)

type TestCronSchedulerFactory struct{}

func NewTestCronSchedulerFactory() fxcron.CronSchedulerFactory {
	return &TestCronSchedulerFactory{}
}

func (f *TestCronSchedulerFactory) Create(options ...gocron.SchedulerOption) (gocron.Scheduler, error) {
	return gocron.NewScheduler()
}
