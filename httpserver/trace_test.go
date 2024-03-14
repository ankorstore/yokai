package httpserver_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestAnnotateTracerProviderWithNonSdkTracerProvider(t *testing.T) {
	t.Parallel()

	tp := noop.NewTracerProvider()

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
