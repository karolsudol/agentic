package jev

import (
	"context"
	"testing"

	"agentic/pkg/contract"
)

func TestLocalJevClient_DecideRoute(t *testing.T) {
	ctx := context.Background()
	client := NewLocalJevClient()
	c := &contract.TaskContract{ID: "test_task"}

	// 1. With failing evidence -> should decide RETRY_REPAIR
	failingEvidence := []contract.Evidence{
		{Claim: "Unit tests", Passed: false, Type: contract.VerifyUnitTest},
	}
	decision, err := client.DecideRoute(ctx, c, failingEvidence, 1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decision.Action != ActionRetryRepair {
		t.Fatalf("expected ActionRetryRepair, got %s", decision.Action)
	}
	if decision.Confidence < 0.90 {
		t.Fatalf("expected confidence >= 0.90, got %f", decision.Confidence)
	}

	// 2. With all passing evidence -> should decide PROCEED_TO_GATE
	passingEvidence := []contract.Evidence{
		{Claim: "Unit tests", Passed: true, Type: contract.VerifyUnitTest},
	}
	decision, err = client.DecideRoute(ctx, c, passingEvidence, 1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decision.Action != ActionProceedToGate {
		t.Fatalf("expected ActionProceedToGate, got %s", decision.Action)
	}

	// 3. Exhausted budget -> should decide HALT_EXHAUSTED
	decision, err = client.DecideRoute(ctx, c, failingEvidence, 3, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decision.Action != ActionHaltExhausted {
		t.Fatalf("expected ActionHaltExhausted, got %s", decision.Action)
	}
}
