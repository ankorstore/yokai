package grpcserver

import (
	"context"
	"time"

	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	HeaderXRequestId  = "x-request-id"
	LogFieldRequestId = "requestID"
)

// GrpcLoggerInterceptor is a gRPC unary and stream server interceptor to produce correlated logs.
type GrpcLoggerInterceptor struct {
	generator  uuid.UuidGenerator
	logger     *log.Logger
	metadata   map[string]string
	exclusions []string
}

// NewGrpcLoggerInterceptor returns a new [GrpcLoggerInterceptor] instance.
func NewGrpcLoggerInterceptor(generator uuid.UuidGenerator, logger *log.Logger) *GrpcLoggerInterceptor {
	return &GrpcLoggerInterceptor{
		generator:  generator,
		logger:     logger,
		metadata:   map[string]string{HeaderXRequestId: LogFieldRequestId},
		exclusions: []string{},
	}
}

// Metadata configures a list of metadata to log from incoming context.
func (i *GrpcLoggerInterceptor) Metadata(metadata map[string]string) *GrpcLoggerInterceptor {
	for k, v := range metadata {
		i.metadata[k] = v
	}

	return i
}

// Exclude configures a list of method names to exclude from logging.
func (i *GrpcLoggerInterceptor) Exclude(methods ...string) *GrpcLoggerInterceptor {
	i.exclusions = append(i.exclusions, methods...)

	return i
}

// UnaryInterceptor handles the unary requests.
//
//nolint:cyclop,dupl,gocognit,nestif
func (i *GrpcLoggerInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		exclude := Contains(i.exclusions, info.FullMethod)

		grpcLogger := i.logger.With().Fields(i.extractLogFieldsFromContextMetadata(ctx)).Logger()

		newCtx := grpcLogger.WithContext(ctx)

		spanContext := trace.SpanContextFromContext(newCtx)

		traceId := ""
		if spanContext.HasTraceID() {
			traceId = spanContext.TraceID().String()
		}

		spanId := ""
		if spanContext.HasSpanID() {
			spanId = spanContext.SpanID().String()
		}

		if !exclude {
			evt := grpcLogger.
				Debug().
				Str("grpcType", "unary").
				Str("grpcMethod", info.FullMethod)

			if traceId != "" {
				evt.Str("traceID", traceId)
			}

			if spanId != "" {
				evt.Str("spanID", spanId)
			}

			evt.Msg("grpc call start")
		}

		now := time.Now()

		resp, err := handler(newCtx, req)

		errStatus := status.Convert(err)

		if !exclude {
			if err != nil {
				evt := grpcLogger.
					Error().
					Err(err).
					Str("grpcType", "unary").
					Str("grpcMethod", info.FullMethod).
					Uint32("grpcCode", uint32(errStatus.Code())).
					Str("grpcStatus", errStatus.Code().String()).
					Str("grpcDuration", time.Since(now).String())

				if traceId != "" {
					evt.Str("traceID", traceId)
				}

				if spanId != "" {
					evt.Str("spanID", spanId)
				}

				evt.Msg("grpc call error")
			} else {
				evt := grpcLogger.
					Info().
					Str("grpcType", "unary").
					Str("grpcMethod", info.FullMethod).
					Uint32("grpcCode", uint32(codes.OK)).
					Str("grpcStatus", codes.OK.String()).
					Str("grpcDuration", time.Since(now).String())

				if traceId != "" {
					evt.Str("traceID", traceId)
				}

				if spanId != "" {
					evt.Str("spanID", spanId)
				}

				evt.Msg("grpc call success")
			}
		} else if err != nil {
			evt := grpcLogger.
				Error().
				Err(err).
				Str("grpcType", "unary").
				Str("grpcMethod", info.FullMethod).
				Uint32("grpcCode", uint32(errStatus.Code())).
				Str("grpcStatus", errStatus.Code().String()).
				Str("grpcDuration", time.Since(now).String())

			if traceId != "" {
				evt.Str("traceID", traceId)
			}

			if spanId != "" {
				evt.Str("spanID", spanId)
			}

			evt.Msg("grpc call error")
		}

		return resp, err
	}
}

// StreamInterceptor handles the stream requests.
//
//nolint:cyclop,dupl,gocognit,nestif
func (i *GrpcLoggerInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		exclude := Contains(i.exclusions, info.FullMethod)

		grpcLogger := i.logger.
			With().
			Fields(i.extractLogFieldsFromContextMetadata(ctx)).
			Logger()

		newCtx := grpcLogger.WithContext(ctx)

		spanContext := trace.SpanContextFromContext(newCtx)

		traceId := ""
		if spanContext.HasTraceID() {
			traceId = spanContext.TraceID().String()
		}

		spanId := ""
		if spanContext.HasSpanID() {
			spanId = spanContext.SpanID().String()
		}

		if !exclude {
			evt := grpcLogger.
				Info().
				Str("grpcType", "server-streaming").
				Str("grpcMethod", info.FullMethod)

			if traceId != "" {
				evt.Str("traceID", traceId)
			}

			if spanId != "" {
				evt.Str("spanID", spanId)
			}

			evt.Msg("grpc call start")
		}

		wrappedStream := &middleware.WrappedServerStream{
			ServerStream:   ss,
			WrappedContext: newCtx,
		}

		now := time.Now()

		err := handler(srv, wrappedStream)

		errStatus := status.Convert(err)

		if !exclude {
			if err != nil {
				evt := grpcLogger.
					Error().
					Err(err).
					Str("grpcType", "server-streaming").
					Str("grpcMethod", info.FullMethod).
					Uint32("grpcCode", uint32(errStatus.Code())).
					Str("grpcStatus", errStatus.Code().String()).
					Str("grpcDuration", time.Since(now).String())

				if traceId != "" {
					evt.Str("traceID", traceId)
				}

				if spanId != "" {
					evt.Str("spanID", spanId)
				}

				evt.Msg("grpc call error")
			} else {
				evt := grpcLogger.
					Info().
					Str("grpcType", "server-streaming").
					Str("grpcMethod", info.FullMethod).
					Uint32("grpcCode", uint32(codes.OK)).
					Str("grpcStatus", codes.OK.String()).
					Str("grpcDuration", time.Since(now).String())

				if traceId != "" {
					evt.Str("traceID", traceId)
				}

				if spanId != "" {
					evt.Str("spanID", spanId)
				}

				evt.Msg("grpc call success")
			}
		} else if err != nil {
			evt := grpcLogger.
				Error().
				Err(err).
				Str("grpcType", "server-streaming").
				Str("grpcMethod", info.FullMethod).
				Uint32("grpcCode", uint32(errStatus.Code())).
				Str("grpcStatus", errStatus.Code().String()).
				Str("grpcDuration", time.Since(now).String())

			if traceId != "" {
				evt.Str("traceID", traceId)
			}

			if spanId != "" {
				evt.Str("spanID", spanId)
			}

			evt.Msg("grpc call error")
		}

		return err
	}
}

func (i *GrpcLoggerInterceptor) extractLogFieldsFromContextMetadata(ctx context.Context) map[string]interface{} {
	ctxMd, _ := metadata.FromIncomingContext(ctx)

	md := make(map[string]interface{})
	for mk, mv := range i.metadata {
		if val, ok := ctxMd[mk]; ok && len(val) > 0 {
			md[mv] = val[0]
		} else if mk == HeaderXRequestId {
			md[mv] = i.generator.Generate()
		}
	}

	return md
}
