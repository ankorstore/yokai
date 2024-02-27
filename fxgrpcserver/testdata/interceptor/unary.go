package interceptor

import (
	"context"

	"github.com/ankorstore/yokai/fxgrpcserver/testdata/service"
	"github.com/ankorstore/yokai/log"
	"google.golang.org/grpc"
)

type UnaryInterceptor struct {
	dependency *service.TestServiceDependency
}

func NewUnaryInterceptor(dependency *service.TestServiceDependency) *UnaryInterceptor {
	return &UnaryInterceptor{
		dependency: dependency,
	}
}

func (i *UnaryInterceptor) HandleUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log.CtxLogger(ctx).Info().Msgf("in unary interceptor of %s", i.dependency.AppName())

		return handler(ctx, req)
	}
}
