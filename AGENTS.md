# VIBE-SDD Agent Collaboration Guide

> **AI 入口文件** — 此文件定义 AI 进入项目时的起点。
> **详细执行步骤和 vic 命令调用 → 见各 SKILL.md**

---

## 系统状态

| 特性 | 状态 |
|------|------|
| 结构化开发 | ✅ .vic-sdd/ SPEC 工作流 |
| AI 自我认知 | ✅ context-tracker (auto_activate) |
| Gate 检查 | ✅ vic spec gate 0-3 |
| 规范约束 | ✅ constitution.yaml |
| 技能系统 | ✅ 5 Skills (Google Cloud Agent Skills 规范) |

---

## 技能一览 (5 个核心)

| Skill | 何时激活 | 职责 |
|-------|---------|------|
| **`context-tracker`** | **每次会话开始 + 每个行动后 + 会话结束** | AI 自我认知、信心度追踪、Blocker 识别 |
| **`spec-workflow`** | 需求模糊 / 架构设计 / SPEC 创建 | 需求分析 → 架构设计 → SPEC 冻结 |
| **`implementation`** | 代码实现 / Bug 修复 / 测试 / SPEC 对齐 | TDD 红绿重构、系统调试、Gate 2/3 检查 |
| **`unified-workflow`** | 功能交付 / 阶段推进 / 提交前 / 追溯检查 | SDD 状态机、Constitution 执行、Traceability |
| **`quick`** | 简单单文件改动（不涉及 SPEC） | Typo 修复、变量重命名、简单注释 |

---

## 决策树：何时用哪个技能？

```
AI 进入项目 → 确认上下文 → 执行工作 → 收尾

Step 1: 确认上下文 (context-tracker, auto_activate)
  → vic status
  → vic spec status
  → vic spec hash
  → vic gate check --blocking
  → 查看 .vic-sdd/ 状态文件
  (详细 → skills/context-tracker/SKILL.md)

Step 2: 判断任务类型
  │
  ├─ 🤔 需求模糊 / 架构未设计 / 需要 SPEC
  │   └─→ spec-workflow
  │       (详细 → skills/spec-workflow/SKILL.md)
  │
  ├─ 💻 代码实现 / Bug 修复 / 写测试 / 检查对齐
  │   └─→ implementation
  │       (详细 → skills/implementation/SKILL.md)
  │
  ├─ 🚀 功能交付 / 阶段推进 / 提交前 / 追溯
  │   └─→ unified-workflow
  │       (详细 → skills/unified-workflow/SKILL.md)
  │
  └─ 🔧 简单改动（单文件、无 SPEC 影响）
      └─→ quick
          (详细 → skills/quick/SKILL.md)
```

---

## 规划前置命令（项目启动 / 规划需求前必看）

> **这些是 AI 在开始任何实质性工作前应该先运行的命令。**
> **详细命令说明和参数 → 各 SKILL.md**

### 会话开始（每次对话第一件事）

```bash
vic status                              # 项目整体状态
vic spec status                         # SPEC 文档状态
vic spec hash                           # 检查 SPEC 是否变更
vic gate check --blocking               # 所有 Gate 状态（阻断性问题）
```

### 规划阶段（开始设计或澄清需求前）

```bash
vic spec list                           # 列出所有 SPEC 文档
vic spec show                           # 显示 SPEC 概要
vic milestone list                       # 项目里程碑
vic task list                           # 剩余任务（如果有）
```

### 状态查询（随时可用）

```bash
vic history --limit 10                  # 最近事件
vic search <关键词>                     # 搜索技术决策和风险
vic deps list                           # 模块依赖概览
vic cost status                         # Token/费用追踪
```

---

## SDD 状态机

```
Ideation → Explore → SpecCheckpoint → Build → Verify → ReleaseReady → Released
    │         │            │             │        │          │            │
    ▼         ▼            ▼             ▼        ▼          ▼            ▼
spec-workflow                   implementation              unified-workflow
                               (Gate 2: 代码对齐)            (Gate 3: 测试覆盖)
                               (Gate 3: 测试覆盖)            (最终交付检查)
```

---

## 质量红线（不可违反）

详见 `skills/context-tracker/SKILL.md` 和 `.vic-sdd/constitution.yaml`

| 规则 ID | 说明 | 触发 |
|---------|------|------|
| `SPEC-FIRST` | 改功能必须先改 SPEC | implementation |
| `SPEC-ALIGNED` | 代码必须对齐 SPEC | Gate 2 |
| `NO-TODO-IN-CODE` | 代码禁止 TODO/FIXME | Gate 0 |
| `NO-CONSOLE-IN-PROD` | 生产代码禁止 console.log | 提交前 |
| `GATE-BEFORE-COMMIT` | 提交前必须过 Gate | unified-workflow |
| `TESTS-REQUIRED` | 新功能必须有测试 | implementation |
| `SELF-AWARENESS` | 每步行动后更新 context | context-tracker |

---

## 信心度（context-tracker 自动计算）

```
confidence = (positive - warnings×0.3 - blockers×0.5) / max_signals

> 0.7    → 🟢 HIGH   → 继续
0.4-0.7  → 🟡 MODERATE → 继续，关注警告
< 0.4    → 🔴 LOW   → 暂停，解决阻塞
blockers >= 2 → 🛑 STOP → 停止，等待人类
```

---

## 目录结构（AI 必读文件）

```
.vic-sdd/
├── SPEC-REQUIREMENTS.md    # 需求规范（先读）
├── SPEC-ARCHITECTURE.md    # 架构规范（先读）
├── PROJECT.md               # 项目状态追踪
├── constitution.yaml        # 不可违反规则（先读）
├── context.yaml            # AI 自我认知状态（context-tracker 维护）
├── agent-prompt.md         # AI 工作流提示（含强制检查清单）
└── status/
    └── spec-hash.json      # SPEC 变更检测

skills/
├── context-tracker/        # AI 自我认知（auto_activate: true）
├── spec-workflow/          # 需求/架构/SPEC 创建
├── implementation/          # 代码/调试/测试/对齐
├── unified-workflow/        # SDD 编排/Constitution/追溯
└── quick/                 # 简单单文件改动
```

---

## 详细文档索引

| 场景 | 文档 |
|------|------|
| 我是谁 / 我该做什么 | AGENTS.md（此文件）|
| 每次行动后如何更新状态 | skills/context-tracker/SKILL.md |
| 需求模糊、架构设计、创建 SPEC | skills/spec-workflow/SKILL.md |
| 写代码、修复 Bug、测试、对齐 SPEC | skills/implementation/SKILL.md |
| 功能交付、阶段推进、Constitution、追溯 | skills/unified-workflow/SKILL.md |
| 简单单文件改动 | skills/quick/SKILL.md |
| CLI 工具完整命令参考 | docs/VIC-CLI-GUIDE.md |

---

> **核心原则**：AGENTS.md 是 AI 的"入口地图"，保持简洁。
> 详细的工作步骤、vic 命令调用链、具体参数 → 在激活对应 Skill 后加载对应 SKILL.md。
> 这样避免上下文爆炸，同时保证每个执行步骤都有据可查。
