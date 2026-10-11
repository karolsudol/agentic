package jev

// Action represents discrete System 1 routing decisions
type Action string

const (
	ActionProceedToGate Action = "PROCEED_TO_GATE"
	ActionRetryRepair   Action = "RETRY_REPAIR"
	ActionHaltExhausted Action = "HALT_EXHAUSTED"
	ActionEscalateHuman Action = "ESCALATE_HUMAN"
)

// RoutingDecision is the schema-guaranteed typed output from Jev
type RoutingDecision struct {
	Action         Action   `json:"action"`
	Confidence     float64  `json:"confidence"` // Calibrated probability (e.g. 0.98)
	Reason         string   `json:"reason"`
	SuggestedFixes []string `json:"suggested_fixes,omitempty"`
}

// VerifierScore represents Jev's evaluation of candidate disproof
type VerifierScore struct {
	Passed          bool    `json:"passed"`
	Confidence      float64 `json:"confidence"`
	DisproofFound   bool    `json:"disproof_found"`
	RejectionReason string  `json:"rejection_reason,omitempty"`
}
