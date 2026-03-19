# Agent Collaboration Guide

**Updated for Vibe-SDD Development Workflow**

## Current System Status

✅ **Multi-Agent Supported**: Yes, via Git branching workflow
✅ **Multi-User Supported**: Yes, via Git collaboration
✅ **Structured Development**: Yes, via .vic-sdd/ SPEC workflow
✅ **Self-Aware AI**: Yes, via Knowledge Boundary + Signal Register
❌ **Real-Time Collaboration**: No, requires Git merge workflow

---

## 核心设计原则

### AI 的自知之明

```
当前 Coding Agent 的根本问题：
─────────────────────────────────────────────────────────────
不是"AI 太笨"，而是"AI 不知道自己不知道"

当 AI 说"我理解了"时，它可能是在：
  1. 真的理解了（基于代码库的事实）
  2. 从模式推断的（可能是错的）
  3. 假设的（完全没验证）
  4. 幻觉的（编造的）
─────────────────────────────────────────────────────────────

VIC-SDD 的目标：
  • 不是监控 AI，而是确保 AI 有"自知之明"
  • 不是给 AI 下命令，而是给 AI 画边界
  • 不是事后检查，而是事前约束
```

---

## 四个核心机制

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                   │
│   1. Knowledge Boundary (认知地图)                               │
│      → AI 知道什么、推测什么、假设什么、不知道什么                    │
│                                                                   │
│   2. Pre-Decision Check (决策前自查)                              │
│      → 重大决策前自动检查边界和约束                                 │
│                                                                   │
│   3. Signal Register (信号注册)                                    │
│      → 用"证据链"代替"进度百分比"                                   │
│                                                                   │
│   4. Exploration Journal (探索日志)                                │
│      → AI 记录思考过程，避免重复探索                                │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 1. Knowledge Boundary（认知地图）

`.vic-sdd/knowledge-boundary.yaml`

```yaml
known:        # 验证过的事实（最高可信度）
inferred:     # 从模式推断的（需要验证）
assumed:      # 假设的（高风险）
unknown:      # 完全不知道的（可能阻塞）
```

**Skill**: `knowledge-boundary`

### 2. Pre-Decision Check（决策前自查）

`.vic-sdd/decision-guardrails.yaml`

```yaml
scope:        # 范围约束
attempts:     # 尝试次数约束
quality:      # 质量红线
signals:      # 信号约束
```

**Skill**: `pre-decision-check`

### 3. Signal Register（信号注册）

`.vic-sdd/signal-register.yaml`

```yaml
signals:
  positive:   # 正面信号
  warnings:   # 警告信号
  blockers:   # 阻塞信号
confidence:   # 信心度计算
```

**Skill**: `signal-register`

### 4. Exploration Journal（探索日志）

`.vic-sdd/exploration-journal.yaml`

```yaml
entries:
  - action: explore   # 开始探索
  - action: tried     # 尝试方法
  - action: decided   # 做出决策
  - action: learned   # 学习教训
```

**Skill**: `exploration-journal`

---

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
Agent-Develop activates self-awareness skills
    ↓
    ├── knowledge-boundary    (查询认知地图)
    ├── pre-decision-check   (决策前检查)
    ├── signal-register      (记录信号)
    └── exploration-journal  (记录探索)
    ↓
Implementation + Tests
    ↓
vic gate pass --gate 4 (Code Compiles)
vic gate pass --gate 5 (Code Aligns SPEC)
    ↓
vic phase advance --to 3
vic gate pass --gate 6-7
    ↓
vic spec merge → PRD.md / ARCH.md / PROJECT.md
```

---

## Self-Aware AI Workflow

```
开始任务
    ↓
┌─────────────────────────────────────────┐
│  1. 照镜子 (knowledge-boundary)          │
│     → 知道什么？不知道什么？              │
│     → 有 unknown/assumed 阻塞？          │
└───────────────┬─────────────────────────┘
                ↓
┌─────────────────────────────────────────┐
│  2. 决策前检查 (pre-decision-check)      │
│     → 范围检查                           │
│     → 质量红线检查                        │
│     → 信号检查                           │
└───────────────┬─────────────────────────┘
                ↓
         检查结果
         ├── ✅ PASS → 继续执行
         ├── ⚠️ WARN → 记录，继续
         ├── 🛑 STOP → 等待人类
         └── 🔴 BLOCK → 解决阻塞
                ↓
