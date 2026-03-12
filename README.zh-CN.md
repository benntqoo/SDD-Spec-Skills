# SDD-Spec Skills

[English README](./README.md)

SDD-Spec Skills 是一套可开源复用的 **严格 Spec-Driven Development（SDD）技能工具集**。
它通过状态机编排与关卡校验，把特性交付过程从“经验驱动”升级为“可追踪、可验证、可发布”。

## LAP 版本标签

- `lap-v1-strict-sdd`：v1 基线，默认对大多数任务采用重关卡严格流程
- `lap-v2-adaptive-sdd`：v2 自适应，按风险分级启用关卡并保留探索轻流程

## LAP v2 差异化设计

LAP v2 保留 v1 的可追踪与发布安全能力，同时削减阻碍快速迭代的过重仪式。

- 任务粒度升级：从 2-5 分钟原子切分，升级为保留架构上下文的有边界纵向切片
- Spec 同步升级：从全程人工同步，升级为检查点同步（`SpecCheckpoint`）并输出差异摘要
- Worktree 策略升级：改为风险分级触发，仅在高风险跨模块或并行开发场景强制使用
- 关卡策略升级：拆分为 `Explore`、`Build`、`Release` 三种模式并配置不同必选校验

### v2 状态流

`Ideation -> Explore -> SpecCheckpoint -> Build -> Verify -> ReleaseReady -> Released`

### v2 状态-技能映射

| 状态 | 主要技能 | 用途 |
|------|----------|------|
| `Ideation` | `spec-architect` | 将模糊需求转换为可执行规格 |
| `Explore` | `spec-architect`, `spec-traceability` | 架构探索，Spec 快照可选 |
| `SpecCheckpoint` | `spec-architect` | Spec 验证与差异摘要同步 |
| `Build` | `spec-to-codebase`, `spec-contract-diff`, `spec-traceability` | 代码生成与聚焦验证 |
| `Verify` | `spec-driven-test`, `spec-traceability` | 契约验证与测试覆盖 |
| `ReleaseReady` | `sdd-release-guard` | 最终发布关卡与回滚就绪 |
| `Released` | - | 功能已交付 |

`Ideation -> Explore -> SpecCheckpoint -> Build -> Verify -> ReleaseReady -> Released`

### v2 模式矩阵

- Explore 模式：本地探索实验、架构草图，Spec 快照可选
- Build 模式：功能实现与聚焦验证，要求进行检查点 Spec 同步
- Release 模式：完整契约校验、追踪通过、发布守门通过

### 快速路径模式

对于简单需求（配置变更、文档修复、Bug 修复），SDD-Spec Skills 支持**快速路径**模式，跳过非必要关卡：

```bash
# 使用快速路径配置模板
python skills/sdd-orchestrator/validate-sdd.py --config skills/sdd-orchestrator/validate-sdd.config.fast-path.json

# 或通过命令行
python skills/sdd-orchestrator/validate-sdd.py --fast-path true --fast-path-skips spec-traceability spec-contract-diff
```

**快速路径特性：**

