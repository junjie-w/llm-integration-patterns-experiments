package function_calling

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/tool"
	"github.com/junjie-w/llm-integration-patterns-experiments/pkg/config"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	client      *openai.Client
	toolRegistry *tool.Registry
}

func NewService(cfg *config.Config, toolRegistry *tool.Registry) *Service {
	client := openai.NewClient(cfg.OpenAIKey)
	return &Service{
		client:      client,
		toolRegistry: toolRegistry,
	}
}

type Request struct {
	Message string `json:"message"`
}

type ToolCallInfo struct {
	Name      string      `json:"name"`
	Arguments interface{} `json:"arguments"`
	Result    interface{} `json:"result"`
}

type Response struct {
	Reply     string        `json:"reply"`
	ToolCalls []ToolCallInfo `json:"tool_calls,omitempty"`
}

func (s *Service) GetCompletion(ctx context.Context, req Request) (*Response, error) {
	tools := make([]openai.Tool, 0)
	for _, t := range s.toolRegistry.List() {
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			},
		})
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful customer support assistant that can use tools to look up information.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: req.Message,
		},
	}

	var toolCalls []ToolCallInfo
	
	for i := 0; i < 3; i++ {
		chatReq := openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
			Tools:    tools,
		}

		resp, err := s.client.CreateChatCompletion(ctx, chatReq)
		if err != nil {
			return nil, fmt.Errorf("failed to get completion: %w", err)
		}

		if len(resp.Choices) == 0 {
			return nil, fmt.Errorf("no completion choices returned")
		}

		assistantMsg := resp.Choices[0].Message
		messages = append(messages, assistantMsg)

		if len(assistantMsg.ToolCalls) == 0 {
			return &Response{
				Reply:     assistantMsg.Content,
				ToolCalls: toolCalls,
			}, nil
		}

		for _, call := range assistantMsg.ToolCalls {
			if call.Function.Name == "" {
				continue
			}

			var args map[string]interface{}
			if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
				continue
			}

			result, err := s.toolRegistry.CallTool(ctx, call.Function.Name, args)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}

			toolCalls = append(toolCalls, ToolCallInfo{
				Name:      call.Function.Name,
				Arguments: args,
				Result:    result,
			})

			resultJSON, _ := json.Marshal(result)
			messages = append(messages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    string(resultJSON),
				ToolCallID: call.ID,
			})
		}
	}

	finalReq := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	}

	finalResp, err := s.client.CreateChatCompletion(ctx, finalReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get final completion: %w", err)
	}

	return &Response{
		Reply:     finalResp.Choices[0].Message.Content,
		ToolCalls: toolCalls,
	}, nil
}