┌─────────────────────────────────────────┐
│  3. 执行任务                             │
│     → 产生信号 (signal-register)          │
│     → 记录探索 (exploration-journal)      │
└───────────────┬─────────────────────────┘
                ↓
         周期性信心度检查
         ├── confidence >= 0.7 → 继续
         ├── 0.4 <= confidence < 0.7 → 关注
         └── confidence < 0.4 → 暂停
```

---

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

---

## Conflict Resolution Workflow

When multiple agents modify the same YAML files:

1. **Git detects conflict** during merge/pull request
2. **Union merge** preserves both versions (may create duplicates)
3. **Run validation script** to detect duplicate IDs
4. **Manual resolution** required to merge similar decisions
5. **Verify** application still works with merged memory

---

## Best Practices

1. **Use Separate Branches**: Each agent gets own branch
2. **Activate Self-Awareness**: Use skills before/during/after tasks
3. **Maintain Knowledge Boundary**: Keep known/inferred/assumed/unknown up-to-date
4. **Record All Signals**: Every meaningful action = one signal
5. **Query Before Acting**: Check journal to avoid duplicate exploration

---

## Directory Structure

```
.vic-sdd/
├── SPEC-REQUIREMENTS.md    # Requirements spec
├── SPEC-ARCHITECTURE.md    # Architecture spec
├── PROJECT.md              # Project status tracking
│
├── knowledge-boundary.yaml  # AI 认知地图
├── decision-guardrails.yaml # 决策约束
├── signal-register.yaml    # 信号注册
├── exploration-journal.yaml  # 探索日志
│
├── status/
│   ├── events.yaml         # Event history
│   └── state.yaml          # Current state
├── tech/
│   └── tech-records.yaml  # Technical decisions
├── risk-zones.yaml         # Risk records
├── project.yaml            # AI quick reference
└── dependency-graph.yaml  # Module dependencies

scripts/
└── verify.sh              # 外部验证脚本
```

---

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

# External Verification
./scripts/verify.sh
```

---

## Skills Reference

### Self-Awareness Skills

| Skill | Purpose | When to Activate |
|-------|---------|------------------|
| `knowledge-boundary` | AI 自知之明 | 开始任务前、遇到不确定时 |
| `pre-decision-check` | 决策前刹车 | 重大决策前 |
| `signal-register` | 证据链进度 | 每个有意义行动后 |
| `exploration-journal` | 思考过程记忆 | 开始探索、尝试、决策时 |

### Development Skills

| Skill | Purpose |
|-------|---------|
| `vibe-think` | Requirements clarification |
| `vibe-architect` | Architecture design |
| `vibe-debug` | Systematic debugging |
| `vibe-design` | Design system |
| `vibe-redesign` | Product redesign |
| `adaptive-planning` | Adaptive replanning |

### SDD Skills

| Skill | Purpose |
|-------|---------|
| `spec-architect` | SDD: 需求凝固 |
| `spec-to-codebase` | SDD: 代码生成 |
| `spec-contract-diff` | SDD: 代码对齐检查 |
| `spec-traceability` | SDD: 追溯验证 |
| `spec-driven-test` | SDD: 契约测试 + TDD |
| `sdd-release-guard` | SDD: 发布守卫 |

---

## 质量红线

违反以下任一条都是不允许的：

| 红线 | 说明 |
|------|------|
| `no_todo_in_code` | 代码里不能有 TODO/FIXME |
| `no_console_in_prod` | 生产代码不能有 console.log |
| `no_hardcoded_secrets` | 不能有硬编码密钥 |
| `tests_required` | 新功能必须有测试 |
| `spec_aligned` | 必须与 SPEC 对齐 |

---

## 信心度阈值

```
confidence = (positive - warnings×0.3 - blockers×0.5) / max_signals

> 0.7    → 🟢 HIGH   → 状态良好，继续推进
0.4-0.7  → 🟡 MODERATE → 可以继续，关注警告
< 0.4    → 🔴 LOW   → 暂停，优先解决警告和阻塞
blockers >= 2 → 🛑 STOP → 停止，等待人类
```

---

> 注：详细CLI命令参考 [VIC-CLI-GUIDE.md](./docs/VIC-CLI-GUIDE.md)
