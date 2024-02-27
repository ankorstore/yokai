package fxgrpcserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/ankorstore/yokai/grpcserver"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestNewResolvedGrpcService(t *testing.T) {
	t.Parallel()

	service := grpcserver.NewGrpcHealthCheckService(nil)
	description := &grpc_health_v1.Health_ServiceDesc

	resolvedService := fxgrpcserver.NewResolvedGrpcService(service, description)

	assert.IsType(t, &fxgrpcserver.ResolvedGrpcServerService{}, resolvedService)
	assert.Equal(t, service, resolvedService.Implementation())
	assert.Equal(t, description, resolvedService.Description())
}
