---
name: sdd-orchestrator
description: Use when coordinating feature delivery with strict SDD state transitions, gate validation, and requirement for single entry point control.
---

# SDD Orchestrator

## Overview
Single controller for strict Spec-Driven Development (SDD) execution. Enforces one-way, auditable, and recoverable feature progression through mandatory state transitions and gate checks.

## When to Use

**Use when:**
- Starting any new feature delivery workflow
- Managing feature state across SDD lifecycle
- Enforcing gate validation before state promotion
- Requiring auditable state transition history

**When NOT to use:**
- Ad-hoc feature implementation without SDD workflow
- Simple documentation updates (use fast-path mode)

## Invocation Alignment

- Always invoke this skill first for any feature change
- Other SDD skills must run only when directed by this skill
- Direct invocation of downstream skills without matching state is invalid

## State Flow

```
Ideation → Explore → SpecCheckpoint → Build → Verify → ReleaseReady → Released
```

## Routing Rules

| Current State | Call Skill | Purpose |
|---------------|------------|---------|
| Ideation/Explore | `spec-architect` | Create spec and contracts |
| SpecCheckpoint | `spec-to-codebase` | Generate implementation |
| Build | `spec-contract-diff` | Detect contract drift |
| Build/Verify | `spec-driven-test` | Run verification tests |
| Verify | `sdd-release-guard` | Final release gates |

**Note:** Call `spec-traceability` after any spec/code/test changes (verification only, no state change).

## Canonical Enums

- **State:** Ideation | Explore | SpecCheckpoint | Build | Verify | ReleaseReady | Released
- **Result:** pass | fail | blocked
- **Compatibility Mode:** backward | forward | strict

## Quick Reference

| Action | Command/Method |
|--------|----------------|
| Validate workflow | `python skills/sdd-orchestrator/validate-sdd.py` |
| Track state | `.sdd-spec/specs/<feature>.state.json` |
| Schema | `skills/sdd-orchestrator/sdd-machine-schema.json` |
| Gate checklist | `skills/sdd-orchestrator/sdd-gate-checklist.json` |

## Gate Governance

**Never:**
- Skip failed gates
- Downgrade compatibility claims silently
- Promote state without output artifacts

**Always:**
- Persist block reasons in state record

## Recovery Rules

- If skill fails → remain in current valid state
- If contracts change → force return to Explore
- If tests fail → set state to Build with failed IDs

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Direct invocation of downstream skills without orchestrator | Always invoke through sdd-orchestrator |
| Skipping failed gates to proceed | Never skip; fix root cause first |
| Promoting state without artifacts | Ensure all required_outputs exist before promotion |
| Silently changing compatibility mode | Document any changes in contract |

## Machine Contracts

- Schema: `skills/sdd-orchestrator/sdd-machine-schema.json`
- Checklist: `skills/sdd-orchestrator/sdd-gate-checklist.json`
- Validation: `python skills/sdd-orchestrator/validate-sdd.py`