package gemini

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"agentic/pkg/contract"
	"agentic/pkg/jev"
)

// newTestClient creates a Client that points at a local httptest server.
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	c := NewClient("test-key", "test-model")
	c.baseURL = server.URL
	return c
}

func TestClient_GenerateContent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"candidates": []map[string]any{
				{
					"content": map[string]any{
						"parts": []map[string]any{
							{"text": "Hello from Gemini!"},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	client := newTestClient(t, handler)
	result, err := client.GenerateContent(context.Background(), "system", "user prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello from Gemini!" {
		t.Fatalf("expected 'Hello from Gemini!', got %q", result)
	}
}

func TestClient_GenerateContent_EmptyKey(t *testing.T) {
	client := NewClient("", "test-model")
	_, err := client.GenerateContent(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
}

func TestClient_DecideRoute_ProceedToGate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decision := jev.RoutingDecision{
			Action:     jev.ActionProceedToGate,
			Confidence: 0.98,
			Reason:     "all checks passed",
		}
		decisionJSON, _ := json.Marshal(decision)
		resp := map[string]any{
			"candidates": []map[string]any{
				{
					"content": map[string]any{
						"parts": []map[string]any{
							{"text": string(decisionJSON)},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	client := newTestClient(t, handler)
	tc := &contract.TaskContract{ID: "test", Objective: "fix bug"}
	evidence := []contract.Evidence{
		{Claim: "syntax ok", Passed: true, Type: contract.VerifySyntax},
	}

	decision, err := client.DecideRoute(context.Background(), tc, evidence, 1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decision.Action != jev.ActionProceedToGate {
		t.Fatalf("expected PROCEED_TO_GATE, got %s", decision.Action)
	}
	if decision.Confidence < 0.90 {
		t.Fatalf("expected high confidence, got %f", decision.Confidence)
	}
}

func TestClient_DecideRoute_RetryRepair(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decision := jev.RoutingDecision{
			Action:         jev.ActionRetryRepair,
			Confidence:     0.95,
			Reason:         "unit test failure detected",
			SuggestedFixes: []string{"fix double-coupon edge case"},
		}
		decisionJSON, _ := json.Marshal(decision)
		resp := map[string]any{
			"candidates": []map[string]any{
				{
					"content": map[string]any{
						"parts": []map[string]any{
							{"text": string(decisionJSON)},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	client := newTestClient(t, handler)
	tc := &contract.TaskContract{ID: "test", Objective: "fix bug"}
	evidence := []contract.Evidence{
		{Claim: "unit test", Passed: false, Type: contract.VerifyUnitTest},
	}

	decision, err := client.DecideRoute(context.Background(), tc, evidence, 1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decision.Action != jev.ActionRetryRepair {
		t.Fatalf("expected RETRY_REPAIR, got %s", decision.Action)
	}
	if len(decision.SuggestedFixes) == 0 {
		t.Fatal("expected suggested fixes, got none")
	}
}

func TestClient_ScoreVerification_Pass(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		score := jev.VerifierScore{
			Passed:        true,
			Confidence:    0.97,
			DisproofFound: false,
		}
		scoreJSON, _ := json.Marshal(score)
		resp := map[string]any{
			"candidates": []map[string]any{
				{
					"content": map[string]any{
						"parts": []map[string]any{
							{"text": string(scoreJSON)},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	client := newTestClient(t, handler)
	tc := &contract.TaskContract{ID: "test", Objective: "fix bug"}
	evidence := []contract.Evidence{
		{Claim: "all tests pass", Passed: true, Type: contract.VerifyUnitTest},
	}

	score, err := client.ScoreVerification(context.Background(), tc, "diff --git ...", evidence)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !score.Passed {
		t.Fatal("expected score to pass")
	}
	if score.DisproofFound {
		t.Fatal("expected no disproof")
	}
}

func TestClient_ScoreVerification_Disproof(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		score := jev.VerifierScore{
			Passed:          false,
			Confidence:      0.96,
			DisproofFound:   true,
			RejectionReason: "unit test failure detected",
		}
		scoreJSON, _ := json.Marshal(score)
		resp := map[string]any{
			"candidates": []map[string]any{
				{
					"content": map[string]any{
						"parts": []map[string]any{
							{"text": string(scoreJSON)},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	client := newTestClient(t, handler)
	tc := &contract.TaskContract{ID: "test", Objective: "fix bug"}
	evidence := []contract.Evidence{
		{Claim: "unit test", Passed: false, Type: contract.VerifyUnitTest},
	}

	score, err := client.ScoreVerification(context.Background(), tc, "diff --git ...", evidence)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.Passed {
		t.Fatal("expected score to fail")
	}
	if !score.DisproofFound {
		t.Fatal("expected disproof to be found")
	}
	if score.RejectionReason == "" {
		t.Fatal("expected rejection reason")
	}
}

func TestClient_APIError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": {"message": "internal server error"}}`))
	})

	client := newTestClient(t, handler)
	_, err := client.GenerateContent(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestClient_MalformedJSON(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not json`))
	})

	client := newTestClient(t, handler)
	_, err := client.GenerateContent(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("expected error for malformed JSON response")
	}
}

func TestClient_EmptyCandidates(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"candidates": []map[string]any{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	client := newTestClient(t, handler)
	_, err := client.GenerateContent(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("expected error for empty candidates")
	}
}

// TestClient_ImplementsJevInterface verifies the compile-time interface assertion.
func TestClient_ImplementsJevInterface(t *testing.T) {
	var c jev.Client = NewClient("key", "model")
	if c == nil {
		t.Fatal("expected non-nil jev.Client implementation")
	}
}
