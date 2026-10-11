package nodes

import (
	"context"
	"fmt"
	"time"

	"agentic/pkg/graph"
	"agentic/pkg/harness"
)

// BuilderNode generates candidate implementation (Mori Rule 7: Separate Builder from Verifier)
type BuilderNode struct {
	gateway *harness.ToolGateway
}

func NewBuilderNode(gateway *harness.ToolGateway) *BuilderNode {
	return &BuilderNode{gateway: gateway}
}

func (b *BuilderNode) Execute(ctx context.Context, state *graph.AgentState) (*graph.AgentState, error) {
	state.Iteration++
	fmt.Printf("\n[BUILDER (Turn %d/%d)]: Generating candidate implementation...\n", state.Iteration, state.MaxIterations)

	// If there were suggested fixes from Jev, apply them
	if state.LastDecision != nil && len(state.LastDecision.SuggestedFixes) > 0 {
		fmt.Printf("[BUILDER]: Addressing prior verification feedback: %v\n", state.LastDecision.SuggestedFixes)
		state.Memory.RecordDecision(
			fmt.Sprintf("Repair turn %d", state.Iteration),
			"Applied targeted fixes rather than restarting scratch conversation",
		)
	} else {
		state.Memory.RecordDecision(
			"Initial candidate patch",
			"Created baseline coupon validation logic",
		)
	}

	// Use Tool Gateway to edit file (Mori Rule 4: Model proposes -> Gateway validates -> Executes)
	input := map[string]any{
		"path":    "services/checkout/coupon.go",
		"content": "package checkout\nfunc ValidateCoupon() bool { return true }",
	}

	_, err := b.gateway.Invoke(ctx, "workspace_edit", input, fmt.Sprintf("edit_%d", state.Iteration))
	if err != nil {
		fmt.Printf("  ⚠ Tool execution intercepted by gateway: %v\n", err)
	}

	// Candidate diff
	state.CurrentDiff = fmt.Sprintf("diff --git a/services/checkout/coupon.go (turn %d)", state.Iteration)
	time.Sleep(30 * time.Millisecond)

	return state, nil
}
