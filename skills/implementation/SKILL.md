---
name: implementation
description: Handles code implementation, debugging, testing, and SPEC alignment.
metadata:
  domain: engineering
  version: "1.0"
  tags: [implementation, debugging, testing, coding, tdd]
  examples:
    - "Implement a new feature"
    - "Fix a bug"
    - "Write tests for code"
    - "Check code vs SPEC alignment"
  priority: critical
  auto_activate: false
---

# Implementation Workflow

## Overview

Handles the complete implementation lifecycle from coding to testing to SPEC alignment. Includes systematic debugging and TDD workflow.

**Merged from:** debugging + qa + spec-contract-diff

## L1: When to Use

| Situation | Use Skill? |
|-----------|------------|
| Implementing new feature | ✅ Yes |
| Fixing a bug | ✅ Yes |
| Writing or running tests | ✅ Yes |
| Checking code vs SPEC alignment | ✅ Yes |
| Designing system architecture | ❌ No (use spec-workflow) |
| Clarifying requirements | ❌ No (use spec-workflow) |

## L1: Auto-Activate Triggers

| Trigger | When |
|---------|------|
| User explicitly invokes | User asks to implement, fix, or test |
| `implement`, `code`, `fix`, `debug` mentioned | Task involves coding |
| Test coverage | User asks to check or improve test coverage |
| SPEC alignment | `vic spec gate 2` or `vic check` called |

## L2: How to Use

### Option A: Feature Implementation (TDD)

1. **Read SPEC**
   - Read SPEC-ARCHITECTURE.md
   - Read SPEC-REQUIREMENTS.md

2. **Start TDD**
   ```bash
   vic tdd start --feature "[feature]"
   ```

3. **RED Phase**: Write failing test
   ```bash
   vic tdd red --test "[test description]"
   ```

4. **GREEN Phase**: Make it pass
   ```bash
   vic tdd green --test "[test description]" --passed
   ```

5. **REFACTOR Phase**: Improve code
   ```bash
   vic tdd refactor
   ```

6. **Check Alignment**
   - Run `vic spec gate 2`
   - If failed: Fix alignment

7. **Check Tests**
   - Run `vic spec gate 3`
   - If failed: Fix tests

### Option B: Bug Fix (Systematic Debugging)

1. **Start Debug Session**
   ```bash
   vic debug start --problem "[description]"
   ```

2. **Survey**: Gather evidence
   ```bash
   vic debug survey
   ```

3. **Pattern**: Find similar issues
   ```bash
   vic debug pattern
   ```

4. **Hypothesis**: Form and test
   ```bash
   vic debug hypothesis --explain "[explanation]"
   ```

5. **Implement**: Fix root cause
   ```bash
   vic debug implement --fix "[fix description]" --root-cause "[root cause]"
   ```

6. **Verify**: Run tests to confirm fix

### Option C: SPEC Alignment Check

1. **Run Gate 2**
   ```bash
   vic spec gate 2
   ```

2. **If failed**:
   - Option A: Update SPEC (preferred)
   - Option B: Fix code alignment

## Vic Commands

此 Skill 激活时，按以下场景调用 vic 命令：

### 代码实现 (TDD)

| 场景 | 命令 | 何时用 |
|------|------|-------|
| 开始 TDD | `vic tdd start --feature "<feature>"` | 启动 TDD 工作流 |
| 红阶段 | `vic tdd red --test "<test>"` | 写失败测试 |
| 绿阶段 | `vic tdd green --test "<test>" --passed` | 测试通过后标记 |
| 重构 | `vic tdd refactor` | 重构阶段 |
| TDD 状态 | `vic tdd status` | 查看当前 TDD 状态 |

### 代码验证

| 场景 | 命令 | 何时用 |
|------|------|-------|
| 检查对齐 (必须) | `vic spec gate 2` | **每次实现完成必跑**，Gate 2 验证 |
| 检查测试 | `vic spec gate 3` | 实现后验证测试覆盖率 |
| 技术选型检查 | `vic check` | 验证实际技术栈与决策记录一致 |
| 依赖影响 | `vic deps impact <module>` | 修改前查看影响范围 |
| 模块依赖 | `vic deps list` | 了解模块依赖结构 |

### Bug 调试 (Systematic Debugging)

| 场景 | 命令 | 何时用 |
|------|------|-------|
| 开始调试 | `vic debug start --problem "<描述>"` | 启动 4 阶段调试流程 |
| 证据收集 | `vic debug survey` | 调试 Step 1：收集证据 |
| 模式识别 | `vic debug pattern` | 调试 Step 2：寻找相似问题 |
| 假设验证 | `vic debug hypothesis --explain "<解释>"` | 调试 Step 3：验证根因假设 |
| 实施修复 | `vic debug implement --fix "<修复>" --root-cause "<根因>"` | 调试 Step 4：实施修复 |

### 实施完成

| 场景 | 命令 | 何时用 |
|------|------|-------|
| 检查 AI Slop | `vic slop scan` | 清理 AI 生成的低质量代码模式 |
| 修复 Slop | `vic slop fix --dry-run=false` | 应用自动修复（预览后再执行） |

## L3: References (Required Reading)

These references are part of the skill, not optional:

### Required (Always Read)
- `references/implementation-guide.md` - Complete usage guide

### Optional (Read if Needed)
- `references/tdd-guide.md` - TDD workflow details
- `references/debugging-guide.md` - Systematic debugging methodology
- `references/troubleshooting.md` - Common issues and fixes
