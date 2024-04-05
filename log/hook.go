package log

import (
	"github.com/rs/zerolog"
)

type CallerInfoHook struct{}

func (h CallerInfoHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	e.Caller(3)
}
