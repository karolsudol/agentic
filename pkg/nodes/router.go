package nodes

import (
	"context"
	"fmt"

	"agentic/pkg/graph"
	"agentic/pkg/jev"
)

// JevEdgeRouter provides the conditional edge routing using Jev's System 1 typed decisions
type JevEdgeRouter struct {
	jevClient jev.Client
}

func NewJevEdgeRouter(client jev.Client) *JevEdgeRouter {
	return &JevEdgeRouter{jevClient: client}
}

func (r *JevEdgeRouter) Route(ctx context.Context, state *graph.AgentState) (string, error) {
	fmt.Printf("\n[JEV ROUTER]: Evaluating routing decision (System 1)...\n")

	decision, err := r.jevClient.DecideRoute(ctx, state.Contract, state.Evidence, state.Iteration, state.MaxIterations)
	if err != nil {
		return "", fmt.Errorf("jev decision failed: %w", err)
	}

	state.LastDecision = decision

	fmt.Printf("  • Jev Action: %s (Confidence: %.2f)\n", decision.Action, decision.Confidence)
	fmt.Printf("  • Reason: %s\n", decision.Reason)

	switch decision.Action {
	case jev.ActionProceedToGate:
		return "human_gate", nil
	case jev.ActionRetryRepair:
		return "builder", nil
	case jev.ActionHaltExhausted:
		return "receipt_node", nil
	default:
		return "receipt_node", nil
	}
}
