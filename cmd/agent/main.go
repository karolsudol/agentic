package main

import (
	"context"
	"fmt"
	"log"

	"agentic/pkg/config"
	"agentic/pkg/contract"
	"agentic/pkg/gemini"
	"agentic/pkg/graph"
	"agentic/pkg/harness"
	"agentic/pkg/jev"
	"agentic/pkg/nodes"
)

func main() {
	ctx := context.Background()

	fmt.Println("=================================================================")
	fmt.Println("  AGENT ARCHITECTURE: Go ADK 2.0 + Jev (TypeSafe AI)")
	fmt.Println("  Harness Gateway + StateGraph + Evidence-Based Loops")
	fmt.Println("=================================================================")

	// 0. CONFIGURATION: Load env-driven settings (.env → environment)
	cfg := config.Load(".env")

	// 1. HARNESS LAYER (Mori Rule 2): Define Bounded Task Contract
	taskContract := &contract.TaskContract{
		ID:        "fix_double_coupon_001",
		Objective: "Fix duplicate coupon stacking in checkout checkout flow",
		Inputs: []string{
			"checkout service code",
			"cart validation logic",
			"coupon integration test suite",
		},
		Constraints: []string{
			"Preserve existing auth context",
			"Do not alter payment gateway schema",
			"Do not write outside workspace",
		},
		Deliverable: "Tested pull request preventing stacked discounts",
		DoneWhen: []string{
			"Go syntax and type-check passes",
			"Double coupon regression tests pass (100% exit code 0)",
			"No new warnings from linter",
		},
		ApprovalRequired: []string{
			"deploy_staging",
			"git_push_main",
		},
	}

	// 2. HARNESS LAYER (Mori Rule 3): Progressive Disclosure Project Map
	projectMap := harness.NewProjectMap()
	projectMap.RegisterArea("docs", "docs/", "Architecture guidelines and constraints", []string{"architecture.md"})
	projectMap.RegisterArea("backend", "services/checkout/", "Checkout microservice", []string{"coupon.go", "cart.go"})
	projectMap.RegisterArea("tests", "tests/checkout/", "Unit and integration regression suites", []string{"coupon_test.go"})

	// 3. HARNESS LAYER (Mori Rule 4 & 8): Tool Gateway & Policy Engine
	policyEngine := harness.NewPolicyEngine(true) // Auto-approve reversible edits within workspace
	gateway := harness.NewToolGateway(policyEngine)
	gateway.Register(&harness.WorkspaceEditTool{})
	gateway.Register(&harness.RunTestsTool{})
	gateway.Register(&harness.DeployStagingTool{})

	// 4. SYSTEM 1 DECISION ENGINE: Jev (TypeSafe AI)
	//    Dual-mode: local (deterministic, zero API cost) or gemini (structured outputs)
	var jevClient jev.Client
	if cfg.IsGeminiMode() && cfg.HasGeminiKey() {
		fmt.Println("[CONFIG]: Jev mode = GEMINI (structured outputs via Gemini API)")
		jevClient = gemini.NewClient(cfg.GeminiAPIKey, cfg.GeminiModel)
	} else {
		fmt.Println("[CONFIG]: Jev mode = LOCAL (deterministic decision engine, zero API cost)")
		jevClient = jev.NewLocalJevClient()
	}

	// 5. GRAPH ENGINE: ADK 2.0 StateGraph Assembly
	builderNode := nodes.NewBuilderNode(gateway)
	verifierNode := nodes.NewVerifierNode(gateway, jevClient)
	supervisorNode := nodes.NewSupervisorNode(gateway)
	humanGateNode := nodes.NewHumanGateNode(policyEngine)
	receiptNode := nodes.NewReceiptNode()
	jevRouter := nodes.NewJevEdgeRouter(jevClient)

	workflow := graph.NewStateGraph()
	workflow.AddNode("supervisor", supervisorNode.Execute)
	workflow.AddNode("builder", builderNode.Execute)
	workflow.AddNode("verifier", verifierNode.Execute)
	workflow.AddNode("human_gate", humanGateNode.Execute)
	workflow.AddNode("receipt_node", receiptNode.Execute)

	// Topology & Transitions (Beam Graph layer)
	workflow.SetEntryPoint("supervisor")
	workflow.AddEdge("supervisor", "builder")
	workflow.AddEdge("builder", "verifier")

	// Dynamic Loop Routing via Jev (Beam Loop layer & Mori Rule 9)
	workflow.AddConditionalEdge("verifier", jevRouter.Route)

	// Post-gate static transition
	workflow.AddEdge("human_gate", "receipt_node")
	workflow.AddEdge("receipt_node", graph.EndNode)

	// Initialize Graph State with durable memory (Mori Rule 5)
	initialState := graph.NewInitialState(taskContract, projectMap, 4)

	// Execute ADK 2.0 Graph Workflow
	finalState, err := workflow.Execute(ctx, initialState)
	if err != nil {
		log.Fatalf("Graph execution failed: %v", err)
	}

	// Print Durable Memory (Durable state outside the conversation transcript)
	fmt.Println("\n[DURABLE STATE (Stored Outside Chat Transcript)]:")
	for _, dec := range finalState.Memory.Decisions {
		fmt.Printf("  • [%s]: %s (Rationale: %s)\n", dec.ID, dec.Choice, dec.Rationale)
	}

	// Print Gateway Audit Traces (Mori Rule 11)
	fmt.Println("\n[HARNESS TOOL TRACES]:")
	for i, tr := range gateway.Traces() {
		fmt.Printf("  %d. %s [%s] -> %s (latency: %v)\n",
			i+1, tr.ToolName, tr.Risk, tr.Output["status"], tr.Duration)
	}

	fmt.Println("\nExecution completed successfully!")
}
