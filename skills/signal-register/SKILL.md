---
name: signal-register
description: Use after completing meaningful actions - to record progress as evidence chains instead of percentages, and to calculate confidence levels.
---

# Signal Register

用"证据"代替"百分比"来衡量进展。

---

## 核心思想

```
❌ 错误：60% 完成
   → 什么是 60%？怎么计算的？谁定义的？

✅ 正确：4 个测试通过，1 个功能实现
   → 有具体的、可验证的证据
```

---

## 信号类型

### 正面信号 (positive)

证明在正确方向的证据：

| 类型 | 示例 |
|------|------|
| `code_created` | "创建 login.ts，包含 150 行" |
| `test_created` | "创建 login.test.ts" |
| `test_passed` | "npm test: 5 passed" |
| `refactoring_done` | "重构完成，代码更清晰" |
| `bug_fixed` | "修复了邮箱验证 bug" |
| `spec_aligned` | "覆盖 SPEC #1, #2, #3" |
| `deps_added` | "添加了 validator 依赖" |
| `docs_updated` | "更新了 API 文档" |

### 警告信号 (warnings)

需要关注，但还没阻塞：

| 类型 | 示例 |
|------|------|
| `assumption_made` | "假设 token 过期 1 小时，未验证" |
| `edge_case_found` | "发现 +suffix 邮箱可能有问题" |
| `complexity_increased` | "圈复杂度从 5 升到 12" |
| `deps_added` | "添加了新依赖" |
| `confidence_dropped` | "信心度从 0.8 降到 0.6" |

### 阻塞信号 (blockers)

必须解决才能继续：

| 类型 | 示例 |
|------|------|
| `unknown_blocking` | "不知道 token 刷新机制" |
| `dependency_blocking` | "等待 backend API 文档" |
| `decision_blocking` | "需要决定使用哪个 auth 方案" |
| `spec_unclear` | "SPEC 对这个场景描述不清" |
| `env_blocking` | "本地环境缺少配置" |

---

## 记录原则

### 1. 每个有意义的行动 → 一个信号

```
❌ 模糊信号：
- "开始了"
- "进展顺利"
- "继续工作"

✅ 具体信号：
- "创建 src/auth/login.ts"
- "修复了类型错误"
- "测试通过: 5 passed"
```

### 2. 信号必须有时间戳

```
✅ - id: S-001
     type: code_created
     content: "创建 login.ts"
     timestamp: "2026-03-19 10:15"
```

### 3. 信号必须具体

```
❌ - "测试通过了"

✅ - "npm test: 5 passed, 0 failed"
    - "测试覆盖率从 60% 升到 80%"
    - "修复了 edge case: user+tag@example.com"
```

---

## 信心度计算

```
confidence = (positive - warnings×0.3 - blockers×0.5) / max_signals
```

**示例：**
```
positive = 4
warnings = 2  → 2 × 0.3 = 0.6
blockers = 1  → 1 × 0.5 = 0.5
max_signals = 10

confidence = (4 - 0.6 - 0.5) / 10 = 0.29
```

### 信心度解读

| 信心度 | 状态 | 行动 |
|--------|------|------|
| > 0.7 | 🟢 HIGH | 状态良好，继续推进 |
| 0.4-0.7 | 🟡 MODERATE | 可以继续，关注警告 |
| < 0.4 | 🔴 LOW | 暂停，优先解决警告和阻塞 |
| blockers >= 2 | 🛑 STOP | 停止，等待人类解决 |

---

## 工作流

### 1. 产生信号

每完成有意义的一步，记录：

```yaml
# 正面信号
positive:
  - id: S-xxx
    type: code_created
    content: "具体描述"
    timestamp: "2026-03-19 10:15"

# 警告信号
warnings:
  - id: W-xxx
    type: assumption_made
    content: "假设了什么"
    timestamp: "2026-03-19 10:20"

# 阻塞信号
blockers:
  - id: B-xxx
    type: unknown_blocking
    content: "需要知道什么才能继续"
    severity: high
    timestamp: "2026-03-19 10:25"
```

### 2. 重新计算信心度

每次添加信号后，更新：

```yaml
confidence:
  positive: 4
  warnings: 2
  blockers: 1
  calculated: 0.55
  status: moderate
```

### 3. 定期检查

```
每完成一个阶段性目标时检查：

如果 blockers >= 2：
  → 停止，更新状态，等待人类

如果 confidence < 0.4：
  → 暂停，优先解决警告

如果 confidence >= 0.7 AND blockers == 0：
  → 状态良好，继续推进
```

---

## 人类查看

人类可以用这个文件了解 AI 在做什么：

```bash
# 查看当前任务
grep -A 4 "^current_task:" .vic-sdd/signal-register.yaml

# 查看所有正面信号
grep -B 1 -A 2 "type:" .vic-sdd/signal-register.yaml | grep -A 2 "positive:"

# 查看阻塞
grep -B 1 -A 3 "severity: high" .vic-sdd/signal-register.yaml

# 查看信心度
grep -A 5 "^confidence:" .vic-sdd/signal-register.yaml
```

---

## 示例

**场景：实现用户登录功能**

```
10:00 - 开始
  → current_task: { id: T03, name: "用户登录" }

10:15 - 创建 login.ts
  → positive: { type: code_created, content: "创建 login.ts (150行)" }

10:20 - 创建测试
  → positive: { type: test_created, content: "创建 login.test.ts" }

10:25 - 测试通过
  → positive: { type: test_passed, content: "npm test: 5 passed" }

10:30 - 发现一个假设
  → warnings: { type: assumption_made, content: "假设 token 过期 1h" }

10:35 - 发现边界情况
  → warnings: { type: edge_case_found, content: "+suffix 邮箱可能有问题" }

10:40 - 遇到未知问题
  → blockers: { type: unknown_blocking, content: "不知道刷新机制" }

10:45 - 计算信心度
  → confidence = (3 - 0.6 - 0.5) / 10 = 0.19
  → status: LOW
  → 🛑 STOP → 等待人类
```

---

## 与其他技能的关系

```
knowledge-boundary → 发现 unknown → 可能是 blocker
pre-decision-check → 决策前检查 → 可能产生信号
exploration-journal → 记录探索过程 → 与信号相关
```

---

## 常见错误

| 错误 | 正确做法 |
|------|----------|
| "进度 60%" | 记录具体信号 |
| "开始了" | 记录具体行动 |
| 不更新信心度 | 每次都重新计算 |
| 忽略 warning | 记录并关注 |
| 不处理 blocker | 停止，等待人类 |
