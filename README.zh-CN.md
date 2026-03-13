# Vibe Integrity

[English README](./README.md)

Vibe Integrity 是一个专门为 AI 辅助开发（vibe coding）设计的 **AI 项目记忆与安全系统**。它能防止 AI 编码助手虚假声称完成，并提供结构化的项目知识以实现 AI 的快速理解。

## 概述

Vibe Integrity 解决了 AI 辅助开发中的两个关键问题：

1. **完成守卫** - 检测 AI 是否虚假声称工作已完成（TODO/FIXME 占位符、空函数、假测试等）
2. **架构记忆** - 提供结构化的项目知识，使 AI 能快速理解项目状态而无需阅读数百个文件

与传统开发方法不同，Vibe Integrity 是 **方法论无关** 的 - 它适用于 TDD、SDD、敏捷或纯 vibe 编程方法。

## 核心概念

### 两大支柱

#### 支柱 1：完成守卫
检测和验证以确保 AI 实际完成了工作。

| 技能 | 目的 |
|------|------|
| `vibe-guard` | 检测 TODO、空函数、假测试 |
| `cascade-check` | 防止修复后的级联错误 |
| `integration-check` | 验证组件集成 |

#### 支柱 2：架构记忆
用于 AI 快速理解的结构化项目知识库。

| 文件 | 目的 |
|------|------|
| `project.yaml` | 项目元信息，技术栈 |
| `dependency-graph.yaml` | 模块依赖关系 |
| `module-map.yaml` | 目录结构 |
| `risk-zones.yaml` | 高风险区域 |
| `tech-records.yaml` | 技术决策记录 |
| `schema-evolution.yaml` | 数据模型演进 |

## AI 快速开始

当 AI 开始在这个项目上工作时，请按此顺序阅读：

```
1. .vibe-integrity/project.yaml
   → 了解项目状态和技术栈

2. .vibe-integrity/risk-zones.yaml  
   → 了解哪些区域是高风险的

3. .vibe-integrity/dependency-graph.yaml
   → 了解模块关系

4. .vibe-integrity/module-map.yaml
   → 查找文件位置

5. .vibe-integrity/tech-records.yaml
   → 理解系统为何如此设计
```

**结果**：AI 能在大约 15 秒内理解项目，而不是 3 分钟。


## Base Workflow

Vibe Integrity 提供完整的 AI 辅助开发工作流程：

```
[用户提出需求] → vibe-design (需求澄清/苏格拉底提问)
                                    ↓
                        [做出架构决策] → vibe-integrity-writer (自动更新 tech-records.yaml)
                                    ↓
                        [完成设计] → 生成 .vibe-integrity/ 更新
                                    ↓
                [实现过程] → vibe-integrity-debug (发现问题时进行根因分析)
                                    ↓
                        [发现新风险/决策] → vibe-integrity-writer (自动更新 risk-zones.yaml/tech-records.yaml)
                                    ↓
                        [实现完成] → vibe-guard (验证完整性)
                                    ↓
                        [验证通过] → 开发完成
```

### 工作流程阶段：

| 阶段 | 主要技能 | 辅助技能 | 目的 |
|-------|---------------|------------------|---------|
| **澄清阶段** | `vibe-design` | `vibe-integrity-writer` | 理解需求，捕获决策 |
| **实现阶段** | 开发者/AI | `vibe-integrity-debug` | 构建功能，调试问题 |
| **验证阶段** | `vibe-guard` | `validate-vibe-integrity.py` | 确保完成和完整性 |
| **发现阶段** | `vibe-integrity-debug` | `vibe-integrity-writer` | 根因分析，更新项目记忆 |

## YAML 文件职责

`.vibe-integrity/` 中的每个 YAML 文件都有特定职责，由不同的技能维护：

| 文件 | 职责 | 维护者 | 触发事件 |
|------|------|--------|----------|
| `project.yaml` | 项目元信息，技术栈，状态 | `vibe-design` | 项目创建，范围变更，技术栈更新 |
| `tech-records.yaml` | 技术决策及其理由 | `vibe-design` → `vibe-integrity-writer` | 架构决策，技术选择，实现模式 |
| `dependency-graph.yaml` | 模块/服务依赖关系 | `vibe-design` → `vibe-integrity-writer` | 添加/移除模块，变更依赖，服务集成 |
| `module-map.yaml` | 目录结构和文件组织 | `vibe-design` → `vibe-integrity-writer` | 重组，新目录，文件迁移 |
| `risk-zones.yaml` | 识别的风险和高风险区域 | `vibe-integrity-debug` → `vibe-integrity-writer` | Bug 发现，安全问题，性能瓶颈，架构缺陷 |
| `schema-evolution.yaml` | 数据模型变更和迁移 | `vibe-design` → `vibe-integrity-writer` | 数据库模式变更，API 模型更新，数据结构演进 |

