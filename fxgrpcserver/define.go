package fxgrpcserver

import "google.golang.org/grpc"

type GrpcServiceDefinition interface {
	ReturnType() string
	Description() *grpc.ServiceDesc
}

type grpcServiceDefinition struct {
	returnType  string
	description *grpc.ServiceDesc
}

func NewGrpcServiceDefinition(returnType string, description *grpc.ServiceDesc) GrpcServiceDefinition {
	return &grpcServiceDefinition{
		returnType:  returnType,
		description: description,
	}
}

func (d *grpcServiceDefinition) ReturnType() string {
	return d.returnType
}

func (d *grpcServiceDefinition) Description() *grpc.ServiceDesc {
	return d.description
}
