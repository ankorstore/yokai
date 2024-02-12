package fxcron

import (
	"github.com/go-co-op/gocron/v2"
)

// CronSchedulerFactory is the interface for [gocron.Scheduler] factories.
type CronSchedulerFactory interface {
	Create(options ...gocron.SchedulerOption) (gocron.Scheduler, error)
}

// DefaultCronSchedulerFactory is the default [CronSchedulerFactory] implementation.
type DefaultCronSchedulerFactory struct{}

// NewDefaultCronSchedulerFactory returns a [DefaultCronSchedulerFactory], implementing [CronSchedulerFactory].
func NewDefaultCronSchedulerFactory() CronSchedulerFactory {
	return &DefaultCronSchedulerFactory{}
}

// Create returns a new [gocron.Scheduler] instance for an optional list of [gocron.SchedulerOption].
func (f *DefaultCronSchedulerFactory) Create(options ...gocron.SchedulerOption) (gocron.Scheduler, error) {
	return gocron.NewScheduler(options...)
}
