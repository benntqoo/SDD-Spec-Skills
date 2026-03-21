# Embedding Design: vic ask

## 背景

Session 3 分析确定了 embedding 在 vic-go 中的定位：

- **不是**：`vic deps scan`（静态 import 分析，地基）
- **不是**：regex 技术检测、gate 检查、workflow state CRUD
- **是**：语义代码理解——自然语言查询代码含义

```
deps scan          ← 地基：静态结构（模块、import 关系、调用图）
    ↓
vic ask + embedding ← 迭代层：语义索引、chunk 持续更新、查询增强
    ↓
未来扩展           ← 如 vic deps impact --semantic
```

## 设计决策（已确认）

| # | 问题 | 决策 |
|---|------|------|
| 1 | 接口 | 独立 `vic ask` 命令 |
| 2 | 存储 | SQLite + 纯 Go 实现（无 CGO） |
| 3 | 分块策略 | 按语义段落（func/class/def）切，不打散函数意图 |
| 4 | init 集成 | `vic init` 后台非阻塞构建 |
| 5 | 增量更新 | `vic ask` 触发式增量更新 + git hook 辅助 |

## 技术选型

### Embedding 模型
- **模型**：`all-MiniLM-L6-v2`（384 维，90MB，CPU-only，无 API key）
- **纯 Go 调用**：通过 `github.com/go-skynet/go-llama.cpp` 或 subprocess 调用 `llama.cpp` 二进制
- **备选**：使用 `github.com/tmc/langchaingo/embeddings` 调用本地 OLLAMA API
- **约束**：必须完全离线、无需 API key、纯 Go 依赖

### 向量存储
- **库**：`modernc.org/sqlite`（纯 Go，无 CGO）
- **路径**：`.vic-sdd/embeddings/index.sqlite`
- **Schema**：
```sql
CREATE TABLE chunks (
    id          INTEGER PRIMARY KEY,
    file_path   TEXT NOT NULL,
    chunk_type  TEXT NOT NULL,  -- func, class, def, struct, module
    chunk_name  TEXT NOT NULL,  -- function/class name
    module_path TEXT NOT NULL,  -- e.g. internal/commands
    start_line  INTEGER,
    end_line    INTEGER,
    code        TEXT NOT NULL,
    doc         TEXT,
    lang        TEXT,
    updated_at  INTEGER         -- Unix timestamp
);
CREATE TABLE vectors (
    chunk_id INTEGER PRIMARY KEY,
    vector   BLOB              -- 384 float32, little-endian
);
CREATE INDEX idx_file ON chunks(file_path);
CREATE INDEX idx_module ON chunks(module_path);
```

### Chunk 提取策略

**Go**：
- 按 `func`、`type`、`const`、`var` 分块
- 每个 export 元素（首字母大写）独立 chunk
- 保留完整签名 + doc + body 摘要

**Python**：
- 按 `class`、`def`、`async def` 分块
- 顶层 `import` 单独 chunk

**TypeScript/JavaScript**：
- 按 `function`、`class`、`const` (exported)、`interface`、`type` 分块

### Index 存储
```
.vic-sdd/embeddings/
├── index.sqlite      -- SQLite 向量数据库
├── manifest.json     -- index 元数据（版本、维度、chunk 数、最后更新时间）
└── model/            -- 嵌入模型文件（可选，简化版用 subprocess）
    └── (model files)
```

## 命令接口

### `vic ask "<自然语言查询>"`
```bash
# 基础查询
vic ask "database connection pooling"

# 输出示例
🔍 Query: database connection pooling
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📁 internal/db/pool.go:14-45 [func:NewConnectionPool]
   ├─ module: internal/db
   └─ "NewConnectionPool creates a connection pool with configurable
       min/max connections, connection timeout, and idle eviction..."

📁 internal/db/pool.go:52-78 [func:Acquire]
   ├─ module: internal/db
   └─ "Acquire returns a connection from the pool..."

📁 internal/repo/user.go:12-30 [func:WithTransaction]
   └─ "Wraps database operations in a transaction..."
```

### `vic deps sync`
```bash
# 增量更新 embedding index
vic deps sync

# 强制全量重建
vic deps sync --full

# 输出
🔄 Syncing embedding index...
   ✅ 142 chunks updated
   ✅ 3 new chunks added
   ✅ 0 chunks removed
   ✅ Completed in 12s
```

## 架构

```
cmd/vic-go/
├── internal/
│   ├── embedding/
│   │   ├── chunker/          # 代码分块（语言感知）
│   │   │   ├── go.go
│   │   │   ├── python.go
│   │   │   ├── typescript.go
│   │   │   └── chunker.go    # 接口定义
│   │   ├── embedder.go      # embedding 模型调用
│   │   ├── store.go          # SQLite vector store
│   │   └── sync.go          # 增量同步逻辑
│   ├── commands/
│   │   ├── ask.go            # vic ask 命令
│   │   └── deps_sync.go      # vic deps sync 命令
│   └── config/
│       └── config.go         # 新增 EmbeddingDir 配置
```

## 关键设计原则

1. **无感知增量更新**：`vic ask` 自动检测变更，必要时触发更新
2. **完全离线**：不需要任何 API key 或网络请求
3. **向后兼容**：`vic deps scan` 继续正常工作
4. **轻量优先**：暴力搜索（brute-force），10k chunks 以内 <1s

## 文件变更清单

| 操作 | 文件 |
|------|------|
| 新增 | `cmd/vic-go/internal/embedding/chunker/chunker.go` |
| 新增 | `cmd/vic-go/internal/embedding/chunker/go.go` |
| 新增 | `cmd/vic-go/internal/embedding/chunker/python.go` |
| 新增 | `cmd/vic-go/internal/embedding/chunker/typescript.go` |
| 新增 | `cmd/vic-go/internal/embedding/embedder.go` |
| 新增 | `cmd/vic-go/internal/embedding/store.go` |
| 新增 | `cmd/vic-go/internal/embedding/sync.go` |
| 新增 | `cmd/vic-go/internal/commands/ask.go` |
| 新增 | `cmd/vic-go/internal/commands/deps_sync.go` |
| 修改 | `cmd/vic-go/internal/commands/init.go` |
| 修改 | `cmd/vic-go/internal/commands/root.go` |
| 修改 | `cmd/vic-go/internal/config/config.go` |
| 修改 | `cmd/vic-go/go.mod` |
