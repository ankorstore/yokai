package fxgrpcserver

import "google.golang.org/grpc"

type ResolvedGrpcService struct {
	implementation any
	description    *grpc.ServiceDesc
}

func NewResolvedGrpcService(implementation any, description *grpc.ServiceDesc) *ResolvedGrpcService {
	return &ResolvedGrpcService{
		implementation: implementation,
		description:    description,
	}
}

func (r *ResolvedGrpcService) Implementation() any {
	return r.implementation
}

func (r *ResolvedGrpcService) Description() *grpc.ServiceDesc {
	return r.description
}
