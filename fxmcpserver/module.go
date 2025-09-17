package fxmcpserver

import (
	"context"
	"github.com/ankorstore/yokai/fxmcpserver/server/stream"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxmcpserver/fxmcpservertest"
	fs "github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/ankorstore/yokai/fxmcpserver/server/stdio"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/server"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const ModuleName = "mcpserver"

// FxMCPServerModule is the MCP server module.
var FxMCPServerModule = fx.Module(
	ModuleName,
	fx.Provide(
		// module fixed dependencies
		ProvideMCPServerRegistry,
		ProvideMCPServer,
		ProvideMCPStreamableHTTPServer,
		ProvideMCPStreamableHTTPTestServer,
		ProvideMCPSSEServer,
		ProvideMCPSSETestServer,
		ProvideMCPStdioServer,
		// module overridable dependencies
		fx.Annotate(
			ProvideDefaultMCPServerHooksProvider,
			fx.As(new(fs.MCPServerHooksProvider)),
		),
		fx.Annotate(
			ProvideDefaultMCPServerFactory,
			fx.As(new(fs.MCPServerFactory)),
		),
		fx.Annotate(
			ProvideDefaultMCPStreamableHTTPServerContextHandler,
			fx.As(new(stream.MCPStreamableHTTPServerContextHandler)),
		),
		fx.Annotate(
			ProvideDefaultMCPStreamableHTTPServerFactory,
			fx.As(new(stream.MCPStreamableHTTPServerFactory)),
		),
		fx.Annotate(
			ProvideDefaultMCPSSEServerContextHandler,
			fx.As(new(sse.MCPSSEServerContextHandler)),
		),
		fx.Annotate(
			ProvideDefaultMCPSSEServerFactory,
			fx.As(new(sse.MCPSSEServerFactory)),
		),
		fx.Annotate(
			ProvideDefaultMCPStdioServerContextHandler,
			fx.As(new(stdio.MCPStdioServerContextHandler)),
		),
		fx.Annotate(
			ProvideDefaultMCPStdioServerFactory,
			fx.As(new(stdio.MCPStdioServerFactory)),
		),
		// module info
		fx.Annotate(
			NewMCPServerModuleInfo,
			fx.As(new(any)),
			fx.ResultTags(`group:"core-module-infos"`),
		),
	),
)

// ProvideDefaultMCPServerHooksProviderParams allows injection of the required dependencies in ProvideDefaultMCPServerHooksProvider.
type ProvideDefaultMCPServerHooksProviderParams struct {
	fx.In
	Registry *prometheus.Registry
	Config   *config.Config
}

// ProvideDefaultMCPServerHooksProvider provides the default server.MCPServerHooksProvider instance.
func ProvideDefaultMCPServerHooksProvider(p ProvideDefaultMCPServerHooksProviderParams) *fs.DefaultMCPServerHooksProvider {
	return fs.NewDefaultMCPServerHooksProvider(p.Registry, p.Config)
}

// ProvideDefaultMCPServerFactoryParams allows injection of the required dependencies in ProvideDefaultMCPServerFactory.
type ProvideDefaultMCPServerFactoryParams struct {
	fx.In
	Config *config.Config
}

// ProvideDefaultMCPServerFactory provides the default server.MCPServerFactory instance.
func ProvideDefaultMCPServerFactory(p ProvideDefaultMCPServerFactoryParams) *fs.DefaultMCPServerFactory {
	return fs.NewDefaultMCPServerFactory(p.Config)
}

// ProvideMCPServerRegistryParams allows injection of the required dependencies in ProvideMCPServerRegistry.
type ProvideMCPServerRegistryParams struct {
	fx.In
	Config            *config.Config
	Tools             []fs.MCPServerTool             `group:"mcp-server-tools"`
	Prompts           []fs.MCPServerPrompt           `group:"mcp-server-prompts"`
	Resources         []fs.MCPServerResource         `group:"mcp-server-resources"`
	ResourceTemplates []fs.MCPServerResourceTemplate `group:"mcp-server-resource-templates"`
}

// ProvideMCPServerRegistry provides the server.MCPServerRegistry.
func ProvideMCPServerRegistry(p ProvideMCPServerRegistryParams) *fs.MCPServerRegistry {
	return fs.NewMCPServerRegistry(
		p.Config,
		p.Tools,
		p.Prompts,
		p.Resources,
		p.ResourceTemplates,
	)
}

// ProvideMCPServerParam allows injection of the required dependencies in ProvideMCPServer.
type ProvideMCPServerParam struct {
	fx.In
	Config   *config.Config
	Provider fs.MCPServerHooksProvider
	Factory  fs.MCPServerFactory
	Registry *fs.MCPServerRegistry
}

// ProvideMCPServer provides the server.MCPServer.
func ProvideMCPServer(p ProvideMCPServerParam) *server.MCPServer {
	srv := p.Factory.Create(server.WithHooks(p.Provider.Provide()))

	p.Registry.Register(srv)

	return srv
}

// ProvideDefaultMCPStreamableHTTPContextHandlerParam allows injection of the required dependencies in ProvideDefaultMCPStreamableHTTPServerContextHandler.
type ProvideDefaultMCPStreamableHTTPContextHandlerParam struct {
	fx.In
	Generator                           uuid.UuidGenerator
	TracerProvider                      trace.TracerProvider
	Logger                              *log.Logger
	MCPStreamableHTTPServerContextHooks []stream.MCPStreamableHTTPServerContextHook `group:"mcp-streamable-http-server-context-hooks"`
}

// ProvideDefaultMCPStreamableHTTPServerContextHandler provides the default sse.MCPStreamableHTTPServerContextHandler instance.
func ProvideDefaultMCPStreamableHTTPServerContextHandler(p ProvideDefaultMCPStreamableHTTPContextHandlerParam) *stream.DefaultMCPStreamableHTTPServerContextHandler {
	textMapPropagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	return stream.NewDefaultMCPStreamableHTTPServerContextHandler(
		p.Generator,
		p.TracerProvider,
		textMapPropagator,
		p.Logger,
		p.MCPStreamableHTTPServerContextHooks...,
	)
}

// ProvideDefaultMCPStreamableHTTPServerFactoryParams allows injection of the required dependencies in ProvideDefaultMCPSSEServerFactory.
type ProvideDefaultMCPStreamableHTTPServerFactoryParams struct {
	fx.In
	Config *config.Config
}

// ProvideDefaultMCPStreamableHTTPServerFactory provides the default sse.MCPStreamableHTTPServerFactory instance.
func ProvideDefaultMCPStreamableHTTPServerFactory(p ProvideDefaultMCPStreamableHTTPServerFactoryParams) *stream.DefaultMCPStreamableHTTPServerFactory {
	return stream.NewDefaultMCPStreamableHTTPServerFactory(p.Config)
}

// ProvideMCPStreamableHTTPServerParam allows injection of the required dependencies in ProvideMCPStreamableHTTPServer.
type ProvideMCPStreamableHTTPServerParam struct {
	fx.In
	LifeCycle                             fx.Lifecycle
	Logger                                *log.Logger
	Config                                *config.Config
	MCPServer                             *server.MCPServer
	MCPStreamableHTTPServerFactory        stream.MCPStreamableHTTPServerFactory
	MCPStreamableHTTPServerContextHandler stream.MCPStreamableHTTPServerContextHandler
}

// ProvideMCPStreamableHTTPServer provides the stream.MCPStreamableHTTPServer.
func ProvideMCPStreamableHTTPServer(p ProvideMCPStreamableHTTPServerParam) *stream.MCPStreamableHTTPServer {
	streamableHTTPServer := p.MCPStreamableHTTPServerFactory.Create(
		p.MCPServer,
		server.WithHTTPContextFunc(p.MCPStreamableHTTPServerContextHandler.Handle()),
	)

	streamableHTTPServerCtx := p.Logger.WithContext(context.Background())

	if p.Config.GetBool("modules.mcp.server.transport.stream.expose") {
		p.LifeCycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				if !p.Config.IsTestEnv() {
					//nolint:errcheck
					go streamableHTTPServer.Start(streamableHTTPServerCtx)
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				if !p.Config.IsTestEnv() {
					return streamableHTTPServer.Stop(ctx)
				}

				return nil
			},
		})
	}

	return streamableHTTPServer
}

