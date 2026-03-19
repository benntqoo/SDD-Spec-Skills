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
    ↓
Bridge → Agent-Develop enters sdd-orchestrator
```

### Phase 3: 立规矩 (Implementation)

```
sdd-orchestrator (activates self-awareness at entry)
    ↓
    ├── knowledge-boundary    (query before routing)
    ├── pre-decision-check   (gate check before transition)
    └── signal-register      (record each state)
    ↓
spec-architect → spec-to-codebase → spec-contract-diff → spec-driven-test
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
│     → 记录探索 (exploration-journal)     │
└───────────────┬─────────────────────────┘
                ↓
         周期性信心度检查
         ├── confidence >= 0.7 → 继续
         ├── 0.4 <= confidence < 0.7 → 关注
         └── confidence < 0.4 → 暂停
```

---

## Self-Awareness Activation Protocol

**This is NOT optional.** Every task execution path MUST follow this protocol.

### The Four-Phase Loop

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│  PHASE 1: BEGIN (Every task starts here)                                │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │ skill:knowledge-boundary                                           │  │
│  │                                                                    │  │
│  │ INPUT : current task description                                   │  │
│  │ OUTPUT: categorized list (known/inferred/assumed/unknown)         │  │
│  │ FILE  : .vic-sdd/knowledge-boundary.yaml (read/write)             │  │
│  │                                                                    │  │
│  │ If unknown blocks this task → STOP, record blocker, ask human    │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│                               ↓                                          │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │ skill:pre-decision-check                                           │  │
│  │                                                                    │  │
│  │ INPUT : task decision points + knowledge-boundary output          │  │
│  │ OUTPUT: PASS / WARN / STOP / BLOCK result                        │  │
│  │ FILES : .vic-sdd/decision-guardrails.yaml (read)                   │  │
│  │         .vic-sdd/signal-register.yaml (read)                      │  │
│  │                                                                    │  │
│  │ If STOP or BLOCK → record signals, do not proceed                 │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│                               ↓                                          │
│  PHASE 2: EXECUTE (delegated to domain skills)                         │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │ skill:signal-register  ← invoked after EVERY meaningful action   │  │
│  │ skill:exploration-journal ← invoked after exploration/decision  │  │
│  │                                                                    │  │
│  │ positive/warnings/blockers → update signal-register.yaml          │  │
│  │ explored/tried/decided → update exploration-journal.yaml          │  │
│  │ RECALCULATE confidence after each signal update                   │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│                               ↓                                          │
│  PHASE 3: CHECKPOINT (after each subtask or milestone)                 │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │ skill:pre-decision-check (re-entry, lightweight)                   │  │
│  │                                                                    │  │
│  │ If confidence < 0.4 → pause, resolve warnings/blockers            │  │
│  │ If blockers >= 2 → STOP, ask human                               │  │
│  │ If all clear → continue to next subtask                           │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│                               ↓                                          │
│  PHASE 4: WRAP-UP (task complete)                                       │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │ skill:knowledge-boundary (final update)                           │  │
│  │ skill:signal-register (final confidence + summary)               │  │
│  │                                                                    │  │
│  │ Move inferred → known (if verified)                              │  │
│  │ Move assumed → inferred or known (if validated or confirmed)     │  │
│  │ Emit final confidence score                                       │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### Calling Order (Enforced)

| # | Skill | When | Read | Write |
|---|-------|------|------|-------|
| 1 | `knowledge-boundary` | Task BEGIN | `knowledge-boundary.yaml` | `knowledge-boundary.yaml` |
| 2 | `pre-decision-check` | Task BEGIN + CHECKPOINT | `decision-guardrails.yaml` + `signal-register.yaml` | `decision-guardrails.yaml` |
| 3 | `signal-register` | After every meaningful action | `signal-register.yaml` | `signal-register.yaml` |
| 4 | `exploration-journal` | After every explore/tried/decided/learned | `exploration-journal.yaml` | `exploration-journal.yaml` |

### Data Flow Between Skills

```
knowledge-boundary.yaml
       │
       ├──→ pre-decision-check (for scope/signal queries)
       │         │
       │         └──→ signal-register.yaml (blockers → signal blockers)
       │                   │
       │                   └──→ pre-decision-check (confidence check)
       │
       └──→ exploration-journal.yaml (known items vs journal findings)
                 │
                 └──→ knowledge-boundary.yaml (journal findings → new known)
