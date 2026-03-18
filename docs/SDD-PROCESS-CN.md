# VIBE-SDD 落地型开发流程规范

> 本规范解决AI开发中的三大核心问题：
> - **防幻觉**：每个决策必须有据可查
> - **防盲目**：每个阶段必须有明确产出
> - **防失序**：进度必须透明可追溯

---

## 1. 核心理念

### 1.1 设计原则

| 原则 | 说明 |
|------|------|
| **强制外部验证** | 关键决策必须经过验证，不能自己说自己对 |
| **产出驱动** | 每个阶段必须有明确产出物，无产出=未完成 |
| **可追溯** | 所有变更必须记录原因和依据 |
| **进度可视** | 每次交互必须声明当前状态 |

### 1.2 精简阶段模型

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   PHASE 0   │───▶│   PHASE 1   │───▶│   PHASE 2   │───▶│   PHASE 3   │
│   需求凝固   │    │   架构设计   │    │   代码实现   │    │   验证发布   │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
       │                  │                  │                  │
       ▼                  ▼                  ▼                  ▼
   Gate 0-1           Gate 2-3           Gate 4-5           Gate 6-7
```

---

## 2. 阶段详解

### PHASE 0: 需求凝固

**目标**：将模糊需求转化为可执行的规格文档

**输入**：用户原始需求（可能模糊）

**输出**：
- `SPEC-REQUIREMENTS.md` - 完整需求文档
- `SPEC-CONTRACT.json` - 接口契约（可选）

**Gate检查**：
| 检查项 | 要求 | 状态 |
|--------|------|------|
| User Story | 至少1条 | □ |
| Acceptance Criteria | 每条Story至少2个 | □ |
| 验收标准可测试 | 是/否明确 | □ |
| 边界条件 | 已识别 | □ |

**必须完成**：
- [ ] User Story 完整（who/what/why）
- [ ] 每个Story对应Acceptance Criteria
- [ ] 识别技术依赖
- [ ] 评估风险点

**记录要求**：
```yaml
# 记录到 .vic-sdd/tech/tech-records.yaml
- id: REQ-001
  title: "[User Story标题]"
  story: "As a [角色], I want [功能], so that [价值]"
  acceptance_criteria:
    - "AC1: [可测试标准]"
    - "AC2: [可测试标准]"
  status: frozen
  created_at: "YYYY-MM-DD"
```

---

### PHASE 1: 架构设计

**目标**：确定技术方案和系统设计

**输入**：SPEC-REQUIREMENTS.md

**输出**：
- `SPEC-ARCHITECTURE.md` - 技术架构文档
- 技术选型决策记录

**Gate检查**：
| 检查项 | 要求 | 状态 |
|--------|------|------|
| 技术栈确定 | 语言/框架/数据库 | □ |
| 模块划分 | 清晰定义 | □ |
| 数据模型 | 包含关系 | □ |
| API设计 | RESTful/GraphQL | □ |
| 安全性考虑 | 认证/授权/加密 | □ |

**必须完成**：
- [ ] 技术选型并记录理由
- [ ] 模块/组件划分
- [ ] 数据库schema（如涉及）
- [ ] API接口定义（如涉及）
- [ ] 识别技术风险

**记录要求**：
```yaml
# 记录到 .vic-sdd/tech/tech-records.yaml
- id: ARCH-001
  title: "[架构决策标题]"
  decision: "选择[技术X]作为[用途]"
  alternatives_considered:
    - "方案A: [描述] - 不选原因"
    - "方案B: [描述] - 不选原因"
  reason: "选择理由（必须有数据支撑或经验依据）"
  impact: low/medium/high
  files_affected: ["file1", "file2"]
  status: approved
  created_at: "YYYY-MM-DD"
```

---

### PHASE 2: 代码实现

**目标**：基于架构实现功能

**输入**：
- SPEC-REQUIREMENTS.md
- SPEC-ARCHITECTURE.md

**输出**：
- 源代码
- 单元测试（如适用）

**Gate检查**：
| 检查项 | 要求 | 状态 |
|--------|------|------|
| 代码编译通过 | 无编译错误 | □ |
| 对齐SPEC | 功能vs需求 | □ |
| 单元测试 | 核心逻辑覆盖 | □ |
| 代码规范 | linter通过 | □ |

**必须完成**：
- [ ] 实现所有Acceptance Criteria
- [ ] 保持代码与架构一致
- [ ] 编写必要的测试
- [ ] 确保代码可运行

**进度追踪**：
```
状态声明格式（每次响应必须包含）：

