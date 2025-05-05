package server

import (
	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPServerTool is the interface for MCP server tools.
type MCPServerTool interface {
	Name() string
	Options() []mcp.ToolOption
	Handle() server.ToolHandlerFunc
}

// MCPServerPrompt is the interface for MCP server prompts.
type MCPServerPrompt interface {
	Name() string
	Options() []mcp.PromptOption
	Handle() server.PromptHandlerFunc
}

// MCPServerResource is the interface for MCP server resources.
type MCPServerResource interface {
	Name() string
	URI() string
	Options() []mcp.ResourceOption
	Handle() server.ResourceHandlerFunc
}

// MCPServerResourceTemplate is the interface for MCP server resource templates.
type MCPServerResourceTemplate interface {
	Name() string
	URI() string
	Options() []mcp.ResourceTemplateOption
	Handle() server.ResourceTemplateHandlerFunc
}

// MCPServerRegistryInfo is the information of the MCPServerRegistry.
type MCPServerRegistryInfo struct {
	Capabilities struct {
		Tools     bool
		Prompts   bool
		Resources bool
	}
	Registrations struct {
		Tools             map[string]string
		Prompts           map[string]string
		Resources         map[string]string
		ResourceTemplates map[string]string
	}
}

// MCPServerRegistry is the registry for MCP tools, prompts, resources and resource templates.
type MCPServerRegistry struct {
	config            *config.Config
	tools             map[string]MCPServerTool
	prompts           map[string]MCPServerPrompt
	resources         map[string]MCPServerResource
	resourceTemplates map[string]MCPServerResourceTemplate
}

// NewMCPServerRegistry returns a new MCPServerRegistry instance.
func NewMCPServerRegistry(
	config *config.Config,
	tools []MCPServerTool,
	prompts []MCPServerPrompt,
	resources []MCPServerResource,
	resourceTemplates []MCPServerResourceTemplate,
) *MCPServerRegistry {
	toolsMap := make(map[string]MCPServerTool, len(tools))
	promptsMap := make(map[string]MCPServerPrompt, len(prompts))
	resourcesMap := make(map[string]MCPServerResource, len(resources))
	resourceTemplatesMap := make(map[string]MCPServerResourceTemplate, len(resourceTemplates))

	for _, tool := range tools {
		toolsMap[tool.Name()] = tool
	}

	for _, prompt := range prompts {
		promptsMap[prompt.Name()] = prompt
	}

	for _, resource := range resources {
		resourcesMap[resource.Name()] = resource
	}

	for _, resourceTemplate := range resourceTemplates {
		resourceTemplatesMap[resourceTemplate.Name()] = resourceTemplate
	}

	return &MCPServerRegistry{
		config:            config,
		tools:             toolsMap,
		prompts:           promptsMap,
		resources:         resourcesMap,
		resourceTemplates: resourceTemplatesMap,
	}
}

// Register registers MCP tools, prompts, resources and resource templates on a provided MCPServer instance.
func (r *MCPServerRegistry) Register(mcpServer *server.MCPServer) {
	if r.config.GetBool("modules.mcp.server.capabilities.tools") {
		for _, tool := range r.tools {
			mcpServer.AddTool(
				mcp.NewTool(tool.Name(), tool.Options()...),
				tool.Handle(),
			)
		}
	}

	if r.config.GetBool("modules.mcp.server.capabilities.prompts") {
		for _, prompt := range r.prompts {
			mcpServer.AddPrompt(
				mcp.NewPrompt(prompt.Name(), prompt.Options()...),
				prompt.Handle(),
			)
		}
	}

	if r.config.GetBool("modules.mcp.server.capabilities.resources") {
		for _, resource := range r.resources {
			mcpServer.AddResource(
				mcp.NewResource(resource.URI(), resource.Name(), resource.Options()...),
				resource.Handle(),
			)
		}

		for _, resourceTemplate := range r.resourceTemplates {
			mcpServer.AddResourceTemplate(
				mcp.NewResourceTemplate(resourceTemplate.URI(), resourceTemplate.Name(), resourceTemplate.Options()...),
				resourceTemplate.Handle(),
			)
		}
	}
}

// Info returns information about the capabilities and the registered MCP tools, prompts, resources and resource templates.
func (r *MCPServerRegistry) Info() MCPServerRegistryInfo {
	toolsInfo := make(map[string]string, len(r.tools))
	for _, tool := range r.tools {
		toolsInfo[tool.Name()] = FuncName(tool.Handle())
	}

	promptsInfo := make(map[string]string, len(r.prompts))
	for _, prompt := range r.prompts {
		promptsInfo[prompt.Name()] = FuncName(prompt.Handle())
	}

	resourcesInfo := make(map[string]string, len(r.resources))
	for _, resource := range r.resources {
		resourcesInfo[resource.Name()] = FuncName(resource.Handle())
	}

	resourceTemplatesInfo := make(map[string]string, len(r.resourceTemplates))
	for _, resourceTemplate := range r.resourceTemplates {
		resourceTemplatesInfo[resourceTemplate.Name()] = FuncName(resourceTemplate.Handle())
	}

	return MCPServerRegistryInfo{
		Capabilities: struct {
			Tools     bool
			Prompts   bool
			Resources bool
		}{
			Tools:     r.config.GetBool("modules.mcp.server.capabilities.tools"),
			Prompts:   r.config.GetBool("modules.mcp.server.capabilities.prompts"),
			Resources: r.config.GetBool("modules.mcp.server.capabilities.resources"),
		},
		Registrations: struct {
			Tools             map[string]string
			Prompts           map[string]string
			Resources         map[string]string
			ResourceTemplates map[string]string
		}{
			Tools:             toolsInfo,
			Prompts:           promptsInfo,
			Resources:         resourcesInfo,
			ResourceTemplates: resourceTemplatesInfo,
		},
	}
}
