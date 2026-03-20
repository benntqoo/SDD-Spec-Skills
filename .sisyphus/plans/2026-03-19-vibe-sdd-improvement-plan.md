# VIBE-SDD 改进计划：使其真正可用

> 生成时间: 2026-03-19
> 状态: ✅ 全部完成 (2026-03-20)

---

## 问题诊断总结

### 审计结果：vic-go CLI 真实状态

| 命令模块 | 状态 | 说明 |
|---------|------|------|
| `init` | ✅ 已实现 | 目录结构初始化 |
| `spec init` | ✅ 已实现 | SPEC 文档生成 |
| `spec status` | ✅ 已实现 | 显示文档存在性 |
| `gate status` | ✅ 已实现 | 显示 Gate 状态 |
| `gate pass` | ✅ 已实现 | 标记 Gate 通过 |
| `gate check` | ⚠️ 部分 | 仅显示状态，无实际检查 |
| `spec gate` | ❌ 空壳 | 声称检查但打印占位符 |
| `check` | ⚠️ 部分 | 仅技术检测，无对齐验证 |
| `record` | ✅ 已实现 | 记录技术决策/风险 |
| `auto` | ✅ 已实现 | 自动执行模式 |
| `tdd` | ✅ 已实现 | TDD 循环 |
| `debug` | ✅ 已实现 | 调试流程 |
| `slop` | ✅ 已实现 | AI 垃圾检测 |
| `qa` | ✅ 已实现 | E2E 测试 |

### 关键发现

1. **Gate 检查是空壳**：`spec gate 0-3` 只打印 "not implemented yet"
2. **技能系统与 CLI 脱节**：19 个 SKILL.md 没有对应的 CLI 命令
3. **流程无法强制**：没有任何机制阻止 AI 直接写代码

---

## 改进计划

### 阶段 1：实现真实 Gate 检查 (高优先级)

#### 1.1 Gate 0：需求完整性检查

**目标**：验证 SPEC-REQUIREMENTS.md 包含必要章节

**实现要点**：
```go
// 检查项
- [ ] 包含 ## User Stories 章节
- [ ] 包含 ## Key Features 章节
- [ ] 每个 Feature 有验收标准
- [ ] 包含非功能性需求
- [ ] 无未完成的 TBD/TODO
```

**文件位置**：`cmd/vic-go/internal/commands/gate0.go` (新增)

#### 1.2 Gate 1：架构完整性检查

**目标**：验证 SPEC-ARCHITECTURE.md 包含技术决策

**实现要点**：
```go
// 检查项
- [ ] 包含 ## Technology Stack
- [ ] 包含 ## System Design
- [ ] 包含 ## API Design
- [ ] 技术选型有选择理由
- [ ] 无 TBD/TODO
```

**文件位置**：`cmd/vic-go/internal/commands/gate1.go` (新增)

#### 1.3 Gate 2：代码对齐检查

**目标**：验证代码实现与 SPEC 一致

**实现要点**：
```go
// 检查项
- [ ] SPEC 声明的技术栈在代码中存在
- [ ] SPEC 定义的 API 端点已实现
- [ ] SPEC 定义的模块已创建
- [ ] 记录的技术决策与代码一致
```

**文件位置**：扩展 `cmd/vic-go/internal/checker/code_analysis.go`

#### 1.4 Gate 3：测试覆盖检查

**目标**：验证关键路径有测试覆盖

**实现要点**：
```go
// 检查项
- [ ] 存在测试文件
- [ ] 关键模块有测试
- [ ] 使用已知测试框架 (jest, pytest, go test 等)
```

**文件位置**：`cmd/vic-go/internal/commands/gate3.go` (新增)

---

### 阶段 2：精简 Skill 系统 (中优先级)

#### 2.1 合并 Self-Awareness 技能 (4 → 1)

**当前**：
- `knowledge-boundary`
- `pre-decision-check`
- `signal-register`
- `exploration-journal`

**合并为**：`context-tracker`
- 单一文件 `.vic-sdd/context.yaml`
- 统一的上下文追踪格式

#### 2.2 合并 Vibe Exploration 技能 (7 → 3)

**当前**：
- `vibe-think` → 保留，改名 `requirements`
- `vibe-architect` → 保留，改名 `architecture`
- `vibe-design` → 保留 `design-review`
- `vibe-redesign` → 合并到 `requirements`
- `vibe-debug` → 合并到 `debugging`
- `vibe-qa` → 合并到 `qa`
- `adaptive-planning` → 合并到 `planning`

**合并后**：
- `requirements` - 需求分析
- `architecture` - 架构设计
- `design-review` - 设计审查

#### 2.3 保留 SDD Core 技能 (7 → 5)

**当前**：
- `sdd-orchestrator` → 保留
- `spec-architect` → 保留
- `spec-to-codebase` → 合并到 `sdd-orchestrator`
- `spec-contract-diff` → 保留
- `spec-driven-test` → 合并到 `qa`
- `spec-traceability` → 保留
- `sdd-release-guard` → 合并到 `sdd-orchestrator`

**合并后**：
- `sdd-orchestrator` - 状态机 + 流程控制
- `spec-architect` - 规范编写
- `spec-contract-diff` - 差异检测
- `spec-traceability` - 追溯追踪
- `qa` - 测试覆盖

#### 2.4 精简结果

