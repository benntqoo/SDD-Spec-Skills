---
name: vibe-architect
description: Use when evaluating technology options, designing system architecture, making tech stack decisions, or creating SPEC-ARCHITECTURE.md.
---

# Vibe Architect

Technical architect tool for technology selection, system architecture design, and SPEC-ARCHITECTURE.md creation.

---

## When to Use

**Use when:**
- Need to select tech stack
- Designing system architecture
- Defining data models
- Creating API contracts
- Determining server-side boundaries
- Evaluating technical trade-offs

**NOT use when:**
- Requirements unclear (use vibe-think)
- Implementing code (use vibe-develop)
- Debugging issues (use vibe-debug)

---

## Core Method

### 1. Technology Selection

```
┌─────────────────────────────────────────┐
│       Technology Evaluation Matrix      │
├─────────────────────────────────────────┤
│                                         │
│  ┌─────────┐    ┌─────────┐           │
│  │ Tech A │ vs │ Tech B │           │
│  └────┬────┘    └────┬────┘           │
│       │               │                  │
│       ▼               ▼                  │
│  ┌─────────────────────────────────┐  │
│  │         Evaluation Dimensions    │  │
│  │  • Learning curve               │  │
│  │  • Community ecosystem          │  │
│  │  • Documentation quality        │  │
│  │  • Maintenance status          │  │
│  │  • Team familiarity           │  │
│  │  • Performance                │  │
│  │  • Security                  │  │
│  │  • Cost                     │  │
│  └─────────────────────────────────┘  │
│                                         │
│  Output: Tech decision → SPEC-ARCH     │
└─────────────────────────────────────────┘
```

### 2. Architecture Design Flow

```
Requirements Analysis
     │
     ▼
┌─────────────────────┐
│   Tech Selection   │ ← SPEC-REQUIREMENTS.md
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   System Design    │ ← Layered architecture, modules
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   Data Model       │ ← ER diagram, relations
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   API Contract     │ ← REST/GraphQL
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   Security Design │
└──────────┬──────────┘
           │
           ▼
     Write SPEC-ARCHITECTURE.md
```

---

## Output

After completion:

1. **SPEC-ARCHITECTURE.md** - Complete technical architecture
2. **Tech Decision Records** - Use `vic record tech`
3. **Risk Identification** - Use `vic record risk`

```bash
# Record tech selection
vic rt --id ARCH-001 \
  --title "Choose PostgreSQL over MongoDB" \
  --decision "Use PostgreSQL as primary database" \
  --category database \
  --reason "Need ACID compliance, complex relationships"

# Record architecture risk
vic rr --id ARCH-RISK-001 \
  --area architecture \
  --desc "Single point of failure in auth service" \
  --impact high
```

---

## Quick Reference

| Phase | Output | Command |
|-------|--------|---------|
| Tech Selection | Evaluated options | `vic rt` |
| Architecture | System diagram | (manual) |
| Data Model | ER diagram | (manual) |
| API Contract | API spec | (manual) |
| Security | Security checklist | `vic rr` |
| Gate Check | Verified | `vic spec gate 1` |

---

## Architecture Diagram Template

```
┌─────────────────────────────────────────────────────────────┐
│                         Client                           │
└─────────────────────────┬───────────────────────────────────┘
                          │ HTTPS
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                     接入层                               │
│   (Next.js / Express / ...)                             │
└─────────────────────────┬───────────────────────────────────┘
                          │
              ┌───────────┴───────────┐
              ▼                       ▼
┌─────────────────────────┐   ┌─────────────────────────────┐
│       业务服务层          │   │       外部服务             │
└───────────┬─────────────┘   └─────────────────────────────┘
            │
            ▼
┌─────────────────────────┐
│       数据层             │
└─────────────────────────┘
```

---

## Phase 1 → Phase 2 Bridge (SPEC → SDD)

When SPEC-ARCHITECTURE.md is complete and Gate 1 passed:

