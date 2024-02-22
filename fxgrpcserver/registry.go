package fxgrpcserver

import (
	"fmt"
	"google.golang.org/grpc"

	"go.uber.org/fx"
)

type GrpcServerRegistry struct {
	options     []grpc.ServerOption
	services    []any
	definitions []GrpcServiceDefinition
}

type FxGrpcServiceRegistryParam struct {
	fx.In
	Options     []grpc.ServerOption     `group:"grpc-server-options"`
	Services    []any                   `group:"grpc-server-services"`
	Definitions []GrpcServiceDefinition `group:"grpc-server-service-definitions"`
}

func NewFxGrpcServerRegistry(p FxGrpcServiceRegistryParam) *GrpcServerRegistry {
	return &GrpcServerRegistry{
		options:     p.Options,
		services:    p.Services,
		definitions: p.Definitions,
	}
}

func (r *GrpcServerRegistry) ResolveGrpcServerOptions() []grpc.ServerOption {
	return r.options
}

func (r *GrpcServerRegistry) ResolveGrpcServerServices() ([]*ResolvedGrpcService, error) {
	var grpcServices []*ResolvedGrpcService

	for _, definition := range r.definitions {
		implementation, err := r.lookupRegisteredServiceImplementation(definition.ReturnType())
		if err != nil {
			return nil, err
		}

		grpcServices = append(grpcServices, NewResolvedGrpcService(implementation, definition.Description()))
	}

	return grpcServices, nil
}

func (r *GrpcServerRegistry) lookupRegisteredServiceImplementation(returnType string) (any, error) {
	for _, implementation := range r.services {
		if GetType(implementation) == returnType {
			return implementation, nil
		}
	}

	return nil, fmt.Errorf("cannot find grpc service implementation for type %s", returnType)
}
