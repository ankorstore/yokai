package tracker

import "sync"

type CronExecutionTracker struct {
	mutex      sync.Mutex
	executions map[string]int
}

func NewCronExecutionTracker() *CronExecutionTracker {
	return &CronExecutionTracker{
		executions: make(map[string]int),
	}
}

func (t *CronExecutionTracker) TrackJobExecution(jobName string) *CronExecutionTracker {
	t.mutex.Lock()

	if executions, ok := t.executions[jobName]; ok {
		t.executions[jobName] = executions + 1
	} else {
		t.executions[jobName] = 1
	}

	t.mutex.Unlock()

	return t
}

func (t *CronExecutionTracker) JobExecutions(jobName string) int {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	jobExecutions := 0
	if executions, ok := t.executions[jobName]; ok {
		jobExecutions = executions
	}

	return jobExecutions
}
