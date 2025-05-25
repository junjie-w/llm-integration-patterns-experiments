package multi_agent

import (
	"context"
	"math"
	"strings"

	"github.com/junjie-w/llm-integration-patterns-experiments/pkg/config"
	"github.com/sashabaranov/go-openai"
)

type Agent interface {
	GetName() string
	
	GetExpertise() string
	
	ProcessMessage(ctx context.Context, message Message) (Message, error)
	
	CanHandle(query string) float64
}

type BaseAgent struct {
	Name      string
	Expertise string
	Client    *openai.Client
	SystemPrompt string
}

func (a *BaseAgent) GetName() string {
	return a.Name
}

func (a *BaseAgent) GetExpertise() string {
	return a.Expertise
}

func (a *BaseAgent) createResponse(ctx context.Context, content string, history []Message) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: a.SystemPrompt,
		},
	}
	
	for _, msg := range history {
		role := openai.ChatMessageRoleUser
		if msg.From == a.Name {
			role = openai.ChatMessageRoleAssistant
		}
		
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}
	
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	})
	
	resp, err := a.Client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo16K,
		Messages:    messages,
		Temperature: 0.7,
	})
	
	if err != nil {
		return "", err
	}
	
	if len(resp.Choices) == 0 {
		return "", nil
	}
	
	return resp.Choices[0].Message.Content, nil
}

type CustomerSupportAgent struct {
	BaseAgent
}

func NewCustomerSupportAgent(cfg *config.Config) *CustomerSupportAgent {
	return &CustomerSupportAgent{
		BaseAgent: BaseAgent{
			Name:      "CustomerSupport",
			Expertise: "general customer support, policies, account issues",
			Client:    openai.NewClient(cfg.OpenAIKey),
			SystemPrompt: `You are a customer support specialist who excels at handling general inquiries,
account issues, and policy questions. If a question is outside your expertise (technical problems or 
order-specific details), say you'll delegate it to a specialist. Be helpful, concise, and friendly.`,
		},
	}
}

func (a *CustomerSupportAgent) CanHandle(query string) float64 {
	query = strings.ToLower(query)
	keywords := []string{"account", "login", "password", "policy", "contact", "feedback", 
		"complaint", "question", "help", "support", "subscription", "membership"}
	
	score := 0.5
	matches := 0
	for _, keyword := range keywords {
		if strings.Contains(query, keyword) {
			score += 0.1
			matches++
		}
	}
	return math.Min(score, 0.9)
}

func (a *CustomerSupportAgent) ProcessMessage(ctx context.Context, message Message) (Message, error) {
	history := []Message{}
	
	response, err := a.createResponse(ctx, message.Content, history)
	if err != nil {
		return Message{}, err
	}
	
	return NewMessage(MessageTypeAnswer, response, a.Name, message.From, message.ThreadID), nil
}

type TechnicalSupportAgent struct {
	BaseAgent
}

func NewTechnicalSupportAgent(cfg *config.Config) *TechnicalSupportAgent {
	return &TechnicalSupportAgent{
		BaseAgent: BaseAgent{
			Name:      "TechnicalSupport",
			Expertise: "technical issues, product functionality, troubleshooting",
			Client:    openai.NewClient(cfg.OpenAIKey),
			SystemPrompt: `You are a technical support specialist who excels at troubleshooting
product issues and providing technical guidance. Focus on clear step-by-step instructions
and technical details. If a question is about order status or general policies, indicate
you'll delegate it to customer support.`,
		},
	}
}

func (a *TechnicalSupportAgent) CanHandle(query string) float64 {
	query = strings.ToLower(query)
	keywords := []string{"troubleshoot", "connect", "pair", "bluetooth", "reset", "not working", 
		"broken", "error", "issue", "problem", "battery", "charging", "firmware", "update", 
		"sound", "audio", "wireless", "manual", "instructions"}
	
	score := 0.3
	matches := 0
	for _, keyword := range keywords {
		if strings.Contains(query, keyword) {
			score += 0.15
			matches++
		}
	}
	if matches >= 3 {
		score += 0.1
	}
	return math.Min(score, 0.95)
}

func (a *TechnicalSupportAgent) ProcessMessage(ctx context.Context, message Message) (Message, error) {
	history := []Message{}
	
	response, err := a.createResponse(ctx, message.Content, history)
	if err != nil {
		return Message{}, err
	}
	
	return NewMessage(MessageTypeAnswer, response, a.Name, message.From, message.ThreadID), nil
}

type OrderSpecialistAgent struct {
	BaseAgent
}

func NewOrderSpecialistAgent(cfg *config.Config) *OrderSpecialistAgent {
	return &OrderSpecialistAgent{
		BaseAgent: BaseAgent{
			Name:      "OrderSpecialist",
			Expertise: "order status, shipping, returns, product availability",
			Client:    openai.NewClient(cfg.OpenAIKey),
			SystemPrompt: `You are an order and shipping specialist who excels at handling order status inquiries,
shipping questions, returns, and product availability. You should focus on order-specific details
and logistics. If a question is about technical issues or general account questions, indicate
you'll delegate it to the appropriate team.`,
		},
	}
}

func (a *OrderSpecialistAgent) CanHandle(query string) float64 {
	query = strings.ToLower(query)
	keywords := []string{"order", "shipping", "delivery", "package", "tracking", 
		"return", "refund", "replacement", "warranty", "receipt", "payment", "status", 
		"expedited", "delay", "address", "inventory", "availability"}
	
	score := 0.4
	matches := 0
	for _, keyword := range keywords {
		if strings.Contains(query, keyword) {
			score += 0.15
			matches++
		}
	}
	if matches >= 2 {
		score += 0.2
	}
	return math.Min(score, 0.99)
}

func (a *OrderSpecialistAgent) ProcessMessage(ctx context.Context, message Message) (Message, error) {
	history := []Message{}

	response, err := a.createResponse(ctx, message.Content, history)
	if err != nil {
		return Message{}, err
	}
	
	return NewMessage(MessageTypeAnswer, response, a.Name, message.From, message.ThreadID), nil
}
