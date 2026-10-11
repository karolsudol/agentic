package contract

import (
	"fmt"
	"strings"
	"time"
)

// VerificationType classifies how evidence was produced (Mori Rule 6: deterministic checks first)
type VerificationType string

const (
	VerifySyntax      VerificationType = "SYNTAX"
	VerifyTypes       VerificationType = "TYPE_CHECK"
	VerifyUnitTest    VerificationType = "UNIT_TEST"
	VerifyIntegration VerificationType = "INTEGRATION_TEST"
	VerifySemantic    VerificationType = "SEMANTIC_REVIEW"
)

// Evidence represents verifiable proof of an outcome (Mori Rule 6)
type Evidence struct {
	Claim     string           `json:"claim"`
	Type      VerificationType `json:"type"`
	Passed    bool             `json:"passed"`
	Details   string           `json:"details"`
	Timestamp time.Time        `json:"timestamp"`
}

// TaskContract turns natural language requests into bounded engineering contracts (Mori Rule 2)
type TaskContract struct {
	ID               string   `json:"id"`
	Objective        string   `json:"objective"`
	Inputs           []string `json:"inputs"`
	Constraints      []string `json:"constraints"`
	Deliverable      string   `json:"deliverable"`
	DoneWhen         []string `json:"done_when"`         // Completion conditions (Mori: "done_when")
	ApprovalRequired []string `json:"approval_required"` // Actions requiring human gate
}

// Receipt is the final verifiable summary delivered to humans (Mori Rule 12)
type Receipt struct {
	TaskID          string     `json:"task_id"`
	Objective       string     `json:"objective"`
	ChangedFiles    []string   `json:"changed_files"`
	Verified        []Evidence `json:"verified"`
	UnverifiedRisks []string   `json:"unverified_risks"`
	ApprovalStatus  string     `json:"approval_status"`
	CompletedAt     time.Time  `json:"completed_at"`
}

func (r *Receipt) Format() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== RUN RECEIPT: %s ===\n", r.TaskID))
	sb.WriteString(fmt.Sprintf("Objective: %s\n", r.Objective))
	sb.WriteString(fmt.Sprintf("Status: %s (at %s)\n", r.ApprovalStatus, r.CompletedAt.Format(time.RFC3339)))

	sb.WriteString("\n[Files Changed]:\n")
	for _, f := range r.ChangedFiles {
		sb.WriteString(fmt.Sprintf("  • %s\n", f))
	}

	sb.WriteString("\n[Verifiable Evidence]:\n")
	for _, ev := range r.Verified {
		status := "PASS"
		if !ev.Passed {
			status = "FAIL"
		}
		sb.WriteString(fmt.Sprintf("  [%s] (%s): %s -> %s\n", status, ev.Type, ev.Claim, ev.Details))
	}

	if len(r.UnverifiedRisks) > 0 {
		sb.WriteString("\n[Open Risks & Unverified Items]:\n")
		for _, risk := range r.UnverifiedRisks {
			sb.WriteString(fmt.Sprintf("  ⚠ %s\n", risk))
		}
	}
	sb.WriteString("===============================\n")
	return sb.String()
}
