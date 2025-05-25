package reasoning_agent

import (
	"time"

	"github.com/google/uuid"
)

type Step struct {
	ID        string    `json:"id"`
	Type      StepType  `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type StepType string

const (
	StepTypeThought StepType = "thought"
	StepTypeAction  StepType = "action"
	StepTypeObservation StepType = "observation"
	StepTypeAnswer  StepType = "answer"
)

type State struct {
	ID            string    `json:"id"`
	UserQuery     string    `json:"user_query"`
	Steps         []Step    `json:"steps"`
	CreatedAt     time.Time `json:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	IsComplete    bool      `json:"is_complete"`
}

func NewState(query string) *State {
	now := time.Now()
	return &State{
		ID:            uuid.New().String(),
		UserQuery:     query,
		Steps:         []Step{},
		CreatedAt:     now,
		LastUpdatedAt: now,
		IsComplete:    false,
	}
}

func (s *State) AddThought(content string) {
	s.addStep(StepTypeThought, content)
}

func (s *State) AddAction(content string) {
	s.addStep(StepTypeAction, content)
}

func (s *State) AddObservation(content string) {
	s.addStep(StepTypeObservation, content)
}

func (s *State) AddAnswer(content string) {
	s.addStep(StepTypeAnswer, content)
	s.IsComplete = true
}

func (s *State) addStep(stepType StepType, content string) {
	now := time.Now()
	step := Step{
		ID:        uuid.New().String(),
		Type:      stepType,
		Content:   content,
		Timestamp: now,
	}
	s.Steps = append(s.Steps, step)
	s.LastUpdatedAt = now
}

func (s *State) GetFormattedHistory() string {
	var result string
	
	for _, step := range s.Steps {
		switch step.Type {
		case StepTypeThought:
			result += "Thought: " + step.Content + "\n\n"
		case StepTypeAction:
			result += "Action: " + step.Content + "\n\n"
		case StepTypeObservation:
			result += "Observation: " + step.Content + "\n\n"
		case StepTypeAnswer:
			result += "Answer: " + step.Content + "\n\n"
		}
	}
	
	return result
}
