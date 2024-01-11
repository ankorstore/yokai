package fxlog

import (
	"io"
	"os"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "log"

// FxLogModule is the [Fx] log module.
//
// [Fx]: https://github.com/uber-go/fx
var FxLogModule = fx.Module(
	ModuleName,
	fx.Provide(
		log.NewDefaultLoggerFactory,
		logtest.NewDefaultTestLogBuffer,
		NewFxLogger,
	),
)

// FxLogParam allows injection of the required dependencies in [NewFxLogger].
type FxLogParam struct {
	fx.In
	Factory log.LoggerFactory
	Buffer  logtest.TestLogBuffer
	Config  *config.Config
}

// NewFxLogger returns a [log.Logger].
func NewFxLogger(p FxLogParam) (*log.Logger, error) {
	var level zerolog.Level
	if p.Config.AppDebug() {
		level = zerolog.DebugLevel
	} else {
		level = log.FetchLogLevel(p.Config.GetString("modules.log.level"))
	}

	var outputWriter io.Writer
	if p.Config.IsTestEnv() {
		outputWriter = p.Buffer
	} else {
		switch log.FetchLogOutputWriter(p.Config.GetString("modules.log.output")) {
		case log.NoopOutputWriter:
			outputWriter = io.Discard
		case log.TestOutputWriter:
			outputWriter = p.Buffer
		case log.ConsoleOutputWriter:
			outputWriter = zerolog.ConsoleWriter{Out: os.Stderr}
		default:
			outputWriter = os.Stdout
		}
	}

	return p.Factory.Create(
		log.WithServiceName(p.Config.AppName()),
		log.WithLevel(level),
		log.WithOutputWriter(outputWriter),
	)
}
