package trace

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DefaultOtlpGrpcTimeout is the default timeout in seconds for the OTLP gRPC connection.
const DefaultOtlpGrpcTimeout = 30

// NewOtlpGrpcClientConnection returns a gRPC connection, and accept a host and a list of [grpc.DialOption].
func NewOtlpGrpcClientConnection(ctx context.Context, host string, dialOptions ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialCtx, cancel := context.WithTimeout(ctx, DefaultOtlpGrpcTimeout*time.Second)
	defer cancel()

	dialContextOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	dialContextOptions = append(dialContextOptions, dialOptions...)

	return grpc.DialContext(dialCtx, host, dialContextOptions...)
}
