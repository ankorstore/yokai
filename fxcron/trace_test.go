package fxcron_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/fxcron"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestAnnotateTracerProviderWithNonSdkTracerProvider(t *testing.T) {
	t.Parallel()

	tp := noop.NewTracerProvider()

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
