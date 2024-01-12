package fxcore

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/rs/zerolog"
)

// FxModuleInfo is the interface to implement by modules to provide their info to fxcore.
type FxModuleInfo interface {
	Name() string
	Data() map[string]any
}

// FxCoreModuleInfo is a module info collector for fxcore.
type FxCoreModuleInfo struct {
	AppName        string
	AppEnv         string
	AppDebug       bool
	AppVersion     string
	LogLevel       string
	LogOutput      string
	TraceProcessor string
	TraceSampler   string
}

// NewFxCoreModuleInfo returns a new [FxCoreModuleInfo].
func NewFxCoreModuleInfo(cfg *config.Config) *FxCoreModuleInfo {
	logLevel, logOutput := "", ""
	if cfg.IsTestEnv() {
		logLevel = zerolog.DebugLevel.String()
		logOutput = log.TestOutputWriter.String()
	} else {
		logLevel = log.FetchLogLevel(cfg.GetString("modules.log.level")).String()
		logOutput = log.FetchLogOutputWriter(cfg.GetString("modules.log.output")).String()
	}

	traceProcessor := ""
	traceSampler := trace.FetchSampler(cfg.GetString("modules.trace.sampler.type")).String()
	if cfg.IsTestEnv() {
		traceProcessor = trace.TestSpanProcessor.String()
	} else {
		traceProcessor = trace.FetchSpanProcessor(cfg.GetString("modules.trace.processor.type")).String()
	}

	return &FxCoreModuleInfo{
		AppName:        cfg.AppName(),
		AppEnv:         cfg.AppEnv(),
		AppDebug:       cfg.AppDebug(),
		AppVersion:     cfg.AppVersion(),
		LogLevel:       logLevel,
		LogOutput:      logOutput,
		TraceProcessor: traceProcessor,
		TraceSampler:   traceSampler,
	}
}

// Name return the name of the module info.
func (i *FxCoreModuleInfo) Name() string {
	return ModuleName
}

// Data return the data of the module info.
func (i *FxCoreModuleInfo) Data() map[string]interface{} {
	return map[string]interface{}{
		"app": map[string]interface{}{
			"name":    i.AppName,
			"env":     i.AppEnv,
			"debug":   i.AppDebug,
			"version": i.AppVersion,
		},
		"log": map[string]interface{}{
			"level":  i.LogLevel,
			"output": i.LogOutput,
		},
		"trace": map[string]interface{}{
			"processor": i.TraceProcessor,
			"sampler":   i.TraceSampler,
		},
	}
}
