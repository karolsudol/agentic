# Agentic: High-Throughput Go Agent Architecture

A production-grade AI agent framework implemented in **Go** following **Google ADK 2.0 principles**, featuring a two-tier cognitive model (**System 2 Frontier Planning + System 1 Jev Typed Decisions**), evidence-based loops, and a secure execution harness.

---

## Architecture Overview

This architecture replaces monolithic, prompt-heavy loops with a decoupled three-layer system designed for reliability, type safety, and high-concurrency deployment.

```
┌─────────────────────────────────────────────────────────────┐
│ 1. THE HARNESS (Environment, Security & Memory Floor)       │
│    • Bounded Task Contracts with done_when criteria         │
│    • Progressive Disclosure via Project Map attention index │
│    • Tool Gateway with 4-Tier Permission Ladder             │
│    • Structured Durable Memory outside chat transcripts     │
├─────────────────────────────────────────────────────────────┤
│ 2. THE GRAPH (Topology & Orchestration Engine)              │
│    • Google ADK 2.0 StateGraph state machine                │
│    • Supervisor & Parallel Fleet Workers (Goroutines)       │
│    • Decoupled Builder vs. Adversarial Verifier             │
│    • Human-in-the-Loop (HITL) Checkpoint Interrupts         │
├─────────────────────────────────────────────────────────────┤
│ 3. THE LOOP (Evidence & Cognitive Engine)                   │
│    • System 2 (Frontier LLM): Planning & Code Generation    │
│    • System 1 (Jev / TypeSafe AI): Calibrated Typed Routing │
│    • Deterministic Evidence Gates (Compilers, Tests, Lints) │
│    • Verifiable Run Receipts for human operators            │
└─────────────────────────────────────────────────────────────┘
```

For complete technical specifications, sequence diagrams, and design rationales, see [Architecture Specification](docs/ARCHITECTURE.md).

---

## Documentation Index

Explore the documentation guides for in-depth information:

* **[Architecture Specification](docs/ARCHITECTURE.md):** Detailed component topologies, Mermaid flowcharts, layer descriptions, and comparisons with traditional agent architectures.
* **[Developer Guide](docs/DEVELOPMENT.md):** Local setup, prerequisites, development workflows, and command references.
* **[Evaluation Framework](docs/EVALS.md):** The evaluation blueprint covering contract adherence, trajectory efficiency, Jev calibration metrics, and test datasets.
* **[Repository Rules & Guidelines](AGENTS.md):** Invariants and coding standards governing AI assistant interactions.
* **[Contributing Guide](CONTRIBUTING.md):** PR workflows, conventional commits, and pre-commit hook setup.

---

## Directory Structure

```
.
├── cmd/
│   └── agent/
│       └── main.go              # End-to-end runnable demonstration
├── pkg/
│   ├── contract/                # Task contracts, Evidence, and Receipt models
│   ├── harness/                 # Tool Gateway, Permission Ladder, ProjectMap, Durable Memory
│   ├── jev/                     # Jev client (System 1 calibrated decision engine)
│   ├── graph/                   # ADK 2.0 StateGraph engine & interrupt handling
│   └── nodes/                   # Supervisor, Builder, Verifier, Router, HumanGate nodes
├── docs/                        # In-depth architectural & developer documentation
└── .agents/                     # Local skills and assistant operating guidelines
```

---

## Quickstart

### Prerequisites
* Go 1.27+
* Make
* Git

### Running the Demonstration
Execute the end-to-end multi-agent workflow:
```bash
make run
```

### Running Tests
Run the full test suite with the race detector enabled:
```bash
make test
```

### Formatting & Linting
```bash
make fmt
make vet
make lint
```
