package fxgrpcserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/ankorstore/yokai/grpcserver"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestNewGrpcServerRegistry(t *testing.T) {
	t.Parallel()

	param := fxgrpcserver.FxGrpcServiceRegistryParam{
		Services:    []any{},
		Definitions: []fxgrpcserver.GrpcServiceDefinition{},
	}

	registry := fxgrpcserver.NewFxGrpcServerRegistry(param)

	assert.IsType(t, &fxgrpcserver.GrpcServerRegistry{}, registry)
}

func TestResolveGrpcServicesSuccess(t *testing.T) {
	t.Parallel()

	service := grpcserver.NewGrpcHealthCheckService(nil)
	description := &grpc_health_v1.Health_ServiceDesc

	param := fxgrpcserver.FxGrpcServiceRegistryParam{
		Services: []any{service},
		Definitions: []fxgrpcserver.GrpcServiceDefinition{
			fxgrpcserver.NewGrpcServiceDefinition(fxgrpcserver.GetType(service), description),
		},
	}

	registry := fxgrpcserver.NewFxGrpcServerRegistry(param)

	resolvedServices, err := registry.ResolveGrpcServerServices()
	assert.NoError(t, err)

	assert.Len(t, resolvedServices, 1)
	assert.Equal(t, service, resolvedServices[0].Implementation())
	assert.Equal(t, description, resolvedServices[0].Description())
}

func TestResolveCheckerProbesRegistrationsFailure(t *testing.T) {
	t.Parallel()

	service := grpcserver.NewGrpcHealthCheckService(nil)
	description := &grpc_health_v1.Health_ServiceDesc

	param := fxgrpcserver.FxGrpcServiceRegistryParam{
		Services: []any{service},
		Definitions: []fxgrpcserver.GrpcServiceDefinition{
			fxgrpcserver.NewGrpcServiceDefinition("invalid", description),
		},
	}

	registry := fxgrpcserver.NewFxGrpcServerRegistry(param)

	_, err := registry.ResolveGrpcServerServices()
	assert.Error(t, err)
	assert.Equal(t, "cannot find grpc service implementation for type invalid", err.Error())
}
