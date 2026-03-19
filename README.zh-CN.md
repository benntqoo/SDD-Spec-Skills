# VIBE-SDD

[English](./README.md)

VIBE-SDD 是一个结合了结构化 SDD (Spec-Driven Development) 与灵活 Vibe Coding 的**Vibe 驱动软件开发系统**。它为 AI 辅助开发提供了完整的流程，包含规范的检查点和文档管理。

## 概述

VIBE-SDD 解决了 AI 辅助开发中的三个关键问题：

1. **规范** - 结构化的需求和架构文档
2. **门禁** - 推进前的质量检查点
3. **记忆** - 项目知识供 AI 快速理解

## 快速开始

```bash
# 初始化项目
vic init --name "My Project" --tech "React,Node,PostgreSQL"

# 初始化 SPEC 文档
vic spec init --name "My Project"

# 记录技术决策
vic rt --id DB-001 --title "Use PostgreSQL" --decision "Primary database" --reason "Need ACID"

# 查看 SPEC 状态
vic spec status

# 运行 Gate 检查
vic spec gate 0  # 需求完整性
vic spec gate 1  # 架构完整性

# 验证
vic validate
```

## 命令

| 命令 | 别名 | 描述 |
|------|------|------|
| `vic init` | - | 初始化 .vic-sdd/ |
| `vic spec init` | - | 初始化 SPEC 文档 |
| `vic spec status` | - | 查看 SPEC 状态 |
| `vic spec gate [0-3]` | - | 运行 Gate 检查 |
| `vic rt` | `record-tech` | 记录技术决策 |
| `vic rr` | `record-risk` | 记录风险 |
| `vic rd` | `record-dep` | 记录依赖 |
| `vic check` | - | 检查代码对齐 |
| `vic validate` | - | 完整验证 |
| `vic status` | - | 查看项目状态 |
| `vic search` | - | 搜索记录 |
| `vic history` | - | 查看历史 |
| `vic export` | - | 导出数据 |
| `vic import` | - | 导入数据 |

完整文档：[cmd/vic/README.md](./cmd/vic/README.md)

## 开发流程

```
定图纸 (需求)              打地基 (架构)              立规矩 (实现)
    │                         │                        │
vibe-think            vibe-architect         sdd-orchestrator
    │                         │                        │
    ▼                         ▼                        ▼
SPEC-REQUIREMENTS.md ─▶ SPEC-ARCHITECTURE.md ─▶ 实现代码
    │                         │                        │
    ▼                         ▼                        ▼
   Gate 0                  Gate 1                  Gate 2 + 3
(需求完整)              (架构完整)              (代码 + 测试)
                                                        │
                                                        ▼
                                              收敛到 PRD/ARCH/PROJECT
```

## 目录结构

```
project/
├── cmd/
│   └── vic/                    # CLI 工具
│       ├── vic                  # 主程序
│       ├── README.md           # 英文文档
│       └── *.py                # 脚本
│
├── skills/                         # 共 19 个 skills
│   │
│   ├── 自我认知 (4):              # AI 自知之明机制
│   │   ├── knowledge-boundary/ # 知道/推断/假设/不知
│   │   ├── pre-decision-check/ # 决策前门禁检查
│   │   ├── signal-register/    # 证据链进度追踪
│   │   └── exploration-journal/ # 探索过程记忆
│   │
│   ├── Vibe 探索 (7):             # 灵活探索流程
│   │   ├── vibe-think/         # 需求澄清
│   │   ├── vibe-architect/      # 技术选型 + 架构设计
│   │   ├── vibe-redesign/      # 产品重新设计
│   │   ├── vibe-design/        # 设计系统
│   │   ├── vibe-debug/         # 系统性调试
│   │   ├── vibe-qa/            # 质量保证
│   │   └── adaptive-planning/  # 自适应重规划
│   │
│   └── SDD 核心 (7):              # 严格的契约驱动交付
│       ├── sdd-orchestrator/    # 状态机 + 门禁执行
│       ├── spec-architect/      # 将需求凝固为契约
│       ├── spec-to-codebase/    # 从契约生成代码
│       ├── spec-contract-diff/  # 检测契约漂移
│       ├── spec-driven-test/    # 契约验证测试（契约先行）
│       ├── spec-traceability/   # 故事→契约→代码→测试追溯
│       └── sdd-release-guard/  # 最终发布门禁
│
│   └── TDD 独立 (1):               # 独立的红绿重构模式
│       └── test-driven-development/ # 单模块逻辑（无需 SPEC）
├── docs/                       # 设计文档
│   ├── VIC-CLI-GUIDE.md       # CLI操作指南
│   └── *.md
│
└── .vic-sdd/                   # 项目记忆与规范
    ├── SPEC-REQUIREMENTS.md    # 需求规范
    ├── SPEC-ARCHITECTURE.md    # 架构规范
    ├── PROJECT.md              # 项目状态
    ├── knowledge-boundary.yaml # AI 认知地图
    ├── decision-guardrails.yaml # 决策约束
    ├── signal-register.yaml    # 证据链进度
    ├── exploration-journal.yaml # 探索日志
    ├── status/
    │   ├── events.yaml          # 事件历史
    │   └── state.yaml          # 当前状态
    ├── tech/
    │   └── tech-records.yaml   # 技术决策
    ├── risk-zones.yaml        # 风险记录
    ├── project.yaml            # AI 快速参考
    └── dependency-graph.yaml   # 模块依赖
```

