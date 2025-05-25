package tool

import (
	"context"
	"fmt"
)

type Tool struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
	Handler     ToolHandler
}

type ToolHandler func(ctx context.Context, args map[string]interface{}) (interface{}, error)

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Register(tool Tool) error {
	if _, exists := r.tools[tool.Name]; exists {
		return fmt.Errorf("tool with name %s already registered", tool.Name)
	}
	
	r.tools[tool.Name] = tool
	return nil
}

func (r *Registry) Get(name string) (Tool, bool) {
	tool, exists := r.tools[name]
	return tool, exists
}

func (r *Registry) List() []Tool {
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

func (r *Registry) CallTool(ctx context.Context, name string, args map[string]interface{}) (interface{}, error) {
	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", name)
	}
	
	return tool.Handler(ctx, args)
}
