package interceptor

import (
	"context"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func TestUnaryInterceptor(traceId string, spanId string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = trace.ContextWithSpanContext(ctx, prepareSpanContext(traceId, spanId))

		return handler(ctx, req)
	}
}

func TestStreamInterceptor(traceId string, spanId string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := trace.ContextWithSpanContext(ss.Context(), prepareSpanContext(traceId, spanId))

		wrappedStream := &middleware.WrappedServerStream{
			ServerStream:   ss,
			WrappedContext: ctx,
		}

		return handler(srv, wrappedStream)
	}
}

func prepareSpanContext(traceId string, spanId string) trace.SpanContext {
	tId, _ := trace.TraceIDFromHex(traceId)
	sId, _ := trace.SpanIDFromHex(spanId)

	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: tId,
		SpanID:  sId,
	})
}
