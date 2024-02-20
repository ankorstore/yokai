package grpcserver_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/grpcserver/testdata/proto"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestCtxLogger(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// client
	client, closer := prepareTestServiceGrpcServerAndClient(
		t,
		logger,
		[]string{},
		map[string]string{"x-meta": "meta"},
		true,
	)
	defer closer()

	// call assertions
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-request-id", testRequestId)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-meta", "data")

	response, err := client.Unary(ctx, &proto.Request{
		ShouldFail: false,
		Message:    "test",
	})
	assert.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, "test", response.Message)

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "unary call",
		"requestID": testRequestId,
		"meta":      "data",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "unary call success",
		"requestID": testRequestId,
		"meta":      "data",
	})
}

func TestCtxTracer(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// tracer
	exporter := tracetest.NewDefaultTestTraceExporter()

	_, err = trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	// client
	client, closer := prepareTestServiceGrpcServerAndClient(
		t,
		logger,
		[]string{},
		map[string]string{},
		false,
	)
	defer closer()

	// call assertions
	response, err := client.Unary(context.Background(), &proto.Request{
		ShouldFail: false,
		Message:    "test",
	})
	assert.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, "test", response.Message)

	// trace assertions
	tracetest.AssertHasTraceSpan(t, exporter, "unary trace")
}
