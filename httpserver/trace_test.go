package httpserver_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type tracerProviderMock struct {
	mock.Mock
}

func (m *tracerProviderMock) Tracer(string, ...trace.TracerOption) trace.Tracer {
	args := m.Called()

	return otel.GetTracerProvider().Tracer(args.String(0))
}

func TestAnnotateTracerProviderWithNonSdkTracerProvider(t *testing.T) {
	t.Parallel()

	tp := new(tracerProviderMock)

	assert.Equal(t, tp, httpserver.AnnotateTracerProvider(tp))
}

func TestNewTracerProviderRequestIdAnnotator(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	annotator := httpserver.NewTracerProviderRequestIdAnnotator()

	assert.IsType(t, &httpserver.TracerProviderRequestIdAnnotator{}, annotator)
	assert.Nil(t, annotator.ForceFlush(ctx))
	assert.Nil(t, annotator.Shutdown(ctx))
}
