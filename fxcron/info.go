package fxcron

import (
	"reflect"

	"github.com/go-co-op/gocron/v2"
)

const NON_AVAILABLE = "n/a"

type FxCronModuleInfo struct {
	scheduler gocron.Scheduler
	registry  *CronJobRegistry
}

func NewFxCronModuleInfo(scheduler gocron.Scheduler, registry *CronJobRegistry) *FxCronModuleInfo {
	return &FxCronModuleInfo{
		scheduler: scheduler,
		registry:  registry,
	}
}

func (i *FxCronModuleInfo) Name() string {
	return ModuleName
}

func (i *FxCronModuleInfo) Data() map[string]interface{} {
	scheduledJobs := i.scheduler.Jobs()

	resolvedJobs, err := i.registry.ResolveCronJobs()
	if err != nil {
		return map[string]interface{}{
			"jobs": map[string]interface{}{
				"scheduled":   NON_AVAILABLE,
				"unscheduled": NON_AVAILABLE,
			},
		}
	}

	scheduledJobsData := make(map[string]interface{})
	unscheduledJobsData := make(map[string]interface{})

	for _, resolvedJob := range resolvedJobs {
		isJobScheduled := false

		for _, scheduledJob := range scheduledJobs {
			if resolvedJob.Implementation().Name() == scheduledJob.Name() {
				isJobScheduled = true

				scheduledJobsData[resolvedJob.Implementation().Name()] = map[string]interface{}{
					"expression": resolvedJob.Expression(),
					"last_run":   i.jobLastRun(scheduledJob),
					"next_run":   i.jobNextRun(scheduledJob),
					"type":       i.jobType(resolvedJob.Implementation()),
				}
			}
		}

		if !isJobScheduled {
			unscheduledJobsData[resolvedJob.Implementation().Name()] = map[string]interface{}{
				"expression": resolvedJob.Expression(),
				"type":       i.jobType(resolvedJob.Implementation()),
			}
		}
	}

	return map[string]interface{}{
		"jobs": map[string]interface{}{
			"scheduled":   scheduledJobsData,
			"unscheduled": unscheduledJobsData,
		},
	}
}

func (i *FxCronModuleInfo) jobLastRun(job gocron.Job) string {
	if run, err := job.LastRun(); err == nil {
		return run.String()
	}

	return NON_AVAILABLE
}

func (i *FxCronModuleInfo) jobNextRun(job gocron.Job) string {
	if run, err := job.NextRun(); err == nil {
		return run.String()
	}

	return NON_AVAILABLE
}

func (i *FxCronModuleInfo) jobType(job CronJob) string {
	return reflect.ValueOf(job).Type().String()
}
