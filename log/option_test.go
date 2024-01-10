package log_test

import (
	"bytes"
	"testing"

	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLoggerOptions(t *testing.T) {
	t.Parallel()

	t.Run("test WithName", func(t *testing.T) {
		t.Parallel()

		o := &log.Options{}
		name := "test"
		opt := log.WithServiceName(name)
		opt(o)
		assert.Equal(t, name, o.ServiceName)
	})

	t.Run("test WithLevel", func(t *testing.T) {
		t.Parallel()

		o := &log.Options{}
		level := zerolog.WarnLevel
		opt := log.WithLevel(level)
		opt(o)
		assert.Equal(t, level, o.Level)
	})

	t.Run("test WithOutputWriter", func(t *testing.T) {
		t.Parallel()

		o := &log.Options{}
		var buf bytes.Buffer
		opt := log.WithOutputWriter(&buf)
		opt(o)
		assert.Equal(t, &buf, o.OutputWriter)
	})
}
