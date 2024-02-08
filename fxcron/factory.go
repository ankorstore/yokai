package fxcron

import (
	"github.com/go-co-op/gocron/v2"
)

type CronSchedulerFactory interface {
	Create(options ...gocron.SchedulerOption) (gocron.Scheduler, error)
}

type DefaultCronSchedulerFactory struct{}

func NewDefaultCronSchedulerFactory() CronSchedulerFactory {
	return &DefaultCronSchedulerFactory{}
}

func (f *DefaultCronSchedulerFactory) Create(options ...gocron.SchedulerOption) (gocron.Scheduler, error) {
	return gocron.NewScheduler(options...)
}