```
┌─────────────────────────────────────────────────────────┐
│  ✅ SPEC-ARCHITECTURE.md complete                       │
│  ✅ Gate 1 (Architecture Completeness) passed           │
│                                                         │
│  Next step: Enter SDD Phase                             │
│                                                         │
│  1. Record completion:                                  │
│     vic spec gate pass --gate 1                        │
│                                                         │
│  2. Transition to SDD workflow:                         │
│     → Run skill:sdd-orchestrator                       │
│       (it will call skill:knowledge-boundary +          │
│        skill:pre-decision-check at entry)              │
│                                                         │
│  3. The orchestrator routes to spec-architect:          │
│     → spec-architect freezes requirements into          │
│       SPEC-REQUIREMENTS.frozen.md + contracts          │
│                                                         │
│  4. Then routes to spec-to-codebase for generation      │
│                                                         │
│  ⚠️ Do NOT skip the orchestrator and call               │
│     downstream SDD skills directly                     │
└─────────────────────────────────────────────────────────┘
```

---

## Self-Awareness Integration

### At Entry

```
1. skill:knowledge-boundary
   → Query: "What do I know/infer/assume/unknown about
     the tech stack, architecture patterns, and this domain?"
   → If unknown blocks architecture decisions → STOP

2. skill:pre-decision-check
   → Check scope, quality hard-lines, signals
   → If STOP/BLOCK → do not proceed, ask human
```

### At Each Tech Decision Point

```
Before committing to any tech selection:
1. skill:pre-decision-check
   → Scope check: Is this in approved/forbidden?
   → Quality check: Any hard-line violations?

After making a tech decision:
1. skill:signal-register
   → Record: type=tech_decided, content="<decision>"

2. skill:exploration-journal
   → Record: action=decided, choice="<decision>",
     alternatives_considered=[...]

3. vic rt (CLI)
   → Record technical decision for traceability
```

### At Completion

```
1. skill:signal-register
   → Final confidence recalculation

2. skill:knowledge-boundary (wrap-up)
   → Move inferred → known (if verified during design)
   → Move assumed → inferred/known (if validated)
```

---

## Required Sections

| Section | Content | Importance |
|---------|---------|------------|
| Tech Stack | Each tech with rationale | ⭐⭐⭐ |
| Architecture | Diagram, module划分 | ⭐⭐⭐ |
| Data Model | Entities, relationships | ⭐⭐⭐ |
| API Design | Contracts, error codes | ⭐⭐ |
| Security | Auth, encryption, protections | ⭐⭐⭐ |
| Server Boundaries | What must be server-side | ⭐⭐ |

---

## Related Skills

| Skill | Relationship |
|-------|--------------|
| `vibe-think` | Requirements input → SPEC-REQUIREMENTS.md |
| `vic CLI` | Record technical decisions |
| `knowledge-boundary` | Self-awareness: query known/inferred/assumed/unknown |
| `pre-decision-check` | Self-awareness: check before tech decisions |
| `signal-register` | Self-awareness: record tech decisions as signals |
| `sdd-orchestrator` | **Phase 1→Phase 2 bridge**: enter SDD after Gate 1 |
| `vibe-debug` | Analyze architecture issues |

---

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Selecting trendy tech, not appropriate | Evaluate based on project needs |
| Skipping trade-off analysis | Always compare 2+ options |
| Not documenting rationale | Record reason for each decision |
| Over-engineering early | Start simple, evolve as needed |
| Ignoring team skills | Consider team familiarity |
| Skipping security design | Address security early |

---

## Quick Checklist

Before architecture design:
- [ ] Requirements complete? (Gate 0 passed)
- [ ] Understanding what needs to build?
- [ ] Activated knowledge-boundary? (known/inferred/assumed/unknown categorized)

Before technology selection:
- [ ] Evaluated alternatives?
- [ ] Considered team familiarity?
- [ ] Considered long-term maintenance?
- [ ] pre-decision-check passed for each major choice?

Before SPEC completion:
- [ ] All sections filled?
- [ ] Architecture diagram clear?
- [ ] Gate 1 check passed?
- [ ] Entered SDD workflow via sdd-orchestrator?
