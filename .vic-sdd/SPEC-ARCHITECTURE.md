# SPEC-ARCHITECTURE: <项目名称>

> 此文档为技术架构规范，定义了项目的技术选型、系统架构和数据模型。
> 详细需求请参考 SPEC-REQUIREMENTS.md。

---

## 元数据

| 字段 | 值 |
|------|-----|
| version | 1.0.0 |
| status | draft / spec / build / verify / done |
| owner | @agent-name |
| created | YYYY-MM-DD |
| updated | YYYY-MM-DD |

---

## 1. 技术选型

### 1.1 技术栈总览

| 层级 | 技术 | 版本 | 状态 |
|------|------|------|------|
| 前端框架 | [React/Vue/其他] | x.x.x | 选中 |
| 状态管理 | [Zustand/Redux/Pinia] | x.x.x | 选中 |
| 后端框架 | [Next.js/Express/其他] | x.x.x | 选中 |
| 数据库 | [PostgreSQL/MySQL/其他] | x.x.x | 选中 |
| ORM | [Prisma/Drizzle/其他] | x.x.x | 选中 |
| 认证 | [JWT/Session/其他] | - | 选中 |
| 部署 | [Vercel/AWS/其他] | - | 选中 |

### 1.2 技术选型评估

#### 前端框架

| 技术 | 优点 | 缺点 | 适用场景 |
|------|------|------|---------|
| React | 生态大、灵活性高 | 学习曲线陡 | 中大型项目 |
| Vue | 上手简单、文档好 | 灵活性略低 | 快速开发 |
| Svelte | 性能好、代码少 | 生态较小 | 轻量项目 |

**最终选择**: [技术名称]
**选择理由**: [简明理由]

#### 数据库

| 技术 | 优点 | 缺点 | 适用场景 |
|------|------|------|---------|
| PostgreSQL | ACID、JSON支持 | 资源需求高 | 事务性数据 |
| MySQL | 成熟稳定 | JSON支持一般 | 通用场景 |
| MongoDB | 灵活、易扩展 | 事务弱 | 文档数据 |

**最终选择**: [技术名称]
**选择理由**: [简明理由]

### 1.3 开发工具

| 工具 | 用途 | 选择 |
|------|------|------|
| 语言 | 开发语言 | TypeScript |
| 包管理 | 依赖管理 | pnpm |
| 代码规范 | Linting | ESLint |
| 代码格式 | Formatting | Prettier |
| 测试 | 单元测试 | Vitest |
| 构建 | 打包 | Vite |

---

## 2. 系统架构

### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         客户端                              │
│         (Web / Mobile / Desktop)                            │
│                                                              │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│   │   页面 A    │  │   页面 B    │  │   页面 C    │       │
│   └──────┬──────┘  └──────┬──────┘  └──────┬──────┘       │
│          │                │                │              │
│          └────────────────┼────────────────┘              │
│                           ▼                                │
│                    ┌─────────────┐                         │
│                    │  状态管理   │                         │
│                    └──────┬──────┘                         │
└───────────────────────────┼────────────────────────────────┘
                            │ HTTPS
                            ▼
┌───────────────────────────────────────────────────────────────┐
│                      接入层                                  │
│                                                              │
│   ┌─────────────────────────────────────────────────────┐   │
│   │                   Next.js / Express                   │   │
│   │   - SSR / CSR                                       │   │
│   │   - API Routes / Controllers                        │   │
│   │   - Middleware (鉴权、日志)                          │   │
│   └─────────────────────────────────────────────────────┘   │
│                           │                                 │
└───────────────────────────┼────────────────────────────────┘
                            │
              ┌─────────────┴─────────────┐
              ▼                           ▼
┌─────────────────────────┐   ┌─────────────────────────────┐
│       业务服务层          │   │        外部服务             │
│                          │   │                             │
│  ┌───────────────────┐  │   │  ┌─────────────┐           │
│  │   Auth Service    │  │   │  │  支付网关    │           │
│  │   - 登录/注册    │  │   │  └─────────────┘           │
│  │   - Token 签发   │  │   │  ┌─────────────┐           │
│  └───────────────────┘  │   │  │  邮件服务    │           │
│  ┌───────────────────┐  │   │  └─────────────┘           │
│  │   User Service   │  │   │  ┌─────────────┐           │
│  │   - 用户资料     │  │   │  │    CDN      │           │
│  │   - 权限管理     │  │   │  └─────────────┘           │
│  └───────────────────┘  │   │                             │
│  ┌───────────────────┐  │   │                             │
│  │  Order Service   │  │   │                             │
│  └───────────────────┘  │   │                             │
└───────────┬─────────────┘   └─────────────────────────────┘
            │
            ▼
