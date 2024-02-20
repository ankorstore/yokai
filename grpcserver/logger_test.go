package grpcserver_test

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/ankorstore/yokai/generate/generatetest/uuid"
	"github.com/ankorstore/yokai/grpcserver"
	"github.com/ankorstore/yokai/grpcserver/grpcservertest"
	"github.com/ankorstore/yokai/grpcserver/testdata/interceptor"
	"github.com/ankorstore/yokai/grpcserver/testdata/proto"
	"github.com/ankorstore/yokai/grpcserver/testdata/service"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	testRequestId = "33084b3e-9b90-926c-af19-3859d70bd296"
	testTraceId   = "c4ca71e03e42c2c3d54293a6e2608bfa"
	testSpanId    = "8d0fdc8a74baaaea"
)

func TestUnarySuccess(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
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
		"level":      "debug",
		"grpcMethod": "/test.Service/Unary",
		"grpcType":   "unary",
		"message":    "grpc call start",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
		"meta":       "data",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "unary call",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"meta":      "data",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "unary call success",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"meta":      "data",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"grpcMethod": "/test.Service/Unary",
		"grpcType":   "unary",
		"grpcStatus": "OK",
		"message":    "grpc call success",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
		"meta":       "data",
	})
}

func TestUnaryFailure(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// client
	client, closer := prepareTestServiceGrpcServerAndClient(
		t,
		logger,
		[]string{},
		map[string]string{},
		true,
	)
	defer closer()

	// call assertions
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-request-id", testRequestId)

	response, err := client.Unary(ctx, &proto.Request{
		ShouldFail: true,
		Message:    "test",
	})
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = failure", err.Error())

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "debug",
		"grpcMethod": "/test.Service/Unary",
		"grpcType":   "unary",
		"message":    "grpc call start",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "unary call",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"message":   "unary call failure",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "error",
		"grpcMethod": "/test.Service/Unary",
		"grpcType":   "unary",
		"grpcStatus": "Internal",
		"error":      "rpc error: code = Internal desc = failure",
		"message":    "grpc call error",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
	})
}

func TestUnaryWithExclusion(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// client
	client, closer := prepareTestServiceGrpcServerAndClient(
		t,
		logger,
		[]string{"/test.Service/Unary"},
		map[string]string{},
		true,
	)
	defer closer()

	// call assertions
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-request-id", testRequestId)

	response, err := client.Unary(ctx, &proto.Request{
		ShouldFail: false,
		Message:    "test",
	})
	assert.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, "test", response.Message)

	// logs assertions
	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "debug",
		"grpcMethod": "/test.Service/Unary",
		"grpcType":   "unary",
		"message":    "grpc call start",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "unary call",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "unary call success",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"grpcMethod": "/test.Service/Unary",
		"grpcType":   "unary",
		"grpcStatus": "OK",
		"message":    "grpc call success",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
	})
}

func TestUnaryWithExclusionAndError(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// client
	client, closer := prepareTestServiceGrpcServerAndClient(
		t,
		logger,
		[]string{"/test.Service/Unary"},
		map[string]string{},
		true,
	)
	defer closer()

	// call assertions
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-request-id", testRequestId)

	response, err := client.Unary(ctx, &proto.Request{
		ShouldFail: true,
		Message:    "test",
	})
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = failure", err.Error())

	// logs assertions
	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "debug",
		"grpcMethod": "/test.Service/Unary",
		"grpcType":   "unary",
		"message":    "grpc call start",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "unary call",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"message":   "unary call failure",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "error",
		"grpcMethod": "/test.Service/Unary",
		"grpcType":   "unary",
		"grpcStatus": "Internal",
		"error":      "rpc error: code = Internal desc = failure",
		"message":    "grpc call error",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
	})
}

func TestBidiSuccess(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
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

	// call
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-request-id", testRequestId)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-meta", "data")

	stream, err := client.Bidi(ctx)
	assert.NoError(t, err)

	wait := make(chan struct{})

	// client send
	go func() {
		err = stream.Send(&proto.Request{
			ShouldFail: false,
			Message:    "this is a test",
		})
		assert.NoError(t, err)

		err = stream.CloseSend()
		assert.NoError(t, err)
	}()

	// client recv
	var responses []*proto.Response
	go func() {
		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}

			assert.NoError(t, err)

			responses = append(responses, resp)
		}

		close(wait)
	}()

	<-wait

	// call assertions
	assert.Len(t, responses, 4)
	assert.Equal(t, "this", responses[0].Message)
	assert.Equal(t, "is", responses[1].Message)
	assert.Equal(t, "a", responses[2].Message)
	assert.Equal(t, "test", responses[3].Message)

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"grpcMethod": "/test.Service/Bidi",
		"grpcType":   "server-streaming",
		"message":    "grpc call start",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
		"meta":       "data",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi call",
		"requestID": testRequestId,
		"meta":      "data",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi recv value this is a test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"meta":      "data",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi send value this",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"meta":      "data",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi send value is",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"meta":      "data",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi send value a",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"meta":      "data",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi send value test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"meta":      "data",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"grpcMethod": "/test.Service/Bidi",
		"grpcType":   "server-streaming",
		"grpcStatus": "OK",
		"message":    "grpc call success",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
		"meta":       "data",
	})
}

