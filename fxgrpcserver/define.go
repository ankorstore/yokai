package fxgrpcserver

import "google.golang.org/grpc"

type GrpcServerServiceDefinition interface {
	ReturnType() string
	Description() *grpc.ServiceDesc
}

type grpcServerServiceDefinition struct {
	returnType  string
	description *grpc.ServiceDesc
}

func NewGrpcServiceDefinition(returnType string, description *grpc.ServiceDesc) GrpcServerServiceDefinition {
	return &grpcServerServiceDefinition{
		returnType:  returnType,
		description: description,
	}
}

func (d *grpcServerServiceDefinition) ReturnType() string {
	return d.returnType
}

func (d *grpcServerServiceDefinition) Description() *grpc.ServiceDesc {
	return d.description
}
