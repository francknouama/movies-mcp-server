package server

import (
	"fmt"

	"movies-mcp-server/internal/interfaces/dto"
)

// ToolHandlerFunc defines the signature for tool handler functions
type ToolHandlerFunc func(id any, arguments map[string]any, sender ResponseSender)

// ResourceHandler defines the signature for resource handler functions
type ResourceHandler func(uri string, sender ResponseSender) 

// PromptHandler defines the signature for prompt handler functions
type PromptHandler func(id any, name string, arguments map[string]interface{}, sender ResponseSender)

// Registry manages the registration and discovery of tools, resources, and prompts
type Registry struct {
	tools     map[string]ToolHandlerFunc
	resources map[string]ResourceHandler
	prompts   map[string]PromptHandler
	schemas   []dto.Tool
	resourceList []dto.Resource
	promptList   []dto.Prompt
}

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	return &Registry{
		tools:     make(map[string]ToolHandlerFunc),
		resources: make(map[string]ResourceHandler),
		prompts:   make(map[string]PromptHandler),
		schemas:   make([]dto.Tool, 0),
		resourceList: make([]dto.Resource, 0),
		promptList:   make([]dto.Prompt, 0),
	}
}

// RegisterTool registers a tool handler with its schema
func (r *Registry) RegisterTool(name string, handler ToolHandlerFunc, schema dto.Tool) {
	r.tools[name] = handler
	r.schemas = append(r.schemas, schema)
}

// RegisterResource registers a resource handler
func (r *Registry) RegisterResource(uri string, handler ResourceHandler, resource dto.Resource) {
	r.resources[uri] = handler
	r.resourceList = append(r.resourceList, resource)
}

// RegisterPrompt registers a prompt handler
func (r *Registry) RegisterPrompt(name string, handler PromptHandler, prompt dto.Prompt) {
	r.prompts[name] = handler
	r.promptList = append(r.promptList, prompt)
}

// GetToolHandler returns a tool handler by name
func (r *Registry) GetToolHandler(name string) (ToolHandlerFunc, bool) {
	handler, exists := r.tools[name]
	return handler, exists
}

// GetResourceHandler returns a resource handler by URI
func (r *Registry) GetResourceHandler(uri string) (ResourceHandler, bool) {
	handler, exists := r.resources[uri]
	return handler, exists
}

// GetPromptHandler returns a prompt handler by name
func (r *Registry) GetPromptHandler(name string) (PromptHandler, bool) {
	handler, exists := r.prompts[name]
	return handler, exists
}

// GetToolSchemas returns all registered tool schemas
func (r *Registry) GetToolSchemas() []dto.Tool {
	return r.schemas
}

// GetResources returns all registered resources
func (r *Registry) GetResources() []dto.Resource {
	return r.resourceList
}

// GetPrompts returns all registered prompts
func (r *Registry) GetPrompts() []dto.Prompt {
	return r.promptList
}

// ToolRegistrar provides methods for registering tools
type ToolRegistrar interface {
	RegisterTools(registry *Registry)
}

// ResourceRegistrar provides methods for registering resources
type ResourceRegistrar interface {
	RegisterResources(registry *Registry)
}

// PromptRegistrar provides methods for registering prompts
type PromptRegistrar interface {
	RegisterPrompts(registry *Registry)
}

// AutoRegister automatically registers tools, resources, and prompts from registrars
func (r *Registry) AutoRegister(registrars ...interface{}) error {
	for _, registrar := range registrars {
		if toolReg, ok := registrar.(ToolRegistrar); ok {
			toolReg.RegisterTools(r)
		}
		if resourceReg, ok := registrar.(ResourceRegistrar); ok {
			resourceReg.RegisterResources(r)
		}
		if promptReg, ok := registrar.(PromptRegistrar); ok {
			promptReg.RegisterPrompts(r)
		}
	}
	return nil
}

// ValidateRegistrations checks that all registrations are valid
func (r *Registry) ValidateRegistrations() error {
	// Check for duplicate tool names
	toolNames := make(map[string]bool)
	for _, schema := range r.schemas {
		if toolNames[schema.Name] {
			return fmt.Errorf("duplicate tool name: %s", schema.Name)
		}
		toolNames[schema.Name] = true
	}

	// Check for duplicate resource URIs
	resourceURIs := make(map[string]bool)
	for _, resource := range r.resourceList {
		if resourceURIs[resource.URI] {
			return fmt.Errorf("duplicate resource URI: %s", resource.URI)
		}
		resourceURIs[resource.URI] = true
	}

	// Check for duplicate prompt names
	promptNames := make(map[string]bool)
	for _, prompt := range r.promptList {
		if promptNames[prompt.Name] {
			return fmt.Errorf("duplicate prompt name: %s", prompt.Name)
		}
		promptNames[prompt.Name] = true
	}

	return nil
}