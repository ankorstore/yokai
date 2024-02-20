package grpcserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/grpcserver"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestWithServerOptions(t *testing.T) {
	t.Parallel()

	opt := grpcserver.DefaultGrpcServerOptions()
	grpcserver.WithServerOptions(&grpc.EmptyServerOption{})(&opt)

	assert.Equal(t, []grpc.ServerOption{&grpc.EmptyServerOption{}}, opt.ServerOptions)
}

func TestWithReflection(t *testing.T) {
	t.Parallel()

	opt := grpcserver.DefaultGrpcServerOptions()
	grpcserver.WithReflection(true)(&opt)

	assert.True(t, opt.Reflection)
}
