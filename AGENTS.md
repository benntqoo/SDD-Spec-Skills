# Agent Collaboration Guide

**Updated for Vibe-SDD Development Workflow**

## Current System Status

✅ **Multi-Agent Supported**: Yes, via Git branching workflow
✅ **Multi-User Supported**: Yes, via Git collaboration
✅ **Structured Development**: Yes, via .vic-sdd/ SPEC workflow
❌ **Real-Time Collaboration**: No, requires Git merge workflow

## Development Workflow

### Phase 1: 定图纸 (Requirements)

```
Agent-Product uses vibe-think
    ↓
SPEC-REQUIREMENTS.md
    ↓
vic spec gate 0 (Gate: Requirements Completeness)
```

### Phase 2: 打地基 (Architecture)

```
Agent-Architect uses vibe-architect
    ↓
SPEC-ARCHITECTURE.md
    ↓
vic spec gate 1 (Gate: Architecture Completeness)
```

### Phase 3: 立规矩 (Implementation)

```
Agent-Develop uses vibe-develop
    ↓
Implementation + Tests
    ↓
vic spec gate 2 (Code Alignment)
vic spec gate 3 (Test Coverage)
    ↓
vic spec merge → PRD.md / ARCH.md / PROJECT.md
```

## Multi-Agent Scenarios

### Scenario 1: Sequential Agents (Recommended)
```
Agent A (design): Completes work → Pushes branch → Creates PR
Agent B (review): Reviews PR → Merges → Continues work
```

### Scenario 2: Parallel Agents (Use Caution)
```
Agent A: Working on branch feature/auth
Agent B: Working on branch feature/database
Both: Use separate branches, merge independently
```

### Scenario 3: Same Branch (Avoid if Possible)
```
⚠️ Risk: Merge conflicts in .vic-sdd/ files
⚠️ Solution: Coordinate via PR reviews, use union merge
```

## Conflict Resolution Workflow

When multiple agents modify the same YAML files:

1. **Git detects conflict** during merge/pull request
2. **Union merge** preserves both versions (may create duplicates)
3. **Run validation script** to detect duplicate IDs
4. **Manual resolution** required to merge similar decisions
5. **Verify** application still works with merged memory

## Best Practices

1. **Use Separate Branches**: Each agent gets own branch
2. **Document Decisions**: Use .vic-sdd/tech/tech-records.yaml for major choices
3. **Run Gate Checks**: `vic spec gate 0-3` before progressing
4. **PR Reviews**: Review .vic-sdd/ changes before merging
5. **Validate After Updates**: Run validation frequently

## Directory Structure

```
.vic-sdd/
├── SPEC-REQUIREMENTS.md    # Requirements spec
├── SPEC-ARCHITECTURE.md    # Architecture spec
├── PROJECT.md              # Project status tracking
├── status/
│   ├── events.yaml         # Event history
│   └── state.yaml          # Current state
├── tech/
│   └── tech-records.yaml  # Technical decisions
├── risk-zones.yaml         # Risk records
├── project.yaml            # AI quick reference
└── dependency-graph.yaml  # Module dependencies
```

## Quick Commands

```bash
# Initialize
vic init
vic spec init

# SPEC Management
vic spec status
vic spec gate 0  # Requirements
vic spec gate 1  # Architecture
vic spec gate 2  # Code alignment
vic spec gate 3  # Test coverage
vic spec merge   # Merge to final docs

# Recording
vic rt --id DB-001 --title "Use PostgreSQL" --decision "Primary DB"
vic rr --id RISK-001 --area auth --desc "JWT handling"

# Validation
vic check
vic validate
```

## Related Skills

| Skill | Purpose |
|-------|---------|
| `vibe-think` | Requirements clarification |
| `vibe-architect` | Architecture design |
| `vibe-develop` | Implementation workflow |
| `vibe-integrity` | Memory and validation |
| `vibe-debug` | Systematic debugging |
