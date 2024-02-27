package fxgrpcserver

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// AsGrpcServerOptions registers a list of grpc server options into Fx.
func AsGrpcServerOptions(options ...grpc.ServerOption) fx.Option {
	var serverOptions []fx.Option

	for _, option := range options {
		serverOptions = append(
			serverOptions,
			fx.Supply(
				fx.Annotate(
					option,
					fx.As(new(grpc.ServerOption)),
					fx.ResultTags(`group:"grpc-server-options"`),
				),
			),
		)
	}

	return fx.Options(serverOptions...)
}

// AsGrpcServerUnaryInterceptor registers a grpc server unary interceptor into Fx.
func AsGrpcServerUnaryInterceptor(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(GrpcServerUnaryInterceptor)),
			fx.ResultTags(`group:"grpc-server-unary-interceptors"`),
		),
	)
}

// AsGrpcServerStreamInterceptor registers a grpc server stream interceptor into Fx.
func AsGrpcServerStreamInterceptor(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(GrpcServerStreamInterceptor)),
			fx.ResultTags(`group:"grpc-server-stream-interceptors"`),
		),
	)
}

// AsGrpcServerService registers a grpc server service into Fx.
func AsGrpcServerService(constructor any, description *grpc.ServiceDesc) fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				constructor,
				fx.As(new(interface{})),
				fx.ResultTags(`group:"grpc-server-services"`),
			),
		),
		fx.Supply(
			fx.Annotate(
				NewGrpcServiceDefinition(GetReturnType(constructor), description),
				fx.As(new(GrpcServerServiceDefinition)),
				fx.ResultTags(`group:"grpc-server-service-definitions"`),
			),
		),
	)
}
