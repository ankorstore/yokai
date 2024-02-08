package fxcron_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/fxcron"
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

	assert.Equal(t, tp, fxcron.AnnotateTracerProvider(tp))
}

func TestNewTracerProviderCronJobAnnotator(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	annotator := fxcron.NewTracerProviderCronJobAnnotator()

	assert.IsType(t, &fxcron.TracerProviderCronJobAnnotator{}, annotator)
	assert.Nil(t, annotator.ForceFlush(ctx))
	assert.Nil(t, annotator.Shutdown(ctx))
}
