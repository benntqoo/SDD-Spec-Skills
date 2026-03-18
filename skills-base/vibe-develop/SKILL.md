---
name: vibe-develop
description: Use when implementing features, managing development workflow, running gate checks, or coordinating milestone deliveries.
---

# Vibe Develop

Development execution tool for code implementation, Gate checks, milestone management, and workflow control.

---

## When to Use

**Use when:**
- Starting feature implementation
- Running Gate checks
- Managing milestones
- Code quality control
- Coordinating development workflow

**NOT use when:**
- Requirements unclear (use vibe-think)
- Technology selection (use vibe-architect)
- Debugging issues (use vibe-debug)

---

## Development Flow

### Simplified SDD State Machine

```
    ┌──────┐     ┌───────────┐     ┌──────┐     ┌────────┐     ┐
    │ Plan │────▶│   Spec    │────▶│Build │────▶│Verify  │────▶│ Done
    └──────┘     └───────────┘     └──────┘     └────────┘     └
       │              │               │              │
       ▼              ▼               ▼              ▼
    Requirements    Contract        Code         Tests
    
    Gate 0:      Gate 1:         Gate 2:      Gate 3:
    Requirements  Contract        Code         Tests
    Complete     Complete        Aligned      Covered
```

### Gate Check Flow

```
Start Development
    │
    ▼
┌─────────────────┐
│ Gate 0 Check    │ ← vic spec gate 0
│ Requirements    │
└────────┬────────┘
         │ PASS
         ▼
┌─────────────────┐
│ Gate 1 Check    │ ← vic spec gate 1
│ Architecture    │
└────────┬────────┘
         │ PASS
         ▼
    Code Implementation
         │
         ▼
┌─────────────────┐
│ Gate 2 Check    │ ← vic spec gate 2
│ Code Alignment  │ ← vic check
└────────┬────────┘
         │ PASS
         ▼
    Test Implementation
         │
         ▼
┌─────────────────┐
│ Gate 3 Check    │ ← vic spec gate 3
│ Test Coverage    │
└────────┬────────┘
         │ PASS
         ▼
    Milestone Acceptance
```

---

## Core Methods

### 1. Small Iteration

```
Don't:
❌ Implement entire feature → Test → Commit

Do:
✅ Implement skeleton → Commit
✅ Implement core logic → Commit  
✅ Implement edge cases → Commit
✅ Add tests → Commit
```

Each small stage:
- Feature runs
- Has corresponding commit
- Can be rolled back

### 2. Human Intervention Points

After each milestone, required:

1. **Code Review** - Check for extractable common modules
2. **Architecture Review** - Any code smells?
3. **Refactoring Opportunities** - Any duplicate code?

```bash
# Check code alignment
vic check

# Run Gate 2
vic spec gate 2

# Run Gate 3
vic spec gate 3

# Check project status
vic spec status
```

### 3. AI Permission Boundaries

**AI CANNOT on its own:**
- Change code style
- Refactor unrelated modules
- Modify UI styles
- Introduce new dependencies
- Change architecture decisions

**Human must confirm:**
- Introducing new dependencies
- Architecture changes
- Major refactoring

---

## Quick Reference

| Gate | Command | Checks |
|------|---------|--------|
| Gate 0 | `vic spec gate 0` | User stories, acceptance criteria, phase plan |
| Gate 1 | `vic spec gate 1` | Tech stack, data model, API design |
| Gate 2 | `vic spec gate 2` | Code matches SPEC |
| Gate 3 | `vic spec gate 3` | Acceptance criteria coverage |

| Command | Purpose |
|---------|---------|
| `vic spec gate [0-3]` | Run specific Gate check |
| `vic check` | Code alignment check |
| `vic validate` | Full validation |
| `vic status` | Project status |

---

## Output

After completion:

1. **Implementation Code** - Aligns with SPEC-ARCHITECTURE.md
2. **Acceptance Tests** - Cover acceptance criteria
3. **Gate Reports** - Pass Gate 2/3
4. **Code Commits** - Meaningful commit messages

```bash
# Record milestone completion
vic rt --id DEV-001 \
  --title "Feature X Implementation Complete" \
  --decision "Implemented as per SPEC" \
  --status completed
```

---

## Related Skills

| Skill | Relationship |
|-------|--------------|
| `vibe-think` | Requirements → SPEC-REQUIREMENTS.md |
| `vibe-architect` | Architecture → SPEC-ARCHITECTURE.md |
| `vibe-integrity` | Validation → `vic check` |
| `vibe-debug` | Issue debugging |

---

## Security Boundaries

**MUST follow:**
- API Keys NEVER in frontend code
- Database MUST have authentication
- Sensitive data MUST be encrypted
- All APIs MUST have authorization
- Input MUST be validated

**MUST NOT trust:**
- Frontend validation
- Client-side hiding
- Obscurity

---

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Skipping Gate checks | Run `vic spec gate [0-3]` before progression |
| Large commits | Small iterations with meaningful messages |
| AI changing code style | Human must approve style changes |
| Skipping tests | Always add regression tests |
| Not recording decisions | Use `vic rt` for changes |
| Ignoring security | Follow security boundaries strictly |

---

## Quick Checklist

Before implementation:
- [ ] Gate 0 passed?
- [ ] Gate 1 passed?
- [ ] Understand acceptance criteria?
- [ ] Know server boundaries?

After each small stage:
- [ ] Code runs?
- [ ] Commit message meaningful?
- [ ] Using vibe-debug for issues?

After milestone:
- [ ] Gate 2 passed?
- [ ] Gate 3 passed?
- [ ] Any code to extract to common modules?
- [ ] Commit message clear?

Before claiming done:
- [ ] All Gates passed?
- [ ] `vic validate` passes?
- [ ] Security boundaries followed?
