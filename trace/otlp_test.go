package trace_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/ankorstore/yokai/trace"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/test/bufconn"
)

func TestNewOtlpGrpcInsecureConnectionSuccess(t *testing.T) {
	t.Parallel()

	listener := bufconn.Listen(1)
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			t.Error(err)
		}
	}()

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := trace.NewOtlpGrpcClientConnection(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(bufDialer),
	)
	assert.NoError(t, err)

	assert.Equal(t, connectivity.Ready, conn.GetState())

	err = conn.Close()
	assert.NoError(t, err)

	grpcServer.GracefulStop()
}

func TestNewOtlpGrpcInsecureConnectionError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Microsecond))
	defer cancel()

	_, err := trace.NewOtlpGrpcClientConnection(ctx, "https://example.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}
