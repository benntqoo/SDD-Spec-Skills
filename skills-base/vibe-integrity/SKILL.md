---
name: vibe-integrity
description: Use when working on a project, recording technical decisions, identifying risks, verifying code alignment, or managing SPEC workflow.
---

# Vibe Integrity

AI Project Memory & Safety System - prevents false completion claims and provides structured project knowledge.

---

## Quick Commands

```bash
# Initialize project
vic init --name "My Project" --tech "React,Node,PostgreSQL"

# Record technical decision
vic rt --id DB-001 --title "Use PostgreSQL" --decision "Primary database" --reason "Need ACID"

# Record risk
vic rr --id RISK-001 --area auth --desc "JWT handling issue"

# SPEC commands
vic spec init --name "My Project"
vic spec status
vic spec gate 0  # Requirements
vic spec gate 1  # Architecture

# Validate
vic check
vic validate
vic status
```

---

## When to Use

| Scenario | Command |
|----------|---------|
| Start new project | `vic init` |
| Make technical decision | `vic rt` |
| Identify risk | `vic rr` |
| Record dependency | `vic rd` |
| AI claims "done" | `vic check` |
| Before commit | `vic validate` |
| Check SPEC status | `vic spec status` |
| Run Gate checks | `vic spec gate [0-3]` |
| Backup memory | `vic export` |

---

## Quick Reference

| Command | Alias | Purpose |
|---------|-------|---------|
| `vic init` | - | Initialize .vic-sdd/ |
| `vic spec init` | - | Initialize SPEC docs |
| `vic rt` | `record-tech` | Record decision |
| `vic rr` | `record-risk` | Record risk |
| `vic rd` | `record-dep` | Record dependency |
| `vic check` | - | Code alignment |
| `vic validate` | - | Full validation |
| `vic spec gate [0-3]` | - | Gate checks |
| `vic status` | - | Project status |
| `vic search` | - | Search records |
| `vic history` | - | Event history |
| `vic export` | - | Export data |
| `vic import` | - | Import data |

---

## Gate Reference

| Gate | Name | Checks |
|------|------|--------|
| Gate 0 | Requirements | User stories, acceptance criteria, phase plan |
| Gate 1 | Architecture | Tech stack, data model, API design |
| Gate 2 | Code Alignment | Code matches SPEC-ARCHITECTURE.md |
| Gate 3 | Test Coverage | Acceptance criteria covered |

---

## Directory Structure

```
.vic-sdd/
├── SPEC-REQUIREMENTS.md    # Requirements spec
├── SPEC-ARCHITECTURE.md    # Architecture spec
├── PROJECT.md              # Project status
├── status/
│   ├── events.yaml         # Event history
│   └── state.yaml          # Current state
├── tech/
│   └── tech-records.yaml  # Technical decisions
├── risk-zones.yaml        # Risk records
├── project.yaml           # AI quick reference
└── dependency-graph.yaml  # Module dependencies
```

---

## Related Skills

| Skill | Purpose |
|-------|---------|
| `vibe-think` | Requirements clarification → SPEC-REQUIREMENTS.md |
| `vibe-architect` | Tech selection → SPEC-ARCHITECTURE.md |
| `vibe-develop` | Implementation, Gate checks |
| `vibe-debug` | Systematic debugging |

---

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Skipping `vic check` before commit | Always validate code alignment |
| Recording vague decisions | Include specific reason and impact |
| Not updating risk status | Use `vic rr --status resolved` when fixed |
| Forgetting SPEC updates | Update SPEC before Gate checks |
| Ignoring Gate failures | Fix issues before proceeding |

---

## Quick Checklist

Before claiming completion:
- [ ] Ran `vic check`?
- [ ] All tech decisions recorded?
- [ ] New risks identified?

Before Gate progression:
- [ ] Gate 0: Requirements complete?
- [ ] Gate 1: Architecture complete?
- [ ] Gate 2: Code aligns with SPEC?
- [ ] Gate 3: Tests cover acceptance criteria?

Before commit:
- [ ] `vic validate` passes?
- [ ] `vic export` for backup?
