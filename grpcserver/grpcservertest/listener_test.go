package grpcservertest_test

import (
	"testing"

	"github.com/ankorstore/yokai/grpcserver/grpcservertest"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
)

func TestNewBufconnListener(t *testing.T) {
	t.Parallel()

	lis := grpcservertest.NewBufconnListener(1024)
	assert.IsType(t, &bufconn.Listener{}, lis)
}
