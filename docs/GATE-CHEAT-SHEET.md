# VIBE-SDD Gate 检查快速参考

> 🚪 4 个核心 Gates + Git Hooks = 质量保证系统

## 核心 Gates

| Gate | 检查内容 | 常见问题 | 解决方案 |
|------|----------|----------|----------|
| **Gate 0** | 需求完整性 | 缺少验收标准 | 为每个功能添加 AC |
| **Gate 1** | 架构完整性 | 缺少决策理由 | 解释为什么选这个技术 |
| **Gate 2** | 代码对齐 | TODO 注释、console.log | 修复或移除违规代码 |
| **Gate 3** | 测试覆盖率 | 缺少测试文件 | 为关键功能编写测试 |

## 常用命令

```bash
# 初始化
vic init --name "Project" --tech "Go"

# 检查 Gates
vic spec gate 0              # 需求完整性
vic spec gate 1 --format json # 架构完整性
vic spec gate 2              # 代码对齐
vic spec gate 3              # 测试覆盖率

# 智能选择
vic gate smart               # 推荐 gates
vic gate smart --execute     # 运行推荐 gates

# Gate 管理
vic gate status              # 查看状态
vic gate pass --gate 0       # 通过 gate
vic gate check --blocking    # 阻塞检查

# Git Hooks
vic hooks install            # 安装 pre-commit hook
vic hooks uninstall          # 卸载
```

## SPEC 模板

### SPEC-REQUIREMENTS.md
```markdown
# SPEC-REQUIREMENTS.md

## User Stories
- [ ] As a user, I can...

## Key Features
1. Feature 1
2. Feature 2

## Acceptance Criteria
### Feature 1
- [ ] Criterion 1
- [ ] Criterion 2

## Non-Functional Requirements
- Performance:
- Security:

## Out of Scope
- Feature X
```

### SPEC-ARCHITECTURE.md
```markdown
# SPEC-ARCHITECTURE.md

## System Design
### Components
- Component 1

## Data Model
### Core Entities
- **Entity**: Description

## Technology Stack
| Layer | Technology | Rationale |
|-------|------------|----------|
| Backend | Go | Performance |
| Database | SQLite | Lightweight |

## Decision Rationale
### Technology Choice
- **Why**: Reason
- **Alternative**: Other option
- **Impact**: Consequence
```

## 宪法规则（Constitution）

| 规则 | 描述 | 违规示例 |
|------|------|----------|
| NO-TODO-IN-CODE | 不允许 TODO 注释 | // TODO: Fix later |
| NO-CONSOLE-IN-PROD | 生产代码禁用 console.log | console.log("debug") |
| NO-HARD-CODED-SECRETS | 不允许硬编码密钥 | password = "secret" |

## 错误解决

### Gate 0 失败
```
❌ Only 5/13 features have acceptance criteria
```
**Fix**: 为每个功能添加验收标准

### Gate 1 失败
```
❌ Missing rationale for tech decisions
```
**Fix**: 在 Technology Stack 中添加 Rationale 列

### Gate 2 失败
```
❌ Found 87 TODO/FIXME/XXX/HACK comments
❌ Console statement should not be in production code
```
**Fix**: 移除或解决所有 TODO，删除 console.log

### Gate 3 失败
```
❌ Only 0/1 critical files have tests
```
**Fix**: 为 main.go、handlers 等关键文件添加测试

## Git Hooks 流程

```bash
# 1. 安装 hooks
vic hooks install

# 2. 开发代码
# ... coding ...

# 3. 提交（自动检查）
git commit -m "feat: add feature"

# 4. 如果被阻止
# 运行检查找出问题
vic spec gate 2
# 修复问题
# 重新提交

# 5. 绕过（紧急情况）
git commit --no-verify -m "hotfix"
```

## 快速检查清单

- [ ] SPEC 文档完整了吗？
- [ ] 每个功能都有验收标准吗？
- [ ] 技术选择有理由吗？
- [ ] 代码里有 TODO 吗？
- [] 生产代码有 console.log 吗？
- [] 关键功能有测试吗？
- [] Git hooks 安装了吗？

## Pro Tips

1. **每次提交前**: `vic gate check --blocking`
2. **新项目先**: `vic spec init`
3. **架构设计**: 先写 SPEC，再写代码
4. **测试驱动**: 先写测试，再写实现
5. **定期审查**: 使用 `vic gate smart` 优化检查

记住：**SPEC 先行，代码跟进，质量保证**