### 自动更新规则：

| 工作流阶段 | 涉及技能 | 更新的 YAML 文件 | 示例 |
|----------------|-----------------|-------------------|---------|
| **需求澄清** | `vibe-design` → `vibe-integrity-writer` | `tech-records.yaml`, `project.yaml` | 用户确认 PostgreSQL 代替 MongoDB → 记录决策 |
| **架构讨论** | `vibe-design` → `vibe-integrity-writer` | `tech-records.yaml`, `dependency-graph.yaml` | 讨论模块边界 → 更新依赖图 |
| **风险识别** | `vibe-integrity-debug` → `vibe-integrity-writer` | `risk-zones.yaml`, `tech-records.yaml` | 发现紧耦合 → 记录为风险区域 |
| **模式变更** | `vibe-design` → `vibe-integrity-writer` | `schema-evolution.yaml`, `tech-records.yaml` | 添加新表 → 记录模式演进 |
| **调试发现** | `vibe-integrity-debug` → `vibe-integrity-writer` | `risk-zones.yaml`, `tech-records.yaml` | 发现架构问题 → 记录洞察 |

### 关键工作流程原则：

1. **AI 维护记忆**：无需手动编辑 YAML - AI 在正常工作流程中自动更新项目记忆
2. **实时记录决策**：决策在做出时立即捕获，而非事后补充
3. **安全更新**：所有 YAML 修改通过 `vibe-integrity-writer` 进行，具备备份和验证
4. **审计追踪**：每次更改都被跟踪且可撤销
5. **交叉引用**：YAML 文件相互引用以保持一致性

## 使用方法

### AI：在进行更改之前

```bash
# 1. 检查风险区
cat .vibe-integrity/risk-zones.yaml

# 2. 检查依赖
cat .vibe-integrity/dependency-graph.yaml

# 3. 检查模式
cat .vibe-integrity/schema-evolution.yaml
```

### AI：在"完成"之后

```bash
# 运行 vibe-guard
python skills/vibe-guard/validate-vibe-guard.py --check
```

### 人类：在进行重大更改之后

```bash
# 更新技术记录
python skills/vibe-integrity/validate-vibe-integrity.py  # 首先检查完整性

# 向 .vibe-integrity/tech-records.yaml 添加新决策
# 向 .vibe-integrity/schema-evolution.yaml 添加新版本  
# 在 .vibe-integrity/dependency-graph.yaml 中反映新的模块关系
```

## 目录结构

```
.vibe-integrity/
├── project.yaml              # 项目元信息
├── dependency-graph.yaml     # 模块依赖关系
├── module-map.yaml          # 目录结构
├── risk-zones.yaml          # 高风险区域
├── tech-records.yaml        # 技术决策记录
└── schema-evolution.yaml   # 数据模型演进

skills/
├── vibe-guard/             # 完成检测
└── vibe-integrity/         # 此技能
    ├── SKILL.md
    ├── validate-vibe-integrity.py
    ├── validate-all.py
    └── template/           # Schema 模板
        ├── project.schema.json
        ├── dependency-graph.schema.json
        ├── module-map.schema.json
        ├── risk-zones.schema.json
        ├── tech-records.schema.json
        └── schema-evolution.schema.json
```


## 目录结构与技能

```
.vibe-integrity/
├── project.yaml              # 项目元信息
├── dependency-graph.yaml     # 模块依赖关系
├── module-map.yaml          # 目录结构
├── risk-zones.yaml          # 高风险区域
├── tech-records.yaml        # 技术决策记录
└── schema-evolution.yaml   # 数据模型演进

skills/
├── vibe-guard/              # 完成检测 (支柱 1)
├── vibe-integrity/          # 验证框架
│   ├── SKILL.md
│   ├── validate-vibe-integrity.py
│   ├── validate-all.py
│   └── template/           # Schema 模板
│       ├── project.schema.json
│       ├── dependency-graph.schema.json
│       ├── module-map.schema.json
│       ├── risk-zones.schema.json
│       ├── tech-records.schema.json
│       └── schema-evolution.schema.json
├── vibe-design/             # 需求澄清 (支柱 2)
│   ├── SKILL.md
│   └── 使用 vibe-integrity-writer 进行更新
├── vibe-integrity-debug/    # 系统调试
│   ├── SKILL.md
│   └── 使用 vibe-integrity-writer 记录发现
└── vibe-integrity-writer/   # YAML 文件写入器 (支柱 2)
    ├── SKILL.md
    └── 处理所有 .vibe-integrity/ 更新
```

