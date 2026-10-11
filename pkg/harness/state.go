package harness

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Decision documents an architectural choice and rationale
type Decision struct {
	ID        string    `json:"id"`
	Choice    string    `json:"choice"`
	Rationale string    `json:"rationale"`
	MadeAt    time.Time `json:"made_at"`
}

// Lesson captures failure modes so future runs do not repeat them (Mori Rule 5)
type Lesson struct {
	Trigger  string `json:"trigger"`
	Failure  string `json:"failure"`
	RuleForm string `json:"rule_form"` // What policy/test now prevents it (Rule 10)
}

// DurableMemory stores persistent system state outside the LLM chat transcript (Mori Rule 5)
type DurableMemory struct {
	mu        sync.RWMutex
	Facts     map[string]string `json:"facts"`     // Stable domain knowledge
	Decisions []Decision        `json:"decisions"` // Choices made and why
	State     map[string]any    `json:"state"`     // Current step, variables
	Lessons   []Lesson          `json:"lessons"`   // Retrospective feedback
}

func NewDurableMemory() *DurableMemory {
	return &DurableMemory{
		Facts:     make(map[string]string),
		Decisions: make([]Decision, 0),
		State:     make(map[string]any),
		Lessons:   make([]Lesson, 0),
	}
}

func (m *DurableMemory) SetFact(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Facts[key] = value
}

func (m *DurableMemory) RecordDecision(choice, rationale string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Decisions = append(m.Decisions, Decision{
		ID:        fmt.Sprintf("dec_%d", len(m.Decisions)+1),
		Choice:    choice,
		Rationale: rationale,
		MadeAt:    time.Now(),
	})
}

func (m *DurableMemory) SetStateVar(key string, val any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.State[key] = val
}

func (m *DurableMemory) RecordLesson(trigger, failure, ruleForm string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Lessons = append(m.Lessons, Lesson{
		Trigger:  trigger,
		Failure:  failure,
		RuleForm: ruleForm,
	})
}

func (m *DurableMemory) Clone() *DurableMemory {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, _ := json.Marshal(m)
	var cloned DurableMemory
	_ = json.Unmarshal(data, &cloned)
	return &cloned
}