| 特性 | 标准模式 | 快速路径 |
|------|----------|----------|
| 必需技能 | 6 个 | 4 个（最低）|
| 追踪矩阵 | 强制 | 可选 |
| 契约差异 | 必需 | 可选 |
MR|| 关卡检查 | 完整 | 简化 |
#JQ|
#JJ|## Vibe Guard - AI 完整性校验器
#RT|
#BQ|Vibe Guard 是一个 **AI 完整性校验器**，用于防止 AI 编码助手虚假声称完成。它能检测常见的幻觉模式，确保代码真正完成后才能进入下一阶段。
#KY|
#QM|### 问题背景
#RT|
#XZ|AI 编码助手经常在以下情况声称完成：
#XZ|- 代码中仍有 TODO/FIXME 占位符
#XZ|- 函数是空的存根
#XZ|- 测试永远通过（假断言）
#XZ|- 构建/验证从未实际运行
#XZ|
#BQ|Vibe Guard 通过自动化检查解决这些问题。
#KY|
#QM|### 三种模式
#RT|
#XZ|Vibe Guard 支持三种模式，平衡速度与严格性：
#XZ|
#XZ|| 模式 | 适用场景 | 阻塞条件 |
#XZ||------|----------|----------|
#XZ|| `vibe` | 快速原型、POC | 构建失败、严重安全问题 |
#XZ|| `standard` | 中小企业项目、团队开发 | 构建 + 安全 + 核心测试 |
#XZ|| `strict` | 企业级、生产环境 | 所有检查失败 |
#KY|
#QM|### 快速使用
#RT|
#XZ|```bash
#XZ|# 以不同模式运行
#XZ|python skills/vibe-guard/validate-vibe-guard.py --mode vibe
#XZ|python skills/vibe-guard/validate-vibe-guard.py --mode standard
#XZ|python skills/vibe-guard/validate-vibe-guard.py --mode strict
#XZ|
#XZ|# 配置（可选）
#XZ|# 创建 .sdd-spec/vibe-guard.config.json
#XZ|```
#KY|
#QM|### 检查类别
#RT|
#XZ|- **完整性**：TODO/FIXME、空函数、存根实现
#XZ|- **安全**：硬编码密钥、SQL 注入、XSS 漏洞
#XZ|- **可执行性**：构建成功、类型检查、代码规范
#XZ|- **测试真实性**：假测试、永远通过的断言、跳过的测试
#KY|
#QM|### 集成方式
#RT|
#XZ|Vibe Guard 可以通过以下方式调用：
#XZ|- **独立运行**：随时手动检查
#XZ|- **通过 Orchestrator**：集成到 SDD 状态转换
#XZ|- **自动触发**：检测完成短语（"done"、"ready"、"complete"）
#JQ|
#JR|## 为什么使用这套工具

## 为什么使用这套工具

MB|- 统一状态流转：`Ideation -> Explore -> SpecCheckpoint -> Build -> Verify -> ReleaseReady -> Released`
- 统一产物约束：规格、契约、测试、追踪矩阵、发布守门报告
- 统一机器校验：`validate-sdd.py` 自动检查技能一致性与关卡完整性
- 兼容多工具目录：支持单层与多层 `skills` 结构

## 技能清单

- `sdd-orchestrator`：状态机入口与路由控制
- `spec-architect`：规格与契约设计
- `spec-to-codebase`：从规格生成实现
- `spec-contract-diff`：契约漂移检测
- `spec-driven-test`：基于规格的测试关卡
- `spec-traceability`：需求-契约-代码-测试追踪
- `sdd-release-guard`：发布前最终守门
- `vibe-guard`：AI 完整性校验器（防止 AI 虚假完成）
- `spec-architect`：规格与契约设计
- `spec-to-codebase`：从规格生成实现
- `spec-contract-diff`：契约漂移检测
- `spec-driven-test`：基于规格的测试关卡
- `spec-traceability`：需求-契约-代码-测试追踪
- `sdd-release-guard`：发布前最终守门
- `vibe-guard`：AI 完整性校验器（防止 AI 虚假完成）

## 产物存储

所有 SDD 产物统一存储在 `.sdd-spec` 目录下，与项目代码分离：

```text
.sdd-spec/
  specs/              # 规格、契约、追踪文件
    <feature>.md
    <feature>.contract.json
    <feature>.traceability.yaml
    <feature>.state.json
    ...
  tests/specs/       # 测试文件
    <feature>.contract.spec.*
    <feature>.acceptance.spec.*
    ...
```

> **注意**：`.sdd-spec` 目录已通过 `.gitignore` 自动忽略版本控制。

## 目录结构