## 核心理念

### 定图纸 (需求)
- 定义用户故事和验收标准
- 规划开发阶段
- 创建 SPEC-REQUIREMENTS.md

### 打地基 (架构)
- 评估技术选型
- 设计系统架构
- 创建 SPEC-ARCHITECTURE.md

### 立规矩 (实现)
- 小步迭代
- 门禁检查推进
- 收敛到 PRD/ARCH/PROJECT

### 自我认知 (Self-Awareness)
VIBE-SDD 通过 4 个机制赋予 AI"自知之明"：
- **Knowledge Boundary** — 知道自己知道什么、推断什么、假设什么、不知道什么
- **Pre-Decision Check** — 重大决策前的门禁检查（范围、质量、信号）
- **Signal Register** — 证据链代替"60% 完成"来追踪进展
- **Exploration Journal** — 记录探索过程，避免重复踩坑

## AI 快速开始

当 AI 在这个项目上开始工作时，请按以下顺序阅读：

```
1. .vic-sdd/PROJECT.md                → 项目状态、里程碑
2. .vic-sdd/SPEC-REQUIREMENTS.md      → 需求、验收标准
3. .vic-sdd/SPEC-ARCHITECTURE.md      → 架构、技术栈
4. .vic-sdd/risk-zones.yaml           → 高风险区域
```

**结果**: AI 能在约 15 秒内理解项目上下文。

## 典型工作流

| 场景 | 命令 |
|------|------|
| 开始新项目 | `vic init` |
| 初始化 SPEC | `vic spec init` |
| 做技术决策 | `vic rt` |
| 发现风险 | `vic rr` |
| 推进阶段 | `vic phase advance` |
| 查看阶段 | `vic phase status` |
| 通过门禁 | `vic gate pass --gate N` |
| AI 说"完成了" | `vic check` |
| 提交前验证 | `vic validate` |
| 备份记忆 | `vic export` |

## 相关 Skills

19 个 Skills 均按 Google 5 Agent Design Patterns 分类：

| 模式 | Skill | 用途 |
|------|-------|------|
| **Generator** | `spec-architect` | 将需求凝固为契约 |
| **Generator** | `spec-to-codebase` | 从契约生成实现代码 |
| **Generator** | `vibe-think` | 需求澄清与权衡分析 |
| **Generator** | `vibe-redesign` | 产品探索 (EXPANSION/SELECTIVE/HOLD/REDUCTION) |
| **Generator** | `vibe-architect` | 技术选型 + 架构设计 |
| **Generator** | `vibe-design` | 设计系统咨询 |
| **Generator** | `test-driven-development` | 单模块逻辑的红绿重构 (TDD 独立模式) |
| **Reviewer** | `spec-contract-diff` | 检测代码与契约的漂移 |
| **Reviewer** | `spec-traceability` | 验证故事→契约→代码→测试链路 |
| **Reviewer** | `spec-driven-test` | 强制 100% 测试覆盖率 |
| **Reviewer** | `vibe-qa` | 端到端质量保证 (Playwright) |
| **Reviewer** | `vibe-design` (Mode 2) | 80项设计审计 + AI Slop 检测 |
| **Reviewer** | `pre-decision-check` | 所有重大决策前的门禁检查 |
| **Reviewer** | `signal-register` | 证据链进度追踪 |
| **Reviewer** | `knowledge-boundary` | 知识完整性审查 |
| **Reviewer** | `exploration-journal` | 探索过程记忆（避免重复踩坑）|
| **Reviewer** | `vibe-debug` | 根因分析 (SURVEY→PATTERN→HYPOTHESIS→IMPLEMENT) |
| **Reviewer** | `adaptive-planning` | 新信息矛盾时重新评估计划 |
| **Pipeline** | `sdd-orchestrator` | 强制执行 SDD 状态机 (Ideation→Released) |
| **Tool Wrapper** | `vic` CLI | 25个命令 — 见 [cmd/vic-go/README.md](./cmd/vic-go/README.md) |

### Schema 文件

Generator 模式产出均通过 JSON Schema 验证：

| Schema | 用途 |
|--------|------|
| `skills/spec-architect/spec-requirements.schema.json` | 验证 SPEC-REQUIREMENTS.md 结构 |
| `skills/spec-architect/spec-architecture.schema.json` | 验证 SPEC-ARCHITECTURE.md 结构 |
| `skills/sdd-orchestrator/sdd-machine-schema.json` | 验证 SDD 报告输出 |
| `skills/sdd-orchestrator/reviewer.interface.yaml` | 统一 Reviewer 接口 (8 个 Reviewer, 20 个 criteria) |

## 安装

```bash
# 依赖
pip install pyyaml

# Linux/macOS - 添加到 PATH
chmod +x cmd/vic/vic
sudo ln -s $(pwd)/cmd/vic/vic /usr/local/bin/vic

# Windows PowerShell
Set-Alias vic "python D:\Code\aaa\cmd\vic\vic"
```

## 许可证

MIT License. See [LICENSE](./LICENSE).
