# Skills Decision Tree

> **快速参考** — 30 秒内找到正确的技能

---

## 入口决策

```
About to plan, commit, or make a major decision?
├─ ✅ Run context-tracker FIRST (每次行动后都要)
└─ Then proceed with your task

What is your task?
│
├─ 🤔 需求模糊 / 架构未设计 / 需要 SPEC
│   └─→ spec-workflow
│
├─ 💻 代码实现 / Bug 修复 / 写测试 / 检查对齐
│   └─→ implementation
│
├─ 🚀 功能交付 / 阶段推进 / 提交前 / 追溯检查
│   └─→ unified-workflow
│
└─ 🔧 简单单文件改动（无 SPEC 影响）
    └─→ quick
```

---

## SDD 状态机子决策树

进入 SDD 流程时（由 unified-workflow 触发）：

```
当前状态?
│
├─ Ideation / Explore
│   └─→ spec-workflow (需求分析 → SPEC 冻结)
│
├─ SpecCheckpoint
│   └─→ spec-workflow (运行 vic spec gate 0/1 验证完整性)
│
├─ Build (实现阶段)
│   ├─ vic spec gate 2 → 检查代码对齐
│   └─ vic check → 验证技术选型
│
├─ Verify (验证阶段)
│   ├─ vic spec gate 3 → 检查测试覆盖
│   └─ vic slop scan → 检查 AI Slop
│
└─ ReleaseReady
    └─→ unified-workflow (最终门控检查)
```

---

## 快速参考卡

| 情形 | 技能 | 关键问题 |
|------|------|---------|
| 会话开始/每次行动后 | `context-tracker` | 信心度多少？blocker？ |
| 需求模糊、架构设计 | `spec-workflow` | SPEC 是否完整？ |
| 写代码、调试、测试 | `implementation` | Gate 2/3 通过了吗？ |
| 提交前、阶段推进 | `unified-workflow` | Constitution 满足？ |
| 简单单文件改动 | `quick` | 真的不需要 SPEC？ |

---

## Skill 职责边界

```
┌─────────────────────────────────────────────────────────────┐
│                    context-tracker                            │
│         (auto_activate: true — 始终激活)                    │
│              信心度、blocker、上下文更新                      │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                     spec-workflow                             │
│          (需求 → 架构 → SPEC 冻结)                           │
│              何时用：需求不清晰或需要架构设计                   │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                     implementation                             │
│           (代码实现 / 调试 / 测试 / Gate 2/3)                 │
│              何时用：SPEC 已冻结，开始实现                      │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                    unified-workflow                           │
│         (SDD 状态机 / Constitution / 提交前 / 追溯)           │
│          何时用：功能交付、阶段推进、pre-commit                │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                       quick                                   │
│              (单文件、无 SPEC 影响)                            │
│              何时用：typo、rename、简单注释                   │
└─────────────────────────────────────────────────────────────┘
```

---

## 常见错误

| 错误 | 正确做法 |
|------|---------|
| 单文件改动却用 SDD | 直接用 `quick` |
| 实现前不检查 SPEC | 先用 `spec-workflow` |
| 调试前不读 context-tracker | 始终激活 `context-tracker` |
| 盲目实现不验证 Gate | `implementation` 后必须跑 `vic spec gate 2` |
| 提交前跳过 Constitution | `unified-workflow` 是强制门控 |

---

## 文件位置

此决策树可在以下位置找到：
- `AGENTS.md` — AI 入口（简化版）
- `skills/SKILL_DECISION.md` — 详细版（本文件）
- 激活任意 Skill 后，对应的 SKILL.md 包含完整的 vic 命令调用链