```text
skills/
  sdd-orchestrator/
    sdd-machine-schema.json
    sdd-gate-checklist.json
    validate-sdd.py
    validate-sdd.config.single-layer.json
    validate-sdd.config.multi-layer.json
  spec-architect/
  spec-to-codebase/
  spec-contract-diff/
  spec-driven-test/
  spec-traceability/
  sdd-release-guard/
  vibe-guard/
    SKILL.md
    vibe-guard.config.json
    validate-vibe-guard.py
```
skills/
  sdd-orchestrator/
    sdd-machine-schema.json
    sdd-gate-checklist.json
    validate-sdd.py
    validate-sdd.config.single-layer.json
    validate-sdd.config.multi-layer.json
  spec-architect/
  spec-to-codebase/
  spec-contract-diff/
  spec-driven-test/
  spec-traceability/
  sdd-release-guard/
```

## 快速开始

1) 默认校验（扫描 `<root>/skills`）：

```bash
python skills/sdd-orchestrator/validate-sdd.py
```

2) 使用单层目录模板：

```bash
python skills/sdd-orchestrator/validate-sdd.py --config skills/sdd-orchestrator/validate-sdd.config.single-layer.json
```

3) 使用多层目录模板：

```bash
python skills/sdd-orchestrator/validate-sdd.py --config skills/sdd-orchestrator/validate-sdd.config.multi-layer.json
```

4) 使用初始化工具创建新项目：

```bash
# 创建新项目结构
python skills/sdd-orchestrator/bootstrap-sdd.py init ./my-project

# 添加新功能
python skills/sdd-orchestrator/bootstrap-sdd.py add my-feature ./my-project

# 添加 skills 目录
python skills/sdd-orchestrator/bootstrap-sdd.py add-skills ./my-project
```


## 示例输出

```text
SDD validation passed
Root: D:\Code\aaa
Skills paths:
- D:\Code\aaa\skills
Schema: D:\Code\aaa\skills\sdd-orchestrator\sdd-machine-schema.json
Checklist: D:\Code\aaa\skills\sdd-orchestrator\sdd-gate-checklist.json
```

出现 `SDD validation passed` 时，表示技能覆盖、状态枚举与关卡清单结构均已通过一致性检查。

## 配置方式

`validate-sdd.py` 支持三类配置来源：命令参数、环境变量、JSON 配置文件。

优先级：

- `root_path`：命令参数 > 环境变量 > 配置文件 > 脚本默认
- `skills_paths`：命令参数 + 环境变量 + 配置文件合并去重

常用参数：

- `--root-path`
- `--skills-path`（可重复传入）
- `--orchestrator-path`
- `--schema-path`
- `--checklist-path`
- `--recursive-search true|false`
- `--config <json>`

环境变量：

- `SDD_VALIDATE_CONFIG`
- `SDD_ROOT_PATH`
- `SDD_SKILLS_PATHS`
- `SDD_ORCHESTRATOR_PATH`
- `SDD_SCHEMA_PATH`
- `SDD_CHECKLIST_PATH`
- `SDD_RECURSIVE_SEARCH`

## 常见失败与排查

- `Unable to resolve sdd-orchestrator path from configured skills paths`
  - 检查 `skills_paths` 是否指向真实技能根目录
  - 检查 `sdd-orchestrator` 是否包含 `sdd-machine-schema.json` 与 `sdd-gate-checklist.json`
- `SKILL.md not found for <skill>`
  - 检查目标技能目录是否存在
  - 多层目录结构请启用 `--recursive-search true`
- `missing schema reference` 或 `missing checklist reference`
  - 检查技能 `SKILL.md` 是否包含 schema 与 checklist 引用
- `State enum mismatch between schema and checklist`
  - 对齐 `sdd-machine-schema.json` 与 `sdd-gate-checklist.json` 的状态枚举
- `Checklist section incomplete for <skill>`
  - 检查 checklist 是否包含 `entry_state`、`required_outputs`、`gate_checks`

## 开源发布建议

- 技能目录统一放在项目根 `skills/`
- 避免使用工具私有路径（例如 `.trae/skills/`）
- 每次发布前执行校验脚本
- `LICENSE` 与 `.gitignore` 与功能变更一起提交

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](./LICENSE)。
