package fxcron

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

type CronJob interface {
	Name() string
	Run(ctx context.Context) error
}

type CronJobRegistry struct {
	cronJobs           []CronJob
	cronJobDefinitions []CronJobDefinition
}

type FxCronJobRegistryParam struct {
	fx.In
	CronJobs            []CronJob           `group:"cron-jobs"`
	CronJobsDefinitions []CronJobDefinition `group:"cron-jobs-definitions"`
}

func NewFxCronJobRegistry(p FxCronJobRegistryParam) *CronJobRegistry {
	return &CronJobRegistry{
		cronJobs:           p.CronJobs,
		cronJobDefinitions: p.CronJobsDefinitions,
	}
}

func (r *CronJobRegistry) ResolveCronJobs() ([]*ResolvedCronJob, error) {
	resolvedCronJobs := []*ResolvedCronJob{}

	for _, definition := range r.cronJobDefinitions {
		implementation, err := r.lookupRegisteredCronJob(definition.ReturnType())
		if err != nil {
			return nil, err
		}

		resolvedCronJobs = append(
			resolvedCronJobs,
			NewResolvedCronJob(implementation, definition.Expression(), definition.Options()...),
		)
	}

	return resolvedCronJobs, nil
}

func (r *CronJobRegistry) lookupRegisteredCronJob(returnType string) (CronJob, error) {
	for _, implementation := range r.cronJobs {
		if GetType(implementation) == returnType {
			return implementation, nil
		}
	}

	return nil, fmt.Errorf("cannot find cron job implementation for type %s", returnType)
}
