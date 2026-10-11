# Contributing to Agentic

Thank you for contributing to the **agentic** project! This guide explains our workflow, conventions, and quality gates.

---

## 1. Development Principles

All contributions must align with the core architectural invariants:
* **Harness First:** Protect tool access with the `ToolGateway` and enforce the 4-tier permission ladder.
* **Deterministic Evidence:** Tests and compilers must verify state before any task completes.
* **Two-Tier Cognition:** Use System 2 (Frontier models) for generation; System 1 (Jev) for calibrated routing decisions.
* **Durable State:** Record choices and lessons in `DurableMemory`, not in raw LLM transcripts.

---

## 2. Local Setup

1. **Clone & Branch:**
   ```bash
   git checkout -b feat/your-feature-name
   ```
2. **Install Hooks & Verify Environment:**
   ```bash
   make setup
   ```
3. **Run Checks Locally:**
   ```bash
   make fmt
   make vet
   make test
   ```

---

## 3. Pull Request Guidelines

* Keep PRs focused on a single responsibility or feature.
* Include unit tests for all new nodes, tools, and graph transitions.
* Update documentation (`README.md`, `docs/`) if architectural changes or new nodes are added.
* Ensure all pre-commit hooks and CI checks pass before requesting review.

---

## 4. Conventional Commits

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:
* `feat:` A new feature or graph node
* `fix:` A bug fix or gate correction
* `docs:` Documentation updates
* `refactor:` Code changes that neither fix a bug nor add a feature
* `test:` Adding or updating tests and evals
