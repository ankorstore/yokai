package fxgrpcserver

import (
	"fmt"

	"google.golang.org/grpc"

	"go.uber.org/fx"
)

// GrpcServerUnaryInterceptor is the interface for grpc server unary interceptors.
type GrpcServerUnaryInterceptor interface {
	HandleUnary() grpc.UnaryServerInterceptor
}

// GrpcServerStreamInterceptor is the interface for grpc server stream interceptors.
type GrpcServerStreamInterceptor interface {
	HandleStream() grpc.StreamServerInterceptor
}

// GrpcServerRegistry is the registry collecting grpc server options, interceptors, services and their definitions.
type GrpcServerRegistry struct {
	options            []grpc.ServerOption
	unaryInterceptors  []GrpcServerUnaryInterceptor
	streamInterceptors []GrpcServerStreamInterceptor
	services           []any
	definitions        []GrpcServerServiceDefinition
}

// FxGrpcServiceRegistryParam allows injection of the required dependencies in [NewFxGrpcServerRegistry].
type FxGrpcServiceRegistryParam struct {
	fx.In
	Options            []grpc.ServerOption           `group:"grpc-server-options"`
	UnaryInterceptors  []GrpcServerUnaryInterceptor  `group:"grpc-server-unary-interceptors"`
	StreamInterceptors []GrpcServerStreamInterceptor `group:"grpc-server-stream-interceptors"`
	Services           []any                         `group:"grpc-server-services"`
	Definitions        []GrpcServerServiceDefinition `group:"grpc-server-service-definitions"`
}

// NewFxGrpcServerRegistry returns as new [GrpcServerRegistry].
func NewFxGrpcServerRegistry(p FxGrpcServiceRegistryParam) *GrpcServerRegistry {
	return &GrpcServerRegistry{
		options:            p.Options,
		unaryInterceptors:  p.UnaryInterceptors,
		streamInterceptors: p.StreamInterceptors,
		services:           p.Services,
		definitions:        p.Definitions,
	}
}

// ResolveGrpcServerOptions resolves a list of grpc server options.
func (r *GrpcServerRegistry) ResolveGrpcServerOptions() []grpc.ServerOption {
	return r.options
}

// ResolveGrpcServerUnaryInterceptors resolves a list of grpc server unary interceptors.
func (r *GrpcServerRegistry) ResolveGrpcServerUnaryInterceptors() []GrpcServerUnaryInterceptor {
	return r.unaryInterceptors
}

// ResolveGrpcServerStreamInterceptors resolves a list of grpc server stream interceptors.
func (r *GrpcServerRegistry) ResolveGrpcServerStreamInterceptors() []GrpcServerStreamInterceptor {
	return r.streamInterceptors
}

// ResolveGrpcServerServices resolves a list of [ResolvedGrpcServerService] from their definitions.
func (r *GrpcServerRegistry) ResolveGrpcServerServices() ([]*ResolvedGrpcServerService, error) {
	var grpcServices []*ResolvedGrpcServerService

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
