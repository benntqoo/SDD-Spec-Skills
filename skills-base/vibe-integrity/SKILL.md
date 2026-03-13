# Vibe Integrity Framework
# ========================
## AI Project Memory & Safety System

---

## Overview

Vibe Integrity Framework is **NOT a development methodology**. It's a **safety net for vibe coding** that provides:

1. **Completion Guard** - Detects AI false completion claims
2. **Architecture Memory** - Project knowledge base for AI quick understanding

**Methodology Agnostic**: Works with TDD, SDD, Agile, or pure vibe coding.

---

## Two Pillars

### Pillar 1: Completion Guard
Detection and validation to ensure AI actually completed the work.

| Skill | Purpose |
|-------|---------|
| `vibe-guard` | Detects TODO, empty functions, fake tests |
| `cascade-check` | Prevents cascading errors after fixes |
| `integration-check` | Validates component integration |

#XR|**See**: `skills-base/vibe-guard/SKILL.md`

### Pillar 2: Architecture Memory
Structured project knowledge base for AI quick understanding.

| File | Purpose |
|------|---------|
| `project.yaml` | Project metadata, tech stack |
| `dependency-graph.yaml` | Module dependencies |
| `module-map.yaml` | Directory structure |
| `risk-zones.yaml` | High-risk areas |
| `tech-records.yaml` | Technical decisions |
| `schema-evolution.yaml` | Data model changes |

---

## AI Quick Start

When AI starts work on this project, read in this order:

```
1. .vibe-integrity/project.yaml
   → Understand project status, tech stack

2. .vibe-integrity/risk-zones.yaml  
   → Know what areas are high-risk

3. .vibe-integrity/dependency-graph.yaml
   → Understand module relationships

4. .vibe-integrity/module-map.yaml
   → Find where files are located

5. .vibe-integrity/tech-records.yaml
   → Understand why system is designed this way
```

**Result**: AI understands project in ~15 seconds instead of 3 minutes.

---

## Directory Structure

```
.vibe-integrity/
├── project.yaml              # Project metadata
├── dependency-graph.yaml     # Module dependencies
├── module-map.yaml          # Directory structure
├── risk-zones.yaml          # Risk areas
├── tech-records.yaml        # Technical decisions
└── schema-evolution.yaml   # Data model changes

skills/
├── vibe-guard/             # Completion detection
└── vibe-integrity/         # This skill
    └── template/           # Schema templates
```

---

## Usage

### For AI: Before Making Changes

```bash
# 1. Check risk zone
cat .vibe-integrity/risk-zones.yaml

# 2. Check dependencies
cat .vibe-integrity/dependency-graph.yaml

# 3. Check schema
cat .vibe-integrity/schema-evolution.yaml
```

### For AI: After "Completing"

```bash
# Run vibe-guard
#SZ|python skills-base/vibe-guard/validate-vibe-guard.py --check
```

### For Humans: After Significant Changes

```bash
# Update tech-records
# Add new decision to .vibe-integrity/tech-records.yaml

# Update schema-evolution  
# Add new version to .vibe-integrity/schema-evolution.yaml

# Update dependency-graph
# Reflect new module relationships
```

---

## Core Principles

### 1. AI Should Read First
**Before any work**, AI reads `.vibe-integrity/` files to understand project state.

### 2. Record Decisions, Not Just Code
Every technical decision should be documented in `tech-records.yaml`.

### 3. Track Schema Changes
Any database schema change should be recorded in `schema-evolution.yaml`.

### 4. Annotate Risk
High-risk areas should be clearly marked in `risk-zones.yaml`.

### 5. Verify After Changes
After any fix, run `vibe-guard` to ensure completeness.

---

## Key Benefit

| Without Vibe Integrity | With Vibe Integrity |
|------------------------|---------------------|
| AI reads 100+ files | AI reads 5 key files |
| 3+ minutes understanding | 15 seconds understanding |
| May miss dependencies | Dependencies documented |
| May break high-risk code | Risk zones marked |
| No project memory | Complete history |

---

## Files Detail

### project.yaml
```yaml
name: my-project
version: 0.1.0
status: mvp
tech_stack:
  frontend: [Vue, Vite]
  backend: [Express, Node]
  database: [SQLite]
```

### dependency-graph.yaml
```yaml
modules:
  auth-service:
    depends_on: [user-repository, jwt-service]
    used_by: [api-routes-auth]
```

### risk-zones.yaml
```yaml
zones:
  auth-service:
    risk_level: critical
    reason: Handles credentials
```

### tech-records.yaml
```yaml
records:
  - id: DB-001
    date: "2026-01-15"
    title: "Choose SQLite for MVP"
    decision: "Use SQLite for fast iteration"
```

### schema-evolution.yaml
```yaml
tables:
  - name: users
    versions:
      - version: "1.0"
        fields:
          - name: email
            type: String
```

---

## Validation
Run validation to ensure integrity:
#PB|#- `python skills-base/vibe-integrity/validate-vibe-integrity.py` - checks .vibe-integrity/ files
#QZ|#- `python skills-base/vibe-integrity/validate-all.py` - runs both vibe-guard and vibe-integrity validations
#TR|#- `python skills-base/vibe-guard/validate-vibe-guard.py --check` - AI completion check

## Related Skills

#PW|#- `skills-base/vibe-guard` - Completion detection
- `superpowers/test-driven-development` - TDD workflow (optional)
#WR|#- `skills-sdd/sdd-orchestrator` - SDD workflow (optional)

**Note**: Vibe Integrity works with ANY development approach.
