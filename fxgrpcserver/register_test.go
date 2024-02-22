package fxgrpcserver_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/proto"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/service"
	"github.com/stretchr/testify/assert"
)

func TestAsGrpcService(t *testing.T) {
	t.Parallel()

	result := fxgrpcserver.AsGrpcServerService(service.NewTestServiceServer, &proto.Service_ServiceDesc)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result))
}