┌─────────────────────────────────────────────────────────────┐
│                       数据层                                │
│                                                              │
│   ┌─────────────┐        ┌─────────────┐                   │
│   │ PostgreSQL  │        │    Redis    │                   │
│   │   (主库)   │        │   (缓存)    │                   │
│   └─────────────┘        └─────────────┘                   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 模块划分

| 模块 | 职责 | 依赖 | 边界 |
|------|------|------|------|
| auth | 用户认证、JWT 签发和验证 | 数据库 | 服务端 |
| user | 用户资料管理、权限 | auth | 服务端 |
| order | 订单管理、状态流转 | user, payment | 服务端 |
| notification | 消息通知、推送 | 外部服务 | 服务端 |

---

## 3. 数据模型

### 3.1 核心实体

#### User (用户)

```typescript
interface User {
  id: UUID;           // 主键
  email: string;      // 唯一索引
  passwordHash: string;
  name: string;
  role: 'admin' | 'user';
  createdAt: DateTime;
  updatedAt: DateTime;
}
```

#### Order (订单)

```typescript
interface Order {
  id: UUID;                    // 主键
  userId: UUID;               // 外键 -> User
  status: OrderStatus;        // 枚举
  totalAmount: Decimal;
  items: OrderItem[];
  createdAt: DateTime;
  updatedAt: DateTime;
}

enum OrderStatus {
  PENDING = 'pending',
  PAID = 'paid',
  CANCELLED = 'cancelled',
  REFUNDED = 'refunded'
}
```

### 3.2 关系图

```
┌─────────┐       ┌─────────┐
│  User   │ 1   ∞ │  Order  │
│ (用户)  │───────│ (订单)  │
└─────────┘       └────┬────┘
                      │
                      │ 1
                      │ 
                      ∞
                 ┌────┴────┐
                 │OrderItem │
                 │(订单项)  │
                 └─────────┘
```

### 3.3 索引设计

| 表 | 索引类型 | 字段 | 说明 |
|----|---------|------|------|
| User | unique | email | 登录查询 |
| User | index | createdAt | 排序 |
| Order | index | userId | 用户订单查询 |
| Order | index | status + createdAt | 状态筛选 |

---

## 4. API 设计

### 4.1 API 目录

| 路由 | 方法 | 功能 | 鉴权 | 状态 |
|------|------|------|------|------|
| /api/auth/login | POST | 登录 | 否 | done |
| /api/auth/register | POST | 注册 | 否 | done |
| /api/auth/refresh | POST | 刷新Token | 是 | done |
| /api/users | GET | 用户列表 | 是 | done |
| /api/users/:id | GET | 用户详情 | 是 | done |

### 4.2 API 契约示例

#### POST /api/auth/login

**请求**
```json
{
  "email": "user@example.com",
  "password": "string"
}
```

**响应 (200)**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expiresIn": 3600,
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "用户名"
    }
  }
}
```

**错误响应**
| 状态码 | 场景 | 响应 |
|--------|------|------|
| 401 | 密码错误 | `{"success": false, "error": "Invalid credentials"}` |
| 404 | 用户不存在 | `{"success": false, "error": "User not found"}` |

### 4.3 错误码规范

```typescript
// 错误码格式: 模块_序号
// 认证模块: AUTH_001, AUTH_002
// 用户模块: USER_001, USER_002

enum ErrorCode {
  // 认证模块
  AUTH_001 = 'Invalid credentials',
  AUTH_002 = 'Token expired',
  AUTH_003 = 'Token invalid',
  
