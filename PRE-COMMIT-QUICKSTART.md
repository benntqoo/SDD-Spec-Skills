# Pre-commit Quick Start
# =====================

## 一次性安装

```bash
# 1. 安装 Python 依赖
pip install -r requirements.txt

# 2. 安装 git hooks
cd /path/to/your/project
pre-commit install

# 3. (可选) 安装 commit-msg hook
pre-commit install --hook-type commit-msg

# 4. 运行所有检查
pre-commit run --all-files
```

## 日常使用

```bash
# 正常提交 (自动运行 hooks)
git commit -m "feat: add new feature"

# 跳过检查 (紧急情况)
git commit --no-verify -m "hotfix: emergency fix"

# 手动运行检查
pre-commit run --all-files
pre-commit run vibe-integrity-check --all-files
```

## 可用的 Hooks

| Hook ID | 描述 | 阻断级别 |
|---------|------|---------|
| `vibe-integrity-check` | AI 完成完整性检查 (vic validate) | 警告/错误 |
| `code-alignment-check` | 代码与决策对齐检查 (vic check) | 警告/错误 |
| `trailing-whitespace` | 尾部空格检查 | 自动修复 |
| `end-of-file-fixer` | 文件末尾换行 | 自动修复 |
| `detect-private-key` | 私钥检测 | 阻断 |
| `commitlint` | 提交消息格式 | 阻断 |

## 配置文件位置

```
your-project/
├── .pre-commit-config.yaml    # 主配置
├── .commitlintrc.yaml         # 提交消息规范 (可选)
└── cmd/vic/
    └── vic                    # CLI 工具
```

## 集成到新项目

1. 复制 `.pre-commit-config.yaml` 到新项目
2. 复制 `skills/` 目录 (或修改 hook 路径)
3. 运行 `pre-commit install`

## 自定义跳过

在 `.pre-commit-config.yaml` 中:

```yaml
# CI 中跳过耗时检查
ci:
  skip: [vibe-guard-check, code-alignment-check]

# 或在 commit 时跳过
# git commit --no-verify
```
