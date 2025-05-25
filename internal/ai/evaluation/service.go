package evaluation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/basic_llm_completion"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/function_calling"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/knowledge_rag"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/multi_agent"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/reasoning_agent"
	"github.com/junjie-w/llm-integration-patterns-experiments/pkg/config"

	"github.com/sashabaranov/go-openai"
)

type Service struct {
	basicService      *basic_llm_completion.Service
	knowledgeService  *knowledge_rag.Service
	functionService   *function_calling.Service
	reasoningService  *reasoning_agent.Service
	multiAgentService *multi_agent.Service
	evalClient        *openai.Client
	reports           map[string]*Report
}

func NewService(
	cfg *config.Config,
	basicService *basic_llm_completion.Service,
	knowledgeService *knowledge_rag.Service,
	functionService *function_calling.Service,
	reasoningService *reasoning_agent.Service,
	multiAgentService *multi_agent.Service,
) *Service {
	return &Service{
		basicService:      basicService,
		knowledgeService:  knowledgeService,
		functionService:   functionService,
		reasoningService:  reasoningService,
		multiAgentService: multiAgentService,
		evalClient:        openai.NewClient(cfg.OpenAIKey),
		reports:           make(map[string]*Report),
	}
}

type Request struct {
	Query        string        `json:"query"`
	PatternTypes []PatternType `json:"pattern_types"`
}

type Response struct {
	ReportID string  `json:"report_id"`
	Report   *Report `json:"report"`
}

func (s *Service) Evaluate(ctx context.Context, req Request) (*Response, error) {
	report := NewReport(req.Query)
	
	for _, patternType := range req.PatternTypes {
		result, err := s.evaluatePattern(ctx, patternType, req.Query)
		if err != nil {
			fmt.Printf("Error evaluating pattern %s: %v\n", patternType, err)
			continue
		}
		
		report.AddResult(result)
	}
	
	s.reports[report.ID] = report
	
	return &Response{
		ReportID: report.ID,
		Report:   report,
	}, nil
}

func (s *Service) GetReport(reportID string) (*Report, error) {
	report, exists := s.reports[reportID]
	if !exists {
		return nil, fmt.Errorf("report with ID %s not found", reportID)
	}
	
	return report, nil
}

func (s *Service) evaluatePattern(ctx context.Context, patternType PatternType, query string) (EvaluationResult, error) {
	result := EvaluationResult{
		ID:          uuid.New().String(),
		PatternType: patternType,
		Query:       query,
	}
	
	startTime := time.Now()
	var response string
	var err error
	
	switch patternType {
	case PatternBasicLLMCompletion:
		resp, err := s.basicService.GetCompletion(ctx, basic_llm_completion.Request{Message: query})
		if err != nil {
			return result, err
		}
		response = resp.Reply
		
	case PatternKnowledgeRAG:
		resp, err := s.knowledgeService.GetCompletion(ctx, knowledge_rag.Request{Message: query, UseVectorSearch: true})
		if err != nil {
			return result, err
		}
		response = resp.Reply
		
	case PatternFunctionCalling:
		resp, err := s.functionService.GetCompletion(ctx, function_calling.Request{Message: query})
		if err != nil {
			return result, err
		}
		response = resp.Reply
		
	case PatternReasoningAgent:
		resp, err := s.reasoningService.Execute(ctx, reasoning_agent.Request{Message: query})
		if err != nil {
			return result, err
		}
		response = resp.Answer
		
	case PatternMultiAgent:
		resp, err := s.multiAgentService.Process(ctx, multi_agent.Request{Message: query})
		if err != nil {
			return result, err
		}
		response = resp.Reply
		
	default:
		return result, fmt.Errorf("unsupported pattern type: %s", patternType)
	}
	
	elapsed := time.Since(startTime)
	
	result.Response = response
	result.ResponseTime = elapsed.Milliseconds()
	
	autoRating, notes, err := s.autoEvaluate(ctx, query, response)
	if err == nil {
		result.AutoRating = autoRating
		result.EvaluationNotes = notes
	}
	
	return result, nil
}

func (s *Service) autoEvaluate(ctx context.Context, query, response string) (float64, string, error) {
	prompt := fmt.Sprintf(`Evaluate this customer support response to the query.
  
Query: %s

Response: %s

Criteria:
- Relevance: Does it address the query directly?
- Accuracy: Is the information correct?
- Completeness: Does it fully answer the question?
- Clarity: Is it easy to understand?
- Helpfulness: Is it actually helpful to the customer?

Rate the response on a scale of 0 to 1 where:
- 0.0-0.2: Poor (completely fails to address the query)
- 0.3-0.4: Below Average (partially addresses but with major issues)
- 0.5-0.6: Average (addresses the query but with some issues)
- 0.7-0.8: Good (addresses the query well with minor issues)
- 0.9-1.0: Excellent (perfectly addresses the query)

Provide your rating as a single number followed by a brief explanation.

Rating: `, query, response)

	resp, err := s.evalClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an objective evaluator of customer support responses. Provide fair ratings based on the given criteria.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})
	
	if err != nil {
		return 0, "", err
	}
	
	if len(resp.Choices) == 0 {
		return 0, "", fmt.Errorf("no evaluation response generated")
	}
	
	content := resp.Choices[0].Message.Content
	
	var rating float64
	var explanation string
	
	_, err = fmt.Sscanf(content, "Rating: %f", &rating)
	if err != nil {
		rating = 0.5
		explanation = content
	} else {
		parts := strings.SplitN(content, "\n", 2)
		if len(parts) > 1 {
			explanation = parts[1]
		} else {
			explanation = "No detailed explanation provided."
		}
	}
	
	return rating, explanation, nil
}
