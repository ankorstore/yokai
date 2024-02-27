package fxgrpcserver

import "google.golang.org/grpc"

// GrpcServerServiceDefinition is the interface for grpc server services definitions.
type GrpcServerServiceDefinition interface {
	ReturnType() string
	Description() *grpc.ServiceDesc
}

type grpcServerServiceDefinition struct {
	returnType  string
	description *grpc.ServiceDesc
}

// NewGrpcServiceDefinition returns a new [GrpcServerServiceDefinition] instance.
func NewGrpcServiceDefinition(returnType string, description *grpc.ServiceDesc) GrpcServerServiceDefinition {
	return &grpcServerServiceDefinition{
		returnType:  returnType,
		description: description,
	}
}

// ReturnType returns the definition return type.
func (d *grpcServerServiceDefinition) ReturnType() string {
	return d.returnType
}

// Description returns the definition service description.
func (d *grpcServerServiceDefinition) Description() *grpc.ServiceDesc {
	return d.description
}
