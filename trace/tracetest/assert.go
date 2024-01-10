package tracetest

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

// AssertHasTraceSpan allows to assert if a trace span exactly matching provided name and attributes can be found.
func AssertHasTraceSpan(tb testing.TB, exporter TestTraceExporter, expectedName string, expectedAttributes ...attribute.KeyValue) bool {
	tb.Helper()

	if !exporter.HasSpan(expectedName, expectedAttributes...) {
		tb.Errorf("cannot find trace span with matching name %s and matching attributes %+v", expectedName, expectedAttributes)

		return false
	}

	return true
}

// AssertHasNotTraceSpan allows to assert if a trace span exactly matching provided name and attributes cannot be found.
func AssertHasNotTraceSpan(tb testing.TB, exporter TestTraceExporter, expectedName string, expectedAttributes ...attribute.KeyValue) bool {
	tb.Helper()

	if exporter.HasSpan(expectedName, expectedAttributes...) {
		tb.Errorf("can find trace span with matching name %s and matching attributes %+v", expectedName, expectedAttributes)

		return false
	}

	return true
}

// AssertContainTraceSpan allows to assert if a trace span partially matching provided name and attributes can be found.
func AssertContainTraceSpan(tb testing.TB, exporter TestTraceExporter, expectedName string, expectedAttributes ...attribute.KeyValue) bool {
	tb.Helper()

	if !exporter.ContainSpan(expectedName, expectedAttributes...) {
		tb.Errorf("cannot find trace span with contained name %s and contained attributes %+v", expectedName, expectedAttributes)

		return false
	}

	return true
}

// AssertContainNotTraceSpan allows to assert if a trace span partially matching provided name and attributes cannot be found.
func AssertContainNotTraceSpan(tb testing.TB, exporter TestTraceExporter, expectedName string, expectedAttributes ...attribute.KeyValue) bool {
	tb.Helper()

	if exporter.ContainSpan(expectedName, expectedAttributes...) {
		tb.Errorf("can find trace span with contained name %s and contained attributes %+v", expectedName, expectedAttributes)

		return false
	}

	return true
}