---
**当前阶段**: PHASE 2 - 代码实现
**当前Gate**: Gate 4 (代码实现中)
**已完成**:
  - [x] REQ-001 / AC1 实现
  - [x] REQ-001 / AC2 实现
**进行中**:
  - [~] REQ-002 / AC1 实现 (60%)
**待完成**:
  - [ ] REQ-002 / AC2 实现
  - [ ] REQ-003 全部
**卡点**: [如有]
---
```

---

### PHASE 3: 验证发布

**目标**：确保实现符合需求，可发布

**输入**：完整代码 + 测试

**输出**：
- 测试报告
- 发布检查清单
- Release Notes

**Gate检查**：
| 检查项 | 要求 | 状态 |
|--------|------|------|
| 功能测试 | 全部AC通过 | □ |
| 代码审查 | 无重大问题 | □ |
| 安全检查 | 无高危漏洞 | □ |
| 文档更新 | 同步代码变更 | □ |

**必须完成**：
- [ ] 所有Acceptance Criteria验证通过
- [ ] 代码审查（如适用）
- [ ] 安全扫描通过
- [ ] 文档已更新

---

## 3. 进度追踪机制

### 3.1 状态文件

在 `.vic-sdd/status/` 下维护：

```
.vic-sdd/status/
├── state.yaml          # 全局状态
├── current-phase.yaml  # 当前阶段
├── gate-status.yaml    # Gate检查状态
└── progress.yaml       # 详细进度
```

### 3.2 状态声明规范

**每次响应必须声明当前状态**，格式：

```yaml
# 当前阶段状态
current_phase: 2  # 0-3
current_gate: 4   # 0-7

# 阶段详情
phase_details:
  phase_0:
    status: completed
    completed_at: "YYYY-MM-DD"
    outputs:
      - SPEC-REQUIREMENTS.md
      - SPEC-CONTRACT.json
  
  phase_1:
    status: completed
    completed_at: "YYYY-MM-DD"
    outputs:
      - SPEC-ARCHITECTURE.md
      - tech-records.yaml (updated)
  
  phase_2:
    status: in_progress
    started_at: "YYYY-MM-DD"
    completion: 60%
    progress:
      - feature: "REQ-001"
        status: completed
      - feature: "REQ-002"
        status: in_progress
        progress: 60%
      - feature: "REQ-003"
        status: pending
  
  phase_3:
    status: pending

# Gate检查状态
gates:
  gate_0: passed  # 需求完整性
  gate_1: passed  # 需求可测试
  gate_2: passed  # 架构完整性
  gate_3: passed  # 技术选型合理
  gate_4: passed  # 代码编译通过
  gate_5: passed  # 代码对齐SPEC
  gate_6: pending  # 功能测试
  gate_7: pending  # 发布检查
```

### 3.3 快速查看命令

```bash
# 查看当前状态
cat .vic-sdd/status/current-phase.yaml

# 查看Gate状态
cat .vic-sdd/status/gate-status.yaml

# 查看详细进度
cat .vic-sdd/status/progress.yaml
```

---

## 4. 决策记录规范

### 4.1 技术决策记录

每个技术决策必须记录：

```yaml
- id: [类型-序号，如 ARCH-001, DB-001 ]
  title: "[决策标题]"
  decision: |
    [决策内容，描述选择了什么方案]
  context: |
    [背景：为什么需要做这个决定]
  alternatives_considered:
    - name: "[方案A]"
      pros: "[优点]"
      cons: "[缺点]"
      rejected_reason: "[为什么被拒绝]"
    - name: "[方案B]"
      pros: "[优点]"
      cons: "[缺点]"
      rejected_reason: "[为什么被拒绝]"
  reason: "[最终选择理由，必须有依据]"
  evidence: |
    [证据：可以是文档链接、性能数据、经验总结等]
  impact: low/medium/high
  files_affected: ["file1", "file2"]
  status: proposed/approved/completed
  created_at: "YYYY-MM-DD"
  decided_by: "agent/human"
```

### 4.2 风险记录

```yaml
- id: RISK-001
  title: "[风险标题]"
  description: "[风险描述]"
  probability: low/medium/high
  impact: low/medium/high
  mitigation: "[应对措施]"
  status: identified/monitoring/mitigated/occurred
  created_at: "YYYY-MM-DD"