| 类别 | 原来 | 精简后 |
|------|------|--------|
| Self-Awareness | 4 | 1 |
| Vibe | 7 | 3 |
| SDD | 7 | 5 |
| TDD | 1 | 1 |
| **总计** | **19** | **10** |

---

### 阶段 3：强制流程钩子 (高优先级)

#### 3.1 pre-commit 钩子

**文件**：`.pre-commit-config.yaml`

```yaml
# 检查 Gate 状态才能提交
- repo: local
  hooks:
    - id: vic-gate-check
      name: VIBE-SDD Gate Check
      entry: vic gate check --phase
      language: system
      pass_frozen: false
```

#### 3.2 Agent 提示词钩子

**文件**：`.vic-sdd/agent-prompt.md`

在每个会话开始时显示：
```
VIBE-SDD 工作流程:
1. vic init (初始化)
2. vic spec init (创建 SPEC)
3. vic spec gate 0 (验证需求)
4. vic spec gate 1 (验证架构)
5. ... 实现功能 ...
6. vic spec gate 2 (验证对齐)
7. vic spec gate 3 (验证测试)
```

#### 3.3 检查点验证

**文件**：`cmd/vic-go/internal/commands/checkpoint.go`

```go
// 在执行 phase advance 前验证
func ValidateCheckpoint(cfg *Config, fromPhase int) error {
    // 必须先通过当前 phase 的所有 gate
    for i := 0; i <= fromPhase; i++ {
        if !passed, _ := gate.ValidateGateCheck(cfg, i); !passed {
            return fmt.Errorf("Phase %d gates not passed", i)
        }
    }
    return nil
}
```

---

### 阶段 4：文档与代码同步 (中优先级)

#### 4.1 SPEC 变更检测

```go
// 监控 SPEC-ARCHITECTURE.md 的技术栈声明
// 如果变更，提示重新检查代码对齐

func CheckTechDrift(cfg *Config) error {
    // 1. 读取 SPEC 声明的技术栈
    // 2. 扫描代码检测实际技术栈
    // 3. 对比差异
    // 4. 报告不一致
}
```

#### 4.2 变更日志

```go
// 当 SPEC 变更时记录
// .vic-sdd/status/change-log.yaml

type ChangeLog struct {
    Date      string
    File      string
    Changes   []string
    Impact    string
    Reviewed  bool
}
```

---

## 实施顺序

```
Week 1:
  ├── 实现 Gate 0 检查
  ├── 实现 Gate 1 检查
  └── 更新 spec gate 命令

Week 2:
  ├── 实现 Gate 2 检查
  ├── 实现 Gate 3 检查
  └── 集成 pre-commit 钩子

Week 3:
  ├── 精简 Skill 系统 (19 → 10)
  ├── 更新 AGENTS.md
  └── 简化 SDD 流程文档

Week 4:
  ├── 实现强制流程钩子
  ├── 实现 SPEC 变更检测
  └── 编写使用指南
```

---

## 验收标准

### CLI 功能
- [x] `vic spec gate 0` 实际检查 SPEC-REQUIREMENTS.md 结构
- [x] `vic spec gate 1` 实际检查 SPEC-ARCHITECTURE.md 结构
- [x] `vic spec gate 2` 实际检查代码与 SPEC 对齐
- [x] `vic spec gate 3` 实际检查测试覆盖
- [x] `vic phase advance` 在 Gate 未通过时拒绝推进

### Skill 系统
- [x] Skill 从 19 个减少到 10 个
- [x] 每个 Skill 有 CLI 命令对应
- [x] SKILL.md 文档简洁明了

### 流程强制
- [x] pre-commit 钩子在 Gate 未通过时阻止提交
- [x] 新会话开始时显示流程提示
- [x] SPEC 变更时自动检测技术漂移

---

## 文件变更清单

### 新增文件
- `cmd/vic-go/internal/commands/gate0.go`
- `cmd/vic-go/internal/commands/gate1.go`
- `cmd/vic-go/internal/commands/gate2.go`
- `cmd/vic-go/internal/commands/gate3.go`
- `cmd/vic-go/internal/commands/checkpoint.go`
- `cmd/vic-go/internal/checker/spec_analyzer.go`
- `.vic-sdd/agent-prompt.md`
- `.sisyphus/plans/2026-03-19-improvement-plan.md`

### 修改文件
- `cmd/vic-go/internal/commands/spec.go` (移除空壳)
- `cmd/vic-go/internal/commands/phase.go` (添加验证)
- `AGENTS.md` (简化流程)
- `.pre-commit-config.yaml` (添加钩子)

### 删除/合并文件
- skills/knowledge-boundary → 合并到 context-tracker
- skills/pre-decision-check → 合并到 context-tracker
- skills/signal-register → 合并到 context-tracker
- skills/exploration-journal → 合并到 context-tracker
- skills/vibe-redesign → 合并到 requirements
- skills/adaptive-planning → 合并到 planning
- skills/vibe-qa → 合并到 qa
- skills/test-driven-development → 合并到 qa
- skills/spec-to-codebase → 合并到 sdd-orchestrator
- skills/sdd-release-guard → 合并到 sdd-orchestrator

---

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| Gate 检查太严格 | 用户体验差 | 提供 `--force` 选项 |
| Skill 合并丢失功能 | 工具能力下降 | 确保合并后的 Skill 功能完整 |
| pre-commit 钩子被跳过 | 流程失效 | 在 AGENTS.md 中强调重要性 |

---

**下一步**：开始实现 Gate 0 检查
