package evaluation

import (
	"time"

	"github.com/google/uuid"
)

type PatternType string

const (
	PatternBasicLLMCompletion PatternType = "basic_llm_completion"
	PatternKnowledgeRAG       PatternType = "knowledge_rag"
	PatternFunctionCalling    PatternType = "function_calling"
	PatternReasoningAgent     PatternType = "reasoning_agent"
	PatternMultiAgent         PatternType = "multi_agent"
)

type EvaluationResult struct {
	ID              string      `json:"id"`
	PatternType     PatternType `json:"pattern_type"`
	Query           string      `json:"query"`
	Response        string      `json:"response"`
	ResponseTime    int64       `json:"response_time_ms"`
	TokensUsed      int         `json:"tokens_used,omitempty"`
	HumanRating     int         `json:"human_rating,omitempty"`
	AutoRating      float64     `json:"auto_rating,omitempty"`
	EvaluationNotes string      `json:"evaluation_notes,omitempty"`
}

type Report struct {
	ID        string             `json:"id"`
	Query     string             `json:"query"`
	Timestamp time.Time          `json:"timestamp"`
	Results   []EvaluationResult `json:"results"`
}

func NewReport(query string) *Report {
	return &Report{
		ID:        uuid.New().String(),
		Query:     query,
		Timestamp: time.Now(),
		Results:   []EvaluationResult{},
	}
}

func (r *Report) AddResult(result EvaluationResult) {
	r.Results = append(r.Results, result)
}
