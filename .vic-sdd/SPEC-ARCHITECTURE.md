# SPEC-ARCHITECTURE: VIBE-SDD CLI

> 此文档为技术架构规范，定义了项目的技术选型、系统架构和数据模型。
> 详细需求请参考 SPEC-REQUIREMENTS.md。

---

## 元数据

| 字段 | 值 |
|------|-----|
| version | 1.0.0 |
| status | spec |
| owner | @sisyphus |
| created | 2026-03-18 |
| updated | 2026-03-18 |

---

## 1. 技术选型

### 1.1 技术栈总览

| 层级 | 技术 | 版本 | 状态 |
|------|------|------|------|
| CLI框架 | Cobra | v1.2.0+ | 选中 |
| 配置文件 | Viper | v1.18.0+ | 选中 |
| 数据存储 | YAML | - | 选中 |
| 测试 | Go testing | 内置 | 选中 |
| 构建 | Make | - | 选中 |

**最终选择**: Go 1.21+  
**选择理由**: 单一二进制，跨平台，无运行时依赖，启动快

### 1.2 开发工具

| 工具 | 用途 | 选择 |
|------|------|------|
| 语言 | 开发语言 | Go |
| 包管理 | 依赖管理 | Go modules |
| CLI框架 | 命令行 | Cobra + Viper |
| 代码规范 | Linting | golangci-lint |
| 测试 | 单元测试 | Go testing |

---

## 2. 系统架构

### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI 客户端                          │
│                           vic                               │
└───────────────────────────┬────────────────────────────────┘
                            │
                            ▼
┌───────────────────────────────────────────────────────────────┐
│                      命令层 (Cobra)                          │
│                                                               │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐         │
│  │  init   │  │ record  │  │ status  │  │  spec   │   ...   │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘         │
└───────┼────────────┼────────────┼────────────┼───────────────┘
        │            │            │            │
        ▼            ▼            ▼            ▼
┌───────────────────────────────────────────────────────────────┐
│                       业务逻辑层                              │
│                                                               │
│  ┌─────────────────────────────────────────────────────┐     │
│  │                  config (配置加载)                    │     │
│  └─────────────────────────────────────────────────────┘     │
│                                                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           │
│  │   commands  │  │   checker   │  │    utils    │           │
│  └─────────────┘  └─────────────┘  └─────────────┘           │
└───────────────────────────┬────────────────────────────────┘
                            │
                            ▼
