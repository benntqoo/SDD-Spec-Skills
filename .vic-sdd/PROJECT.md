# PROJECT: <项目名称>

> 此文档为项目状态追踪，记录了项目的整体进度、里程碑和开发状态。
> 需求规范请参考 SPEC-REQUIREMENTS.md，技术架构请参考 SPEC-ARCHITECTURE.md。

---

## 元数据

| 字段 | 值 |
|------|-----|
| version | 1.0.0 |
| status | planning / active / completed |
| owner | @agent-name |
| started | YYYY-MM-DD |
| target_release | YYYY-MM-DD |

---

## 1. 项目概览

### 1.1 项目信息

| 项目 | 内容 |
|------|------|
| 项目名称 | VIBE-SDD |
| 项目类型 | CLI工具 |
| 核心价值 | AI驱动的Spec-Driven Development开发流程工具，解决AI开发中的幻觉、盲目、失序问题 |
| 目标用户 | 使用AI辅助开发的团队/个人 |

### 1.2 技术栈

| 类别 | 技术 |
|------|------|
| CLI开发 | Go 1.21+ |
| 配置存储 | YAML |
| 验证脚本 | Python 3.10+ |
| 构建系统 | Make |

### 1.3 团队

| 角色 | Agent | 职责 |
|------|-------|------|
| Product | @agent-product | 需求分析、验收标准 |
| Architect | @agent-architect | 技术选型、架构设计 |
| Developer | @agent-develop | 代码实现、测试 |
| Lead | @agent-sisyphus | 审核、里程碑验收 |

---

## 2. 开发阶段

### 2.1 Phase 概览

| Phase | 名称 | 状态 | 开始 | 结束 | 完成度 |
|-------|------|------|------|------|--------|
| Phase 1 | MVP | in_progress | YYYY-MM-DD | - | 30% |
| Phase 2 | 扩展功能 | pending | - | - | 0% |
| Phase 3 | 优化 | pending | - | - | 0% |

### 2.2 当前 Phase 详情

#### Phase 1: MVP

**目标**: [核心价值交付]

**时间线**:
- 开始: YYYY-MM-DD
- 预计结束: YYYY-MM-DD
- 实际进度: 30%

**进度详情**:

| 功能 | 状态 | 验收标准 | 完成度 |
|------|------|---------|--------|
| F-001 用户登录 | done | 3/3 | 100% |
| F-002 用户注册 | in_progress | 2/4 | 50% |
| F-003 仪表盘 | pending | 0/5 | 0% |

---

## 3. 里程碑

### 3.1 里程碑列表

| 里程碑 | 名称 | 状态 | 日期 | 说明 |
|--------|------|------|------|------|
| M1 | 项目初始化 | completed | YYYY-MM-DD | 初始化完成 |
| M2 | MVP 功能完成 | pending | - | Phase 1 功能完成 |
| M3 | Beta 发布 | pending | - | 内测版本 |
| M4 | 正式发布 | pending | - | 生产环境 |

### 3.2 里程碑详情

#### M1: 项目初始化

| 字段 | 内容 |
|------|------|
| 状态 | ✅ completed |
| 完成日期 | YYYY-MM-DD |
| 完成条件 | [x] 仓库初始化<br>[x] 技术栈确定<br>[x] 基础架构搭建 |
| 产出 | SPEC-REQUIREMENTS.md<br>SPEC-ARCHITECTURE.md |

#### M2: MVP 功能完成

| 字段 | 内容 |
|------|------|
| 状态 | ⏳ pending |
| 目标日期 | YYYY-MM-DD |
| 完成条件 | [ ] Phase 1 所有功能完成<br>[ ] 所有验收标准通过<br>[ ] 无阻断性问题 |

---

## 4. 状态追踪

### 4.1 SPEC 状态

| 文档 | 版本 | 状态 | 最后更新 |
|------|------|------|----------|
| SPEC-REQUIREMENTS.md | 1.0.0 | spec | YYYY-MM-DD |
| SPEC-ARCHITECTURE.md | 1.0.0 | spec | YYYY-MM-DD |

### 4.2 Gate 状态

| Gate | 名称 | Feature | 状态 | 检查时间 |
|------|------|---------|------|----------|
| Gate 0 | 需求完整性 | user-auth | pass | YYYY-MM-DD |
| Gate 1 | 契约完整性 | user-auth | pass | YYYY-MM-DD |
| Gate 2 | 代码对齐 | user-auth | pending | - |
| Gate 3 | 测试覆盖 | user-auth | pending | - |

---

## 5. 问题追踪

### 5.1 阻断性问题

| ID | 问题 | 严重性 | 状态 | 负责人 |
|----|------|--------|------|--------|
| BLOCK-001 | [问题描述] | critical | open | @agent |

### 5.2 待解决技术债

| ID | 描述 | 优先级 | 阶段 |
|----|------|--------|------|
| TECH-001 | [技术债描述] | medium | Phase 2 |

---

## 6. 发布管理

### 6.1 版本规划

| 版本 | 名称 | 功能 | 目标日期 |
|------|------|------|----------|
| 0.1.0 | MVP | Phase 1 功能 | YYYY-MM-DD |
| 0.2.0 | Beta | Phase 2 功能 | YYYY-MM-DD |
| 1.0.0 | Release | 正式发布 | YYYY-MM-DD |

### 6.2 发布检查清单

- [ ] 所有 Gate 通过
- [ ] 无阻断性问题
- [ ] 安全扫描通过
- [ ] 性能测试通过
- [ ] 文档已更新

---

## 7. 风险登记

### 7.1 已识别风险

| ID | 风险 | 影响 | 概率 | 应对措施 |
|----|------|------|------|---------|
| R-001 | 第三方服务不可用 | 高 | 低 | 降级策略 |
| R-002 | 技术选型变更 | 中 | 中 | 预留方案 |

---

## 8. 会议和沟通

### 8.1 会议记录

| 日期 | 类型 | 参与者 | 议题 | 结论 |
|------|------|--------|------|------|
| YYYY-MM-DD | 需求评审 | @agent-product, @agent-architect | Phase 1 范围 | 确认范围 |

---

## 9. 变更历史

| 日期 | 变更内容 | 变更人 | 原因 |
|------|---------|--------|------|
| YYYY-MM-DD | 创建项目 | - | 初始版本 |
| YYYY-MM-DD | [变更内容] | [变更人] | [原因] |

---

## 附录

### A. 相关链接

- 仓库: [URL]
- 文档: [URL]
- 监控: [URL]

### B. 快速命令

```bash
# 查看项目状态
vic status

# 查看 SPEC 状态
vic spec status

# 查看里程碑
vic milestone list

# 运行 Gate 检查
vic spec gate 0
vic spec gate 1
vic spec gate 2
vic spec gate 3
```
