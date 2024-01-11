package plugin

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// TracerName is the tracer name for ORM operations.
const TracerName = "orm"

// OrmRegisterer is the interface for compatibility with [Gorm plugins].
//
// [Gorm plugins]: https://gorm.io/docs/write_plugins.html
type OrmRegisterer interface {
	Register(name string, fn func(*gorm.DB)) error
}

// OrmTracerPlugin is the Gorm tracing plugin.
type OrmTracerPlugin struct {
	tracerProvider trace.TracerProvider
	withValues     bool
}

// NewOrmTracerPlugin returns a new [OrmTracerPlugin].
func NewOrmTracerPlugin(tracerProvider trace.TracerProvider, withValues bool) gorm.Plugin {
	return &OrmTracerPlugin{
		tracerProvider: tracerProvider,
		withValues:     withValues,
	}
}

// Name returns the plugin name.
func (p *OrmTracerPlugin) Name() string {
	return TracerName
}

// Initialize is called upon plugin initialization.
func (p *OrmTracerPlugin) Initialize(db *gorm.DB) error {
	cb := db.Callback()

	hooks := []struct {
		callback OrmRegisterer
		hook     func(*gorm.DB)
		name     string
	}{
		{cb.Create().Before("orm:create"), p.before("orm.Create"), "before:create"},
		{cb.Create().After("orm:create"), p.after(), "after:create"},

		{cb.Query().Before("orm:query"), p.before("orm.Query"), "before:select"},
		{cb.Query().After("orm:query"), p.after(), "after:select"},

		{cb.Delete().Before("orm:delete"), p.before("orm.Delete"), "before:delete"},
		{cb.Delete().After("orm:delete"), p.after(), "after:delete"},

		{cb.Update().Before("orm:update"), p.before("orm.Update"), "before:update"},
		{cb.Update().After("orm:update"), p.after(), "after:update"},

		{cb.Row().Before("orm:row"), p.before("orm.Row"), "before:row"},
		{cb.Row().After("orm:row"), p.after(), "after:row"},

		{cb.Raw().Before("orm:raw"), p.before("orm.Raw"), "before:raw"},
		{cb.Raw().After("orm:raw"), p.after(), "after:raw"},
	}

	var firstErr error

	for _, h := range hooks {
		if err := h.callback.Register("otel:"+h.name, h.hook); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("callback register %s failed: %w", h.name, err)
		}
	}

	return firstErr
}

func (p *OrmTracerPlugin) before(spanName string) func(*gorm.DB) {
	return func(tx *gorm.DB) {
		tx.Statement.Context, _ = p.tracerProvider.Tracer(TracerName).Start(
			tx.Statement.Context,
			spanName,
			trace.WithSpanKind(trace.SpanKindClient),
		)
	}
}

func (p *OrmTracerPlugin) after() func(*gorm.DB) {
	return func(tx *gorm.DB) {
		span := trace.SpanFromContext(tx.Statement.Context)
		if !span.IsRecording() {
			return
		}
		defer span.End()

		var attrs []attribute.KeyValue

		if sys := semconv.DBSystemKey.String(tx.Dialector.Name()); sys.Valid() {
			attrs = append(attrs, sys)
		}

		vars := tx.Statement.Vars
		if !p.withValues {
			vars = make([]interface{}, len(tx.Statement.Vars))

			for i := 0; i < len(vars); i++ {
				vars[i] = "?"
			}
		}

		query := tx.Dialector.Explain(tx.Statement.SQL.String(), vars...)

		attrs = append(attrs, semconv.DBStatementKey.String(query))
		if tx.Statement.Table != "" {
			attrs = append(attrs, semconv.DBSQLTableKey.String(tx.Statement.Table))
		}

		span.SetAttributes(attrs...)

		//nolint:errorlint
		switch tx.Error {
		case nil,
			gorm.ErrRecordNotFound,
			driver.ErrSkip,
			io.EOF,
			sql.ErrNoRows:
		default:
			span.RecordError(tx.Error)
			span.SetStatus(codes.Error, tx.Error.Error())
		}
	}
}
