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
vibe-think            vibe-architect            vibe-develop
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
├── skills-base/                # Skills 定义
│   ├── vibe-think/             # 需求澄清
│   ├── vibe-architect/          # 架构设计
│   ├── vibe-develop/           # 开发流程
│   ├── vibe-integrity/         # 记忆与验证
│   └── vibe-debug/             # 调试
│
├── docs/                       # 设计文档
│   └── *.md
│
└── .vic-sdd/                   # 项目记忆与规范
    ├── SPEC-REQUIREMENTS.md    # 需求规范
    ├── SPEC-ARCHITECTURE.md    # 架构规范
    ├── PROJECT.md              # 项目状态
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
| 推进前检查 | `vic spec gate [0-3]` |
| AI 说"完成了" | `vic check` |
| 提交前验证 | `vic validate` |
| 备份记忆 | `vic export` |

## 相关 Skills

| Skill | 用途 |
|-------|------|
| `vibe-think` | 需求澄清 |
| `vibe-architect` | 架构设计 |
| `vibe-develop` | 开发流程 |
| `vibe-integrity` | 记忆与验证 |
| `vibe-debug` | 系统性调试 |

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
