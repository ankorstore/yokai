package fxgrpcserver

import (
	"github.com/ankorstore/yokai/config"
	"google.golang.org/grpc"
)

// FxGrpcServerModuleInfo is a module info collector for fxcore.
type FxGrpcServerModuleInfo struct {
	Port     int
	Services map[string]grpc.ServiceInfo
}

// NewFxGrpcServerModuleInfo returns a new [FxGrpcServerModuleInfo].
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

// Name return the name of the module info.
func (i *FxGrpcServerModuleInfo) Name() string {
	return ModuleName
}

// Data return the data of the module info.
func (i *FxGrpcServerModuleInfo) Data() map[string]interface{} {
	return map[string]interface{}{
		"port":     i.Port,
		"services": i.Services,
	}
}
