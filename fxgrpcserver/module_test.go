package fxgrpcserver_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/factory"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/interceptor"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/probes"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/proto"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/service"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/grpcserver/grpcservertest"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

var (
	// request headers parts.
	testRequestId   = "33084b3e-9b90-926c-af19-3859d70bd296"
	testTraceId     = "c4ca71e03e42c2c3d54293a6e2608bfa"
	testSpanId      = "8d0fdc8a74baaaea"
	testTraceParent = fmt.Sprintf("00-%s-%s-01", testTraceId, testSpanId)
)

//nolint:maintidx
func TestModule(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("METRICS_NAMESPACE", "foo")
	t.Setenv("METRICS_SUBSYSTEM", "bar")
	t.Setenv("APP_ENV", "test")

	var grpcServer *grpc.Server
	var connFactory grpcservertest.TestBufconnConnectionFactory
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxgenerate.FxGenerateModule,
		fxmetrics.FxMetricsModule,
		fxhealthcheck.FxHealthcheckModule,
		fxgrpcserver.FxGrpcServerModule,
		fx.Provide(service.NewTestServiceDependency),
		fx.Options(
			fxgrpcserver.AsGrpcServerUnaryInterceptor(interceptor.NewUnaryInterceptor),
			fxgrpcserver.AsGrpcServerStreamInterceptor(interceptor.NewStreamInterceptor),
			fxgrpcserver.AsGrpcServerService(service.NewTestServiceServer, &proto.Service_ServiceDesc),
		),
		fx.Populate(&grpcServer, &connFactory, &logBuffer, &traceExporter, &metricsRegistry),
	).RequireStart().RequireStop()

	defer func() {
		grpcServer.GracefulStop()
	}()

	// conn preparation
	conn, err := connFactory.Create(
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)

	// client preparation
	client := proto.NewServiceClient(conn)

	// context preparation
	unaryCtx := context.Background()
	unaryCtx = metadata.AppendToOutgoingContext(unaryCtx, "x-request-id", testRequestId)
	unaryCtx = metadata.AppendToOutgoingContext(unaryCtx, "traceparent", testTraceParent)
	unaryCtx = metadata.AppendToOutgoingContext(unaryCtx, "x-foo", "foo")

	// unary call assertions
	response, err := client.Unary(unaryCtx, &proto.Request{
		ShouldFail: false,
		Message:    "test",
	})
	assert.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, "test received on test", response.Message)

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "debug",
		"system":     "grpcserver",
		"service":    "test",
		"grpcMethod": "/test.Service/Unary",
		"grpcType":   "unary",
		"message":    "grpc call start",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
		"foo":        "foo",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"system":    "grpcserver",
		"service":   "test",
		"message":   "in unary interceptor of test",
		"requestID": testRequestId,
		"foo":       "foo",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"system":    "grpcserver",
		"service":   "test",
		"message":   "unary call on test",
		"requestID": testRequestId,
		"foo":       "foo",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"system":    "grpcserver",
		"service":   "test",
		"message":   "unary call success on test",
		"requestID": testRequestId,
		"foo":       "foo",
	})

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"system":     "grpcserver",
		"service":    "test",
		"grpcStatus": "OK",
		"message":    "grpc call success",
		"requestID":  testRequestId,
		"foo":        "foo",
	})

	tracetest.AssertHasTraceSpan(t, traceExporter, "unary trace on test")
	tracetest.AssertHasNotTraceSpan(t, traceExporter, "test.Service/Unary")

	expectedUnaryMetric := `
		# HELP foo_bar_grpc_server_started_total Total number of RPCs started on the server.
		# TYPE foo_bar_grpc_server_started_total counter
		foo_bar_grpc_server_started_total{grpc_method="Unary",grpc_service="test.Service",grpc_type="unary"} 1
	`

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedUnaryMetric),
		"foo_bar_grpc_server_started_total",
	)
	assert.NoError(t, err)

	// log buffer reset
	logBuffer.Reset()

	// context preparation
	bidiCtx := context.Background()
	bidiCtx = metadata.AppendToOutgoingContext(bidiCtx, "x-request-id", testRequestId)
	bidiCtx = metadata.AppendToOutgoingContext(bidiCtx, "traceparent", testTraceParent)
	bidiCtx = metadata.AppendToOutgoingContext(bidiCtx, "x-foo", "foo")

	// bidi call assertions
	stream, err := client.Bidi(bidiCtx)
	assert.NoError(t, err)

	wait := make(chan struct{})

	go func() {
		err = stream.Send(&proto.Request{
			ShouldFail: false,
			Message:    "this is a test",
		})
		assert.NoError(t, err)

		err = stream.CloseSend()
		assert.NoError(t, err)
	}()

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

	assert.Len(t, responses, 4)
	assert.Equal(t, "this", responses[0].Message)
	assert.Equal(t, "is", responses[1].Message)
	assert.Equal(t, "a", responses[2].Message)
	assert.Equal(t, "test", responses[3].Message)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"system":     "grpcserver",
		"service":    "test",
		"grpcMethod": "/test.Service/Bidi",
		"grpcType":   "server-streaming",
		"message":    "grpc call start",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
		"foo":        "foo",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"system":    "grpcserver",
		"message":   "in stream interceptor of test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"foo":       "foo",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"system":    "grpcserver",
		"message":   "bidi call on test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"foo":       "foo",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"system":    "grpcserver",
		"message":   "bidi recv value this is a test on test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"foo":       "foo",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"system":    "grpcserver",
		"message":   "bidi send value this on test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"foo":       "foo",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"system":    "grpcserver",
		"message":   "bidi send value is on test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"foo":       "foo",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"system":    "grpcserver",
		"message":   "bidi send value a on test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"foo":       "foo",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"system":    "grpcserver",
		"message":   "bidi send value test on test",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"foo":       "foo",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":      "info",
		"system":     "grpcserver",
		"grpcStatus": "OK",
		"message":    "grpc call success",
		"requestID":  testRequestId,
		"traceID":    testTraceId,
		"foo":        "foo",
	})

	tracetest.AssertHasTraceSpan(t, traceExporter, "bidi trace on test")
	tracetest.AssertHasTraceSpan(t, traceExporter, "test.Service/Bidi")

	expectedBidiMetric := `
		# HELP foo_bar_grpc_server_handled_total Total number of RPCs completed on the server, regardless of success or failure.
		# TYPE foo_bar_grpc_server_handled_total counter
		foo_bar_grpc_server_handled_total{grpc_code="OK",grpc_method="Unary",grpc_service="test.Service",grpc_type="unary"} 1
		foo_bar_grpc_server_handled_total{grpc_code="OK",grpc_method="Bidi",grpc_service="test.Service",grpc_type="bidi_stream"} 1
	`

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedBidiMetric),
		"foo_bar_grpc_server_handled_total",
	)
	assert.NoError(t, err)
}