// ProvideMCPStreamableHTTPTestServerParam allows injection of the required dependencies in ProvideMCPStreamableHTTPTestServer.
type ProvideMCPStreamableHTTPTestServerParam struct {
	fx.In
	Config                                *config.Config
	MCPServer                             *server.MCPServer
	MCPStreamableHTTPServerContextHandler stream.MCPStreamableHTTPServerContextHandler
}

// ProvideMCPStreamableHTTPTestServer provides the fxmcpservertest.MCPStreamableHTTPTestServer.
func ProvideMCPStreamableHTTPTestServer(p ProvideMCPStreamableHTTPTestServerParam) *fxmcpservertest.MCPStreamableHTTPTestServer {
	return fxmcpservertest.NewMCPStreamableHTTPTestServer(p.Config, p.MCPServer, p.MCPStreamableHTTPServerContextHandler)
}

// ProvideDefaultMCPSSEContextHandlerParam allows injection of the required dependencies in ProvideDefaultMCPSSEServerContextHandler.
type ProvideDefaultMCPSSEContextHandlerParam struct {
	fx.In
	Generator                uuid.UuidGenerator
	TracerProvider           trace.TracerProvider
	Logger                   *log.Logger
	MCPSSEServerContextHooks []sse.MCPSSEServerContextHook `group:"mcp-sse-server-context-hooks"`
}

// ProvideDefaultMCPSSEServerContextHandler provides the default sse.MCPSSEServerContextHandler instance.
func ProvideDefaultMCPSSEServerContextHandler(p ProvideDefaultMCPSSEContextHandlerParam) *sse.DefaultMCPSSEServerContextHandler {
	textMapPropagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	return sse.NewDefaultMCPSSEServerContextHandler(
		p.Generator,
		p.TracerProvider,
		textMapPropagator,
		p.Logger,
		p.MCPSSEServerContextHooks...,
	)
}

