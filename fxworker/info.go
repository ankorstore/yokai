package fxworker

import (
	"time"

	"github.com/ankorstore/yokai/worker"
)

// FxWorkerModuleInfo is a module info collector for fxcore.
type FxWorkerModuleInfo struct {
	pool *worker.WorkerPool
}

// NewFxWorkerModuleInfo returns a new [FxWorkerModuleInfo].
func NewFxWorkerModuleInfo(pool *worker.WorkerPool) *FxWorkerModuleInfo {
	return &FxWorkerModuleInfo{
		pool: pool,
	}
}

// Name return the name of the module info.
func (i *FxWorkerModuleInfo) Name() string {
	return ModuleName
}

// Data return the data of the module info.
func (i *FxWorkerModuleInfo) Data() map[string]interface{} {
	data := map[string]interface{}{}

	for name, execution := range i.pool.Executions() {
		var events []map[string]string
		for _, event := range execution.Events() {
			events = append(events, map[string]string{
				"execution": event.ExecutionId(),
				"message":   event.Message(),
				"time":      event.Timestamp().Format(time.DateTime),
			})
		}

		data[name] = map[string]interface{}{
			"status": execution.Status().String(),
			"events": events,
		}
	}

	return map[string]interface{}{
		"workers": data,
	}
}
