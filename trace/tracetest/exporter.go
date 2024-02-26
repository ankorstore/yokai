package tracetest

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// TestTraceExporter is the interface for test trace span exporters.
type TestTraceExporter interface {
	trace.SpanExporter
	Exporter() trace.SpanExporter
	Reset() TestTraceExporter
	Spans() tracetest.SpanStubs
	Span(name string) (tracetest.SpanStub, error)
	HasSpan(expectedName string, expectedAttributes ...attribute.KeyValue) bool
	ContainSpan(expectedName string, expectedAttributes ...attribute.KeyValue) bool
	Dump()
}

// DefaultTestTraceExporter is the default [TestTraceExporter] implementation.
type DefaultTestTraceExporter struct {
	*tracetest.InMemoryExporter
}

// NewDefaultTestTraceExporter returns a [DefaultTestTraceExporter], implementing [TestTraceExporter].
func NewDefaultTestTraceExporter() TestTraceExporter {
	return &DefaultTestTraceExporter{
		tracetest.NewInMemoryExporter(),
	}
}

// Exporter returns the in memory internal exporter.
func (e *DefaultTestTraceExporter) Exporter() trace.SpanExporter {
	return e.InMemoryExporter
}

// Reset resets the in memory internal exporter.
func (e *DefaultTestTraceExporter) Reset() TestTraceExporter {
	e.InMemoryExporter.Reset()

	return e
}

// Spans get the [tracetest.SpanStubs] from the in memory internal exporter.
func (e *DefaultTestTraceExporter) Spans() tracetest.SpanStubs {
	return e.InMemoryExporter.GetSpans()
}

// Span get a specific [tracetest.SpanStub] from the in memory internal exporter by name.
func (e *DefaultTestTraceExporter) Span(name string) (tracetest.SpanStub, error) {
	for _, span := range e.InMemoryExporter.GetSpans() {
		if span.Name == name {
			return span, nil
		}
	}

	return tracetest.SpanStub{}, fmt.Errorf("span with name %s cannot be found", name)
}

// HasSpan return true if a trace span from the in memory internal buffer is exactly matching provided name and attributes.
func (e *DefaultTestTraceExporter) HasSpan(expectedName string, expectedAttributes ...attribute.KeyValue) bool {
	for _, span := range e.InMemoryExporter.GetSpans() {
		if span.Name == expectedName {
			if len(expectedAttributes) == 0 {
				return true
			}

			allMatch := true
			for _, expectedAttribute := range expectedAttributes {
				found := false

				for _, spanAttribute := range span.Attributes {
					if spanAttribute.Key == expectedAttribute.Key {
						found = true

						allMatch = allMatch && spanAttribute.Value == expectedAttribute.Value
					}
				}

				if !found {
					return false
				}
			}

			if allMatch {
				return true
			}
		}
	}

	return false
}

// ContainSpan return true if a trace span from the in memory internal buffer is partially matching provided name and attributes.
//
//nolint:cyclop,gocognit,exhaustive
func (e *DefaultTestTraceExporter) ContainSpan(expectedName string, expectedAttributes ...attribute.KeyValue) bool {
	for _, span := range e.InMemoryExporter.GetSpans() {
		if span.Name == expectedName {
			if len(expectedAttributes) == 0 {
				return true
			}

			allMatch := true
			for _, expectedAttribute := range expectedAttributes {
				found := false

				for _, spanAttribute := range span.Attributes {
					if spanAttribute.Key == expectedAttribute.Key {
						found = true

						switch spanAttribute.Value.Type() {
						case attribute.STRING:
							allMatch = allMatch && strings.Contains(
								spanAttribute.Value.AsString(),
								expectedAttribute.Value.AsString(),
							)
						default:
							allMatch = allMatch && spanAttribute.Value == expectedAttribute.Value
						}
					}
				}

				if !found {
					return false
				}
			}

			if allMatch {
				return true
			}
		}
	}

	return false
}

// Dump prints the [tracetest.SpanStubs] snapshots from the in memory internal exporter, for debugging purposes.
func (e *DefaultTestTraceExporter) Dump() {
	for _, span := range e.Spans().Snapshots() {
		//nolint:forbidigo
		fmt.Printf("%v\n", span)
	}
}
