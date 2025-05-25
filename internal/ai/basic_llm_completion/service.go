package basic_llm_completion

import (
	"context"
	"fmt"

	"github.com/junjie-w/llm-integration-patterns-experiments/pkg/config"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	client *openai.Client
}

func NewService(cfg *config.Config) *Service {
	client := openai.NewClient(cfg.OpenAIKey)
	return &Service{
		client: client,
	}
}

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Reply string `json:"reply"`
}

func (s *Service) GetCompletion(ctx context.Context, req Request) (*Response, error) {
	chatReq := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful customer support assistant. Provide concise and accurate responses.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: req.Message,
			},
		},
		Temperature: 0.7,
		MaxTokens:   150,
	}

	resp, err := s.client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no completion choices returned")
	}

	return &Response{
		Reply: resp.Choices[0].Message.Content,
	}, nil
}