func TestModuleHealthCheck(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_ENV", "test")

	var grpcServer *grpc.Server
	var connFactory grpcservertest.TestBufconnConnectionFactory
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxgenerate.FxGenerateModule,
		fxmetrics.FxMetricsModule,
		fxhealthcheck.FxHealthcheckModule,
		fxgrpcserver.FxGrpcServerModule,
		fx.Options(
			fxhealthcheck.AsCheckerProbe(probes.NewSuccessProbe),
			fxhealthcheck.AsCheckerProbe(probes.NewFailureProbe, healthcheck.Liveness, healthcheck.Readiness),
		),
		fx.Populate(&grpcServer, &connFactory, &logBuffer, &traceExporter, &metricsRegistry),
	).RequireStart().RequireStop()

	defer func() {
		grpcServer.GracefulStop()
	}()

	// conn preparation
	conn, err := connFactory.Create(
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)

	// client preparation
	client := grpc_health_v1.NewHealthClient(conn)

	// context preparation
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", testRequestId)
	ctx = metadata.AppendToOutgoingContext(ctx, "traceparent", testTraceParent)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-foo", "foo")

	// startup call assertions
	response, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "test::startup"})
	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, response.Status)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "info",
		"system":    "grpcserver",
		"kind":      "startup",
		"caller":    "test::startup",
		"message":   "grpc health check success",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"foo":       "foo",
	})

	assert.True(t, traceExporter.HasSpan("grpc.health.v1.Health/Check"))

	// liveness call assertions
	response, err = client.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "test::liveness"})
	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING, response.Status)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"system":    "grpcserver",
		"kind":      "liveness",
		"caller":    "test::liveness",
		"message":   "grpc health check failure",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"foo":       "foo",
	})

	assert.True(t, traceExporter.HasSpan("grpc.health.v1.Health/Check"))

	// readiness call assertions
	response, err = client.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "test::readiness"})
	assert.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_NOT_SERVING, response.Status)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "error",
		"system":    "grpcserver",
		"kind":      "readiness",
		"caller":    "test::readiness",
		"message":   "grpc health check failure",
		"requestID": testRequestId,
		"traceID":   testTraceId,
		"foo":       "foo",
	})

	assert.True(t, traceExporter.HasSpan("grpc.health.v1.Health/Check"))

	// metrics
	expectedMetrics := `
		# HELP grpc_server_handled_total Total number of RPCs completed on the server, regardless of success or failure.
		# TYPE grpc_server_handled_total counter
		grpc_server_handled_total{grpc_code="OK",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 3
	`

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetrics),
		"grpc_server_handled_total",
	)
	assert.NoError(t, err)
}

func TestModuleDecoration(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_ENV", "test")

	// reflection is enabled, but the custom factory should ignore this
	t.Setenv("REFLECTION_ENABLED", "true")

	var grpcServer *grpc.Server

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxgenerate.FxGenerateModule,
		fxmetrics.FxMetricsModule,
		fxhealthcheck.FxHealthcheckModule,
		fxgrpcserver.FxGrpcServerModule,
		fx.Decorate(factory.NewTestGrpcServerFactory),
		fx.Populate(&grpcServer),
	).RequireStart().RequireStop()

	info := grpcServer.GetServiceInfo()

	// the service info returns only the health check, since decorated factory ignored reflection option
	assert.Len(t, info, 1)
	assert.Contains(t, fmt.Sprintf("%+v", info), "grpc.health.v1.Health")
	assert.NotContains(t, fmt.Sprintf("%+v", info), "grpc.reflection.v1alpha.ServerReflection")
}
