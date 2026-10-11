package nodes

import (
	"context"
	"fmt"
	"sync"
	"time"

	"agentic/pkg/graph"
	"agentic/pkg/harness"
)

// SubTask represents a decomposed unit of work for parallel execution (Makarevich's Fleet pattern)
type SubTask struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Area        string `json:"area"`
}

// SupervisorNode orchestrates the fleet and progressive disclosure
type SupervisorNode struct {
	gateway *harness.ToolGateway
}

func NewSupervisorNode(gateway *harness.ToolGateway) *SupervisorNode {
	return &SupervisorNode{gateway: gateway}
}

func (s *SupervisorNode) Execute(ctx context.Context, state *graph.AgentState) (*graph.AgentState, error) {
	fmt.Printf("\n[SUPERVISOR]: Analyzing Task Contract '%s'\n", state.Contract.Objective)
	state.Memory.RecordDecision(
		"Adopt Spec-Driven Decomposition",
		fmt.Sprintf("Decomposed objective based on %d constraints and %d done_when gates", len(state.Contract.Constraints), len(state.Contract.DoneWhen)),
	)

	// Progressive disclosure: identify relevant project areas (Mori Rule 3)
	fmt.Println("[SUPERVISOR]: Inspecting Project Map attention index:")
	fmt.Print(state.ProjectMap.Summary())

	// Fleet Fan-Out: Run subtasks in parallel using Go goroutines (Makarevich Fleet pattern)
	subtasks := []SubTask{
		{ID: "task_api", Description: "Implement bounded checkout API validation", Area: "backend"},
		{ID: "task_test", Description: "Write regression test suite for double-coupon edge case", Area: "tests"},
	}

	fmt.Printf("[SUPERVISOR]: Spawning fleet of %d parallel workers (goroutines)...\n", len(subtasks))
	var wg sync.WaitGroup
	resultsChan := make(chan string, len(subtasks))

	for _, st := range subtasks {
		wg.Add(1)
		go func(task SubTask) {
			defer wg.Done()
			// Each worker runs concurrently in its own goroutine
			time.Sleep(20 * time.Millisecond) // Simulated parallel work
			resultsChan <- fmt.Sprintf("Completed %s in area %s", task.ID, task.Area)
		}(st)
	}

	wg.Wait()
	close(resultsChan)

	for res := range resultsChan {
		fmt.Printf("  ✓ Worker: %s\n", res)
	}

	state.Memory.SetStateVar("subtasks_completed", len(subtasks))
	return state, nil
}
