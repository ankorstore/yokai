package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"runtime"
)

type CallerInfoHook struct{}

func (h CallerInfoHook) Run(e *zerolog.Event, l zerolog.Level, msg string) {
	_, file, line, ok := runtime.Caller(0)
	if ok {
		e.Str("caller", fmt.Sprintf("%s:%d", file, line))
	}
}
