package worker_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/worker"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestAnnotateTracerProviderWithNonSdkTracerProvider(t *testing.T) {
	t.Parallel()

	tp := noop.NewTracerProvider()

	assert.Equal(t, tp, worker.AnnotateTracerProvider(tp))
}

func TestNewTracerProviderCronJobAnnotator(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	annotator := worker.NewTracerProviderWorkerAnnotator()

	assert.IsType(t, &worker.TracerProviderWorkerAnnotator{}, annotator)
	assert.Nil(t, annotator.ForceFlush(ctx))
	assert.Nil(t, annotator.Shutdown(ctx))
}
