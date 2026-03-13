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

## 验证

运行验证以确保完整性：

```bash
python skills/vibe-integrity/validate-vibe-integrity.py  # 检查 .vibe-integrity/ 文件
python skills/vibe-integrity/validate-all.py             # 运行 vibe-guard 和 vibe-integrity 双重验证
python skills/vibe-guard/validate-vibe-guard.py --check  # AI 完成检查
```

## 相关技能

- `vibe-guard` - 完成检测
- `superpowers/test-driven-development` - TDD 工作流（可选）
- `sdd-orchestrator` - SDD 工作流（可选）

**注意**：Vibe Integrity 适用于 ANY 开发方法。您可以单独使用 Vibe Integrity，或者将其与 SDD、TDD、敏捷或任何其他方法结合使用。上述列出的 SDD 和 TDD 技能是可选的附加功能，供希望在仍然受益于 Vibe Integrity 的完成守卫和项目记忆的同时遵循这些特定方法的团队使用。

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