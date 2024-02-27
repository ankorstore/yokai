package interceptor

import (
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/service"
	"github.com/ankorstore/yokai/log"
	"google.golang.org/grpc"
)

type StreamInterceptor struct {
	dependency *service.TestServiceDependency
}

func NewStreamInterceptor(dependency *service.TestServiceDependency) *StreamInterceptor {
	return &StreamInterceptor{
		dependency: dependency,
	}
}

func (i *StreamInterceptor) HandleStream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.CtxLogger(ss.Context()).Info().Msgf("in stream interceptor of %s", i.dependency.AppName())

		return handler(srv, ss)
	}
}
