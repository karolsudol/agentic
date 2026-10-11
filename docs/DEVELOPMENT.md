# Developer Environment Guide

This document describes the recommended developer environment, toolchain requirements, and workflows for developing and testing agents locally.

---

## 1. Prerequisites

* **Go:** Version 1.27 or higher (`go version`)
* **Make:** For standardized command automation (`make --version`)
* **Git:** Version 2.30 or higher
* **Optional / Recommended:**
  * `golangci-lint` for fast multi-linter execution (`brew install golangci-lint` on macOS)
  * `pre-commit` for automated git hook validation (`brew install pre-commit` or `pip install pre-commit`)

---

## 2. Directory Layout & Conventions

```
.
├── cmd/               # Application binaries (e.g., cmd/agent/main.go)
├── pkg/
│   ├── contract/      # Bounded task contracts, DoneWhen criteria, and Receipts
│   ├── harness/       # The Harness (Gateway, permissions, state, project map)
│   ├── jev/           # System 1 calibrated decision engine client
│   ├── graph/         # ADK 2.0 StateGraph engine
│   ├── nodes/         # Graph nodes (Supervisor, Builder, Verifier, Router, HumanGate)
│   └── evals/         # Evaluation harness & benchmark suites (upcoming)
├── docs/              # In-depth architectural & developer documentation
└── .agents/           # Custom local skills and rules for AI assistants
```

---

## 3. Daily Development Workflow

### Running the Workflow
To execute the baseline agent demonstration:
```bash
make run
```

### Running Unit Tests
All packages include unit tests with race detection:
```bash
make test
```

### Formatting & Linting
Ensure code matches repository formatting standards:
```bash
make fmt
make vet
make lint
```

### Cleaning Up
To remove compiled binaries and test profile outputs:
```bash
make clean
```
