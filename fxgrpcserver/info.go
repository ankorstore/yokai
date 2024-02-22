package fxgrpcserver

import (
	"github.com/ankorstore/yokai/config"
	"google.golang.org/grpc"
)

type FxGrpcServerModuleInfo struct {
	Port     int
	Services map[string]grpc.ServiceInfo
}

func NewFxGrpcServerModuleInfo(grpcServer *grpc.Server, cfg *config.Config) *FxGrpcServerModuleInfo {
	port := cfg.GetInt("modules.grpc.server.port")
	if port == 0 {
		port = DefaultPort
	}

	return &FxGrpcServerModuleInfo{
		Port:     port,
		Services: grpcServer.GetServiceInfo(),
	}
}

func (i *FxGrpcServerModuleInfo) Name() string {
	return ModuleName
}

func (i *FxGrpcServerModuleInfo) Data() map[string]interface{} {
	return map[string]interface{}{
		"port":     i.Port,
		"services": i.Services,
	}
}