┌───────────────────────────────────────────────────────────────┐
│                        数据层                                 │
│                                                               │
│   ┌─────────────┐        ┌─────────────┐                     │
│   │  .vic-sdd/  │        │    YAML     │                     │
│   │  (文件存储)  │        │   文件操作   │                     │
│   └─────────────┘        └─────────────┘                     │
└───────────────────────────────────────────────────────────────┘
```

### 2.2 模块划分

| 模块 | 职责 | 依赖 | 边界 |
|------|------|------|------|
| commands | CLI命令实现 | Cobra | 入口 |
| config | 配置加载和管理 | Viper | 基础设施 |
| checker | 代码对齐检查 | 无 | 业务 |
| utils | 文件/YAML工具 | 无 | 基础设施 |
| types | 类型定义 | 无 | 基础设施 |

---

## 3. 目录结构

```
cmd/vic-go/
├── main.go                 # 入口文件
├── Makefile               # 构建配置
├── go.mod                # 依赖管理
├── go.sum                # 依赖锁定
├── README.md             # 项目文档
│
├── internal/
│   ├── commands/         # CLI命令实现
│   │   ├── root.go       # 根命令
│   │   ├── init.go       # vic init
│   │   ├── record.go     # vic record
│   │   ├── spec.go       # vic spec
│   │   ├── check.go      # vic check
│   │   ├── status.go     # vic status
│   │   ├── misc.go       # 其他命令
│   │   └── ...
│   │
│   ├── config/           # 配置管理
│   │   └── config.go
│   │
│   ├── checker/          # 代码检查
│   │   └── code_analysis.go
│   │
│   ├── types/            # 类型定义
│   │   └── types.go
│   │
│   └── utils/            # 工具函数
│       ├── file.go
│       └── yaml.go
```

## 4. API/命令设计

### 4.1 命令目录

| 命令 | 方法 | 功能 | 状态 |
|------|------|------|------|
| vic init | - | 初始化项目 | done |
| vic rt / vic record tech | - | 记录技术决策 | done |
| vic rr / vic record risk | - | 记录风险 | done |
| vic status | - | 显示状态 | done |
| vic check | - | 代码对齐检查 | done |
| vic validate | - | 完整验证 | done |
| vic spec init | - | 初始化SPEC | done |
| vic spec gate | - | Gate检查 | done |

## 4. 命令设计

### 4.1 命令目录

| 命令 | 功能 | 状态 |
|------|------|------|
| vic init | 初始化项目 | done |
| vic record tech (vic rt) | 记录技术决策 | done |
| vic record risk (vic rr) | 记录风险 | done |
| vic status | 显示状态 | done |
| vic check | 代码对齐检查 | done |
| vic validate | 完整验证 | done |
| vic spec init | 初始化SPEC | done |
| vic spec gate | Gate检查 | done |
| vic fold | 事件折叠 | done |
| vic search | 搜索 | done |
| vic history | 历史记录 | done |
| vic export | 导出数据 | done |
| vic import | 导入数据 | done |

---

## 5. 核心模块详解

### 5.1 commands 模块

负责所有CLI命令的实现：

| 命令文件 | 功能 | 依赖 |
|---------|------|------|
| init.go | vic init - 初始化项目 | utils, config |
| record.go | vic record - 记录决策/风险 | utils, config |
| spec.go | vic spec - SPEC文档管理 | utils, config |
| check.go | vic check - 代码对齐检查 | checker |
| status.go | vic status - 状态显示 | config |
| fold.go | vic fold - 事件折叠 | utils |
| search.go | vic search - 搜索 | utils |
| export.go / import.go | 数据导出导入 | utils |

### 5.2 checker 模块

代码对齐检查功能：

| 功能 | 说明 |
|------|------|
| code_analysis.go | 分析代码结构与SPEC一致性 |
| 规则引擎 | 可扩展的检查规则 |

### 5.3 utils 模块

| 功能 | 说明 |
|------|------|
| file.go | 文件操作工具 |
| yaml.go | YAML读写工具 |

---

## 6. 配置文件

### 6.1 环境变量

| 变量 | 默认值 | 说明 |
|------|-------|------|
| VIC_DIR | .vic-sdd | VIC目录名 |
| VIC_PROJECT_DIR | 当前目录 | 项目目录 |
| VIC_OUTPUT | plain | 输出格式 |
| VIC_VERBOSE | false | 详细输出 |

---

## 7. 部署和运维

### 7.1 构建

```bash
cd cmd/vic-go
make build    # 构建当前平台
make build-all # 跨平台构建
make install  # 安装到PATH
```

### 7.2 发布

| 平台 | 格式 |
|------|------|
| Linux | vic-{version}-linux-amd64.tar.gz |
| macOS | vic-{version}-darwin-amd64.tar.gz |
| Windows | vic-{version}-windows-amd64.zip |

---

## 8. 变更历史

| 日期 | 变更内容 | 变更人 | 原因 |
|------|---------|--------|------|
| 2026-03-18 | 初始版本 | @sisyphus | 创建SPEC-ARCHITECTURE |
| 2026-03-18 | 补充CLI架构 | @sisyphus | Phase 1 架构设计 |

---

## 附录

### A. 相关文档

- SPEC-REQUIREMENTS.md - 需求规范
- PROJECT.md - 项目状态追踪
- docs/SDD-PROCESS-CN.md - SDD流程规范

### B. 参考资料

- Cobra文档: https://cobra.dev
- Viper文档: https://github.com/spf13/viper

### C. 术语表

| 术语 | 定义 |
|------|------|
| CLI | Command-Line Interface 命令行工具 |
| Cobra | Go语言CLI框架 |
| Viper | Go语言配置管理库 |
| Gate | 阶段检查点 |
| SPEC | 需求/架构规范文档 |
