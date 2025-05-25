package multi_agent

import (
	"context"
	"fmt"

	"github.com/junjie-w/llm-integration-patterns-experiments/pkg/config"
)

type Service struct {
	coordinator *Coordinator
}

func NewService(cfg *config.Config) *Service {
	coordinator := NewCoordinator(cfg)
	
	coordinator.RegisterAgent(NewCustomerSupportAgent(cfg))
	coordinator.RegisterAgent(NewTechnicalSupportAgent(cfg))
	coordinator.RegisterAgent(NewOrderSpecialistAgent(cfg))
	
	return &Service{
		coordinator: coordinator,
	}
}

type Request struct {
	Message        string `json:"message"`
	ConversationID string `json:"conversation_id,omitempty"`
}

type Response struct {
	ConversationID string   `json:"conversation_id"`
	Reply          string   `json:"reply"`
	Agents         []string `json:"agents"`
	Complete       bool     `json:"complete"`
}

func (s *Service) Process(ctx context.Context, req Request) (*Response, error) {
	if req.ConversationID != "" {
		return nil, fmt.Errorf("continuing conversations not yet implemented")
	}
	
	conversation, err := s.coordinator.StartConversation(ctx, req.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to process request: %w", err)
	}
	
	var finalReply string
	if len(conversation.Messages) > 0 {
		finalReply = conversation.Messages[len(conversation.Messages)-1].Content
	}
	
	agentNames := make([]string, 0)
	agentMap := make(map[string]bool)
	
	for _, msg := range conversation.Messages {
		if msg.From != "User" && msg.From != "Coordinator" && !agentMap[msg.From] {
			agentNames = append(agentNames, msg.From)
			agentMap[msg.From] = true
		}
	}
	
	return &Response{
		ConversationID: conversation.ID,
		Reply:          finalReply,
		Agents:         agentNames,
		Complete:       conversation.IsComplete,
	}, nil
}
