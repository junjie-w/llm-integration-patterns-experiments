package reasoning_agent

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
	memory      *Memory
}

func NewService(cfg *config.Config, toolRegistry *tool.Registry) *Service {
	return &Service{
		client:      openai.NewClient(cfg.OpenAIKey),
		toolRegistry: toolRegistry,
		memory:      NewMemory(),
	}
}

type Request struct {
	Message  string `json:"message"`
	AgentID  string `json:"agent_id,omitempty"`
}

type Response struct {
	AgentID  string      `json:"agent_id"`
	Answer   string      `json:"answer"`
	Complete bool        `json:"complete"`
	Steps    []StepInfo  `json:"steps"`
}

type StepInfo struct {
	Type    string      `json:"type"`
	Content string      `json:"content"`
}

func (s *Service) Execute(ctx context.Context, req Request) (*Response, error) {
	var state *State
	var err error

	if req.AgentID != "" {
		state, err = s.memory.GetState(req.AgentID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve agent state: %w", err)
		}
	} else {
		state = NewState(req.Message)
	}

	maxIterations := 5
	
	tools := s.getToolsForOpenAI()

	for i := 0; i < maxIterations && !state.IsComplete; i++ {
		messages := s.buildMessages(state)
		
		resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo16K,
			Messages:    messages,
			Tools:       tools,
			Temperature: 0.2,
			MaxTokens:   1000,
		})
		
		if err != nil {
			return nil, fmt.Errorf("failed to get completion: %w", err)
		}
		
		if len(resp.Choices) == 0 {
			return nil, fmt.Errorf("no completion choices returned")
		}
		
		assistantMsg := resp.Choices[0].Message
		
		if len(assistantMsg.ToolCalls) > 0 {
			if assistantMsg.Content != "" {
				thought := assistantMsg.Content
				state.AddThought(thought)
			}
			
			for _, call := range assistantMsg.ToolCalls {
				action := fmt.Sprintf("Using tool '%s' with args: %s", 
					call.Function.Name, call.Function.Arguments)
				state.AddAction(action)
				
				var args map[string]interface{}
				if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
					observation := fmt.Sprintf("Error parsing arguments: %v", err)
					state.AddObservation(observation)
					continue
				}
				
				result, err := s.toolRegistry.CallTool(ctx, call.Function.Name, args)
				var observation string
				if err != nil {
					observation = fmt.Sprintf("Error calling tool: %v", err)
				} else {
					resultJSON, _ := json.MarshalIndent(result, "", "  ")
					observation = string(resultJSON)
				}
				state.AddObservation(observation)
				
				messages = append(messages, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    observation,
					ToolCallID: call.ID,
				})
			}
		} else {
			state.AddAnswer(assistantMsg.Content)
		}
	}
	
	s.memory.SaveState(state)
	
	return &Response{
		AgentID:  state.ID,
		Answer:   s.getFinalAnswer(state),
		Complete: state.IsComplete,
		Steps:    s.formatSteps(state),
	}, nil
}

func (s *Service) buildMessages(state *State) []openai.ChatCompletionMessage {
	systemPrompt := `You are a customer support agent that solves problems step-by-step.
ALWAYS follow this exact process:
1. THINK: First, always start with your reasoning process. Analyze what information you need and your approach.
2. ACT: Only after thinking, use available tools to gather necessary information.
3. OBSERVE: Review the results from your actions.
4. REPEAT steps 1-3 until you have enough information.
5. ANSWER: Provide a clear, complete answer to the customer.

You must explicitly include your reasoning as "Thought: ...your reasoning..." before any tool use.
When you have a final answer, provide it directly without using "Thought:".`

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: state.UserQuery,
		},
	}
	
	if len(state.Steps) > 0 {
		history := state.GetFormattedHistory()
		
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Here is the conversation history so far:\n\n" + history,
		})
	}
	
	return messages
}

func (s *Service) getToolsForOpenAI() []openai.Tool {
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
	return tools
}

func (s *Service) getFinalAnswer(state *State) string {
	for _, step := range state.Steps {
		if step.Type == StepTypeAnswer {
			return step.Content
		}
	}
	
	var lastObservation string
	for i := len(state.Steps) - 1; i >= 0; i-- {
		if state.Steps[i].Type == StepTypeObservation {
			lastObservation = state.Steps[i].Content
			break
		}
	}
	
	if lastObservation != "" {
		return "Based on the information I found: " + lastObservation
	}
	
	return "I couldn't find a specific answer to your question."
}

func (s *Service) formatSteps(state *State) []StepInfo {
    steps := make([]StepInfo, 0, len(state.Steps))
    for _, step := range state.Steps {
        steps = append(steps, StepInfo{
            Type:    string(step.Type),
            Content: step.Content,
        })
    }
    return steps
}
