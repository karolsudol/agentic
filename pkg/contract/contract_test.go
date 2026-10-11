package contract

import (
	"strings"
	"testing"
	"time"
)

func TestReceiptFormatting(t *testing.T) {
	receipt := &Receipt{
		TaskID:       "task_100",
		Objective:    "Test double coupon fix",
		ChangedFiles: []string{"services/coupon.go"},
		Verified: []Evidence{
			{
				Claim:     "Syntax valid",
				Type:      VerifySyntax,
				Passed:    true,
				Details:   "exit code 0",
				Timestamp: time.Now(),
			},
		},
		UnverifiedRisks: []string{"Potential mobile cache staleness"},
		ApprovalStatus:  "APPROVED",
		CompletedAt:     time.Now(),
	}

	formatted := receipt.Format()
	if !strings.Contains(formatted, "task_100") {
		t.Fatalf("expected formatted receipt to contain task ID, got: %s", formatted)
	}
	if !strings.Contains(formatted, "[PASS] (SYNTAX)") {
		t.Fatalf("expected formatted receipt to show PASS syntax check, got: %s", formatted)
	}
	if !strings.Contains(formatted, "Potential mobile cache staleness") {
		t.Fatalf("expected formatted receipt to show open risk, got: %s", formatted)
	}
}
