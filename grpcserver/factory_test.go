package grpcserver_test

import (
	"context"
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
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	reflection "google.golang.org/grpc/reflection/grpc_reflection_v1"
)

func TestDefaultGrpcServerFactory(t *testing.T) {
	t.Parallel()

	factory := grpcserver.NewDefaultGrpcServerFactory()

	assert.IsType(t, &grpcserver.DefaultGrpcServerFactory{}, factory)
	assert.Implements(t, (*grpcserver.GrpcServerFactory)(nil), factory)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	// factory
	factory := grpcserver.NewDefaultGrpcServerFactory()

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// logger interceptor
	loggerInterceptor := grpcserver.NewGrpcLoggerInterceptor(
		uuid.NewTestUuidGenerator("default-request-id"),
		logger,
	)

	// panic handler
	panicHandler := grpcserver.NewGrpcPanicRecoveryHandler()

	// server
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		interceptor.TestUnaryInterceptor(testTraceId, testSpanId),
		loggerInterceptor.UnaryInterceptor(),
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(panicHandler.Handle(true))),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		interceptor.TestStreamInterceptor(testTraceId, testSpanId),
		loggerInterceptor.StreamInterceptor(),
		recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(panicHandler.Handle(true))),
	}

	server, err := factory.Create(
		grpcserver.WithServerOptions(
			grpc.ChainUnaryInterceptor(unaryInterceptors...),
			grpc.ChainStreamInterceptor(streamInterceptors...),
		),
		grpcserver.WithReflection(true),
	)

	server.RegisterService(
		&proto.Service_ServiceDesc,
		service.NewTestServiceServer(),
	)

	// bufconn listener preparation
	lis := grpcservertest.NewBufconnListener(1024 * 1024)

	go func() {
		//nolint:errcheck
		server.Serve(lis)
	}()

	defer func() {
		err = lis.Close()
		assert.NoError(t, err)

		server.Stop()
	}()

	// context
	ctx := logger.WithContext(context.Background())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", testRequestId)

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

	// service client
	serviceClient := proto.NewServiceClient(conn)

	// service client 1st call assertion
	_, err = serviceClient.Unary(ctx, &proto.Request{
		ShouldPanic: true,
		Message:     "first attempt",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rpc error: code = Internal desc = internal grpc server error, panic = first attempt")

	// service client 2nd call assertion (server should still work due to recovered panic)
	_, err = serviceClient.Unary(ctx, &proto.Request{
		ShouldPanic: true,
		Message:     "second attempt",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rpc error: code = Internal desc = internal grpc server error, panic = second attempt")

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"panic":     "first attempt",
		"message":   "grpc recovered from panic",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"panic":     "second attempt",
		"message":   "grpc recovered from panic",
		"requestID": testRequestId,
		"traceID":   testTraceId,
	})

	// reflection client
	reflectionClient := reflection.NewServerReflectionClient(conn)

	// reflection client call assertion
	stream, err := reflectionClient.ServerReflectionInfo(context.Background(), grpc.WaitForReady(true))
	assert.NoError(t, err)

	err = stream.Send(&reflection.ServerReflectionRequest{
		MessageRequest: &reflection.ServerReflectionRequest_ListServices{},
	})
	assert.NoError(t, err)

	resp, err := stream.Recv()
	assert.NoError(t, err)

	switch resp.MessageResponse.(type) {
	case *reflection.ServerReflectionResponse_ListServicesResponse:
		expectedServices := []string{
			"test.Service",
			"grpc.reflection.v1.ServerReflection",
			"grpc.reflection.v1alpha.ServerReflection",
		}

		reflectionServices := resp.GetListServicesResponse().Service

		assert.Len(t, reflectionServices, len(expectedServices))

		for _, e := range reflectionServices {
			assert.Contains(t, expectedServices, e.Name)
		}
	default:
		t.Errorf("invalid type, expected %s", resp.MessageResponse)
	}

	// logs assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"grpcMethod": "/grpc.reflection.v1.ServerReflection/ServerReflectionInfo",
		"grpcType":   "server-streaming",
		"message":    "grpc call start",
		"requestID":  "default-request-id",
		"traceID":    testTraceId,
	})
}
