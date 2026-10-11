package harness

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Tool represents a capability governed by the harness (Mori Rule 4)
type Tool interface {
	Name() string
	Risk() RiskLevel
	Preconditions(ctx context.Context, input map[string]any) error
	Execute(ctx context.Context, input map[string]any) (map[string]any, error)
}

// ToolTrace records deterministic tool execution history (Mori Rule 11)
type ToolTrace struct {
	ToolName   string         `json:"tool_name"`
	Risk       RiskLevel      `json:"risk"`
	Input      map[string]any `json:"input"`
	Output     map[string]any `json:"output"`
	Error      string         `json:"error,omitempty"`
	Duration   time.Duration  `json:"duration"`
	ExecutedAt time.Time      `json:"executed_at"`
}

// ToolGateway wraps tool execution with policy authorization, precondition validation, and audit traces
type ToolGateway struct {
	mu     sync.Mutex
	policy *PolicyEngine
	tools  map[string]Tool
	traces []ToolTrace
}

func NewToolGateway(policy *PolicyEngine) *ToolGateway {
	return &ToolGateway{
		policy: policy,
		tools:  make(map[string]Tool),
		traces: make([]ToolTrace, 0),
	}
}

func (g *ToolGateway) Register(tool Tool) {
	g.tools[tool.Name()] = tool
}

// Invoke executes a tool through the harness gateway (Model Proposes -> Gateway Validates -> Policy Authorizes -> Tool Executes -> Harness Records)
func (g *ToolGateway) Invoke(ctx context.Context, toolName string, input map[string]any, actionID string) (map[string]any, error) {
	tool, exists := g.tools[toolName]
	if !exists {
		return nil, fmt.Errorf("tool '%s' not registered in harness", toolName)
	}

	start := time.Now()

	// 1. Gateway Precondition Validation
	if err := tool.Preconditions(ctx, input); err != nil {
		g.recordTrace(toolName, tool.Risk(), input, nil, err, time.Since(start))
		return nil, fmt.Errorf("precondition failed for tool %s: %w", toolName, err)
	}

	// 2. Policy Authorization Ladder
	authorized, err := g.policy.Authorize(toolName, tool.Risk(), actionID)
	if !authorized || err != nil {
		g.recordTrace(toolName, tool.Risk(), input, nil, err, time.Since(start))
		return nil, fmt.Errorf("authorization rejected: %w", err)
	}

	// 3. Tool Execution
	output, execErr := tool.Execute(ctx, input)
	duration := time.Since(start)

	// 4. Trace Recording
	g.recordTrace(toolName, tool.Risk(), input, output, execErr, duration)

	if execErr != nil {
		return nil, fmt.Errorf("tool %s execution error: %w", toolName, execErr)
	}

	return output, nil
}

func (g *ToolGateway) recordTrace(name string, risk RiskLevel, in, out map[string]any, err error, d time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	g.traces = append(g.traces, ToolTrace{
		ToolName:   name,
		Risk:       risk,
		Input:      in,
		Output:     out,
		Error:      errMsg,
		Duration:   d,
		ExecutedAt: time.Now(),
	})
}

func (g *ToolGateway) Traces() []ToolTrace {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]ToolTrace, len(g.traces))
	copy(out, g.traces)
	return out
}
