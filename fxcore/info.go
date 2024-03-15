package fxcore

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// FxExtraInfo is the struct used by modules or apps to provide their extra info to the core.
type FxExtraInfo interface {
	Name() string
	Value() string
}

// fxExtraInfo is the default [FxExtraInfo] implementation.
type fxExtraInfo struct {
	name  string
	value string
}

// NewFxExtraInfo returns a new FxExtraInfo.
func NewFxExtraInfo(name string, value string) FxExtraInfo {
	return &fxExtraInfo{
		name:  name,
		value: value,
	}
}

// Name returns the name of the [fxExtraInfo].
func (i *fxExtraInfo) Name() string {
	return i.name
}

// Value returns the value of the [fxExtraInfo].
func (i *fxExtraInfo) Value() string {
	return i.value
}

// FxModuleInfo is the interface to implement by modules to provide their info to the core.
type FxModuleInfo interface {
	Name() string
	Data() map[string]any
}

// FxCoreModuleInfo is a module info collector for the core.
type FxCoreModuleInfo struct {
	AppName        string
	AppEnv         string
	AppDebug       bool
	AppVersion     string
	LogLevel       string
	LogOutput      string
	TraceProcessor string
	TraceSampler   string
	ExtraInfos     map[string]string
}

// FxCoreModuleInfoParam allows injection of the required dependencies in [NewFxCoreModuleInfo].
type FxCoreModuleInfoParam struct {
	fx.In
	Config     *config.Config
	ExtraInfos []FxExtraInfo `group:"core-extra-infos"`
}

// NewFxCoreModuleInfo returns a new [FxCoreModuleInfo].
func NewFxCoreModuleInfo(p FxCoreModuleInfoParam) *FxCoreModuleInfo {
	logLevel, logOutput := "", ""
	if p.Config.IsTestEnv() {
		logLevel = zerolog.DebugLevel.String()
		logOutput = log.TestOutputWriter.String()
	} else {
		logLevel = log.FetchLogLevel(p.Config.GetString("modules.log.level")).String()
		logOutput = log.FetchLogOutputWriter(p.Config.GetString("modules.log.output")).String()
	}

	traceProcessor := ""
	traceSampler := trace.FetchSampler(p.Config.GetString("modules.trace.sampler.type")).String()
	if p.Config.IsTestEnv() {
		traceProcessor = trace.TestSpanProcessor.String()
	} else {
		traceProcessor = trace.FetchSpanProcessor(p.Config.GetString("modules.trace.processor.type")).String()
	}

	extraInfos := make(map[string]string)
	for _, info := range p.ExtraInfos {
		extraInfos[info.Name()] = info.Value()
	}

	return &FxCoreModuleInfo{
		AppName:        p.Config.AppName(),
		AppEnv:         p.Config.AppEnv(),
		AppDebug:       p.Config.AppDebug(),
		AppVersion:     p.Config.AppVersion(),
		LogLevel:       logLevel,
		LogOutput:      logOutput,
		TraceProcessor: traceProcessor,
		TraceSampler:   traceSampler,
		ExtraInfos:     extraInfos,
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
		"extra": i.ExtraInfos,
	}
}
