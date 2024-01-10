package log_test

import (
	"testing"

	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestFetchLogLevel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		level    string
		expected zerolog.Level
	}{
		{
			name:     "Trace Level",
			level:    "trace",
			expected: zerolog.TraceLevel,
		},
		{
			name:     "Debug Level",
			level:    "debug",
			expected: zerolog.DebugLevel,
		},
		{
			name:     "Info Level",
			level:    "info",
			expected: zerolog.InfoLevel,
		},
		{
			name:     "Warning Level",
			level:    "warning",
			expected: zerolog.WarnLevel,
		},
		{
			name:     "Error Level",
			level:    "error",
			expected: zerolog.ErrorLevel,
		},
		{
			name:     "Fatal Level",
			level:    "fatal",
			expected: zerolog.FatalLevel,
		},
		{
			name:     "Panic Level",
			level:    "panic",
			expected: zerolog.PanicLevel,
		},
		{
			name:     "No Level",
			level:    "no-level",
			expected: zerolog.NoLevel,
		},
		{
			name:     "Disabled",
			level:    "disabled",
			expected: zerolog.Disabled,
		},
		{
			name:     "Default Level",
			level:    "unknown",
			expected: zerolog.InfoLevel,
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, log.FetchLogLevel(c.level))
	}
}

func TestLogOutputWriterAsString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, log.Stdout, log.StdoutOutputWriter.String())
	assert.Equal(t, log.Noop, log.NoopOutputWriter.String())
	assert.Equal(t, log.Test, log.TestOutputWriter.String())
	assert.Equal(t, log.Console, log.ConsoleOutputWriter.String())
}

func TestFetchLogOutputWriter(t *testing.T) {
	t.Parallel()

	assert.Equal(t, log.StdoutOutputWriter, log.FetchLogOutputWriter(log.Stdout))
	assert.Equal(t, log.NoopOutputWriter, log.FetchLogOutputWriter(log.Noop))
	assert.Equal(t, log.TestOutputWriter, log.FetchLogOutputWriter(log.Test))
	assert.Equal(t, log.ConsoleOutputWriter, log.FetchLogOutputWriter(log.Console))

	// default fallback on stdout
	assert.Equal(t, log.StdoutOutputWriter, log.FetchLogOutputWriter("random"))
}
