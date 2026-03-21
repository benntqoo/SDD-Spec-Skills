---
name: quick
description: Handles simple, single-file changes that don't require full SDD workflow.
metadata:
  domain: engineering
  version: "1.0"
  tags: [quick, simple, single-file, trivial, typo, rename]
  examples:
    - "Fix a typo"
    - "Rename a variable"
    - "Add a comment"
    - "Simple refactor (single file)"
  priority: medium
  auto_activate: false
---

# Quick Workflow

## Overview

Handles trivial tasks that don't need full SDD workflow. Use when:
- Single file change
- Clear scope
- No SPEC update needed

**This skill is for low-risk, simple changes only.**

## L1: When to Use

| Situation | Use Skill? |
|-----------|------------|
| Fix a typo | ✅ Yes |
| Rename a variable | ✅ Yes |
| Add comments | ✅ Yes |
| Simple single-file refactor | ✅ Yes |
| Multi-file changes | ❌ No (use implementation) |
| New feature | ❌ No (use spec-workflow) |
| Complex bug fix | ❌ No (use implementation) |
| Architecture changes | ❌ No (use spec-workflow) |

## L1: Auto-Activate Triggers

| Trigger | When |
|---------|------|
| User explicitly invokes | User asks for simple change |
| `quick`, `typo`, `rename` mentioned | Task described as simple |
| Single file scope | Only one file affected |

**Escalation triggers** (when to NOT use this skill):
- Multi-file changes → Use `implementation` skill
- SPEC impact → Use `spec-workflow` skill
- Complex logic → Use `implementation` skill

## L2: How to Use

### Step 1: Verify Scope

1. **Confirm Single File**
   - Only one file affected
   - No SPEC impact

2. **Confirm Simplicity**
   - No complex logic changes
   - No test changes needed
   - No dependency changes

### Step 2: Make Change

1. **Edit File**
   - Make the minimal change
   - Don't add unnecessary changes

2. **Run Diagnostics**
   - Check for LSP errors
   - Fix any issues

### Step 3: Quick Verification

1. **Build if applicable**
   - Ensure code compiles

2. **Commit if clean**
   - Simple commit message
   - Example: "fix: typo in README.md"

### When to Escalate

If the change turns out to be more complex:
- Multi-file changes → use `implementation` skill
- SPEC impact → use `spec-workflow` skill
- Complex logic → use `implementation` skill

## Vic Commands

此 Skill 激活时，按以下场景调用 vic 命令：

| 场景 | 命令 | 何时用 |
|------|------|-------|
| 确认单文件 | `vic deps list` | 验证改动确实只影响单个模块 |
| AI Slop 检查 | `vic slop scan --type code` | 改动后扫描是否有 AI Slop 代码模式 |
| AI Slop 修复 | `vic slop fix --dry-run=false` | 应用自动修复（默认 dry-run，需手动确认） |
| 诊断检查 | `vic check --category <category>` | 确认改动不影响技术选型 |
| 状态确认 | `vic status` | 确认项目整体状态正常 |

## L3: References (Required Reading)

These references are part of the skill, not optional:

### Required (Always Read)
- `references/quick-guide.md` - Complete usage guide

### Optional (Read if Needed)
- `references/examples.md` - Examples of quick vs non-quick tasks
- `references/escalation.md` - Escalation criteria
