---
name: eval-runner
description: Run agent evaluations, analyze trajectory traces, and verify contract pass rates.
---

# Agent Evaluation Runner Skill

Use this skill when developing, testing, or benchmarking agent workflows against golden evaluation datasets.

## Workflow

1. **Verify Test Environment:**
   Run `make test` to ensure base unit and integration tests are clean.

2. **Run Evaluation Benchmarks:**
   Execute the evaluation harness:
   ```bash
   make eval
   ```

3. **Analyze Run Receipts:**
   Inspect the generated run receipt:
   * Verify all `done_when` criteria are marked `PASS`.
   * Check `ChangedFiles` to ensure only expected files were modified.
   * Verify no negative constraints were violated.

4. **Review Trajectory Traces:**
   Inspect the execution trace for:
   * Total iteration count (should be <= max allowed budget).
   * Tool precondition failures (should be 0).
   * Jev System 1 routing decisions (appropriate action with >0.90 confidence).
