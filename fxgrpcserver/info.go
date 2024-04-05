package fxgrpcserver

import (
	"github.com/ankorstore/yokai/config"
	"google.golang.org/grpc"
)

// FxGrpcServerModuleInfo is a module info collector for fxcore.
type FxGrpcServerModuleInfo struct {
	Address  string
	Services map[string]grpc.ServiceInfo
}

// NewFxGrpcServerModuleInfo returns a new [FxGrpcServerModuleInfo].
func NewFxGrpcServerModuleInfo(grpcServer *grpc.Server, cfg *config.Config) *FxGrpcServerModuleInfo {
	address := cfg.GetString("modules.grpc.server.address")
	if address == "" {
		address = DefaultAddress
	}

	return &FxGrpcServerModuleInfo{
		Address:  address,
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
		"address":  i.Address,
		"services": i.Services,
	}
}