### 技能职责：

| 技能 | 支柱 | 用途 | 输出 |
|-------|--------|---------|----------|
| `vibe-guard` | 完成守卫 | 检测虚假完成声明 | 验证报告 |
| `vibe-design` | 架构记忆 | 澄清需求，捕获决策 | 更新 .vibe-integrity/ 文件 |
| `vibe-integrity-debug` | 架构记忆 | 根因分析，风险识别 | 更新 .vibe-integrity/ 文件 |
| `vibe-integrity-writer` | 架构记忆 | 安全 YAML 更新与验证 | 更新 .vibe-integrity/ 文件 |
| `vibe-integrity` | 两者 | 验证结构与完整性 | 验证报告 |

## 验证

运行验证以确保完整性：

```bash
python skills/vibe-integrity/validate-vibe-integrity.py  # 检查 .vibe-integrity/ 文件
python skills/vibe-integrity/validate-all.py             # 运行 vibe-guard 和 vibe-integrity 双重验证
python skills/vibe-guard/validate-vibe-guard.py --check  # AI 完成检查
```

## 相关技能

- `vibe-guard` - 完成检测
- `vibe-design` - 需求澄清和设计助手，使用 vibe-integrity-writer 更新项目内存
- `vibe-integrity-debug` - 系统调试助手，确保修复前进行根因分析
- `superpowers/test-driven-development` - TDD 工作流（可选）
- `sdd-orchestrator` - SDD 工作流（可选）
- `vibe-integrity-writer` - 专门用于安全更新 .vibe-integrity/ YAML 文件的技能（由其他技能调用）

**注意**：Vibe Integrity 适用于 ANY 开发方法。您可以单独使用 Vibe Integrity，或者将其与 SDD、TDD、敏捷或任何其他方法结合使用。上述列出的 SDD 和 TDD 技能是可选的附加功能，供希望在仍然受益于 Vibe Integrity 的完成守卫和项目记忆的同时遵循这些特定方法的团队使用。

**工作流程总结**:
- **vibe-design**: 澄清需求，做出决策，调用 writer 更新 YAML
- **vibe-integrity-debug**: 执行根因分析，识别风险，调用 writer 更新 YAML
- **vibe-integrity-writer**: 安全更新 .vibe-integrity/ YAML 文件，具备备份和验证

## 多智能体协作

Vibe Integrity 支持单智能体和多智能体使用，具有特定的协作机制：

### 当前协作能力

| 机制 | 描述 | 用例 |
|-----------|-------------|----------|
| **Git 分支** | 不同智能体/会话使用独立分支 | 多智能体并行开发 |
| **Union 合并** | `.gitattributes` 中的 `.vibe-integrity/*.yaml merge=union` | 避免数据丢失 |
| **验证脚本** | 检测重复 ID 和不一致性 | 早期冲突检测 |
| **宽限期** | 10 分钟去重窗口 | 防止冗余写入 |
| **会话隔离** | 不同 git worktree 用于不同工作流 | 避免工作流混淆 |

### 文件锁定与并发

**当前实现：**
- ✅ 修改前创建备份
- ✅ 原子批量操作
- ✅ 更新后验证
- ✅ **文件级锁定**（已在 vibe-integrity-writer 中实现）
- ✅ **智能体身份追踪**（所有更新包含智能体元数据）
- ✅ **冲突检测**（检测重复 ID 和相似决策）
- ❌ 无分布式锁定（用于多实例部署）

**建议：** 对于多智能体场景，使用新的文件锁定和冲突检测功能。每个智能体仍应尽可能使用独立分支，但系统现在支持带适当锁定的并发写入。

**当前实现：**
- ✅ 修改前创建备份
- ✅ 原子批量操作
- ✅ 更新后验证
- ❌ 无文件级锁定（建议顺序使用）
- ❌ 无分布式锁定（用于多实例部署）

**建议：** 对于多智能体场景，使用独立分支并通过 Git 工作流协调，而非同时写入相同文件。

### 冲突解决

当多个智能体修改相同的 .vibe-integrity/ 文件时：

1. **Git 会在合并/PR 时检测冲突**
2. **Union 合并**保留两个版本（可能产生重复）
3. **验证脚本**识别重复 ID
4. **手动解决**以合并相似决策

