# VIBE-SDD 文档索引

> 📚 所有 VIBE-SDD 相关文档的完整列表

## 🚀 核心文档

### 开始使用
- [**VIC CLI Guide**](./VIC-CLI-GUIDE.md) - 完整的 CLI 使用指南
- **[Gate Guide](./GATE-GUIDE.md)** - Gate 检查系统完整指南
- **[Gate Cheat Sheet](./GATE-CHEAT-SHEET.md)** - 快速参考卡片
- **[README.md](../README.md)** - 项目概述和快速开始

## 🛠️ 开发文档

### Phase 完成记录
- **[FIX-LOG.md](./FIX-LOG.md)** - 详细的问题修复和功能实现记录
  - Phase 1: P0 修复（已完成）
  - Phase 2: 输出层统一（已完成）
  - Phase 3: Gate 检查增强（已完成）
  - Phase 4: 测试覆盖（已完成）
  - Phase 5: Git Hooks 集成（已完成）
  - Phase 6: Gate 检查验证（已完成）

### 项目规划
- **[TODO.md](./TODO.md)** - 待办事项和计划
- **[GAP-ANALYSIS-CODE-INDEX.md](./GAP-ANALYSIS-CODE-INDEX.md)** - 代码索引和分析
- **[GAP-ANALYSIS-PROJECT-OVERVIEW.md](./GAP-ANALYSIS-PROJECT-OVERVIEW.md)** - 项目概览

## 📋 特性文档

### Gate 系统
- **[Gate Guide](./GATE-GUIDE.md)**
  - Gate 0-3 详细说明
  - 使用方法和最佳实践
  - 故障排除指南
- **[Gate Cheat Sheet](./GATE-CHEAT-SHEET.md)**
  - 快速命令参考
  - SPEC 模板
  - 常见问题解决

### Git Hooks
- **[FIX-LOG.md#2026-04-01---phase-5-git-hooks-集成完成](./FIX-LOG.md#2026-04-01---phase-5-git-hooks-集成完成)**
  - Hook 安装和使用
  - 阻止机制说明
  - 绕过方法

### 输出系统
- **[output/output.go](../cmd/vic-go/internal/output/output.go)** - 统一输出层实现
  - 支持 JSON/YAML/Plain 格式
  - 错误码定义
  - Result 和 Error 结构

### 测试覆盖
- **[code_analysis_test.go](../cmd/vic-go/internal/commands/code_analysis_test.go)** - 分析器测试
- **[deps_test.go](../cmd/vic-go/internal/deps/deps_test.go)** - 依赖管理测试
- **[config_test.go](../cmd/vic-go/internal/config/config_test.go)** - 配置管理测试
- **[types_test.go](../cmd/vic-go/internal/types/types_test.go)** - 类型定义测试
- **[yaml_test.go](../cmd/vic-go/internal/utils/yaml_test.go)** - YAML 工具测试

## 🔧 技术文档

### 代码结构
```
cmd/vic-go/
├── main.go                    # CLI 入口点
├── internal/
│   ├── commands/              # 命令实现
│   │   ├── root.go           # 根命令
│   │   ├── spec.go           # SPEC 相关命令
│   │   ├── gate.go           # Gate 管理命令
│   │   ├── gate0.go          # Gate 0 实现
│   │   ├── gate1.go          # Gate 1 实现
│   │   ├── gate2.go          # Gate 2 实现
│   │   ├── gate3.go          # Gate 3 实现
│   │   ├── hooks.go          # Git Hooks 命令
│   │   └── ...
│   ├── config/               # 配置管理
│   ├── checker/              # 代码检查器
│   ├── types/                # 类型定义
│   ├── utils/                # 工具函数
│   ├── embedding/            # 嵌入搜索
│   └── output/               # 输出系统
```

### 核心组件
- **[GateReport](../cmd/vic-go/internal/commands/gate_report.go)** - Gate 报告生成器
- **[CodeScanner](../cmd/vic-go/internal/commands/code_scanner.go)** - 代码扫描器
- **[Embedding Store](../cmd/vic-go/internal/embedding/store.go)** - 向量存储
- **[Config](../cmd/vic-go/internal/config/config.go)** - 配置管理

## 📊 工具和脚本

### 测试脚本
- **[test-all-gates.sh](../cmd/vic-go/test-all-gates.sh)** - 测试所有 Gate 功能
- **[hook-test-summary.md](../cmd/vic-go/hook-test-summary.md)** - Hook 测试总结

### CI/CD
- **[.pre-commit-config.yaml](../.pre-commit-config.yaml)** - Pre-commit 配置
  ```yaml
  repos:
    - repo: local
      hooks:
        - id: vic-gate-check
          name: VIBE-SDD Gate Check
          entry: vic gate check --blocking
          language: system
          pass_filenames: false
          always_run: true
  ```

## 🎯 快速路径

### 新用户
1. 阅读 [README.md](../README.md) 了解项目
2. 查看 [Gate Cheat Sheet](./GATE-CHEAT-SHEET.md) 快速上手
3. 运行 `vic init --name "Project" --tech "Go"` 初始化

### 开发者
1. 查看 [Gate Guide](./GATE-GUIDE.md) 了解详细用法
2. 使用 `vic spec gate [0-3]` 检查质量
3. 安装 Git Hooks: `vic hooks install`

### 高级用户
1. 阅读 [FIX-LOG.md](./FIX-LOG.md) 了解实现细节
2. 查看 TODO.md 了解未来计划
3. 自定义配置和扩展功能

## 📝 文档更新

### 最近更新 (2026-04-01)
- ✅ Phase 5: Git Hooks 集成完成
- ✅ Phase 6: Gate 检查功能验证完成
- 🆕 创建了完整的 Gate 系统文档
- 🆕 添加了快速参考卡片
- 🆕 更新了 README 和索引

### 文档维护
- 所有修复和功能实现都会记录在 [FIX-LOG.md](./FIX-LOG.md)
- 待办事项在 [TODO.md](./TODO.md) 中跟踪
- 文档使用 Markdown 格式，易于维护和扩展

## 🔗 相关链接

- [GitHub Repository](https://github.com/vic-sdd/vic)
- [CLAUDE.md](../CLAUDE.md) - Claude Code 工作指导
- [skills/](../skills/) - AI 技能文档