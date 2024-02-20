package grpcserver

import (
	"google.golang.org/grpc"
)

// Options are options for the [GrpcServerFactory] implementations.
type Options struct {
	ServerOptions []grpc.ServerOption
	Reflection    bool
}

// DefaultGrpcServerOptions are the default options used in the [DefaultGrpcServerFactory].
func DefaultGrpcServerOptions() Options {
	return Options{
		ServerOptions: []grpc.ServerOption{},
		Reflection:    false,
	}
}

// GrpcServerOption are functional options for the [GrpcServerFactory] implementations.
type GrpcServerOption func(o *Options)

// WithServerOptions is used to configure a list of [grpc.ServerOption].
func WithServerOptions(s ...grpc.ServerOption) GrpcServerOption {
	return func(o *Options) {
		o.ServerOptions = s
	}
}

// WithReflection is used to enable gRPC server reflection.
func WithReflection(r bool) GrpcServerOption {
	return func(o *Options) {
		o.Reflection = r
	}
}
