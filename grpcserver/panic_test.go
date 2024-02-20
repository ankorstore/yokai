package grpcserver_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/grpcserver"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/stretchr/testify/assert"
)

func TestHandleWithoutDebug(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// context
	ctx := logger.WithContext(context.Background())

	// handler assertion
	handler := grpcserver.NewGrpcPanicRecoveryHandler()
	err = handler.Handle(false)(ctx, nil)
	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = internal grpc server error", err.Error())

	// log assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"panic":   "%!s(<nil>)",
		"message": "grpc recovered from panic",
	})
}

func TestHandleWithDebug(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// context
	ctx := logger.WithContext(context.Background())

	// handler assertion
	handler := grpcserver.NewGrpcPanicRecoveryHandler()
	err = handler.Handle(true)(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rpc error: code = Internal desc = internal grpc server error, panic = %!s(<nil>), stack = goroutine")

	// log assertion
	logtest.AssertContainLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"panic":   "%!s(<nil>)",
		"stack":   "goroutine",
		"message": "grpc recovered from panic",
	})
}
