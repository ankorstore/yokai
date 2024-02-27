package interceptor

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"google.golang.org/grpc"
)

type UnaryInterceptor struct {
	config *config.Config
}

func NewUnaryInterceptor(cfg *config.Config) *UnaryInterceptor {
	return &UnaryInterceptor{
		config: cfg,
	}
}

func (i *UnaryInterceptor) HandleUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log.CtxLogger(ctx).Info().Msgf("in unary interceptor of %s", i.config.AppName())

		return handler(ctx, req)
	}
}
