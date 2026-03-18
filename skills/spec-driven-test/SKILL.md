---
name: spec-driven-test
description: Use when implementation exists at Build or Verify state and you need to create verification tests from spec contracts.
---

# Spec-Driven Test

## Overview
Converts spec artifacts into mandatory verification suites. Blocks promotion when any contract path lacks executable proof.

## When to Use

**Use when:**
- Current state is Build or Verify
- Contract.json exists
- Traceability.yaml exists
- Contract diff check has passed

**When NOT to use:**
- Contract diff not yet run (run spec-contract-diff first)
- No contract file (return to spec-architect)

## Invocation Alignment

- If `sdd-orchestrator` is present, it is the only state-transition entry
- Direct invocation is limited to read-only test planning or gap analysis
- Direct invocation must not promote state; control returns to `sdd-orchestrator`

## State Transition

| Input State | Success | Failure |
|-------------|---------|---------|
| Build/Verify | Verify | Build |

## Entry Conditions (All Must Be True)

1. Current state is `Build` or `Verify`
2. `.sdd-spec/specs/<feature>.contract.json` exists
3. `.sdd-spec/specs/<feature>.traceability.yaml` exists
4. Contract diff check has passed

## Quick Reference

| Gate Check | Coverage Required |
|------------|-------------------|
| Traceability | 100% for stories and acceptance |
| Contract Operations | 100% coverage |
| Error Codes | 100% coverage |
| Flaky Markers | None allowed |
| Test Result | Must be pass |

## Test Construction Rules

- One acceptance criterion → at least one executable test
- One contract operation → success test AND failure test
- Each error code → at least one assertion
- Tests verify runtime behavior through public interfaces only
- External dependencies may be isolated, business logic cannot be replaced

## Required Outputs

- `.sdd-spec/tests/specs/<feature>.contract.spec.*`
- `.sdd-spec/tests/specs/<feature>.acceptance.spec.*`
- `.sdd-spec/specs/<feature>.test.report.json`

## Failure Policy

**On failure:**
- Emit test.report.json with failed IDs
- Keep state at Build
- Block release guard invocation

## Handoff Contract

**When successful:**
- Emit target state `Verify`
- Publish test report path
- Publish coverage summary
- Next skill: sdd-release-guard

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Missing failure tests | Create both success and failure test cases |
| Only testing happy path | Cover all error scenarios defined in contract |
| Isolating business logic | Test through public interfaces, don't mock business logic |
| Not covering error codes | 100% error code coverage required |

## Machine Contracts

Report structure must conform to `skills/sdd-orchestrator/sdd-machine-schema.json`.
Gate checklist defined in `skills/sdd-orchestrator/sdd-gate-checklist.json`.