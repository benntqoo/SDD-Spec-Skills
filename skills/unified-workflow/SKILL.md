---
name: unified-workflow
description: Orchestrates SDD workflow, enforces Constitution rules, and maintains traceability.
metadata:
  domain: governance
  version: "1.0"
  tags: [sdd, orchestration, constitution, traceability, gates, workflow]
  examples:
    - "Start a new feature delivery"
    - "Advance SDD phase"
    - "Check before commit"
    - "Verify requirements-to-code mapping"
  priority: critical
  auto_activate: false
---

# Unified Workflow

## Overview

Single controller for SDD workflow, Constitution enforcement, and traceability tracking. Manages the complete feature delivery lifecycle.

**Merged from:** constitution-check + sdd-orchestrator + spec-traceability

## L1: When to Use

| Situation | Use Skill? |
|-----------|------------|
| Start new feature delivery | ✅ Yes |
| Advance SDD phase | ✅ Yes |
| Before git commit | ✅ Yes |
| Check requirements traceability | ✅ Yes |
| During implementation | ❌ No (use implementation) |
| Clarifying requirements | ❌ No (use spec-workflow) |

## L1: Auto-Activate Triggers

| Trigger | When |
|---------|------|
| User explicitly invokes | User asks to manage workflow |
| `vic auto start` | Starting autonomous mode |
| `vic gate check` | Checking gate compliance |
| Pre-commit | Before git commit |
| Phase advancement | Moving to next SDD phase |

## L2: How to Use

### Workflow: Feature Delivery

1. **Start Delivery**
   ```bash
   vic auto start
   ```

2. **Check Constitution**
   - Read .vic-sdd/constitution.yaml
   - Verify all rules are satisfied
   - If blockers found: resolve before continuing

3. **Manage SDD Phases**
   SDD State Machine:
   ```
   Ideation → Explore → SpecCheckpoint → Build → Verify → ReleaseReady → Released
   ```

4. **Gate Checks at Each Phase**
   | Phase | Gate | Check |
   |-------|------|-------|
   | Ideation | Gate 0 | Requirements completeness |
   | Explore | Gate 1 | Architecture completeness |
   | Build | Gate 2 | Code alignment |
   | Verify | Gate 3 | Test coverage |

5. **Traceability Check**
   - Verify: User Story → SPEC Contract → Code → Tests
   - Each requirement has implementation
   - Each implementation has tests

6. **End Delivery**
   ```bash
   vic auto stop
   ```

### Workflow: Pre-Commit Check

1. **Run Constitution Check**
   - Read .vic-sdd/constitution.yaml
   - Check each principle

2. **Run Gate Checks**
   ```bash
   vic gate check --blocking
   ```

3. **Fix Issues if Any**
   - Resolve blockers
   - Update SPEC if needed

### Workflow: Traceability Check

1. **Read Traceability Map**
   - User Story → SPEC Contract → Code → Tests

2. **Verify Mapping**
   - Each requirement has implementation
   - Each implementation has tests

3. **Update if Needed**
   - Add missing mappings
   - Remove orphaned code

## Vic Commands

此 Skill 激活时，按以下场景调用 vic 命令：

### 自主模式

| 场景 | 命令 | 何时用 |
|------|------|-------|
| 启动自主模式 | `vic auto start` | 开始 autonomous 开发 |
| 查看状态 | `vic auto status` | 监控自主模式进度、Token 消耗 |
| 暂停 | `vic auto pause` | 临时停止自主模式 |
| 恢复 | `vic auto resume` | 从暂停状态恢复 |
| 停止 | `vic auto stop` | 结束自主模式会话 |

### 阶段推进

| 场景 | 命令 | 何时用 |
|------|------|-------|
| 推进阶段 | `vic phase advance --to <N>` | 从当前阶段推进到下一阶段（自动跑 Gate） |
| 查看当前阶段 | `vic phase show` | 确认当前 SDD 阶段 |

### 提交前检查（强制）

| 场景 | 命令 | 何时用 |
|------|------|-------|
| 阻断性检查 (必须) | `vic gate check --blocking` | **提交前必跑**，所有 blocking 项必须清零 |
| SPEC Hash | `vic spec hash` | 检测 SPEC 是否变更，变更则不能提交 |
| Constitution 摘要 | `vic cost status` | 确认成本在预算内 |

### 功能交付 / Traceability

| 场景 | 命令 | 何时用 |
|------|------|-------|
| 查看里程碑 | `vic milestone list` | 确认交付在哪个里程碑 |
| 追溯链路 | `vic history --limit 10` | 验证 User Story → SPEC → Code → Test 链路 |
| 导出数据 | `vic export --output backup.json` | 交付前备份 vic-sdd 数据 |

## L3: References (Required Reading)

These references are part of the skill, not optional:

### Required (Always Read)
- `references/unified-workflow-guide.md` - Complete usage guide

### Optional (Read if Needed)
- `references/sdd-state-machine.md` - SDD state machine details
- `references/constitution-rules.md` - Constitution rule definitions
- `references/traceability-patterns.md` - Traceability patterns
