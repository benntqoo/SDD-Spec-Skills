---
name: spec-traceability
description: Use when you need to verify complete linkage between requirements, contracts, code, and tests, or when traceability gaps block release.
---

# Spec Traceability

## Overview
**Verification-only skill** - does NOT change state. Validates traceability completeness and blocks progression if gates fail. Maintains complete and verifiable linkage across requirements, contracts, implementation, and tests.

## When to Use

**Use when:**
- User stories or acceptance criteria change
- Contract operations or schemas change
- Public code interfaces change
- Contract or acceptance tests change
- Release requires traceability verification

**When NOT to use:**
- Only need to verify single mapping (acceptable but not efficient)

**Note:** This is a verification-only skill. It validates completeness but does not promote state. Always returns control to sdd-orchestrator.

## Invocation Alignment

- If `sdd-orchestrator` is present, it is the only state-transition entry
- Direct invocation is limited to read-only traceability analysis
- Direct invocation must not promote state; control returns to `sdd-orchestrator`

## Entry Conditions

Run when any of these change:
- User stories or acceptance criteria
- Contract operations or schemas
- Public code interfaces
- Contract or acceptance tests

## Quick Reference

| Consistency Rule | Check |
|-----------------|-------|
| No orphan story | Has acceptance criteria |
| No orphan acceptance | Has contract mapping |
| No orphan contract operation | Has code entry point |
| No orphan code entry point | Has at least one test |

**Status values:** draft | implemented | verified

## Gate Rules

Matrix passes only when:
- Mapping completeness is 100%
- No duplicate IDs across rows
- All verified rows have passing test case references

## Required Outputs

- `.sdd-spec/specs/<feature>.traceability.yaml`
- `.sdd-spec/specs/<feature>.traceability.report.json`

## Report Requirements

**traceability.report.json must include:**
- feature, state_before, state_after, skill, timestamp
- result, blocking_reasons
- completeness (0-100%)
- orphan_items list

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Incomplete mapping | Ensure 100% completeness before release |
| Orphan items in matrix | Identify and add missing links |
| Verifying without sdd-orchestrator | Always route through orchestrator |
| Blocking without clear reasons | Document specific orphan items |

## Machine Contracts

Report structure must conform to `skills/sdd-orchestrator/sdd-machine-schema.json`.
Gate checklist defined in `skills/sdd-orchestrator/sdd-gate-checklist.json`.