### 多智能体最佳实践

1. **使用独立分支**：每个智能体使用自己的功能分支
2. **协调架构决策**：在共享记录中记录重要决策
3. **定期验证**：频繁运行验证脚本
4. **PR 审查**：合并前审查 .vibe-integrity/ 变更
5. **冲突解决**：使用验证脚本检测和修复重复

### 多实例考虑

对于多服务器部署：

- **文件锁定**：在 vibe-integrity-writer 中实现文件级锁
- **分布式锁定**：多服务器使用 Redis 等工具
- **租约机制**：实现 TTL 防止死锁
- **降级策略**：锁不可用时优雅降级

详见 [MULTI_AGENT_COLLABORATION.md](skills-base/vibe-integrity/MULTI_AGENT_COLLABORATION.md) 获取详细分析。
- **vibe-guard**: 实现完成后验证 AI 完成情况
- **vibe-integrity**: 验证 .vibe-integrity/ 目录结构

**工作流程分离指南**：为了避免在使用多种方法时产生混淆：
- 在澄清阶段使用 `vibe-design` 来理解需求并记录决策
- 在实现声明完成后使用 `vibe-guard` 来验证完成情况
- 在调查问题时使用 `vibe-integrity-debug` 进行根因分析后再修复
- 当遵循 TDD 或 SDD 时，考虑使用单独的工作会话或在每个会话开始时明确声明工作流程
- Vibe Design 在决策制定过程中会自动更新 `.vibe-integrity/` 文件，减少手动文档编写的负担

**重要提示**：在没有明确上下文切换的情况下，请不要将 `vibe-design` 与活动的 TDD/SDD 实施会话混合使用，以防止 AI 对正在遵循的工作流程产生混淆。

## 快速开始

1) 运行默认验证（扫描 `<root>/skills`）：

```bash
python skills/vibe-integrity/validate-all.py
```

2) 在您的项目中初始化 Vibe Integrity：

```bash
# 创建带模板文件的 .vibe-integrity 目录
python skills/vibe-integrity/validate-vibe-integrity.py --init

# 或手动复制模板文件：
cp -r skills/vibe-integrity/template/* .vibe-integrity/
```

3) 为您的项目自定义文件：
   - 编辑 `.vibe-integrity/project.yaml` 以填写您的项目详情
   - 更新 `.vibe-integrity/tech-records.yaml` 以包含您的技术决策
   - 自定义 `.vibe-integrity/risk-zones.yaml` 以适用于您项目的风险区域

## 示例输出

一次成功的验证运行看起来像这样：

```text
Vibe Integrity 验证通过
根目录: D:\Code\aaa
已检查的文件:
- .vibe-integrity/project.yaml ✓
- .vibe-integrity/dependency-graph.yaml ✓
- .vibe-integrity/module-map.yaml ✓
- .vibe-integrity/risk-zones.yaml ✓
- .vibe-integrity/tech-records.yaml ✓
- .vibe-integrity/schema-evolution.yaml ✓

Vibe Guard 验证:
- TODO/FIXME 检查: 通过
- 空函数检查: 通过
- 假测试检查: 通过
- 构建成功: 通过
- 类型检查: 通过
- 代码规范检查: 通过
- 安全检查: 通过
- 测试真实性: 通过

所有验证均已通过
```

如果显示 `Vibe Integrity 验证通过`，则表示所有文件均存在且结构有效。

## 配置

Vibe Integrity 使用 `.vibe-integrity/` 目录中的 YAML 文件进行配置。

### project.yaml
```yaml
name: my-project
version: 0.1.0
status: mvp
description: "我的惊人项目"
created_at: 2026-01-15
last_updated: 2026-03-12
tech_stack:
  前端: [Vue, Vite]
  后端: [Express, Node]
  数据库: [SQLite]
```

### tech-records.yaml
```yaml
records:
  - id: DB-001
    日期: "2026-01-15"
    类别: database
    标题: "选择 SQLite 作为 MVP"
    决定: "使用 SQLite 实现快速迭代"
    原因: "MVP 阶段优先考虑速度而非可扩展性"
    影响: 低
    状态: 已完成
```

## 常见操作

### 初始化新项目结构
```bash
python skills/vibe-integrity/validate-vibe-integrity.py --init
```

### 验证完整性
```bash
python skills/vibe-integrity/validate-all.py
```

### AI 完成检查
```bash
python skills/vibe-guard/validate-vibe-guard.py --check
```

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](./LICENSE)。