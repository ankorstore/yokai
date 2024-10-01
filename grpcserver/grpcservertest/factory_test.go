package grpcservertest_test

import (
	"testing"

	"github.com/ankorstore/yokai/grpcserver/grpcservertest"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestDefaultTestBufconnConnectionFactory(t *testing.T) {
	t.Parallel()

	lis := grpcservertest.NewBufconnListener(1024)
	factory := grpcservertest.NewDefaultTestBufconnConnectionFactory(lis)

	t.Run("implements TestBufconnConnectionFactory", func(t *testing.T) {
		t.Parallel()

		assert.Implements(t, (*grpcservertest.TestBufconnConnectionFactory)(nil), factory)
	})

	t.Run("creates a connection", func(t *testing.T) {
		t.Parallel()

		conn, err := factory.Create(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		assert.NoError(t, err)
		assert.IsType(t, &grpc.ClientConn{}, conn)
	})
}
