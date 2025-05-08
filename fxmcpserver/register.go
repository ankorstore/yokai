package fxmcpserver

import (
	"github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"go.uber.org/fx"
)

// AsMCPServerTool registers an MCP tool.
func AsMCPServerTool(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(server.MCPServerTool)),
			fx.ResultTags(`group:"mcp-server-tools"`),
		),
	)
}

// AsMCPServerTools registers several MCP tools.
func AsMCPServerTools(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsMCPServerTool(constructor))
	}

	return fx.Options(options...)
}

// AsMCPServerPrompt registers an MCP prompt.
func AsMCPServerPrompt(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(server.MCPServerPrompt)),
			fx.ResultTags(`group:"mcp-server-prompts"`),
		),
	)
}

// AsMCPServerPrompts registers several MCP prompts.
func AsMCPServerPrompts(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsMCPServerPrompt(constructor))
	}

	return fx.Options(options...)
}

// AsMCPServerResource registers an MCP resource.
func AsMCPServerResource(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(server.MCPServerResource)),
			fx.ResultTags(`group:"mcp-server-resources"`),
		),
	)
}

// AsMCPServerResources registers several MCP resources.
func AsMCPServerResources(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsMCPServerResource(constructor))
	}

	return fx.Options(options...)
}

// AsMCPServerResourceTemplate registers an MCP resource template.
func AsMCPServerResourceTemplate(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(server.MCPServerResourceTemplate)),
			fx.ResultTags(`group:"mcp-server-resource-templates"`),
		),
	)
}

// AsMCPServerResourceTemplates registers several MCP resource templates.
func AsMCPServerResourceTemplates(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsMCPServerResourceTemplate(constructor))
	}

	return fx.Options(options...)
}

// AsMCPSSEServerContextHook registers an MCP SSE server context hook.
func AsMCPSSEServerContextHook(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(sse.MCPSSEServerContextHook)),
			fx.ResultTags(`group:"mcp-sse-server-context-hooks"`),
		),
	)
}

// AsMCPSSEServerContextHooks registers several MCP SSE server context hook.
func AsMCPSSEServerContextHooks(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsMCPSSEServerContextHook(constructor))
	}

	return fx.Options(options...)
}
