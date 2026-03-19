---
name: knowledge-boundary
description: Use when starting new tasks, encountering uncertainties, or making important decisions - activates self-awareness of what AI knows, infers, assumes, or doesn't know.
---

# Knowledge Boundary

让 AI 有"自知之明"——知道自己知道什么、不知道什么。

---

## 激活时机

**必须激活 when:**
- 开始新任务前
- 遇到不确定的问题时
- 做出重要决策前
- 听到/看到不熟悉的概念时

**可以跳过 when:**
- 简单的语法修复
- 明显的小改动
- 完全基于已知代码的操作

---

## 核心概念

```
Knowledge Boundary 的四个区域：

┌─────────────────────────────────────────────────────────┐
│                      UNKNOWN                            │
│         我完全不知道的（可能是阻塞点）                       │
│                                                          │
│    ┌─────────────────────────────────────────────────┐   │
│    │                   ASSUMED                       │   │
│    │          我假设的（高风险，需要验证）              │   │
│    │                                                  │   │
│    │    ┌─────────────────────────────────────────┐   │   │
│    │    │                INFERRED                │   │   │
│    │    │      我从模式推断的（需要验证）          │   │   │
│    │    │                                          │   │   │
│    │    │    ┌─────────────────────────────────┐   │   │   │
│    │    │    │              KNOWN              │   │   │   │
│    │    │    │    我验证过的事实（最高可信度）   │   │   │   │
│    │    │    └─────────────────────────────────┘   │   │   │
│    │    └─────────────────────────────────────────┘   │   │
│    └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

---

## 工作流

### 1. 照镜子（查询）

在做任何事之前，先读 `.vic-sdd/knowledge-boundary.yaml`

```
对照你的问题：

❓ "X 是怎么工作的？"

→ X 在 known 里？→ 直接用这个事实
→ X 在 inferred 里？→ 评估 confidence，考虑验证
→ X 在 assumed 里？→ ⚠️ 高风险，先验证
→ X 在 unknown 里？→ ⚠️ 阻塞，需要先解决
→ X 都不在？→ 你是怎么知道 X 的？标记为 inferred
```

### 2. 分类处理

**如果在 UNKNOWN：**
```
→ 先停下来
→ 尝试解决 unknown：
  - 查代码库
  - 读文档
  - 问人类
→ 如果解决不了 → 更新 knowledge-boundary 为 blocker，等待
→ 如果解决了 → 移到 known
```

**如果在 ASSUMED：**
```
→ 尝试验证（grep、读源码、问人）
→ 验证成功 → 移到 known
→ 验证失败 → 移到 inferred 或 unknown
```

**如果在 INFERRED：**
```
→ 如果 confidence > 0.8 → 可以继续，但记录 warning
→ 如果 confidence < 0.8 → 尝试提升（验证）
```

### 3. 记录（更新镜子）

每次验证/发现后，更新边界：

```yaml
# 如果验证了新事实
known:
  - id: K-xxx
    fact: "你的新发现"
    source: "schemas/users.sql:15"
    verified: 2026-03-19
    by: "grep -n 'email' schemas/users.sql"

# 如果做出了推断
inferred:
  - id: I-xxx
    claim: "你的推断"
    reasoning: "基于什么"
    confidence: 0.7
    needs_verification: true

# 如果发现了不知道的
unknown:
  - id: U-xxx
    question: "你的问题"
    blocks: ["T03", "T05"]  # 阻塞的任务
    priority: high
    resolution_path: "如何解决"
```

---

## 镜子检查原则

**在做以下事前必须照镜子：**
- 实现新功能
- 修改现有逻辑
- 相信自己的"最佳实践"
- 给用户建议
- 引入新依赖
- 做出技术选型

**镜子不完整时的处理：**

| 情况 | 处理 |
|------|------|
| inferred 堆积 > 5 个 | 先验证，再继续 |
| assumed 有 high risk | 必须验证才能继续 |
| unknown 阻塞关键任务 | 暂停，等待人类 |
| 有 3+ high risk assumed | 触发人类干预 |

---

## 快速检查命令

```bash
# 查看所有未验证的推断
grep -A 3 "needs_verification: true" .vic-sdd/knowledge-boundary.yaml

# 查看高风险假设
grep -B 1 -A 2 "risk: high" .vic-sdd/knowledge-boundary.yaml

# 查看阻塞未知
grep -B 1 -A 2 "blocks:" .vic-sdd/knowledge-boundary.yaml
```

---

## 示例

**场景：用户要求实现"社交登录"**

```
1. 照镜子：
   → 搜索 "social login", "OAuth", "SSO"
   → 发现：known 里没有相关内容
   → 发现：unknown 里也没有

2. 分类：
   → 这是一个 ASSUMED 风险领域
   → 我不知道项目是否支持 OAuth

3. 检查 scope guardrails：
   → decision-guardrails.yaml 里有 forbidden: "OAuth2/第三方登录"

4. 决策：
   → ❌ 这是 forbidden 范围
   → 必须停下来，询问人类
```

**场景：需要知道数据库 schema**

```
1. 照镜子：
   → 搜索 "schema", "database", "tables"
   → 发现 known 有：
     - K-001: "Users table exists in schemas/users.sql"
     - K-002: "Users table has email column"
   
2. 应用：
   → 直接使用这些事实
   → 不需要重新验证
```

---

## 常见错误

| 错误 | 正确做法 |
|------|----------|
| 不查镜子就开始做 | 先照镜子 |
| 把假设当事实 | 标记为 assumed |
| 忽略 unknown | 标记为 blocker |
| 推断不标 confidence | 必须标 confidence |
| 验证了但不记录 | 验证后必须更新镜子 |
