package interceptor

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"google.golang.org/grpc"
)

type StreamInterceptor struct {
	config *config.Config
}

func NewStreamInterceptor(cfg *config.Config) *StreamInterceptor {
	return &StreamInterceptor{
		config: cfg,
	}
}

func (i *StreamInterceptor) HandleStream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.CtxLogger(ss.Context()).Info().Msgf("in stream interceptor of %s", i.config.AppName())

		return handler(srv, ss)
	}
}