  // 用户模块
  USER_001 = 'User not found',
  USER_002 = 'Email already exists',
}
```

---

## 5. 安全设计

### 5.1 鉴权授权

#### 认证流程

```
1. 用户提交 credentials
2. 服务端验证
3. 生成 JWT (access + refresh)
4. 返回 token
5. 客户端存储
6. 后续请求携带 access token
7. 服务端验证 token
8. 放行或拒绝
```

#### 权限模型

| 角色 | 权限 |
|------|------|
| admin | 全部权限 |
| user | 自身资源读写 |

### 5.2 安全清单

| 项 | 实现 | 说明 |
|---|------|------|
| HTTPS | ✅ | 全站强制 |
| HSTS | ✅ | 启用 |
| JWT 过期 | ✅ | 1小时 |
| 刷新Token | ✅ | 7天 |
| 密码哈希 | ✅ | bcrypt (12 rounds) |
| SQL注入 | ✅ | ORM 参数化查询 |
| XSS | ✅ | React 默认防护 |
| CSRF | ✅ | SameSite Cookie |
| CORS | ✅ | 白名单配置 |
| 速率限制 | ✅ | 100req/15min |

### 5.3 敏感数据处理

| 数据 | 存储 | 传输 | 日志 |
|------|------|------|------|
| 密码 | bcrypt 哈希 | - | 禁止 |
| Token | - | HttpOnly Cookie | 禁止 |
| 信用卡 | 第三方处理 | 加密 | 禁止 |
| PII | 加密存储 | HTTPS | 脱敏 |

---

## 6. 服务端边界

### 6.1 必须服务端处理

| 场景 | 原因 |
|------|------|
| 用户认证 | 安全性 |
| 支付处理 | 敏感操作 |
| 权限校验 | 业务规则 |
| 业务逻辑 | 数据一致性 |
| 第三方API调用 | 密钥安全 |

### 6.2 可客户端处理

| 场景 | 说明 |
|------|------|
| 表单验证 | 提升体验 |
| UI 状态 | 响应速度 |
| 缓存同步 | 离线支持 |
| 简单计算 | 减轻服务端 |

### 6.3 服务间通信

```
┌──────────────┐     ┌──────────────┐
│   Service A  │────▶│   Service B  │
│              │ HTTP│              │
└──────────────┘     └──────────────┘
        │
        │ 可选: 消息队列
        ▼
┌──────────────┐
│   Service C  │
└──────────────┘
```

---

## 7. 部署和运维

### 7.1 环境

| 环境 | 用途 | URL |
|------|------|-----|
| Dev | 开发自测 | dev.local |
| Staging | 预发布测试 | staging.example.com |
| Prod | 生产环境 | example.com |

### 7.2 部署流程

```
Git Push → CI Build → Test → Deploy to Staging → Manual Verify → Deploy to Prod
```

### 7.3 监控

| 工具 | 用途 |
|------|------|
| Sentry | 错误追踪 |
| Vercel Analytics | 性能监控 |
| Datadog | 基础设施 |

### 7.4 备份

| 数据 | 策略 |
|------|------|
| 数据库 | 每日全量 + 增量 |
| 文件 | 每日同步 |

---

## 8. 目录结构

```
project/
├── src/
│   ├── components/      # 公共组件
│   ├── pages/           # 页面 (或 routes/)
│   ├── services/        # 服务层
│   ├── utils/          # 工具函数
│   ├── hooks/          # 自定义 Hooks
│   ├── types/          # 类型定义
│   └── styles/         # 样式
│
├── server/
│   ├── api/            # API 路由
│   ├── services/       # 业务服务
│   ├── models/         # 数据模型
│   ├── middleware/     # 中间件
│   └── utils/         # 服务端工具
│
├── prisma/             # ORM 配置
├── tests/             # 测试文件
└── docs/              # 文档
```

---

## 9. 变更历史

| 日期 | 变更内容 | 变更人 | 原因 |
|------|---------|--------|------|
| YYYY-MM-DD | 创建文档 | - | 初始版本 |
| YYYY-MM-DD | [变更内容] | [变更人] | [原因] |

---

## 附录

### A. 相关文档

- SPEC-REQUIREMENTS.md - 需求规范
- PROJECT.md - 项目状态追踪

### B. 参考资料

- [技术文档链接]
- [架构设计参考]

### C. 术语表

| 术语 | 定义 |
|------|------|
| JWT | JSON Web Token |
| ORM | Object-Relational Mapping |
| ACID | Atomicity, Consistency, Isolation, Durability |
