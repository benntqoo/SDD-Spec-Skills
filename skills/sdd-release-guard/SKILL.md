---
name: sdd-release-guard
description: Use when feature is at ReleaseReady state and you need final release gates and rollback readiness verification before production deployment.
---

# SDD Release Guard

## Overview
Applies final release governance and decides whether a feature can move from ReleaseReady to Released. Validates all gates pass and rollback is ready.

## When to Use

**Use when:**
- Current state is ReleaseReady
- Contract diff report result is pass
- Test report result is pass
- Traceability completeness is 100%

**When NOT to use:**
- Still in Verify state (need to pass tests first)
- Contract or test gates not yet passed

## Invocation Alignment

- If `sdd-orchestrator` is present, it is the only state-transition entry
- Direct invocation is limited to read-only release-readiness analysis
- Direct invocation must not promote state; control returns to `sdd-orchestrator`

## Entry Conditions (All Must Be True)

1. Current state is `ReleaseReady`
2. Contract diff report result is `pass`
3. Test report result is `pass`
4. Traceability completeness is 100%

## Quick Reference

| Gate | Check |
|------|-------|
| Contract Diff | Pass |
| Test Report | Pass |
| Traceability | 100% complete |
| Feature Flag | Plan present |
| Observability | Checks present |
| Rollback | Plan and threshold present |

## Gate Set

Release decision is `pass` only if all checks pass:
1. Compatibility check
2. Contract test gate
3. Acceptance test gate
4. Traceability gate
5. Feature flag availability
6. Rollback command and threshold completeness

## Required Outputs

- `.sdd-spec/specs/<feature>.release.guard.json`

**Required fields:**
- feature, state_before, state_after, skill, timestamp
- result, blocking_reasons
- gate_results, feature_flag_plan
- observability_checks, rollback_plan
- release_decision

## Failure Policy

**On any gate failure:**
- Set `release_decision` to `fail`
- Keep state at ReleaseReady
- Record blocked gates with explicit reason

## Success Policy

**On all gates pass:**
- Set `release_decision` to `pass`
- Promote state to Released
- Record release timestamp and guard version

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Skipping gates | Never skip any gate check |
| Approving without rollback plan | Require complete rollback plan |
| Missing feature flag plan | Require feature flag strategy |
| Not checking observability | Require observability checks |
| Premature promotion | Stay at ReleaseReady until all pass |

## Machine Contracts

Report structure must conform to `skills/sdd-orchestrator/sdd-machine-schema.json`.
Gate checklist defined in `skills/sdd-orchestrator/sdd-gate-checklist.json`.