package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ankorstore/yokai/fxgrpcserver/testdata/proto"
	"github.com/ankorstore/yokai/grpcserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestService struct {
	proto.UnimplementedServiceServer
	dependency *TestServiceDependency
}

func NewTestServiceServer(dependency *TestServiceDependency) *TestService {
	return &TestService{
		dependency: dependency,
	}
}

func (t *TestService) Unary(ctx context.Context, in *proto.Request) (*proto.Response, error) {
	appName := t.dependency.AppName()

	ctx, span := grpcserver.CtxTracer(ctx).Start(ctx, fmt.Sprintf("unary trace on %s", appName))
	defer span.End()

	logger := grpcserver.CtxLogger(ctx)
	logger.Info().Msgf("unary call on %s", appName)

	if in.ShouldFail {
		logger.Error().Msgf("unary call failure on %s", appName)

		return nil, status.Error(codes.Internal, "failure")
	}

	if in.ShouldPanic {
		logger.Error().Msgf("unary call panic on %s", appName)

		panic(in.Message)
	}

	logger.Info().Msgf("unary call success on %s", appName)

	return &proto.Response{
		Success: true,
		Message: fmt.Sprintf("%s received on %s", in.Message, appName),
	}, nil
}

func (t *TestService) Bidi(stream proto.Service_BidiServer) error {
	appName := t.dependency.AppName()

	ctx, span := grpcserver.CtxTracer(stream.Context()).Start(stream.Context(), fmt.Sprintf("bidi trace on %s", appName))
	defer span.End()

	logger := grpcserver.CtxLogger(ctx)
	logger.Info().Msgf("bidi call on %s", appName)

	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			logger.Error().Err(err).Msgf("bidi recv failed on %s", appName)

			return err
		}

		logger.Info().Msgf("bidi recv value %s on %s", req.Message, appName)

		if req.ShouldFail {
			logger.Error().Msgf("bidi call failed on %s", appName)

			return status.Error(codes.Internal, "failure")
		}

		if req.ShouldPanic {
			logger.Error().Msgf("unary call panic on %s", appName)

			panic(req.Message)
		}

		split := strings.Split(req.Message, " ")

		for _, word := range split {
			logger.Info().Msgf("bidi send value %s on %s", word, appName)

			err = stream.Send(&proto.Response{
				Success: true,
				Message: word,
			})

			time.Sleep(1 * time.Millisecond)

			if err != nil {
				logger.Error().Err(err).Msgf("bidi send failed on %s", appName)

				return err
			}
		}
	}
}
