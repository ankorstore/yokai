package grpcservertest

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var _ TestBufconnConnectionFactory = (*DefaultTestBufconnConnectionFactory)(nil)

type TestBufconnConnectionFactory interface {
	Create(opts ...grpc.DialOption) (*grpc.ClientConn, error)
}

type DefaultTestBufconnConnectionFactory struct {
	lis *bufconn.Listener
}

func NewDefaultTestBufconnConnectionFactory(lis *bufconn.Listener) *DefaultTestBufconnConnectionFactory {
	return &DefaultTestBufconnConnectionFactory{
		lis: lis,
	}
}

func (f *DefaultTestBufconnConnectionFactory) Create(opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialOptions := append(
		[]grpc.DialOption{
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return f.lis.Dial()
			}),
		},
		opts...,
	)

	return grpc.NewClient(fmt.Sprintf("passthrough://%s", f.lis.Addr().String()), dialOptions...)
}
