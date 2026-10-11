# AGENTS.md — Repository Rules & Operating Guidelines

Welcome to the **agentic** codebase. This repository implements a high-throughput, production-grade AI agent architecture in **Go** following **Google ADK 2.0 principles**, **Jev (TypeSafe AI)**, and **Harness Engineering**.

All AI coding assistants (Antigravity, Gemini, Copilot, Claude) working in this repository MUST strictly abide by the following engineering invariants.

---

## 1. Architectural Invariants

1. **The Three Layers (Harness, Loop, Graph):**
   * **Harness:** Never invoke tools or mutate filesystem/network directly from an LLM node without passing through `pkg/harness/ToolGateway`.
   * **Graph:** All workflows must be represented as nodes and edges in `pkg/graph/StateGraph`. Never implement uncontrolled recursive prompt loops.
   * **Loop:** Stopping conditions MUST be evidence-driven (`done_when` checks via compilers, tests, linters). Never stop on model self-confidence.

2. **Two-Tier Cognitive Separation:**
   * **System 2 (Frontier models):** Used for code generation, architectural synthesis, and complex planning.
   * **System 1 (Jev):** Used for calibrated, typed routing decisions, parameter guardrails, and verification scoring.
   * **Deterministic Code:** Used for syntax checks, type checks, unit tests, and policy authorization.

3. **Memory Outside the Transcript:**
   * Long-term state and architectural choices must be recorded in `pkg/harness/DurableMemory` (`Facts`, `Decisions`, `State`, `Lessons`). Do not rely on LLM context history to retain critical choices.

4. **Independent Builder vs. Verifier:**
   * The builder agent must NEVER verify its own output without an independent verifier node attempting adversarial disproof.

---

## 2. Go Coding Standards

* **Language Version:** Go 1.27+
* **Formatting:** All Go code must be formatted with `gofmt` and pass `go vet`.
* **Error Handling:** Always return explicit, wrapped errors (`fmt.Errorf("...: %w", err)`). Never discard errors with `_`.
* **Context Propagation:** All I/O, node handlers, and tool executions MUST accept and propagate `context.Context` as their first parameter.
* **Concurrency:** Use goroutines and channels safely; always ensure termination via `context.Done()` or `sync.WaitGroup`.
* **Strong Typing:** Never use untyped strings for routing actions; use typed enums defined in `pkg/jev/types.go` and `pkg/contract/contract.go`.

---

## 3. Security & Safety

* **Permission Ladder:** Respect the 4 risk tiers (`LowRisk`, `Reversible`, `ExternalEffect`, `Irreversible`). Sensitive or external operations MUST require human approval gates.
* **Secrets:** Never hardcode API keys, tokens, or credentials in any file. Use environment variables with `.env.example` templates.
* **Workspace Boundaries:** Tools must validate paths to ensure operations remain bounded within the workspace.

---

## 4. Verification & Commit Requirements

* Before declaring any task finished:
  1. Run `make test` or `go test ./...` and ensure zero failures.
  2. Run `make lint` or `go vet ./...` and ensure zero warnings.
  3. Generate or update the **Run Receipt** detailing changed files and verified evidence.
* Provide clean, conventional git commit commands (`feat:`, `fix:`, `refactor:`, `test:`, `docs:`) for the user to execute.
