package fxgrpcserver_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/interceptor"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/proto"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestAsGrpcServerOptions(t *testing.T) {
	t.Parallel()

	result := fxgrpcserver.AsGrpcServerOptions(grpc.MaxRecvMsgSize(10))

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result))
}

func TestAsGrpcServerUnaryInterceptor(t *testing.T) {
	t.Parallel()

	result := fxgrpcserver.AsGrpcServerUnaryInterceptor(interceptor.NewUnaryInterceptor)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", result))
}

func TestAsGrpcServerStreamInterceptor(t *testing.T) {
	t.Parallel()

	result := fxgrpcserver.AsGrpcServerStreamInterceptor(interceptor.NewStreamInterceptor)

	assert.Equal(t, "fx.provideOption", fmt.Sprintf("%T", result))
}

func TestAsGrpcServerService(t *testing.T) {
	t.Parallel()

	result := fxgrpcserver.AsGrpcServerService(service.NewTestServiceServer, &proto.Service_ServiceDesc)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result))
}
