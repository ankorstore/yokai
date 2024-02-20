package fxgrpcserver

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func AsGrpcService(constructor any, description *grpc.ServiceDesc) fx.Option {
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
