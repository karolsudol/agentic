package harness

import (
	"context"
	"testing"
)

func TestToolGateway_PreconditionValidation(t *testing.T) {
	ctx := context.Background()
	policy := NewPolicyEngine(true)
	gateway := NewToolGateway(policy)
	gateway.Register(&WorkspaceEditTool{})

	// Missing required "path" should fail precondition
	_, err := gateway.Invoke(ctx, "workspace_edit", map[string]any{}, "action_1")
	if err == nil {
		t.Fatal("expected error due to missing path precondition, got nil")
	}

	// Valid path should succeed
	output, err := gateway.Invoke(ctx, "workspace_edit", map[string]any{"path": "file.go"}, "action_2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output["status"] != "success" {
		t.Fatalf("expected success status, got: %v", output["status"])
	}
}

func TestPolicyEngine_PermissionLadder(t *testing.T) {
	policy := NewPolicyEngine(true)

	// Low risk is automatically authorized
	auth, err := policy.Authorize("inspect", LowRisk, "act_1")
	if !auth || err != nil {
		t.Fatalf("expected LowRisk to be auto authorized, err: %v", err)
	}

	// ExternalEffect requires approval
	auth, err = policy.Authorize("deploy", ExternalEffect, "deploy_act")
	if auth || err == nil {
		t.Fatal("expected ExternalEffect to be blocked pending approval")
	}

	// Grant approval and verify
	policy.GrantApproval("deploy_act")
	auth, err = policy.Authorize("deploy", ExternalEffect, "deploy_act")
	if !auth || err != nil {
		t.Fatalf("expected approval to be granted after explicit grant, err: %v", err)
	}
}
