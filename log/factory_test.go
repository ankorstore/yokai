package log_test

import (
	"testing"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLoggerFactory(t *testing.T) {
	t.Parallel()

	factory := log.NewDefaultLoggerFactory()

	assert.IsType(t, &log.DefaultLoggerFactory{}, factory)
	assert.Implements(t, (*log.LoggerFactory)(nil), factory)
}

func TestCreateSuccess(t *testing.T) {
	t.Parallel()

	testLogBuffer := logtest.NewDefaultTestLogBuffer()

	factory := log.NewDefaultLoggerFactory()
	logger, err := factory.Create(
		log.WithServiceName("test logger"),
		log.WithLevel(zerolog.InfoLevel),
		log.WithOutputWriter(testLogBuffer),
	)

	assert.NoError(t, err)
	assert.IsType(t, &log.Logger{}, logger)

	logger.Info().Msg("some message")

	logtest.AssertHasLogRecord(t, testLogBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test logger",
		"message": "some message",
	})

	logtest.AssertHasNotLogRecord(t, testLogBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test logger",
		"message": "some invalid message",
	})

	logtest.AssertContainLogRecord(t, testLogBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test logger",
		"message": "ome mess",
	})

	logtest.AssertContainNotLogRecord(t, testLogBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test logger",
		"message": "ome invalid mess",
	})
}
