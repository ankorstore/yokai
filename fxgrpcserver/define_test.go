package fxgrpcserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/ankorstore/yokai/fxgrpcserver/testdata/proto"
	"github.com/stretchr/testify/assert"
)

func TestNewGrpcServerServiceDefinition(t *testing.T) {
	t.Parallel()

	definition := fxgrpcserver.NewGrpcServiceDefinition("*TestService", &proto.Service_ServiceDesc)

	assert.Implements(t, (*fxgrpcserver.GrpcServerServiceDefinition)(nil), definition)
	assert.Equal(t, "*TestService", definition.ReturnType())
	assert.Equal(t, &proto.Service_ServiceDesc, definition.Description())
}
