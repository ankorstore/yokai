package grpcserver

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GrpcServerFactory is the interface for [grpc.Server] factories.
type GrpcServerFactory interface {
	Create(options ...GrpcServerOption) (*grpc.Server, error)
}

// DefaultGrpcServerFactory is the default [GrpcServerFactory] implementation.
type DefaultGrpcServerFactory struct{}

// NewDefaultGrpcServerFactory returns a [DefaultGrpcServerFactory], implementing [GrpcServerFactory].
func NewDefaultGrpcServerFactory() GrpcServerFactory {
	return &DefaultGrpcServerFactory{}
}

// Create returns a new [grpc.Server], and accepts an optional list of [GrpcServerOption].
func (f *DefaultGrpcServerFactory) Create(options ...GrpcServerOption) (*grpc.Server, error) {
	appliedOpts := DefaultGrpcServerOptions()
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	grpcServer := grpc.NewServer(appliedOpts.ServerOptions...)

	if appliedOpts.Reflection {
		reflection.Register(grpcServer)
	}

	return grpcServer, nil
}
