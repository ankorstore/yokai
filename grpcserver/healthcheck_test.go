package grpcserver_test

import (
	"context"
	"net"
	"testing"

	"github.com/ankorstore/yokai/generate/generatetest/uuid"
	"github.com/ankorstore/yokai/grpcserver"
	"github.com/ankorstore/yokai/grpcserver/grpcservertest"
	"github.com/ankorstore/yokai/grpcserver/testdata/probes"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestCheckSuccess(t *testing.T) {
	// checker
	checker, err := healthcheck.NewDefaultCheckerFactory().Create(
		healthcheck.WithProbe(probes.NewSuccessProbe()),
	)
	assert.NoError(t, err)

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// startup call assertions
	client, closer := prepareHealthCheckServiceGrpcServerAndClient(t, checker, logger)

	response, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "test::startup"})
	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, response.Status)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"kind":    "startup",
		"caller":  "test::startup",
		"message": "grpc health check success",
	})

	closer()

	// liveness call assertions
	client, closer = prepareHealthCheckServiceGrpcServerAndClient(t, checker, logger)

	response, err = client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "test::liveness"})
	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, response.Status)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"kind":    "liveness",
		"caller":  "test::liveness",
		"message": "grpc health check success",
	})

	closer()

	// readiness call assertions
	client, closer = prepareHealthCheckServiceGrpcServerAndClient(t, checker, logger)

	response, err = client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "test::readiness"})
	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, response.Status)

	// readiness call logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"kind":    "readiness",
		"caller":  "test::readiness",
		"message": "grpc health check success",
	})

	closer()
}

func TestCheckFailure(t *testing.T) {
	// checker
	checker, err := healthcheck.NewDefaultCheckerFactory().Create(
		healthcheck.WithProbe(probes.NewSuccessProbe()),
		healthcheck.WithProbe(probes.NewFailureProbe()),
	)
	assert.NoError(t, err)

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// client
	client, closer := prepareHealthCheckServiceGrpcServerAndClient(t, checker, logger)
	defer closer()

	// call assertions
	response, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "test"})
	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING, response.Status)

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "error",
		"kind":         "startup",
		"caller":       "test",
		"successProbe": "success: true, message: some success",
		"failureProbe": "success: false, message: some failure",
		"message":      "grpc health check failure",
	})
}

func TestWatch(t *testing.T) {
	t.Parallel()

	// checker
	checker, err := healthcheck.NewDefaultCheckerFactory().Create(
		healthcheck.WithProbe(probes.NewSuccessProbe()),
	)
	assert.NoError(t, err)

	// logger
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// client
	client, closer := prepareHealthCheckServiceGrpcServerAndClient(t, checker, logger)
	defer closer()

	// call assertions
	stream, err := client.Watch(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "test"})
	assert.NoError(t, err)

	_, err = stream.Recv()
	assert.Error(t, err)
	assert.Equal(t, "rpc error: code = Unimplemented desc = watch is not implemented", err.Error())

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"caller":  "test",
		"message": "grpc health watch not implemented",
	})
}

func prepareHealthCheckServiceGrpcServerAndClient(t *testing.T, checker *healthcheck.Checker, logger *log.Logger) (grpc_health_v1.HealthClient, func()) {
	t.Helper()

	// context preparation
	ctx := logger.WithContext(context.Background())

	// bufconn listener preparation
	lis := grpcservertest.NewBufconnListener(1024 * 1024)

	// gRPC server preparation
	loggerInterceptor := grpcserver.NewGrpcLoggerInterceptor(uuid.NewTestUuidGenerator("test"), logger)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor.UnaryInterceptor()),
		grpc.StreamInterceptor(loggerInterceptor.StreamInterceptor()),
	)

	server.RegisterService(
		&grpc_health_v1.Health_ServiceDesc,
		grpcserver.NewGrpcHealthCheckService(checker),
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

	client := grpc_health_v1.NewHealthClient(conn)

	// bufconn / server closer preparation
	closer := func() {
		err = lis.Close()
		assert.NoError(t, err)

		server.Stop()
	}

	return client, closer
}