```

### Integration with SDD Skills

SDD skills MUST activate self-awareness at these points:

| SDD Skill | Activation Point | Skills to Call |
|-----------|-----------------|----------------|
| `sdd-orchestrator` | Entry + each state transition | knowledge-boundary, pre-decision-check, signal-register |
| `spec-architect` | Entry | knowledge-boundary, pre-decision-check |
| `spec-to-codebase` | Entry + each file generated | signal-register |
| `spec-contract-diff` | Entry + each diff found | signal-register, exploration-journal |
| `spec-driven-test` | Entry + each test created | signal-register |
| `sdd-release-guard` | Final gate | pre-decision-check (final), signal-register (final summary) |

### Integration with Vibe Skills

Vibe skills MUST activate self-awareness at these points:

| Vibe Skill | Activation Point | Skills to Call |
|------------|-----------------|----------------|
| `vibe-think` | Entry | knowledge-boundary, exploration-journal |
| `vibe-architect` | Entry + each tech decision | knowledge-boundary, pre-decision-check, signal-register |
| `vibe-design` | Entry | knowledge-boundary, pre-decision-check |
| `vibe-debug` | Entry + each attempted fix | knowledge-boundary, exploration-journal, signal-register |
| `vibe-qa` | Entry | signal-register |

### Confidence Check Trigger Points

Calculate confidence AFTER adding signals. Check at these milestones:

1. **After Phase 1 (BEGIN)** — ensure task is viable
2. **After every 3-5 signals** — prevent drift
3. **Before any state promotion** (SPEC→SDD, SDD gate pass)
4. **After completing a skill** — before handing off
5. **Before asking human** — document current state first

### What NOT to Do

| Forbidden | Instead |
|-----------|---------|
| Skip knowledge-boundary when unfamiliar with domain | Stop, query, categorize |
| Proceed when confidence < 0.4 | Pause, resolve blockers |
| Continue when blockers >= 2 | STOP, ask human |
| Record vague signals like "progress" | Record specific evidence |
| Skip pre-decision-check before major decisions | Always check |
| Leave assumed items unverified indefinitely | Either verify or mark as warning |

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

### Self-Awareness Skills (Always Active)

| Skill | Purpose | When to Activate |
|-------|---------|------------------|
| `knowledge-boundary` | AI 自知之明：knows/infers/assumes/unknown | Task BEGIN + wrap-up |
| `pre-decision-check` | 决策前刹车：scope/quality/signals check | Task BEGIN + CHECKPOINT |
| `signal-register` | 证据链进度：positive/warnings/blockers → confidence | After every meaningful action |
| `exploration-journal` | 思考过程记忆：explore/tried/decided/learned | After every exploration/decision |

### SDD Skills (Feature Delivery)

| Skill | Purpose | State |
|-------|---------|-------|
| `sdd-orchestrator` | SDD entry point, state machine, gate enforcement | Entry of Phase 2 |
| `spec-architect` | Freezes requirements into contracts | Ideation/Explore |
| `spec-to-codebase` | Generates implementation from contracts | SpecCheckpoint |
| `spec-contract-diff` | Detects drift between spec and code | Build |
| `spec-driven-test` | Builds and enforces test gates from contracts | Build/Verify |
| `spec-traceability` | Story-to-contract-to-code-to-test mapping | Any state |
| `sdd-release-guard` | Final SDD release gates | ReleaseReady |

### Vibe Skills (Exploration & QA)

| Skill | Purpose |
|-------|---------|
| `vibe-think` | Requirements clarification, user story discovery |
| `vibe-architect` | Tech selection, architecture design, SPEC-ARCHITECTURE.md |
| `vibe-redesign` | Product redesign, scope re-evaluation |
| `vibe-design` | Design system, UI/UX specifications |
| `vibe-debug` | Systematic debugging with root cause analysis |
| `vibe-qa` | Quality assurance, verification against specs |
| `adaptive-planning` | Adaptive replanning when scope changes |

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
