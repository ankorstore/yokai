package fxgrpcserver

import (
	"fmt"
	"google.golang.org/grpc"

	"go.uber.org/fx"
)

type GrpcServerUnaryInterceptor interface {
	HandleUnary() grpc.UnaryServerInterceptor
}

type GrpcServerStreamInterceptor interface {
	HandleStream() grpc.StreamServerInterceptor
}

type GrpcServerRegistry struct {
	options            []grpc.ServerOption
	unaryInterceptors  []GrpcServerUnaryInterceptor
	streamInterceptors []GrpcServerStreamInterceptor
	services           []any
	definitions        []GrpcServerServiceDefinition
}

type FxGrpcServiceRegistryParam struct {
	fx.In
	Options            []grpc.ServerOption           `group:"grpc-server-options"`
	UnaryInterceptors  []GrpcServerUnaryInterceptor  `group:"grpc-server-unary-interceptors"`
	StreamInterceptors []GrpcServerStreamInterceptor `group:"grpc-server-stream-interceptors"`
	Services           []any                         `group:"grpc-server-services"`
	Definitions        []GrpcServerServiceDefinition `group:"grpc-server-service-definitions"`
}

func NewFxGrpcServerRegistry(p FxGrpcServiceRegistryParam) *GrpcServerRegistry {
	return &GrpcServerRegistry{
		options:            p.Options,
		unaryInterceptors:  p.UnaryInterceptors,
		streamInterceptors: p.StreamInterceptors,
		services:           p.Services,
		definitions:        p.Definitions,
	}
}

func (r *GrpcServerRegistry) ResolveGrpcServerOptions() []grpc.ServerOption {
	return r.options
}

func (r *GrpcServerRegistry) ResolveGrpcServerUnaryInterceptors() []GrpcServerUnaryInterceptor {
	return r.unaryInterceptors
}

func (r *GrpcServerRegistry) ResolveGrpcServerStreamInterceptors() []GrpcServerStreamInterceptor {
	return r.streamInterceptors
}

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
