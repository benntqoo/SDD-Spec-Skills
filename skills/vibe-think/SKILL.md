---
name: vibe-think
description: Use when requirements are ambiguous, user provides unclear descriptions, or multiple options need trade-off evaluation before implementation.
---

# Vibe Think

Structured thinking and enhanced questioning tool for requirements clarification.

---

## When to Use

**Use when:**
- Requirements are ambiguous or incomplete
- Multiple方案需要评估权衡
- Technical decisions need analysis
- Need deeper understanding of the problem
- User provides vague descriptions

**NOT use when:**
- Requirements are clear
- Already in implementation phase
- Only need simple confirmation

---

## Core Method

### 1. Clarifying Questions

Ask one focused question at a time:

```
❌ Bad:
"Does this feature need auth? What tech stack? Which database?"

✅ Good:
"Does this feature involve sensitive data?"
   ↓ (answer)
"Does user need to login first?"
   ↓ (answer)
"Any preference for auth method? JWT or Session?"
```

### 2. Progressive Understanding

```
┌─────────────┐
│ User Description│
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Clarifying │ ← What question to ask?
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Explore    │ ← What options available?
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Trade-off  │ ← Pros/cons of each?
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Recommended │ ← Which to choose?
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Record      │ ← vic record tech
└─────────────┘
```

### 3. Trade-off Analysis Template

```markdown
## Option Comparison

| Dimension | Option A | Option B |
|-----------|----------|-----------|
| Dev Cost | Low | High |
| Maintenance | High | Low |
| Performance | Medium | High |
| Risk | Medium | Low |
| Recommendation | ⭐⭐⭐ | ⭐⭐⭐⭐ |

**Recommendation**: Option B, because...
```

---

## Quick Reference

| Step | Action | Command |
|------|---------|---------|
| 1 | Ask clarifying question | (manual) |
| 2 | Explore options | (manual) |
| 3 | Analyze trade-offs | (manual) |
| 4 | Record decision | `vic rt --id X --title Y --decision Z` |
| 5 | Record risk | `vic rr --id X --area Y --desc Z` |

---

## Output

After completion:

1. **Understanding Summary** - Verbal or written
2. **Decision Record** - Use `vic record tech`
3. **Risk Identification** - Use `vic record risk` (if applicable)

```bash
# Record discussion result
vic rt --id ARCH-001 \
  --title "Choose X over Y" \
  --decision "Use X because..." \
  --reason "Trade-off analysis showed..."
```

---

## Example Dialog

```
User: "I want a search feature"

AI (vibe-think):
"What content should be searchable? User input or structured queries?"
↓
User: "User input keywords"

AI: "How large is the data? Hundreds or millions?"
↓
User: "Thousands now, potentially millions later"

AI: "Need real-time or delay acceptable?"
↓
User: "Real-time preferred, but 1-2s delay OK"

AI:
"Based on your needs, I analyzed two options:

| Dimension | SQLite LIKE | Elasticsearch |
|------|-------------|---------------|
| Dev Cost | Low (1 day) | High (3-5 days) |
| Ops Cost | None | Medium |
| Performance (1M) | Poor | Excellent |
| Latency | <100ms | <50ms |

**Recommendation**: Use SQLite LIKE for MVP, migrate to Elasticsearch when data exceeds 100K.

Want me to record this decision?"
↓
User: "Yes"

AI:
vic rt --id SEARCH-001 \
  --title "Search Strategy" \
  --decision "SQLite LIKE for MVP, Elasticsearch later" \
  --reason "Data volume currently low, defer infra complexity"
```

---

## Related Skills

| Skill | Relationship |
|-------|--------------|
| `vic CLI` | Record decisions |
| `vibe-debug` | Analyze problem root cause |
| `vibe-integrity` | Read project context |
| `vibe-architect` | Tech selection → SPEC-ARCHITECTURE.md |
| `vibe-develop` | Implementation → Gate checks |

---

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Asking too many questions at once | Ask one focused question at a time |
| Skipping trade-off analysis | Always document pros/cons |
| Not recording decisions | Use `vic rt` for every decision |
| Accepting vague requirements | Keep probing until clear |
| Jumping to solutions | Understand problem first |

---

## Quick Checklist

Before clarification:
- [ ] Do I understand the core problem?
- [ ] What are the constraints?
- [ ] What options exist?
- [ ] What are trade-offs?
- [ ] Is my recommendation justified?

When recording:
- [ ] Decision ID follows naming convention?
- [ ] Reason clearly stated?
- [ ] Related risks identified?
