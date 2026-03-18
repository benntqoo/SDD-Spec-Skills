---
name: spec-architect
description: Use when requirements are ambiguous, data models are missing, or acceptance criteria are incomplete and you need executable SDD specs.
---

# Spec Architect

## Overview
Converts fuzzy requirements into executable, auditable, and gate-ready specifications for strict SDD delivery. Produces spec, contract, and traceability artifacts in one pass.

## When to Use

**Use when:**
- Requirements are ambiguous or inconsistent
- Data models or API contracts are missing
- Acceptance criteria are incomplete
- Existing implementation changed without contract updates
- Starting any new feature in SDD workflow

**When NOT to use:**
- Clear requirements with existing contracts (use spec-to-codebase directly)
- Maintenance work on existing specs

## Invocation Alignment

- If `sdd-orchestrator` is present, it is the only state-transition entry
- Direct invocation is limited to read-only analysis or refinement planning
- Direct invocation must not promote state; control returns to `sdd-orchestrator`

## State Transition

| Input State | Success | Failure |
|-------------|---------|---------|
| Ideation/Explore | SpecCheckpoint | Explore |

## Required Outputs

- `.sdd-spec/specs/<feature>.md` - Full specification
- `.sdd-spec/specs/<feature>.contract.json` - API contract
- `.sdd-spec/specs/<feature>.traceability.yaml` - Requirement mapping
- `.sdd-spec/specs/<feature>.risk.yaml` - Risk assessment
- `.sdd-spec/specs/<feature>.spec.report.json` - Execution report

## Quick Reference

| Check | Required |
|-------|----------|
| Story → Acceptance | 1:1 mapping |
| Acceptance → Contract | 1:1+ mapping |
| Contract Operation | Error scenarios defined |
| Compatibility | Explicit and testable |
| TBD/TODO | None allowed |

## Gate Checks

All must pass before promoting to SpecCheckpoint:
1. Every user story maps to at least one acceptance criterion
2. Every acceptance criterion maps to at least one contract operation
3. Every contract operation defines at least one error scenario
4. Compatibility commitments are explicit and testable
5. No unresolved TBD or TODO

## Hard Rules

**Never:**
- Embed implementation code in spec files
- Leave unnamed acceptance criteria
- Permit contract fields without types
- Mark breaking changes as non-breaking

## Spec Content Requirements

**spec.md must include:**
1. Objective and scope boundaries
2. User stories and acceptance criteria with unique IDs
3. Domain data models
4. API or interface contracts
5. Error taxonomy and retry/idempotency policy
6. Business invariants and non-functional constraints
7. Backward compatibility commitments
8. Rollback and feature-flag strategy

**contract.json must include:**
- Contract version
- Operations with input/output schemas
- Error codes and messages
- Compatibility mode (backward/forward/strict)
- Breaking change indicator

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Leaving TBD/TODO in contracts | Remove or explicitly document as TODO items requiring separate tracking |
| Unnamed acceptance criteria | Assign unique IDs to all acceptance criteria |
| Missing error scenarios | Define at least one error scenario per contract operation |
| Vague compatibility claims | Specify exact mode (backward/forward/strict) with rationale |

## Machine Contracts

Report structure must conform to `skills/sdd-orchestrator/sdd-machine-schema.json`.
Gate checklist defined in `skills/sdd-orchestrator/sdd-gate-checklist.json`.