```

---

## 5. 强制检查点

### 5.1 进入下一阶段的条件

| 从 → 到 | 必须满足 |
|---------|---------|
| Phase 0 → 1 | Gate 0, 1 通过 |
| Phase 1 → 2 | Gate 2, 3 通过 |
| Phase 2 → 3 | Gate 4, 5 通过 |
| Phase 3 → 完成 | Gate 6, 7 通过 |

### 5.2 禁止事项

| 禁止 | 说明 |
|------|------|
| 跳过Gate检查 | 每个Gate必须明确通过/失败状态 |
| 无记录决策 | 技术决策必须写入tech-records.yaml |
| 模糊需求进入实现 | Phase 0未完成禁止进入Phase 2 |
| 变更不记录 | 任何需求/架构变更必须更新文档 |

---

## 6. 工具支持

### 6.1 建议的CLI命令

```bash
# 初始化项目
sdd init --name "项目名" --phase 0

# 记录决策
sdd decision --id ARCH-001 --title "选择技术X" --decision "描述"

# 记录风险
sdd risk --title "风险描述" --impact high --mitigation "措施"

# 阶段推进
sdd phase advance --to 1 --gate-status pass

# 查看状态
sdd status

# 检查Gate
sdd gate check --gate 0
```

### 6.2 手动查看

```bash
# 查看当前阶段
cat .vic-sdd/status/current-phase.yaml

# 查看技术决策
cat .vic-sdd/tech/tech-records.yaml

# 查看风险
cat .vic-sdd/risk-zones.yaml

# 查看项目状态
cat .vic-sdd/PROJECT.md
```

---

## 7. 完整流程示例

### 示例：添加用户认证功能

#### Step 1: 需求凝固 (Phase 0)

```
用户输入: "我需要添加用户登录注册功能"

AI响应（必须包含状态声明）:

---
**当前阶段**: PHASE 0 - 需求凝固
**当前Gate**: Gate 0-1 (需求分析)
**待完成**:
  - [ ] 编写User Story
  - [ ] 定义Acceptance Criteria
  - [ ] 识别技术依赖
---

输出文档: SPEC-REQUIREMENTS.md
```

#### Step 2: 架构设计 (Phase 1)

```
---
**当前阶段**: PHASE 1 - 架构设计
**当前Gate**: Gate 2-3 (技术选型)
**已完成**:
  - [x] User Story: REQ-001
  - [x] Acceptance Criteria: 4条
**待完成**:
  - [ ] 技术选型决策
  - [ ] 模块划分
---

输出:
- SPEC-ARCHITECTURE.md
- tech-records.yaml (更新)
```

#### Step 3: 代码实现 (Phase 2)

```
---
**当前阶段**: PHASE 2 - 代码实现
**当前Gate**: Gate 4-5 (编码中)
**已完成**:
  - [x] 技术选型: JWT
  - [x] 模块: auth/controller, auth/service, auth/model
**进行中**:
  - [~] 登录API实现 (80%)
  - [~] 注册API实现 (60%)
**待完成**:
  - [ ] 单元测试
  - [ ] 集成测试
---

输出: 源代码 + 测试
```

#### Step 4: 验证发布 (Phase 3)

```
---
**当前阶段**: PHASE 3 - 验证发布
**当前Gate**: Gate 6-7 (测试验证)
**已完成**:
  - [x] 登录功能
  - [x] 注册功能
  - [x] 单元测试 (覆盖率85%)
**待完成**:
  - [ ] 集成测试
  - [ ] 安全扫描
  - [ ] 代码审查
---

输出: 测试报告 + Release Notes
```

---

## 8. 附录

### 8.1 Gate检查清单

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

### 8.2 状态码

| 状态 | 说明 |
|------|------|
| pending | 未开始 |
| in_progress | 进行中 |
| completed | 已完成 |
| blocked | 被阻断 |
| passed | 检查通过 |
| failed | 检查失败 |

### 8.3 快速参考

```
# 核心原则
1. 产出驱动：无产出=未完成
2. 记录可追溯：决策必记录
3. 进度透明：每次响应必声明状态
4. Gate必过：禁止跳阶段

# 状态声明模板
---
**当前阶段**: PHASE X - [名称]
**当前Gate**: Gate Y
**已完成**: [清单]
**进行中**: [清单 + 进度%]
**待完成**: [清单]
**卡点**: [如有]
---
```

---

**文档版本**: 1.0.0  
**创建日期**: 2026-03-18  
**维护者**: Sisyphus AI Agent
