package fxlog

import (
	"strings"

	"github.com/ankorstore/yokai/log"
	"go.uber.org/fx/fxevent"
)

// FxEventLogger is a logger compatible with [Fx], decorating [log.Logger].
//
// [Fx]: https://github.com/uber-go/fx
type FxEventLogger struct {
	logger *log.Logger
}

// NewFxEventLogger returns a new [NewFxEventLogger] from a provided [log.Logger].
func NewFxEventLogger(logger *log.Logger) fxevent.Logger {
	return &FxEventLogger{
		logger: logger,
	}
}

// LogEvent logs a [fxevent.Event].
func (l *FxEventLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logger.
			Info().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logger.
				Warn().
				Err(e.Err).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStart hook failed")
		} else {
			l.logger.
				Info().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		l.logger.
			Info().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logger.
				Warn().
				Err(e.Err).
				Str("callee", e.FunctionName).
				Str("callee", e.CallerName).
				Msg("OnStop hook failed")
		} else {
			l.logger.
				Info().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStop hook executed")
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.logger.
				Warn().
				Err(e.Err).
				Str("type", e.TypeName).
				Msg("supplied")
		} else {
			l.logger.
				Info().
				Str("type", e.TypeName).
				Msg("supplied")
		}
	case *fxevent.Provided:
		for _, rType := range e.OutputTypeNames {
			l.logger.
				Info().
				Str("type", rType).
				Str("constructor", e.ConstructorName).
				Msg("provided")
		}
		if e.Err != nil {
			l.logger.
				Error().
				Err(e.Err).
				Msg("error encountered while applying options")
		}
	case *fxevent.Invoking:
	case *fxevent.Invoked:
		if e.Err != nil {
			l.logger.
				Error().
				Err(e.Err).
				Str("stack", e.Trace).
				Str("function", e.FunctionName).
				Msg("invoke failed")
		} else {
			l.logger.
				Info().
				Str("function", e.FunctionName).
				Msg("invoked")
		}
	case *fxevent.Stopping:
		l.logger.
			Info().
			Str("signal", strings.ToUpper(e.Signal.String())).
			Msg("received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.logger.
				Error().
				Err(e.Err).
				Msg("stop failed")
		}
	case *fxevent.RollingBack:
		l.logger.
			Error().
			Err(e.StartErr).
			Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.logger.
				Error().
				Err(e.Err).
				Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.logger.
				Error().
				Err(e.Err).
				Msg("start failed")
		} else {
			l.logger.
				Info().
				Msg("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.logger.
				Error().
				Err(e.Err).
				Msg("custom logger initialization failed")
		} else {
			l.logger.
				Info().
				Str("function", e.ConstructorName).
				Msg("initialized custom fxevent.Logger")
		}
	}
}
