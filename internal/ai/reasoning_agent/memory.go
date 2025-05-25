package reasoning_agent

import (
	"fmt"
	"sync"
)

type Memory struct {
	states map[string]*State
	mu     sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		states: make(map[string]*State),
	}
}

func (m *Memory) SaveState(state *State) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states[state.ID] = state
}

func (m *Memory) GetState(id string) (*State, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	state, exists := m.states[id]
	if !exists {
		return nil, fmt.Errorf("agent state with ID %s not found", id)
	}
	
	return state, nil
}

func (m *Memory) DeleteState(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.states, id)
}

func (m *Memory) ListStates() []*State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	states := make([]*State, 0, len(m.states))
	for _, state := range m.states {
		states = append(states, state)
	}
	
	return states
}
