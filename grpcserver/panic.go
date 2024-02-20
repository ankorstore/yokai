package grpcserver

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcPanicRecoveryHandler is used to recover panics with the [recovery] interceptor.
//
// [recovery]: https://github.com/grpc-ecosystem/go-grpc-middleware/tree/main/interceptors/recovery
type GrpcPanicRecoveryHandler struct{}

// NewGrpcPanicRecoveryHandler returns a new [GrpcPanicRecoveryHandler] instance.
func NewGrpcPanicRecoveryHandler() *GrpcPanicRecoveryHandler {
	return &GrpcPanicRecoveryHandler{}
}

// Handle handles the panic recovery.
func (h *GrpcPanicRecoveryHandler) Handle(withDebug bool) recovery.RecoveryHandlerFuncContext {
	return func(ctx context.Context, pnc any) error {
		evt := CtxLogger(ctx).Error().Str("panic", fmt.Sprintf("%s", pnc))

		if withDebug {
			evt.Str("stack", string(debug.Stack()))
		}

		evt.Msg("grpc recovered from panic")

		if withDebug {
			return status.Errorf(codes.Internal, "internal grpc server error, panic = %s, stack = %s", pnc, debug.Stack())
		} else {
			return status.Error(codes.Internal, "internal grpc server error")
		}
	}
}