// ProvideDefaultMCPSSEServerFactoryParams allows injection of the required dependencies in ProvideDefaultMCPSSEServerFactory.
type ProvideDefaultMCPSSEServerFactoryParams struct {
	fx.In
	Config *config.Config
}

// ProvideDefaultMCPSSEServerFactory provides the default sse.MCPSSEServerFactory instance.
func ProvideDefaultMCPSSEServerFactory(p ProvideDefaultMCPServerFactoryParams) *sse.DefaultMCPSSEServerFactory {
	return sse.NewDefaultMCPSSEServerFactory(p.Config)
}

// ProvideMCPSSEServerParam allows injection of the required dependencies in ProvideMCPSSEServer.
type ProvideMCPSSEServerParam struct {
	fx.In
	LifeCycle                  fx.Lifecycle
	Logger                     *log.Logger
	Config                     *config.Config
	MCPServer                  *server.MCPServer
	MCPSSEServerFactory        sse.MCPSSEServerFactory
	MCPSSEServerContextHandler sse.MCPSSEServerContextHandler
}

// ProvideMCPSSEServer provides the sse.MCPSSEServer.
func ProvideMCPSSEServer(p ProvideMCPSSEServerParam) *sse.MCPSSEServer {
	sseServer := p.MCPSSEServerFactory.Create(
		p.MCPServer,
		server.WithSSEContextFunc(p.MCPSSEServerContextHandler.Handle()),
	)

	sseServerCtx := p.Logger.WithContext(context.Background())

	if p.Config.GetBool("modules.mcp.server.transport.sse.expose") {
		p.LifeCycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				if !p.Config.IsTestEnv() {
					//nolint:errcheck
					go sseServer.Start(sseServerCtx)
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				if !p.Config.IsTestEnv() {
					return sseServer.Stop(ctx)
				}

				return nil
			},
		})
	}

	return sseServer
}

// ProvideMCPSSETestServerParam allows injection of the required dependencies in ProvideMCPSSETestServer.
type ProvideMCPSSETestServerParam struct {
	fx.In
	Config                     *config.Config
	MCPServer                  *server.MCPServer
	MCPSSEServerContextHandler sse.MCPSSEServerContextHandler
}

// ProvideMCPSSETestServer provides the fxmcpservertest.MCPSSETestServer.
func ProvideMCPSSETestServer(p ProvideMCPSSEServerParam) *fxmcpservertest.MCPSSETestServer {
	return fxmcpservertest.NewMCPSSETestServer(p.Config, p.MCPServer, p.MCPSSEServerContextHandler)
}

// ProvideDefaultMCPStdioContextHandlerParam allows injection of the required dependencies in ProvideDefaultMCPStdioServerContextHandler.
type ProvideDefaultMCPStdioContextHandlerParam struct {
	fx.In
	Generator      uuid.UuidGenerator
	TracerProvider trace.TracerProvider
	Logger         *log.Logger
}

// ProvideDefaultMCPStdioServerContextHandler provides the default stdio.MCPStdioServerContextHandler instance.
func ProvideDefaultMCPStdioServerContextHandler(p ProvideDefaultMCPStdioContextHandlerParam) *stdio.DefaultMCPStdioServerContextHandler {
	return stdio.NewDefaultMCPStdioServerContextHandler(p.Generator, p.TracerProvider, p.Logger)
}

// ProvideDefaultMCPStdioServerFactory provides the default stdio.MCPStdioServerFactory instance.
func ProvideDefaultMCPStdioServerFactory() *stdio.DefaultMCPStdioServerFactory {
	return stdio.NewDefaultMCPStdioServerFactory()
}

// ProvideMCPStdioServerParam allows injection of the required dependencies in ProvideMCPStdioServer.
type ProvideMCPStdioServerParam struct {
	fx.In
	LifeCycle                    fx.Lifecycle
	Logger                       *log.Logger
	Config                       *config.Config
	MCPServer                    *server.MCPServer
	MCPStdioServerFactory        stdio.MCPStdioServerFactory
	MCPStdioServerContextHandler stdio.MCPStdioServerContextHandler
}

// ProvideMCPStdioServer provides the stdio.MCPStdioServer.
func ProvideMCPStdioServer(p ProvideMCPStdioServerParam) *stdio.MCPStdioServer {
	stdioServer := p.MCPStdioServerFactory.Create(
		p.MCPServer,
		server.WithStdioContextFunc(p.MCPStdioServerContextHandler.Handle()),
	)

	stdioServerCtx := p.Logger.WithContext(context.Background())

	if p.Config.GetBool("modules.mcp.server.transport.stdio.expose") {
		p.LifeCycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				if !p.Config.IsTestEnv() {
					//nolint:errcheck
					go stdioServer.Start(stdioServerCtx)
				}

				return nil
			},
		})
	}

	return stdioServer
}
