package fxgrpcserver

import "google.golang.org/grpc"

// ResolvedGrpcServerService is an interface for the resolved grpc server services.
type ResolvedGrpcServerService struct {
	implementation any
	description    *grpc.ServiceDesc
}

// NewResolvedGrpcService returns a new [ResolvedGrpcServerService].
func NewResolvedGrpcService(implementation any, description *grpc.ServiceDesc) *ResolvedGrpcServerService {
	return &ResolvedGrpcServerService{
		implementation: implementation,
		description:    description,
	}
}

// Implementation return the resolved grpc server service implementation.
func (r *ResolvedGrpcServerService) Implementation() any {
	return r.implementation
}

// Description return the resolved grpc server service description.
func (r *ResolvedGrpcServerService) Description() *grpc.ServiceDesc {
	return r.description
}
