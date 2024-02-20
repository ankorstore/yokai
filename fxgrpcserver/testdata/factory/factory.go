package factory

import (
	"github.com/ankorstore/yokai/grpcserver"
	"google.golang.org/grpc"
)

type TestGrpcServerFactory struct{}

func NewTestGrpcServerFactory() grpcserver.GrpcServerFactory {
	return &TestGrpcServerFactory{}
}

func (f *TestGrpcServerFactory) Create(options ...grpcserver.GrpcServerOption) (*grpc.Server, error) {
	return grpc.NewServer(), nil
}
