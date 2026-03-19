---
name: exploration-journal
description: Use when starting new tasks, encountering problems, or making decisions - to record exploration process and avoid repeating failed approaches.
---

# Exploration Journal

记录 AI 的"思考过程"，避免重复探索。

---

## 目的

```
为什么要记录探索过程？

对于 AI：
- 避免重复探索同一个问题
- 记住之前失败的方法
- 保持决策连贯性

对于人类：
- 可以追溯 AI 的思考过程
- 理解 AI 为什么做出某个决定
- 在 AI 卡住时提供帮助
```

---

## 记录类型

### 1. explore - 开始探索

当你开始探索一个领域时：

```yaml
- id: E-001
  timestamp: 2026-03-19 10:00
  action: explore
  topic: "如何实现邮箱验证"
  goal: "找到最佳验证方案"
  findings:
    - "package.json 已有 validator 库"
    - "src/utils/validation.ts 已有 isEmail 实现"
    - "但是不支持 +suffix 格式"
  status: complete
```

### 2. tried - 尝试某个方法

当你尝试但失败或成功时：

```yaml
# 成功
- id: E-002
  timestamp: 2026-03-19 10:15
  action: tried
  approach: "使用 validator.isEmail()"
  result: success
  reason: "支持 edge cases，开箱即用"
  evidence: "测试通过: user+tag@example.com ✓"

# 失败
- id: E-003
  timestamp: 2026-03-19 10:30
  action: tried
  approach: "用正则表达式验证邮箱"
  result: failed
  reason: "正则无法正确处理 edge cases"
  evidence: "测试失败: user+tag@example.com ✗"
```

### 3. decided - 做出决策

当你最终做出选择时：

```yaml
- id: E-004
  timestamp: 2026-03-19 10:45
  action: decided
  choice: "使用 validator.isEmail()"
  alternatives_considered:
    - "正则表达式: rejected (edge cases)"
    - "自定义函数: rejected (overkill)"
  reason: "利用现有依赖，支持 edge cases"
  confidence: 0.8
```

### 4. learned - 学习到教训

当你从错误中学到东西时：

```yaml
- id: E-005
  timestamp: 2026-03-19 11:00
  action: learned
  lesson: "正则表达式不适合复杂验证"
  context: "尝试用正则验证邮箱"
  implication: "未来遇到验证问题优先考虑专用库"
```

---

## 查询时机

### 开始新任务前

```
问：
"我之前探索过这个领域吗？"
"有什么失败的尝试需要避免？"
"之前的决策是什么？"

命令：
grep -B 2 -A 5 "topic:" .vic-sdd/exploration-journal.yaml
grep "result: failed" .vic-sdd/exploration-journal.yaml
grep "action: decided" .vic-sdd/exploration-journal.yaml
```

### 遇到问题时

```
问：
"我之前尝试过什么方法？"
"有没有类似问题的解决方案？"

命令：
grep -B 1 -A 3 "action: tried" .vic-sdd/exploration-journal.yaml
```

### 做出决策前

```
问：
"我考虑过哪些替代方案？"
"之前的决策理由是什么？"

命令：
grep -B 2 -A 5 "action: decided" .vic-sdd/exploration-journal.yaml
```

---

## 工作流

### 1. 开始探索时

```yaml
- id: E-xxx
  timestamp: now
  action: explore
  topic: "我要探索的问题"
  goal: "我想找到的答案"
  status: exploring
```

### 2. 尝试方法时

```yaml
- id: E-xxx
  timestamp: now
  action: tried
  approach: "尝试的方法"
  result: success/failed/blocked
  reason: "为什么"
  evidence: "什么证据"
```

### 3. 做出决策时

```yaml
- id: E-xxx
  timestamp: now
  action: decided
  choice: "最终选择"
  alternatives_considered:
    - "选项1: rejected - 原因"
    - "选项2: rejected - 原因"
  reason: "为什么选这个"
```

---

## 示例

**场景：实现邮箱验证**

```
10:00 - 开始探索
→ 记录：探索邮箱验证领域
→ findings: validator 库、现有实现、edge cases

10:15 - 尝试正则
→ 记录：正则失败
→ reason: 无法处理 +suffix
→ 这个记录了，未来不会重复尝试

10:30 - 尝试 validator 库
→ 记录：validator 成功
→ 决定：使用 validator

10:45 - 更新探索状态
→ status: complete
→ 记录决策理由
```

**场景：修复 bug**

```
11:00 - 开始探索
→ topic: "为什么登录失败"
→ goal: "找到根本原因"

11:15 - 尝试方法1：检查密码哈希
→ result: failed
→ reason: 哈希正确

11:30 - 尝试方法2：检查 token 生成
→ result: failed
→ reason: token 生成正确

11:45 - 尝试方法3：检查 token 验证
→ result: success
→ reason: token 验证时用了错误的 secret

12:00 - 决策
→ 修复：使用正确的 secret
→ lesson: 验证 token 时 secret 必须匹配
```

---

## 快速命令

```bash
# 查看最近的探索
tail -50 .vic-sdd/exploration-journal.yaml

# 查看所有失败尝试
grep -B 2 -A 3 "result: failed" .vic-sdd/exploration-journal.yaml

# 查看所有决策
grep -B 1 -A 5 "action: decided" .vic-sdd/exploration-journal.yaml

# 查看某个主题的探索
grep -B 2 -A 10 "topic:.*邮箱" .vic-sdd/exploration-journal.yaml
```

---

## 好处

```
不记录探索过程：
  - 3 天后又尝试了同样的错误方法
  - 忘记了之前的决策理由
  - 需要从头开始理解代码

记录探索过程：
  - 可以回溯思考路径
  - 知道什么方法失败了
  - 理解为什么做某个决定
  - 新接手的人也能理解
```

---

## 与其他技能的关系

```
knowledge-boundary → 发现未知 → 触发探索
signal-register → 记录进展 → 探索过程也是信号
pre-decision-check → 决策前 → 查询 journal
```

---

## 常见错误

| 错误 | 正确做法 |
|------|----------|
| 不记录探索 | 开始探索就记录 |
| 只记成功 | 成功失败都记录 |
| 不记失败原因 | 失败原因最重要 |
| 不查 journal | 开始任务前先查询 |
| 重复尝试 | 查询避免重复 |
