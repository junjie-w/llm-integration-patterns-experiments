package knowledge_rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/embeddings"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/store/document"
	"github.com/junjie-w/llm-integration-patterns-experiments/pkg/config"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	client          *openai.Client
	docRepo         *document.Repository
	embeddingService *embeddings.Service
}

func NewService(cfg *config.Config, docRepo *document.Repository, embeddingService *embeddings.Service) *Service {
	client := openai.NewClient(cfg.OpenAIKey)
	return &Service{
		client:          client,
		docRepo:         docRepo,
		embeddingService: embeddingService,
	}
}

type Request struct {
	Message string `json:"message"`
	UseVectorSearch bool `json:"use_vector_search"`
}

type Response struct {
	Reply string `json:"reply"`
	Sources []document.Document `json:"sources,omitempty"`
}

func (s *Service) GetCompletion(ctx context.Context, req Request) (*Response, error) {
	var relevantDocs []document.Document
	
	if req.UseVectorSearch {
		results, err := s.embeddingService.FindSimilarDocuments(
			ctx, 
			req.Message, 
			s.docRepo.List(), 
			3,
		)
		if err != nil {
			return nil, fmt.Errorf("vector search failed: %w", err)
		}
		
		for _, result := range results {
			if result.Score > 0.7 {
				relevantDocs = append(relevantDocs, result.Document)
			}
		}
	} else {
		relevantDocs = s.docRepo.SearchByKeyword(req.Message)
	}
	
	context := s.formatContext(relevantDocs)
	
	chatReq := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful customer support assistant. Answer based on the provided context when relevant, or say you don't know. Keep responses concise and accurate.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Context:\n%s\n\nQuestion: %s", context, req.Message),
			},
		},
		Temperature: 0.7,
		MaxTokens:   300,
	}

	resp, err := s.client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no completion choices returned")
	}

	return &Response{
		Reply:   resp.Choices[0].Message.Content,
		Sources: relevantDocs,
	}, nil
}

func (s *Service) formatContext(docs []document.Document) string {
	if len(docs) == 0 {
		return "No relevant information found."
	}
	
	var builder strings.Builder
	
	for i, doc := range docs {
		builder.WriteString(fmt.Sprintf("[%d] %s\n%s\n\n", i+1, doc.Title, doc.Content))
	}
	
	return builder.String()
}
