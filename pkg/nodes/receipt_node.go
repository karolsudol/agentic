package nodes

import (
	"context"
	"fmt"
	"time"

	"agentic/pkg/contract"
	"agentic/pkg/graph"
)

// ReceiptNode compiles the final verifiable receipt (Mori Rule 12)
type ReceiptNode struct{}

func NewReceiptNode() *ReceiptNode {
	return &ReceiptNode{}
}

func (r *ReceiptNode) Execute(ctx context.Context, state *graph.AgentState) (*graph.AgentState, error) {
	fmt.Printf("\n[RECEIPT NODE]: Generating final run receipt...\n")

	status := "APPROVED & COMPLETED"
	if !state.HumanApproved && state.RequiresHuman {
		status = "WAITING_HUMAN_APPROVAL"
	} else if state.LastDecision != nil && state.LastDecision.Action == "HALT_EXHAUSTED" {
		status = "HALTED_BUDGET_EXHAUSTED"
	}

	state.Receipt = &contract.Receipt{
		TaskID:          state.Contract.ID,
		Objective:       state.Contract.Objective,
		ChangedFiles:    []string{"services/checkout/coupon.go", "tests/checkout/coupon_test.go"},
		Verified:        state.Evidence,
		UnverifiedRisks: []string{"Legacy mobile client v1.2 might require endpoint compatibility test"},
		ApprovalStatus:  status,
		CompletedAt:     time.Now(),
	}

	fmt.Print(state.Receipt.Format())
	return state, nil
}
