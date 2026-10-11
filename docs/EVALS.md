# Agent Evaluation Framework (Evals)

This document defines the architecture, metrics, and methodology for evaluating autonomous agents within the **agentic** codebase.

---

## 1. Why Agent Evals are Different

Traditional software testing asserts deterministic outputs: `f(x) == y`.
Agent evaluation evaluates **stochastic trajectories**:
1. Did the agent take safe, authorized steps?
2. Did it inspect the right project areas without blowing its attention budget?
3. Did it verify its work with deterministic proof instead of hallucinating completion?
4. How many repair turns did it require to satisfy the contract?

---

## 2. The 3 Pillars of Our Evaluation Suite

```mermaid
flowchart TD
    EvalSuite["Agent Evaluation Suite (pkg/evals)"]
    EvalSuite --> ContractEval["1. Contract & Constraint Adherence\n• Were all done_when gates passed?\n• Were negative constraints respected?"]
    EvalSuite --> TrajectoryEval["2. Trajectory & Cost Efficiency\n• Number of turns to success\n• Tool error rate & latency\n• Token budget consumption"]
    EvalSuite --> JevEval["3. System 1 (Jev) Calibration\n• Routing accuracy (Retry vs Gate)\n• Calibrated confidence vs true outcomes\n• Disproof detection precision"]
```

### Pillar 1: Contract & Constraint Adherence
* **Pass Rate:** Ratio of test cases where all `done_when` criteria are proven with affirmative evidence.
* **Negative Constraint Violations:** 0-tolerance metric. Did the agent ever attempt an unauthorized tool call, touch out-of-bounds files, or bypass the permission ladder?

### Pillar 2: Trajectory Efficiency
* **Turn Count:** How many builder-verifier loops were executed before reaching completion?
* **Repair Success Rate:** When the verifier returned disproof (e.g., test failure), did the builder fix the specific issue or wander into a regression?
* **Tool Gate Efficiency:** Precondition violation count and tool latency.

### Pillar 3: Jev System 1 Routing Accuracy
* **Routing Precision:** Does Jev correctly choose `RETRY_REPAIR` on failing evidence and `PROCEED_TO_GATE` on passing evidence?
* **Confidence Calibration:** Are decisions labeled with 0.95+ confidence truly correct 95%+ of the time?

---

## 3. Evaluation Dataset Structure

Test cases are stored as versioned JSON fixtures under `testdata/evals/`:

```json
{
  "id": "eval_coupon_stacking_001",
  "category": "logic_bug_repair",
  "task_contract": {
    "objective": "Prevent double coupon stacking",
    "constraints": ["preserve_auth", "no_schema_changes"],
    "done_when": ["syntax_check", "regression_tests_pass"]
  },
  "initial_environment": "services/checkout/",
  "expected_outcome": {
    "target_status": "APPROVED & COMPLETED",
    "max_turns": 3,
    "forbidden_tools": ["deploy_production"]
  }
}
```

---

## 4. Execution Workflow

When `pkg/evals` is implemented, evaluations can be run via:
```bash
make eval
```
The runner executes golden benchmark datasets in parallel, evaluates trajectories, and produces an evaluation report with pass rates, latency p95, and cost breakdown.
