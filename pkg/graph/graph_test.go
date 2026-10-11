package graph

import (
	"context"
	"testing"

	"agentic/pkg/contract"
	"agentic/pkg/harness"
)

func TestStateGraph_ExecutionAndInterrupt(t *testing.T) {
	ctx := context.Background()
	g := NewStateGraph()

	executedNodes := make([]string, 0)

	g.AddNode("node_a", func(ctx context.Context, state *AgentState) (*AgentState, error) {
		executedNodes = append(executedNodes, "node_a")
		return state, nil
	})

	g.AddNode("node_b", func(ctx context.Context, state *AgentState) (*AgentState, error) {
		executedNodes = append(executedNodes, "node_b")
		return state, nil
	})

	g.AddNode("protected_node", func(ctx context.Context, state *AgentState) (*AgentState, error) {
		executedNodes = append(executedNodes, "protected_node")
		return state, nil
	})

	g.SetEntryPoint("node_a")
	g.AddEdge("node_a", "node_b")
	g.AddEdge("node_b", "protected_node")
	g.AddEdge("protected_node", EndNode)

	// Add HITL Interrupt on protected_node
	g.AddInterrupt("protected_node")

	c := &contract.TaskContract{ID: "test"}
	pm := harness.NewProjectMap()
	state := NewInitialState(c, pm, 3)

	// 1. Initial execution should pause at protected_node
	resState, err := g.Execute(ctx, state)
	if err != nil {
		t.Fatalf("unexpected error during graph execution: %v", err)
	}

	if !resState.RequiresHuman {
		t.Fatal("expected graph to pause and require human approval before protected_node")
	}

	if len(executedNodes) != 2 || executedNodes[0] != "node_a" || executedNodes[1] != "node_b" {
		t.Fatalf("expected node_a and node_b to execute before interrupt, got: %v", executedNodes)
	}

	// 2. Grant approval and resume
	resState.HumanApproved = true
	resState, err = g.Execute(ctx, resState)
	if err != nil {
		t.Fatalf("unexpected error during resumed execution: %v", err)
	}

	if !resState.ExecutionEnded {
		t.Fatal("expected resumed graph execution to finish at EndNode")
	}

	if len(executedNodes) != 3 || executedNodes[2] != "protected_node" {
		t.Fatalf("expected protected_node to execute after approval, got: %v", executedNodes)
	}
}
