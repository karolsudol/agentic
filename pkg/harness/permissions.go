package harness

import (
	"fmt"
)

// RiskLevel defines the permission ladder (Mori Rule 8)
type RiskLevel int

const (
	// LowRisk: Read, inspect, search -> Automatic
	LowRisk RiskLevel = iota
	// Reversible: Edit workspace files, run tests, format -> Automatic + Trace
	Reversible
	// ExternalEffect: Send network requests, staging deploy, mutate remote resources -> Approval required
	ExternalEffect
	// Irreversible: Production deploy, delete database, rotate root credentials -> Hard gate / human required
	Irreversible
)

func (r RiskLevel) String() string {
	switch r {
	case LowRisk:
		return "LOW_RISK (Auto)"
	case Reversible:
		return "REVERSIBLE (Auto+Trace)"
	case ExternalEffect:
		return "EXTERNAL_EFFECT (Approval Required)"
	case Irreversible:
		return "IRREVERSIBLE (Hard Gate)"
	default:
		return "UNKNOWN"
	}
}

// PolicyEngine enforces permissions outside the model's instructions (Mori Rule 8)
type PolicyEngine struct {
	autoApproveReversible bool
	grantedApprovals      map[string]bool
}

func NewPolicyEngine(autoApproveReversible bool) *PolicyEngine {
	return &PolicyEngine{
		autoApproveReversible: autoApproveReversible,
		grantedApprovals:      make(map[string]bool),
	}
}

func (p *PolicyEngine) Authorize(toolName string, risk RiskLevel, actionID string) (bool, error) {
	switch risk {
	case LowRisk:
		return true, nil
	case Reversible:
		if p.autoApproveReversible {
			return true, nil
		}
		if p.grantedApprovals[actionID] {
			return true, nil
		}
		return false, fmt.Errorf("tool %s requires authorization for reversible modification", toolName)
	case ExternalEffect, Irreversible:
		if p.grantedApprovals[actionID] {
			return true, nil
		}
		return false, fmt.Errorf("policy gate: tool %s [%s] is blocked pending explicit human approval (actionID: %s)", toolName, risk, actionID)
	default:
		return false, fmt.Errorf("unknown risk tier for tool: %s", toolName)
	}
}

func (p *PolicyEngine) GrantApproval(actionID string) {
	p.grantedApprovals[actionID] = true
}
