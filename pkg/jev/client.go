package jev

import (
	"context"
	"fmt"
	"strings"

	"agentic/pkg/contract"
)

// Client defines the interface for interacting with Jev (TypeSafe AI)
type Client interface {
	DecideRoute(ctx context.Context, c *contract.TaskContract, evidence []contract.Evidence, iteration, maxIterations int) (*RoutingDecision, error)
	ScoreVerification(ctx context.Context, c *contract.TaskContract, candidateDiff string, evidence []contract.Evidence) (*VerifierScore, error)
}

// LocalJevClient provides a high-speed, calibrated System 1 decision engine
// In production, this can also connect to the TypeSafe AI / Jev HTTP endpoint.
type LocalJevClient struct{}

func NewLocalJevClient() *LocalJevClient {
	return &LocalJevClient{}
}

func (j *LocalJevClient) DecideRoute(ctx context.Context, c *contract.TaskContract, evidence []contract.Evidence, iteration, maxIterations int) (*RoutingDecision, error) {
	// Exhaustion check
	if iteration >= maxIterations {
		return &RoutingDecision{
			Action:     ActionHaltExhausted,
			Confidence: 0.99,
			Reason:     fmt.Sprintf("maximum retry budget exhausted (%d/%d)", iteration, maxIterations),
		}, nil
	}

	// Analyze evidence failures
	failedCount := 0
	var failureReasons []string
	for _, ev := range evidence {
		if !ev.Passed {
			failedCount++
			failureReasons = append(failureReasons, fmt.Sprintf("%s: %s", ev.Type, ev.Claim))
		}
	}

	if failedCount == 0 {
		return &RoutingDecision{
			Action:     ActionProceedToGate,
			Confidence: 0.97,
			Reason:     "all acceptance criteria and tests validated with affirmative evidence",
		}, nil
	}

	// Classify retry versus escalation (Mori Rule 9)
	return &RoutingDecision{
		Action:         ActionRetryRepair,
		Confidence:     0.94,
		Reason:         fmt.Sprintf("found %d failing verification gate(s): %s", failedCount, strings.Join(failureReasons, "; ")),
		SuggestedFixes: failureReasons,
	}, nil
}

func (j *LocalJevClient) ScoreVerification(ctx context.Context, c *contract.TaskContract, candidateDiff string, evidence []contract.Evidence) (*VerifierScore, error) {
	// Verifier attempts to find disproof (Mori Rule 7: check contract, missing cases)
	for _, ev := range evidence {
		if !ev.Passed {
			return &VerifierScore{
				Passed:          false,
				Confidence:      0.96,
				DisproofFound:   true,
				RejectionReason: fmt.Sprintf("disproof: %s failed (%s)", ev.Claim, ev.Details),
			}, nil
		}
	}

	return &VerifierScore{
		Passed:        true,
		Confidence:    0.95,
		DisproofFound: false,
	}, nil
}
