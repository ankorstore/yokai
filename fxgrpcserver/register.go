package fxgrpcserver

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

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
				fx.As(new(GrpcServiceDefinition)),
				fx.ResultTags(`group:"grpc-server-service-definitions"`),
			),
		),
	)
}

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
