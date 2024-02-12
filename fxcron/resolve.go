package fxcron

import "github.com/go-co-op/gocron/v2"

// ResolvedCronJob represents a resolved cron job, with its expression and execution options.
type ResolvedCronJob struct {
	implementation CronJob
	expression     string
	options        []gocron.JobOption
}

// NewResolvedCronJob returns a new [ResolvedCronJob] instance.
func NewResolvedCronJob(implementation CronJob, expression string, options ...gocron.JobOption) *ResolvedCronJob {
	return &ResolvedCronJob{
		implementation: implementation,
		expression:     expression,
		options:        options,
	}
}

// Implementation returns the [ResolvedCronJob] cron job implementation.
func (r *ResolvedCronJob) Implementation() CronJob {
	return r.implementation
}

// Expression returns the [ResolvedCronJob] cron job expression.
func (r *ResolvedCronJob) Expression() string {
	return r.expression
}

// Options returns the [ResolvedCronJob] cron job execution options.
func (r *ResolvedCronJob) Options() []gocron.JobOption {
	return r.options
}
