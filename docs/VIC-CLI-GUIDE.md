# VIC CLI 操作指南

> 本文档是 VIC CLI 命令的完整参考手册。
> 由原 vibe-integrity Skill 降级而来。

---

## 快速开始

```bash
# 初始化项目
vic init --name "My Project" --tech "Go,PostgreSQL"

# 记录技术决策
vic rt --id DB-001 --title "Use PostgreSQL" --decision "Primary database" --reason "Need ACID"

# 查看状态
vic status

# 运行Gate检查
vic spec gate 0
```

---

## 命令总览

### 基础命令

| 命令 | 别名 | 功能 |
|------|------|------|
| `vic init` | - | 初始化项目 |
| `vic status` | - | 显示项目状态 |
| `vic check` | - | 代码对齐检查 |
| `vic validate` | - | 完整验证 |

### 记录命令

| 命令 | 别名 | 功能 |
|------|------|------|
| `vic record tech` | `vic rt` | 记录技术决策 |
| `vic record risk` | `vic rr` | 记录风险 |
| `vic record dep` | `vic rd` | 记录依赖 |

### SPEC命令

| 命令 | 别名 | 功能 |
|------|------|------|
| `vic spec init` | - | 初始化SPEC文档 |
| `vic spec status` | - | 查看SPEC状态 |
| `vic spec gate` | - | 运行Gate检查 |

### Phase/Gate命令

| 命令 | 别名 | 功能 |
|------|------|------|
| `vic phase status` | - | 查看当前阶段 |
| `vic phase advance` | - | 推进阶段 |
| `vic phase check` | - | 检查阶段要求 |
| `vic gate status` | - | 查看所有Gate |
| `vic gate pass` | - | 标记Gate通过 |
| `vic gate check` | - | 检查Gate |

### 其他命令

| 命令 | 别名 | 功能 |
|------|------|------|
| `vic search` | - | 搜索记录 |
| `vic history` | - | 查看历史 |
| `vic export` | - | 导出数据 |
| `vic import` | - | 导入数据 |
| `vic fold` | - | 折叠事件到状态 |

---

## Phase流程

```
Phase 0: 需求凝固     → Gate 0, Gate 1
     ↓
Phase 1: 架构设计     → Gate 2, Gate 3
     ↓
Phase 2: 代码实现     → Gate 4, Gate 5
     ↓
Phase 3: 验证发布     → Gate 6, Gate 7
```

### 阶段推进示例

```bash
# 1. 需求凝固完成后，通过Gate 0-1
vic gate pass --gate 0 --notes "需求完整"
vic gate pass --gate 1

# 2. 推进到架构设计阶段
vic phase advance --to 1

# 3. 架构设计完成后，通过Gate 2-3
vic gate pass --gate 2 --notes "技术栈确定"
vic gate pass --gate 3

# 4. 推进到代码实现阶段
vic phase advance --to 2

# ... 继续
```

---

## Gate参考

| Gate | 名称 | 检查内容 |
|------|------|---------|
| Gate 0 | 需求完整性 | User Story完整，Acceptance Criteria覆盖 |
| Gate 1 | 需求可测试 | 所有AC可验证，边界条件识别 |
| Gate 2 | 架构完整性 | 技术栈、模块、数据模型完整 |
| Gate 3 | 技术选型合理 | 选型有据可依，风险识别 |
| Gate 4 | 代码可编译 | 无编译错误，依赖完整 |
| Gate 5 | 代码对齐SPEC | 实现覆盖所有AC |
| Gate 6 | 功能测试通过 | 所有AC验证通过 |
| Gate 7 | 发布就绪 | 安全/性能/文档检查通过 |

---

## 目录结构

```
.vic-sdd/
├── SPEC-REQUIREMENTS.md    # 需求规范
├── SPEC-ARCHITECTURE.md    # 架构规范
├── PROJECT.md              # 项目状态
├── status/
│   ├── events.yaml         # 事件历史
│   ├── state.yaml          # 当前状态
│   ├── phase.yaml          # Phase状态
│   └── gate-status.yaml    # Gate状态
├── tech/
│   └── tech-records.yaml  # 技术决策
├── risk-zones.yaml         # 风险记录
├── project.yaml            # 项目元数据
└── dependency-graph.yaml  # 依赖图
```

---

## 常见用法

### 记录技术决策

```bash
vic rt --id DB-001 \
  --title "选择 PostgreSQL" \
  --decision "使用 PostgreSQL 作为主数据库" \
  --reason "需要 ACID 事务支持" \
  --category database \
  --impact high
```

### 记录风险

```bash
vic rr --id RISK-001 \
  --area auth \
  --desc "JWT token 过期处理不完善" \
  --impact medium
```

### 检查代码对齐

```bash
vic check
vic check --category database
vic check --json
```

### 查看Phase状态

```bash
vic phase status
vic phase check
vic phase advance --to 1
```

### 查看Gate状态

```bash
vic gate status
vic gate check --phase 0
vic gate pass --gate 0 --notes "需求完整"
```

---

## 环境变量

| 变量 | 默认值 | 说明 |
|------|-------|------|
| `VIC_DIR` | `.vic-sdd` | VIC目录名 |
| `VIC_PROJECT_DIR` | 当前目录 | 项目目录 |
| `VIC_OUTPUT` | `plain` | 输出格式 |
| `VIC_VERBOSE` | `false` | 详细输出 |

---

## 相关文档

- [SDD-PROCESS-CN.md](../docs/SDD-PROCESS-CN.md) - SDD流程规范
- [SPEC-REQUIREMENTS.md](./SPEC-REQUIREMENTS.md) - 需求规范
- [SPEC-ARCHITECTURE.md](./SPEC-ARCHITECTURE.md) - 架构规范

---

**版本**: 1.0.0  
**来源**: 由 vibe-integrity Skill 降级
**更新**: 2026-03-18
