---
name: pre-decision-check
description: Use before making any significant decision - implementation, dependency introduction, architecture changes, technical choices, or recommendations to users.
---

# Pre-Decision Check

重大决策前的"刹车检查"——防止 AI 在错误方向上狂奔。

---

## 激活时机

**在做出以下决策前必须激活：**
- 实现新功能
- 引入新依赖
- 修改核心逻辑
- 做出技术选型
- 给用户建议方案
- 添加/修改 API
- 修改数据库 schema

**简单操作可以跳过：**
- 修复明显的 typo
- 重命名变量
- 简单的格式调整
- 完全基于已知代码的小改动

---

## 检查流程

```
决策点
   ↓
┌─────────────────┐
│  1. 范围检查      │  → 在 approved？Forbidden？需要扩展？
└────────┬────────┘
         ↓
┌─────────────────┐
│  2. 尝试次数检查  │  → 超过最大次数？
└────────┬────────┘
         ↓
┌─────────────────┐
│  3. 质量红线检查  │  → 违反任何红线？
└────────┬────────┘
         ↓
┌─────────────────┐
│  4. 信号检查      │  → 阻塞太多？信心度太低？
└────────┬────────┘
         ↓
      结果
```

---

## 1. 范围检查

读取 `.vic-sdd/decision-guardrails.yaml` 中的 scope 部分：

```
✅ 在 approved 列表里
   → 继续

❌ 在 forbidden 列表里
   → 立即停止
   → 询问人类："这个在 forbidden 范围内，是否要扩展？"

⚠️ 不在任一列表里
   → 评估是否是 scope 的合理扩展
   → 如果是合理扩展且 < max_expansion
     → 记录扩展内容，继续
   → 如果超出边界
     → 停止，询问人类
```

**示例：**
```
用户要求实现"记住登录状态"
scope:
  approved: ["用户名密码登录", "JWT token 管理"]
  forbidden: ["OAuth2", "多因素认证"]
  max_expansion: 0.2

检查：
- "记住登录状态" ≈ session/token 管理
- 在 approved 范围内 → ✅ 继续
```

---

## 2. 尝试次数检查

读取 attempts 部分：

```
如果当前任务的尝试次数 >= max_on_task：
  → 停止："尝试次数过多"
  → 更新 signal-register 为 blocker

如果尝试次数 >= escalation_after：
  → 标记 warning
  → 考虑是否需要人类介入
```

---

## 3. 质量红线检查

读取 quality.hard_lines 部分：

```
检查决策是否会违反任何红线：

❌ no_todo_in_code
   → 你的实现里有没有 TODO/FIXME？
   → 有 → 不允许，必须先解决

❌ no_console_in_prod
   → 你的实现里有没有 console.log？
   → 有 → 不允许，必须删除

❌ no_hardcoded_secrets
   → 有没有硬编码的密钥/密码/token？
   → 有 → 不允许，必须用环境变量

❌ tests_required
   → 这个功能有测试吗？
   → 没有 → 不允许，必须先写测试

❌ spec_aligned
   → 这个功能与 SPEC-REQUIREMENTS.md 对齐吗？
   → 没有 → 不允许，必须先确认对齐
```

---

## 4. 信号检查

读取 `signal-register.yaml`：

```
如果 blocker 数量 >= max_blockers：
  → 停止："有太多未解决的阻塞"
  → 列出所有 blocker，等待人类解决

如果 信心度 < min_confidence：
  → 暂停："需要更多验证"
  → 先验证 inferred 项
  → 提升信心度后再继续
```

---

## 检查结果

| 结果 | 含义 | 行动 |
|------|------|------|
| ✅ PASS | 所有检查通过 | 继续执行 |
| ⚠️ WARN | 有警告但可继续 | 记录警告，继续 |
| 🛑 STOP | 不允许执行 | 停止，询问人类 |
| 🔴 BLOCK | 有阻塞无法继续 | 列出阻塞，等待人类 |

---

## 决策记录

每次决策后，更新 `.vic-sdd/decision-guardrails.yaml`：

```yaml
decisions:
  - timestamp: 2026-03-19 10:30
    decision: "使用 validator.isEmail() 验证邮箱"
    checked:
      scope: approved ✅
      attempts: 2/5 ✅
      quality: 通过 ✅
      signals: 信心度 0.7 ✅
    result: proceeded
    blockers_resolved: []
    
  - timestamp: 2026-03-19 11:00
    decision: "引入 OAuth2 依赖"
    checked:
      scope: forbidden ❌
    result: stopped
    blockers_resolved: ["需要人类批准 OAuth2"]
```

---

## 快速检查命令

```bash
# 查看范围
grep -A 10 "^scope:" .vic-sdd/decision-guardrails.yaml

# 查看质量红线
grep "no_" .vic-sdd/decision-guardrails.yaml

# 查看信号状态
grep -A 10 "^signals:" .vic-sdd/decision-guardrails.yaml
```

---

## 示例

**场景：想要引入一个新的 npm 包**

```
1. 范围检查：
   → 新依赖不在 forbidden 里
   → 需要评估是否合理

2. 尝试检查：
   → 当前 1/5 次 → ✅

3. 质量红线检查：
   → 引入依赖 → 需要在 decision 里记录
   → 需要确保代码能处理依赖失败

4. 信号检查：
   → blockers: 1 → ✅
   → confidence: 0.7 → ✅

5. 决策：
   → ✅ 可以继续
   → 记录到 exploration-journal
```

**场景：发现代码里有 TODO**

```
1. 范围检查：N/A（不是范围决策）

2. 尝试检查：N/A

3. 质量红线检查：
   → no_todo_in_code = true
   → 发现 TODO → ❌ 不允许

4. 信号检查：N/A

5. 决策：
   → 🛑 STOP
   → 必须解决 TODO 才能继续
   → 要么实现，要么创建 issue，要么删除
```

---

## 常见错误

| 错误 | 正确做法 |
|------|----------|
| 不检查就开始做 | 必须先检查 |
| 跳过 WARN | 记录警告 |
| 忽略 forbidden | 必须停止，询问 |
| 违反质量红线 | 绝对不允许 |
| 不记录决策 | 每次决策都要记录 |
