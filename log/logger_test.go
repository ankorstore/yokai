package log_test

import (
	"testing"

	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLoggerConversion(t *testing.T) {
	t.Parallel()

	// convert to Logger
	zeroLogger := zerolog.New(nil)
	logger := log.FromZerolog(zeroLogger)
	assert.NotNil(t, logger)
	assert.IsType(t, &log.Logger{}, logger)

	// convert back to zerolog
	backToZeroLogger := logger.ToZerolog()
	assert.NotNil(t, backToZeroLogger)
	assert.IsType(t, &zerolog.Logger{}, backToZeroLogger)
	assert.Equal(t, &zeroLogger, backToZeroLogger)
}
