package fxgrpcserver

import "google.golang.org/grpc"

type ResolvedGrpcServerService struct {
	implementation any
	description    *grpc.ServiceDesc
}

func NewResolvedGrpcService(implementation any, description *grpc.ServiceDesc) *ResolvedGrpcServerService {
	return &ResolvedGrpcServerService{
		implementation: implementation,
		description:    description,
	}
}

func (r *ResolvedGrpcServerService) Implementation() any {
	return r.implementation
}

func (r *ResolvedGrpcServerService) Description() *grpc.ServiceDesc {
	return r.description
}
