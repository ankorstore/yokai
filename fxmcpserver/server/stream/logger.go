package stream

import "github.com/ankorstore/yokai/log"

type Logger interface {
	Infof(format string, v ...any)
	Errorf(format string, v ...any)
}

type MCPStreamableHTTPServerLogger struct {
	logger *log.Logger
}

func NewMCPStreamableHTTPServerLogger(logger *log.Logger) *MCPStreamableHTTPServerLogger {
	return &MCPStreamableHTTPServerLogger{
		logger: logger,
	}
}

func (l *MCPStreamableHTTPServerLogger) Infof(format string, v ...any) {
	l.logger.Info().Msgf(format, v...)
}

func (l *MCPStreamableHTTPServerLogger) Errorf(format string, v ...any) {
	l.logger.Error().Msgf(format, v...)
}
