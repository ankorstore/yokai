package grpcserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/ankorstore/yokai/healthcheck"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// GrpcHealthCheckService is a default gRPC health check server implementation working with the [healthcheck.Checker].
type GrpcHealthCheckService struct {
	grpc_health_v1.UnimplementedHealthServer
	checker *healthcheck.Checker
}

// NewGrpcHealthCheckService returns a new [GrpcHealthCheckService] instance.
func NewGrpcHealthCheckService(checker *healthcheck.Checker) *GrpcHealthCheckService {
	return &GrpcHealthCheckService{
		checker: checker,
	}
}

// Check performs checks on the registered [healthcheck.CheckerProbe].
func (s *GrpcHealthCheckService) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	logger := CtxLogger(ctx)

	serviceName := strings.ToLower(in.Service)

	var kind healthcheck.ProbeKind
	switch {
	case strings.Contains(serviceName, healthcheck.Liveness.String()):
		kind = healthcheck.Liveness
	case strings.Contains(serviceName, healthcheck.Readiness.String()):
		kind = healthcheck.Readiness
	default:
		kind = healthcheck.Startup
	}

	result := s.checker.Check(ctx, kind)
	if !result.Success {
		evt := logger.Error()
		evt.
			Str("kind", kind.String()).
			Str("caller", serviceName)

		for probeName, probeResult := range result.ProbesResults {
			evt.Str(probeName, fmt.Sprintf("success: %v, message: %s", probeResult.Success, probeResult.Message))
		}

		evt.Msg("grpc health check failure")

		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	logger.
		Info().
		Str("kind", kind.String()).
		Str("caller", serviceName).
		Msg("grpc health check success")

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is not implemented.
func (s *GrpcHealthCheckService) Watch(in *grpc_health_v1.HealthCheckRequest, watchServer grpc_health_v1.Health_WatchServer) error {
	CtxLogger(watchServer.Context()).
		Warn().
		Str("caller", in.Service).
		Msg("grpc health watch not implemented")

	return status.Error(codes.Unimplemented, "watch is not implemented")
}
