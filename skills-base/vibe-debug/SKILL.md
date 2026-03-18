---
name: vibe-debug
description: Use when tests fail, bugs occur, unexpected behavior happens, or error messages are unclear and root cause analysis is needed.
---

# Vibe Debug

Systematic debugging methodology for root cause analysis.

---

## When to Use

**Use when:**
- Tests fail
- Bugs need fixing
- Unexpected behavior
- Error messages unclear
- Attempted fixes but problem persists

**NOT use when:**
- Syntax errors (fix directly)
- Simple config issues
- Clear-cut problems
- Writing new code

---

## Core Method: 4-Phase Analysis

```
┌─────────────────────────────────────────────────────────┐
│                    4-Phase Debug                        │
├─────────────────────────────────────────────────────────┤
│                                                          │
│   ┌─────────────┐     ┌─────────────┐                  │
│   │  1. Root   │ ──▶ │  2. Pattern│                  │
│   │  Cause     │     │  Analysis   │                  │
│   │  Survey    │     │             │                  │
│   └──────┬──────┘     └──────┬──────┘                  │
│          │                    │                          │
│          ▼                    ▼                          │
│   ┌─────────────┐     ┌─────────────┐                  │
│   │  3. Hypoth │ ──▶ │  4. Implement │                  │
│   │  esis Test │     │  Fix        │                  │
│   └─────────────┘     └─────────────┘                  │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### Phase 1: Root Cause Survey

**Rule: Never fix symptom without understanding root cause.**

Steps:
1. Reproduce the error
2. Gather evidence (logs, stack traces, environment)
3. Question assumptions
4. Identify what's NOT causing the issue

### Phase 2: Pattern Analysis

Look for patterns:
- Similar issues in codebase?
- Known anti-patterns?
- Recent changes that might relate?
- Platform/library specific issues?

### Phase 3: Hypothesis Testing

Form testable hypotheses:
- "It's caused by X because Y"
- Test each hypothesis minimally
- Measure, don't guess

### Phase 4: Implement Fix

Only after root cause confirmed:
- Fix root cause, not symptoms
- Add regression test
- Document the fix

---

## Quick Reference

| Phase | Action | Output |
|-------|--------|--------|
| 1 | Reproduce + Gather evidence | Error details |
| 2 | Search patterns | Possible causes |
| 3 | Form + test hypotheses | Confirmed root cause |
| 4 | Implement fix | Regression test |

| Command | Purpose |
|---------|---------|
| `vic rr` | Record discovered risk |
| `vic check` | Verify fix alignment |

---

## Example Debug Flow

```
User: "Login always fails"

AI (vibe-think):
"Show me the error message and when it occurs"

Phase 1 - Root Cause Survey:
- Error: "Invalid credentials"
- When: Every login attempt
- Environment: Production only?

Phase 2 - Pattern Analysis:
- Search for similar auth issues
- Check recent changes to auth code

Phase 3 - Hypothesis Testing:
- Hypothesis 1: Password hashing mismatch → Test locally
- Hypothesis 2: Token expired → Check timestamps
- Hypothesis 3: Database connection → Test DB

Phase 4 - Implement Fix:
- Root cause found: Hash algorithm version mismatch
- Fix: Migrate password hashes
- Test: Add regression test for auth
```

---

## Related Skills

| Skill | Relationship |
|-------|--------------|
| `vibe-think` | Questioning techniques |
| `vibe-integrity` | Record findings with `vic rr` |
| `vibe-develop` | Implement fix with Gate checks |

---

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Fixing symptoms not root cause | Always ask "why" 3 times |
| Random changes hoping something works | Form hypotheses first |
| Not reproducing the issue | Always reproduce before fixing |
| Skipping regression tests | Add test before declaring fix |
| Not recording the issue | Use `vic rr` for tracking |
| Changing code without understanding | Read code first, then hypothesize |

---

## Quick Checklist

Before proposing fix:
- [ ] Reproduced the error?
- [ ] Gathered evidence (logs, traces)?
- [ ] Questioned assumptions?
- [ ] Identified what NOT causing it?
- [ ] Found similar patterns?
- [ ] Formed testable hypothesis?
- [ ] Tested hypothesis minimally?
- [ ] Root cause confirmed?
- [ ] Will add regression test?
- [ ] Recorded in `vic rr`?

---

**Golden Rule: Never fix a symptom without understanding the root cause.**
