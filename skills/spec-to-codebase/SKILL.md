---
name: spec-to-codebase
description: Use when you have validated spec and contract artifacts at SpecCheckpoint state and need to generate implementation aligned with contracts.
---

# Spec-to-Codebase

## Overview
Builds deterministic code changes from spec artifacts and enforces contract-preserving generation. Transforms validated specs into working implementation skeletons.

## When to Use

**Use when:**
- Current state is SpecCheckpoint
- Contract.json exists with explicit breaking_change policy
- Traceability file exists with story IDs
- Need to generate implementation from validated specs

**When NOT to use:**
- Spec contracts are incomplete (return to spec-architect)
- Code already diverged from spec (return to spec-architect)

## Invocation Alignment

- If `sdd-orchestrator` is present, it is the only state-transition entry
- Direct invocation is limited to read-only analysis or generation planning
- Direct invocation must not promote state; control returns to `sdd-orchestrator`

## State Transition

| Input State | Success | Failure |
|-------------|---------|---------|
| SpecCheckpoint | Build | SpecCheckpoint |

## Entry Conditions (All Must Be True)

1. Current state is `SpecCheckpoint`
2. `.sdd-spec/specs/<feature>.contract.json` exists
3. `breaking_change` policy in contract is explicit
4. Traceability file exists and contains story IDs

## Quick Reference

| Gate Check | Required |
|------------|----------|
| Contract → Code Mapping | All operations mapped |
| Required Fields | All input/output present |
| Unresolved Items | None allowed |
| Compatibility | Result must be pass |

## Generation Process

1. Read spec as source of truth:
   - `.sdd-spec/specs/<feature>.md`
   - `.sdd-spec/specs/<feature>.contract.json`
   - `.sdd-spec/specs/<feature>.traceability.yaml`

2. Map contracts to existing project structure first
3. Extend existing modules before creating new files
4. Generate minimal compile-ready skeletons
5. Preserve backward-compatible signatures unless breaking changes declared

## Non-Negotiable Constraints

**Never:**
- Rewrite unrelated files
- Invent fields not in contract
- Remove fields without explicit breaking-change declaration
- Generate mock or fake business data in production path
- Bypass existing project lint/type rules

## Required Outputs

- Updated implementation files aligned with contract operations
- `.sdd-spec/specs/<feature>.codegen.report.json`

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Generating too much code | Only generate skeletons, not full implementations |
| Rewriting existing working code | Extend existing modules, don't replace |
| Adding fields not in contract | Stick strictly to contract schema |
| Ignoring breaking_change flag | Respect compatibility mode in contract |

## Machine Contracts

Report structure must conform to `skills/sdd-orchestrator/sdd-machine-schema.json`.
Gate checklist defined in `skills/sdd-orchestrator/sdd-gate-checklist.json`.