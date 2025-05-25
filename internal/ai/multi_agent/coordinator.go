package multi_agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/junjie-w/llm-integration-patterns-experiments/pkg/config"
	"github.com/sashabaranov/go-openai"
)

type Coordinator struct {
	agents       []Agent
	client       *openai.Client
	conversations map[string]*Conversation
}

func NewCoordinator(cfg *config.Config) *Coordinator {
	return &Coordinator{
		agents:       []Agent{},
		client:       openai.NewClient(cfg.OpenAIKey),
		conversations: make(map[string]*Conversation),
	}
}

func (c *Coordinator) RegisterAgent(agent Agent) {
	c.agents = append(c.agents, agent)
}

func (c *Coordinator) StartConversation(ctx context.Context, query string) (*Conversation, error) {
	conversation := NewConversation(query)
	
	bestAgent, err := c.selectAgent(ctx, query)
	if err != nil {
		return nil, err
	}
	
	initialMsg := NewMessage(
		MessageTypeQuestion,
		query,
		"User",
		bestAgent.GetName(),
		conversation.ID,
	)
	
	conversation.AddMessage(initialMsg)
	
	response, err := bestAgent.ProcessMessage(ctx, initialMsg)
	if err != nil {
		return nil, err
	}
	
	conversation.AddMessage(response)
	
	if c.shouldDelegate(response.Content) {
		topic := c.extractDelegationTopic(response.Content)
		
		nextAgent, err := c.selectAgentForTopic(ctx, topic)
		if err != nil {
			return nil, err
		}
		
		delegationMsg := NewMessage(
			MessageTypeDelegation,
			fmt.Sprintf("I need help with this user query about %s: %s", topic, query),
			response.From,
			nextAgent.GetName(),
			conversation.ID,
		)
		
		conversation.AddMessage(delegationMsg)
		
		nextResponse, err := nextAgent.ProcessMessage(ctx, delegationMsg)
		if err != nil {
			return nil, err
		}
		
		conversation.AddMessage(nextResponse)
		
		finalResponse, err := c.synthesizeFinalAnswer(ctx, conversation)
		if err != nil {
			return nil, err
		}
		
		finalMsg := NewMessage(
			MessageTypeAnswer,
			finalResponse,
			"Coordinator",
			"User",
			conversation.ID,
		)
		conversation.AddMessage(finalMsg)
	}
	
	conversation.IsComplete = true
	
	c.conversations[conversation.ID] = conversation
	
	return conversation, nil
}

func (c *Coordinator) shouldDelegate(content string) bool {
    lowerContent := strings.ToLower(content)
    
    delegationPhrases := []string{
        "delegate",
        "transfer",
        "specialist",
        "someone else",
        "different department",
        "not my expertise",
        "would need to",
        "technical support",
        "technical team",
        "technical assistance",
    }
    
    technicalKeywords := []string{
        "connect", "connection", "pairing", "bluetooth", 
        "wifi", "wireless", "software", "update", "driver",
        "troubleshoot", "not working", "doesn't work",
    }
    
    for _, phrase := range delegationPhrases {
        if strings.Contains(lowerContent, phrase) {
            return true
        }
    }
    
    for _, keyword := range technicalKeywords {
        if strings.Contains(lowerContent, keyword) {
            return true
        }
    }
    
    return false
}

func (c *Coordinator) extractDelegationTopic(content string) string {	
	lowerContent := strings.ToLower(content)
	
	if strings.Contains(lowerContent, "technical") {
		return "technical support"
	} else if strings.Contains(lowerContent, "order") || strings.Contains(lowerContent, "shipping") {
		return "order management"
	} else {
		return "customer support"
	}
}

func (c *Coordinator) selectAgentForTopic(ctx context.Context, topic string) (Agent, error) {
	topic = strings.ToLower(topic)
	
	for _, agent := range c.agents {
		expertise := strings.ToLower(agent.GetExpertise())
		if strings.Contains(expertise, topic) {
			return agent, nil
		}
	}
	
	if len(c.agents) > 0 {
		return c.agents[0], nil
	}
	
	return nil, fmt.Errorf("no agent found for topic: %s", topic)
}

func (c *Coordinator) selectAgent(ctx context.Context, query string) (Agent, error) {
    lowerQuery := strings.ToLower(query)
    
    technicalPatterns := []string{"connect", "pair", "troubleshoot", "not working", 
                                "error", "issue", "problem", "broken", "fix", "help with"}
                                
    orderPatterns := []string{"order", "delivery", "shipping", "package", 
                             "tracking", "return", "refund"}
    
    for _, pattern := range technicalPatterns {
        if strings.Contains(lowerQuery, pattern) {
            for _, agent := range c.agents {
                if strings.Contains(strings.ToLower(agent.GetExpertise()), "technical") {
                    return agent, nil
                }
            }
        }
    }
    
    for _, pattern := range orderPatterns {
        if strings.Contains(lowerQuery, pattern) {
            for _, agent := range c.agents {
                if strings.Contains(strings.ToLower(agent.GetExpertise()), "order") {
                    return agent, nil
                }
            }
        }
    }
    
    var bestAgent Agent
    var highestScore float64
    
    for _, agent := range c.agents {
        score := agent.CanHandle(query)
        if score > highestScore {
            highestScore = score
            bestAgent = agent
        }
    }
    
    if bestAgent == nil {
        return nil, fmt.Errorf("no suitable agent found for query")
    }
    
    return bestAgent, nil
}

func (c *Coordinator) synthesizeFinalAnswer(ctx context.Context, conversation *Conversation) (string, error) {
	var prompt strings.Builder
	prompt.WriteString("Synthesize a clear, helpful response to the user based on these agent interactions:\n\n")
	
	prompt.WriteString("User query: " + conversation.Query + "\n\n")
	
	for _, msg := range conversation.Messages {
		if msg.Type == MessageTypeAnswer {
			prompt.WriteString(fmt.Sprintf("%s's answer: %s\n\n", msg.From, msg.Content))
		}
	}
	
	prompt.WriteString("Create a unified, helpful response that incorporates the relevant information from all agents.")
	
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a coordinator that synthesizes information from multiple expert agents into coherent, helpful responses.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt.String(),
			},
		},
	})
	
	if err != nil {
		return "", err
	}
	
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}
	
	return resp.Choices[0].Message.Content, nil
}

func (c *Coordinator) GetConversation(id string) (*Conversation, bool) {
	conv, exists := c.conversations[id]
	return conv, exists
}
