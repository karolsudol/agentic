package harness

import (
	"context"
	"fmt"
)

// WorkspaceEditTool allows modifying files in the workspace (Reversible)
type WorkspaceEditTool struct{}

func (t *WorkspaceEditTool) Name() string    { return "workspace_edit" }
func (t *WorkspaceEditTool) Risk() RiskLevel { return Reversible }

func (t *WorkspaceEditTool) Preconditions(ctx context.Context, input map[string]any) error {
	path, ok := input["path"].(string)
	if !ok || path == "" {
		return fmt.Errorf("precondition error: 'path' must be a non-empty string")
	}
	return nil
}

func (t *WorkspaceEditTool) Execute(ctx context.Context, input map[string]any) (map[string]any, error) {
	path := input["path"].(string)
	return map[string]any{
		"status": "success",
		"path":   path,
		"bytes":  420,
	}, nil
}

// RunTestsTool executes test suites deterministically
type RunTestsTool struct{}

func (t *RunTestsTool) Name() string    { return "run_tests" }
func (t *RunTestsTool) Risk() RiskLevel { return LowRisk }

func (t *RunTestsTool) Preconditions(ctx context.Context, input map[string]any) error {
	return nil
}

func (t *RunTestsTool) Execute(ctx context.Context, input map[string]any) (map[string]any, error) {
	return map[string]any{
		"exit_code": 0,
		"passed":    true,
	}, nil
}

// DeployStagingTool requires authorization on the permission ladder
type DeployStagingTool struct{}

func (t *DeployStagingTool) Name() string    { return "deploy_staging" }
func (t *DeployStagingTool) Risk() RiskLevel { return ExternalEffect }

func (t *DeployStagingTool) Preconditions(ctx context.Context, input map[string]any) error {
	return nil
}

func (t *DeployStagingTool) Execute(ctx context.Context, input map[string]any) (map[string]any, error) {
	return map[string]any{
		"status": "deployed_to_staging",
	}, nil
}
