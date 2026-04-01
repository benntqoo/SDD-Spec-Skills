# VIBE-SDD Gate 检查完整指南

> 本文档详细说明 VIBE-SDD 的质量门控系统，包括所有 Gate 的检查内容、使用方法和最佳实践。

## 概览

VIBE-SDD 使用 4 个主要 Gate 来确保项目的质量：

| Gate | 名称 | 检查内容 | 阶段 |
|------|------|----------|------|
| Gate 0 | 需求完整性 | SPEC-REQUIREMENTS.md 完整性 | Phase 0 |
| Gate 1 | 架构完整性 | SPEC-ARCHITECTURE.md 完整性 | Phase 1 |
| Gate 2 | 代码对齐 | 代码与 SPEC 一致性 | Phase 2 |
| Gate 3 | 测试覆盖率 | 测试文件和覆盖情况 | Phase 3 |

## 快速开始

### 1. 初始化项目
```bash
vic init --name "My Project" --tech "Go,PostgreSQL"
```

### 2. 编辑 SPEC 文件
```bash
# 编辑需求文档
vim .vic-sdd/SPEC-REQUIREMENTS.md

# 编辑架构文档
vim .vic-sdd/SPEC-ARCHITECTURE.md
```

### 3. 运行 Gate 检查
```bash
# 检查单个 gate
vic spec gate 0
vic spec gate 1 --format json

# 使用智能选择
vic gate smart

# 安装 git hooks（自动检查）
vic hooks install
```

### 4. 提交代码
```bash
# 提交时会自动检查 gates
git commit -m "feature: new functionality"

# 绕过检查（紧急情况）
git commit --no-verify -m "message"
```

## 详细 Gate 说明

### Gate 0: 需求完整性

**目的**: 确保需求文档完整且可测试

**检查内容**:
- ✅ User Stories Section
- ✅ Key Features Section
- ✅ Acceptance Criteria
- ✅ Non-Functional Requirements
- ✅ Out of Scope
- ❌ Features Have Criteria（必须每个功能都有验收标准）

**通过条件**: 至少 6/7 项检查通过

**示例 SPEC-REQUIREMENTS.md**:
```markdown
# SPEC-REQUIREMENTS.md

## User Stories
- [x] As a user, I can see "Hello, World!" displayed
- [x] As a developer, I can implement basic Go programs

## Key Features
1. Basic "Hello, World!" program
2. Simple console output
3. No external dependencies

## Acceptance Criteria
### Feature 1: Basic Hello World
- [x] Program successfully displays "Hello, World!"
- [x] No compilation errors
- [x] Runs without external dependencies

## Non-Functional Requirements
- Performance: < 100ms response time
- Security: No input validation required for demo

## Out of Scope
- Database integration
- User authentication
```

**常见问题**:
```
❌ Only 5/13 features have acceptance criteria
   → Add criteria for each feature
```

### Gate 1: 架构完整性

**目的**: 确保架构设计完整且合理

**检查内容**:
- ✅ Technology Stack Section（技术栈）
- ✅ System Design Section（系统设计）
- ✅ API Design Section（API设计）
- ✅ Data Model Section（数据模型）
- ✅ Security Section（安全）
- ❌ Decision Rationale（决策理由）

**通过条件**: 至少 5/6 项检查通过

**示例 SPEC-ARCHITECTURE.md**:
```markdown
# SPEC-ARCHITECTURE.md

## System Design
### Components
- VIC CLI Tool - Command-line interface
- Core Libraries - Go libraries for SDD operations

### Data Flow
User input → VIC CLI → Core libraries → Configuration → Gate checks

## Data Model
### Core Entities
- **Project**: Contains all project metadata
- **Phase**: SDD phase with gates
- **Gate**: Individual gate check

## Technology Stack
| Layer | Technology | Rationale |
|-------|------------|----------|
| CLI Tool | Go | Performance and single binary |
| Database | SQLite | Lightweight embedded database |

## Decision Rationale
### Go for CLI Development
- **Why**: Compiled performance, single binary deployment
- **Alternative**: Python (slower startup)
- **Impact**: Faster execution, easier distribution
```

### Gate 2: 代码对齐

**目的**: 确保代码与 SPEC 保持一致

**检查内容**:
- ✅ Tech Stack Alignment（技术栈对齐）
- ✅ API Endpoint Alignment（API对齐）
- ✅ Module Structure（模块结构）
- ✅ Security Implementation（安全实现）
- ❌ Code TODO Comments（TODO 注释）
- ❌ Constitution Rules（宪法规则）

**宪法规则**:
- NO-TODO-IN-CODE: 不允许 TODO/FIXME/XXX/HACK 注释
- NO-CONSOLE-IN-PROD: 生产代码不允许 console.log
- NO-HARD-CODED-SECRETS: 不允许硬编码密钥

