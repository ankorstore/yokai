package fxgrpcserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestNewFxGrpcServerModuleInfo(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("./testdata/config"),
	)
	assert.NoError(t, err)

	grpcServer := &grpc.Server{}

	info := fxgrpcserver.NewFxGrpcServerModuleInfo(grpcServer, cfg)
	assert.IsType(t, &fxgrpcserver.FxGrpcServerModuleInfo{}, info)

	assert.Equal(t, fxgrpcserver.ModuleName, info.Name())
	assert.Equal(
		t,
		map[string]interface{}{
			"port":     fxgrpcserver.DefaultPort,
			"services": map[string]grpc.ServiceInfo{},
		},
		info.Data(),
	)
}
