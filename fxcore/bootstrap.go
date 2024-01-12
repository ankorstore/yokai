package fxcore

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/fxlog"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// Bootstrapper is the application bootstrapper, that can load a list [fx.Option] and run your application.
//
//nolint:containedctx
type Bootstrapper struct {
	context context.Context
	options []fx.Option
}

// NewBootstrapper returns a new [Bootstrapper].
func NewBootstrapper() *Bootstrapper {
	return &Bootstrapper{
		context: context.Background(),
		options: []fx.Option{
			FxCoreModule,
		},
	}
}

// WithContext is used to pass a parent [context.Context].
func (b *Bootstrapper) WithContext(ctx context.Context) *Bootstrapper {
	b.context = ctx

	return b
}

// WithOptions is used to pass a list of [fx.Option].
func (b *Bootstrapper) WithOptions(options ...fx.Option) *Bootstrapper {
	b.options = append(b.options, options...)

	return b
}

// BootstrapApp boostrap the application, accepting optional bootstrap options.
func (b *Bootstrapper) BootstrapApp(options ...fx.Option) *fx.App {
	return fx.New(
		fx.Supply(fx.Annotate(b.context, fx.As(new(context.Context)))),
		fx.WithLogger(fxlog.NewFxEventLogger),
		fx.Options(b.options...),
		fx.Options(options...),
	)
}

// BootstrapTestApp boostrap the application in test mode, accepting a testing context and optional bootstrap options.
func (b *Bootstrapper) BootstrapTestApp(tb testing.TB, options ...fx.Option) *fxtest.App {
	tb.Helper()

	tb.Setenv("APP_ENV", "test")

	return fxtest.New(
		tb,
		fx.Supply(fx.Annotate(b.context, fx.As(new(context.Context)))),
		fx.NopLogger,
		fx.Options(b.options...),
		fx.Options(options...),
	)
}

// RunApp runs the application, accepting optional runtime options.
func (b *Bootstrapper) RunApp(options ...fx.Option) {
	b.BootstrapApp(options...).Run()
}

// RunTestApp runs the application in test mode, accepting a testing context and optional runtime options.
func (b *Bootstrapper) RunTestApp(tb testing.TB, options ...fx.Option) {
	tb.Helper()

	b.BootstrapTestApp(tb, options...).RequireStart().RequireStop()
}
