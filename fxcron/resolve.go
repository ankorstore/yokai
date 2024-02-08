package fxcron

import "github.com/go-co-op/gocron/v2"

type ResolvedCronJob struct {
	implementation CronJob
	expression     string
	options        []gocron.JobOption
}

func NewResolvedCronJob(implementation CronJob, expression string, options ...gocron.JobOption) *ResolvedCronJob {
	return &ResolvedCronJob{
		implementation: implementation,
		expression:     expression,
		options:        options,
	}
}

func (r *ResolvedCronJob) Implementation() CronJob {
	return r.implementation
}

func (r *ResolvedCronJob) Expression() string {
	return r.expression
}

func (r *ResolvedCronJob) Options() []gocron.JobOption {
	return r.options
}
