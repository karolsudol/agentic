---
name: graph-debugger
description: Debug and trace ADK 2.0 StateGraph transitions, node handlers, and HITL interrupts.
---

# ADK 2.0 StateGraph Debugger Skill

Use this skill when investigating state transitions, conditional routing edge errors, or stuck graph loops.

## Debugging Checklist

1. **Inspect Graph Topology:**
   * Verify all nodes have registered handlers in `StateGraph.AddNode`.
   * Ensure `SetEntryPoint` is explicitly defined.
   * Check that every branch terminates at `__END__` or cycles back intentionally.

2. **Check State Immutability & Cloning:**
   * Ensure nodes mutate state predictably and return valid pointer references.
   * Check that `DurableMemory` correctly records decisions at each step.

3. **Verify Jev Routing Edge Logic:**
   * Inspect the `Action` returned by Jev in the conditional edge router.
   * If an unexpected loop occurs, check whether the verifier is correctly clearing or appending evidence.

4. **Verify HITL Interrupts:**
   * Confirm that nodes registered with `AddInterrupt` pause execution when `HumanApproved == false`.
