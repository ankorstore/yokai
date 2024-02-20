package service

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/ankorstore/yokai/grpcserver"
	"github.com/ankorstore/yokai/grpcserver/testdata/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestService struct {
	proto.UnimplementedServiceServer
}

func NewTestServiceServer() *TestService {
	return &TestService{}
}

func (t *TestService) Unary(ctx context.Context, in *proto.Request) (*proto.Response, error) {
	ctx, span := grpcserver.CtxTracer(ctx).Start(ctx, "unary trace")
	defer span.End()

	logger := grpcserver.CtxLogger(ctx)
	logger.Info().Msg("unary call")

	if in.ShouldFail {
		logger.Error().Msg("unary call failure")

		return nil, status.Error(codes.Internal, "failure")
	}

	if in.ShouldPanic {
		logger.Error().Msg("unary call panic")

		panic(in.Message)
	}

	logger.Info().Msg("unary call success")

	return &proto.Response{
		Success: true,
		Message: in.Message,
	}, nil
}

func (t *TestService) Bidi(stream proto.Service_BidiServer) error {
	ctx, span := grpcserver.CtxTracer(stream.Context()).Start(stream.Context(), "bidi trace")
	defer span.End()

	logger := grpcserver.CtxLogger(ctx)
	logger.Info().Msg("bidi call")

	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			logger.Error().Err(err).Msg("bidi recv failed")

			return err
		}

		logger.Info().Msgf("bidi recv value %s", req.Message)

		if req.ShouldFail {
			logger.Error().Msg("bidi call failed")

			return status.Error(codes.Internal, "failure")
		}

		if req.ShouldPanic {
			logger.Error().Msg("unary call panic")

			panic(req.Message)
		}

		split := strings.Split(req.Message, " ")

		for _, word := range split {
			logger.Info().Msgf("bidi send value %s", word)

			err = stream.Send(&proto.Response{
				Success: true,
				Message: word,
			})

			time.Sleep(1 * time.Millisecond)

			if err != nil {
				logger.Error().Err(err).Msg("bidi send failed")

				return err
			}
		}
	}
}
