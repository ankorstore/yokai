package trace

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DefaultOtlpGrpcTimeout is the default timeout in seconds for the OTLP gRPC connection.
const DefaultOtlpGrpcTimeout = 5

// NewOtlpGrpcClientConnection returns a gRPC connection, and accept a host and a list of [DialOption]
//
// [DialOption]: https://github.com/grpc/grpc-go
func NewOtlpGrpcClientConnection(ctx context.Context, host string, dialOptions ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialCtx, cancel := context.WithTimeout(ctx, DefaultOtlpGrpcTimeout*time.Second)
	defer cancel()

	if len(dialOptions) == 0 {
		dialOptions = []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		}
	}

	return grpc.DialContext(dialCtx, host, dialOptions...)
}
