package graph

import (
	"agentic/pkg/contract"
	"agentic/pkg/harness"
	"agentic/pkg/jev"
)

// AgentState represents the typed state flowing through the ADK 2.0 graph
type AgentState struct {
	Contract       *contract.TaskContract `json:"contract"`
	Memory         *harness.DurableMemory `json:"memory"`
	ProjectMap     *harness.ProjectMap    `json:"project_map"`
	CurrentDiff    string                 `json:"current_diff"`
	Evidence       []contract.Evidence    `json:"evidence"`
	LastDecision   *jev.RoutingDecision   `json:"last_decision"`
	Receipt        *contract.Receipt      `json:"receipt"`
	Iteration      int                    `json:"iteration"`
	MaxIterations  int                    `json:"max_iterations"`
	RequiresHuman  bool                   `json:"requires_human"`
	HumanApproved  bool                   `json:"human_approved"`
	ExecutionEnded bool                   `json:"execution_ended"`
}

func NewInitialState(c *contract.TaskContract, pm *harness.ProjectMap, maxIterations int) *AgentState {
	return &AgentState{
		Contract:      c,
		Memory:        harness.NewDurableMemory(),
		ProjectMap:    pm,
		Evidence:      make([]contract.Evidence, 0),
		Iteration:     0,
		MaxIterations: maxIterations,
	}
}