**通过条件**: 至少 70/77 项检查通过

**常见问题**:
```
❌ Found 87 TODO/FIXME/XXX/HACK comments in code
   → Most critical: code_scanner.go:16

❌ Console statement should not be in production code
   → Details: problematic.js:2

❌ Hardcoded secrets detected
   → Details: test.js:3
```

### Gate 3: 测试覆盖率

**目的**: 确保关键功能有适当的测试

**检查内容**:
- ✅ Test Files Exist（测试文件存在）
- ✅ Test Framework Configured（测试框架配置）
- ✅ Module Test Coverage（模块测试覆盖）
- ❌ Critical Path Coverage（关键路径覆盖）

**通过条件**: 至少 3/4 项检查通过

**检测的测试框架**:
- Go, JavaScript, TypeScript, Java, Python, Rust

## 高级功能

### Smart Gate Selection

智能分析变更并推荐需要检查的 Gates：

```bash
vic gate smart --output json
```

**示例输出**:
```json
{
  "change_type": "feature_addition",
  "risk_level": "medium",
  "risk_score": 1.50,
  "gates_required": [2, 3],
  "gates_skipped": [0, 1],
  "recommended_skill": "implementation"
}
```

**变更类型**:
- `feature_addition`: 新功能 - 高风险
- `bug_fix`: Bug 修复 - 中等风险
- `documentation`: 文档更新 - 低风险

### Phase Management

系统根据当前阶段自动调整检查范围：

- **Phase 0**: 只检查 Gate 0
- **Phase 1**: 检查 Gates 0-1
- **Phase 2**: 检查 Gates 0-3
- **Phase 3**: 检查所有 Gates (0-7)

```bash
# 查看当前阶段状态
vic gate status

# 手动进入下一阶段
vic phase advance --to 1
```

## Git Hooks 集成

### 安装 Hooks

```bash
vic hooks install
```

### 卸载 Hooks

```bash
vic hooks uninstall
```

### Hook 行为

1. **自动运行**: 每次 `git commit` 前执行
2. **检查当前及之前阶段**: 只检查已通过阶段的 Gates
3. **阻止失败**: 如果 Gate 未通过，阻止提交
4. **绕过机制**: 使用 `--no-verify` 绕过检查

### 最佳实践

```bash
# 1. 修复所有 gate 问题
vic spec gate 0
vic spec gate 1
vic spec gate 2
vic spec gate 3

# 2. 通过 gates
vic gate pass --gate 0
vic gate pass --gate 1
vic gate pass --gate 2
vic gate pass --gate 3

# 3. 安装 hooks
vic hooks install

# 4. 正常提交
git commit -m "feat: implementation complete"
```

## 故障排除

### 常见错误

#### 1. "SPEC-REQUIREMENTS.md not found"
```bash
# 解决方案
vic spec init
```

#### 2. "Gate check failed: required gates not passed"
```bash
# 解决方案
vic spec gate <number>
# 然后修复问题
```

#### 3. "vic.exe not found" (Git Hook)
```bash
# 确保 vic.exe 在项目根目录
ls vic.exe
```

### 调试技巧

1. **使用详细输出**
```bash
vic --verbose spec gate 0
```

2. **检查 JSON 输出**
```bash
vic spec gate 0 --format json
```

3. **查看 Gate 状态**
```bash
vic gate status
```

## 最佳实践

### 1. 渐进式验证
- 在每个阶段结束时运行对应的 Gates
- 确保 Gate 通过后再进入下一阶段

### 2. 定期检查
- 每次重要的代码更改后运行 `vic spec gate 2`
- 提交前使用 `vic gate check --blocking`

### 3. SPEC 先行
- 始终先更新 SPEC，再实现代码
- 使用 Constitution 规则指导编码实践

### 4. 测试驱动
- 为关键功能编写测试
- Gate 3 会自动检测缺失的测试

### 5. 智能选择
- 对于大型更改，使用 `vic gate smart` 优化检查流程
- 根据风险级别调整审查力度

## 性能优化

### 大型项目

1. **选择性检查**
```bash
# 只检查必要的 gates
vic gate check --phase 2
```

2. **缓存结果**
- Gate 结果会被缓存，重复检查更快

3. **并行检查**
- Smart Gate 自动并行执行推荐的检查

### CI/CD 集成

```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]
jobs:
  gates:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup Go
        uses: actions/setup-go@v2
      - name: Run Gate Checks
        run: |
          ./vic gate check --blocking
```

## 总结

VIBE-SDD Gate 系统提供了完整的质量保证流程：

1. **需求阶段**: 确保需求清晰可测试
2. **设计阶段**: 确保架构完整合理
3. **实现阶段**: 确保代码质量和测试覆盖
4. **发布阶段**: 确保功能完整和发布就绪

通过合理使用这些 Gates，可以显著提高软件质量和开发效率。