func TestBidiFailure(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// client
	client, closer := prepareTestServiceGrpcServerAndClient(
		t,
		logger,
		[]string{},
		map[string]string{},
		true,
	)
	defer closer()

	// call
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-request-id", testRequestId)

	stream, err := client.Bidi(ctx)
	assert.NoError(t, err)

	wait := make(chan struct{})

	// client send
	go func() {
		err = stream.Send(&proto.Request{
			ShouldFail: true,
			Message:    "this is a test",
		})
		assert.NoError(t, err)

		err = stream.CloseSend()
		assert.NoError(t, err)
	}()

	// client recv
	go func() {
		resp, err := stream.Recv()
		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = failure", err.Error())

		close(wait)
	}()

	<-wait

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"grpcMethod": "/test.Service/Bidi",
		"grpcType":   "server-streaming",
		"message":    "grpc call start",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi call",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi recv value this is a test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "error",
		"grpcMethod": "/test.Service/Bidi",
		"grpcType":   "server-streaming",
		"grpcStatus": "Internal",
		"error":      "rpc error: code = Internal desc = failure",
		"message":    "grpc call error",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
	})
}

func TestBidiWithExclusion(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// client
	client, closer := prepareTestServiceGrpcServerAndClient(
		t,
		logger,
		[]string{"/test.Service/Bidi"},
		map[string]string{},
		true,
	)
	defer closer()

	// call
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-request-id", testRequestId)

	stream, err := client.Bidi(ctx)
	assert.NoError(t, err)

	wait := make(chan struct{})

	// client send
	go func() {
		err = stream.Send(&proto.Request{
			ShouldFail: false,
			Message:    "this is a test",
		})
		assert.NoError(t, err)

		err = stream.CloseSend()
		assert.NoError(t, err)
	}()

	// client recv
	var responses []*proto.Response
	go func() {
		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}

			assert.NoError(t, err)

			responses = append(responses, resp)
		}

		close(wait)
	}()

	<-wait

	// call assertions
	assert.Len(t, responses, 4)
	assert.Equal(t, "this", responses[0].Message)
	assert.Equal(t, "is", responses[1].Message)
	assert.Equal(t, "a", responses[2].Message)
	assert.Equal(t, "test", responses[3].Message)

	// logs assertions
	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"grpcMethod": "/test.Service/Bidi",
		"grpcType":   "server-streaming",
		"message":    "grpc call start",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi call",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi recv value this is a test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi send value this",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi send value is",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi send value a",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi send value test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"grpcMethod": "/test.Service/Bidi",
		"grpcType":   "server-streaming",
		"grpcStatus": "OK",
		"message":    "grpc call success",
		"requestID":  testRequestId,
	})
}

func TestBidiWithExclusionAndError(t *testing.T) {
	t.Parallel()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// client
	client, closer := prepareTestServiceGrpcServerAndClient(
		t,
		logger,
		[]string{"/test.Service/Bidi"},
		map[string]string{},
		true,
	)
	defer closer()

	// call
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-request-id", testRequestId)

	stream, err := client.Bidi(ctx)
	assert.NoError(t, err)

	wait := make(chan struct{})

	// client send
	go func() {
		err = stream.Send(&proto.Request{
			ShouldFail: true,
			Message:    "this is a test",
		})
		assert.NoError(t, err)

		err = stream.CloseSend()
		assert.NoError(t, err)
	}()

	// client recv
	go func() {
		resp, err := stream.Recv()
		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = failure", err.Error())

		close(wait)
	}()

	<-wait

	// logs assertions
	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"grpcMethod": "/test.Service/Bidi",
		"grpcType":   "server-streaming",
		"message":    "grpc call start",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi call",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"message":   "bidi recv value this is a test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "error",
		"grpcMethod": "/test.Service/Bidi",
		"grpcType":   "server-streaming",
		"grpcStatus": "Internal",
		"error":      "rpc error: code = Internal desc = failure",
		"message":    "grpc call error",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
	})
}

func prepareTestServiceGrpcServerAndClient(t *testing.T, logger *log.Logger, exclusions []string, metadata map[string]string, withTestInterceptors bool) (proto.ServiceClient, func()) {
	t.Helper()

	// context preparation
	ctx := logger.WithContext(context.Background())

	// bufconn listener preparation
	lis := grpcservertest.NewBufconnListener(1024 * 1024)

	// gRPC server preparation
	loggerInterceptor := grpcserver.NewGrpcLoggerInterceptor(uuid.NewTestUuidGenerator("test"), logger)

	if len(exclusions) != 0 {
		loggerInterceptor.Exclude(exclusions...)
	}

	if len(metadata) != 0 {
		loggerInterceptor.Metadata(metadata)
	}

	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	if withTestInterceptors {
		unaryInterceptors = []grpc.UnaryServerInterceptor{
			interceptor.TestUnaryInterceptor(testTraceId, testSpanId),
			loggerInterceptor.UnaryInterceptor(),
		}

		streamInterceptors = []grpc.StreamServerInterceptor{
			interceptor.TestStreamInterceptor(testTraceId, testSpanId),
			loggerInterceptor.StreamInterceptor(),
		}
	} else {
		unaryInterceptors = []grpc.UnaryServerInterceptor{
			loggerInterceptor.UnaryInterceptor(),
		}

		streamInterceptors = []grpc.StreamServerInterceptor{
			loggerInterceptor.StreamInterceptor(),
		}
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	server.RegisterService(
		&proto.Service_ServiceDesc,
		service.NewTestServiceServer(),
	)

	go func() {
		//nolint:errcheck
		server.Serve(lis)
	}()

	// gRPC client preparation
	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)

	client := proto.NewServiceClient(conn)

	// bufconn / server closer preparation
	closer := func() {
		err = lis.Close()
		assert.NoError(t, err)

		server.Stop()
	}

	return client, closer
}
