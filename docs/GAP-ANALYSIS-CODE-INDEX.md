# 代码索引与同步机制 - 差距分析报告

> 生成日期: 2026-03-29
> 分析范围: vic ask 语义搜索、vic deps sync 索引同步、自动化程度

---

## 1. 执行摘要

### 1.1 核心问题

**VIBE-SDD 的语义搜索功能存在"索引过时"风险**：当用户修改代码后，如果不手动运行 `vic deps sync`，`vic ask` 返回的结果可能是旧代码，导致 AI 基于错误上下文做出决策。

### 1.2 竞品做法

| 产品 | 同步策略 | 自动化程度 |
|------|----------|------------|
| Cursor | 实时索引 + 增量更新 | 全自动 |
| GitHub Copilot | 云端索引 + 文件监听 | 全自动 |
| Sourcegraph | 定时全量 + 实时增量 | 全自动 |
| **VIBE-SDD (当前)** | 手动触发 `vic deps sync` | ❌ 无自动化 |

### 1.3 改进目标

```
当前状态                          目标状态
┌─────────────────┐              ┌─────────────────┐
│  手动 vic sync  │     ───▶     │  Hook自动触发   │
│  手动 vic ask   │              │  实时索引更新   │
│  可能返回旧数据  │              │  永远是最新的   │
└─────────────────┘              └─────────────────┘
     评分: 3/10                       评分: 9/10
```

---

## 2. 当前实现深度分析

### 2.1 代码架构

```
cmd/vic-go/internal/
├── commands/
│   ├── ask.go          # 语义搜索入口
│   ├── deps_sync.go    # 索引同步命令
│   └── hash.go         # SPEC变更检测
└── embedding/
    ├── store.go        # SQLite向量存储
    ├── embedder.go     # Ollama嵌入生成
    ├── sync.go         # 增量同步逻辑
    └── chunker/        # 代码分块器
        ├── go.go
        ├── python.go
        └── typescript.go
```

### 2.2 问题清单

#### 问题 1: ask.go 中 sync 未实际执行

**位置**: `cmd/vic-go/internal/commands/ask.go:67-69`

```go
// 运行增量同步以获取变更
sync := embedding.NewSync(cfg.ProjectDir, cfg.EmbeddingDir, cfg.EmbeddingIndexFile)
_ = sync // sync 在内部检查 Ollama 可用性
```

**问题**: `sync` 对象被创建后直接被忽略了，`IncrementalSync()` 从未被调用。

**修复方案**:
```go
sync := embedding.NewSync(cfg.ProjectDir, cfg.EmbeddingDir, cfg.EmbeddingIndexFile)
if _, _, _, err := sync.IncrementalSync(); err != nil {
    fmt.Println("⚠️  Index sync failed, results may be stale")
}
```

---

#### 问题 2: 删除文件检测逻辑错误

**位置**: `cmd/vic-go/internal/embedding/sync.go:240-248`

```go
// 处理已删除的文件
for _, filePath := range changedFiles {
    if !seenFiles[filePath] {
        deleted, err := store.DeleteChunksByFile(filePath)
        // ...
    }
}
```

**问题**: `changedFiles` 只包含"修改时间晚于 LastBuild 的文件"，不包含"已删除的文件"。因此这段代码永远不会执行。

**正确逻辑应该是**:
```go
// 从索引中获取所有已索引的文件
indexedFiles := store.GetAllIndexedFiles()
for _, filePath := range indexedFiles {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        // 文件已被删除
        store.DeleteChunksByFile(filePath)
        removed++
    }
}
```

---

#### 问题 3: 只使用 mtime 检测，无内容哈希

**位置**: `cmd/vic-go/internal/embedding/sync.go:183-187`

```go
mtime := info.ModTime().Unix()
if mtime > manifest.LastBuild {
    changedFiles = append(changedFiles, path)
}
```

**问题**:
- `touch file.go` 会导致不必要的重新索引
- `git checkout` 恢复文件也会触发重建
- 无法区分"内容改变"和"元数据改变"

**改进方案**:
```go
type FileHashCache struct {
    FilePath string `json:"file_path"`
    Hash     string `json:"hash"`      // SHA256(content)
    Mtime    int64  `json:"mtime"`
}

// 检测逻辑
content, _ := os.ReadFile(path)
hash := sha256.Sum256(content)
cached, exists := hashCache[path]
if !exists || cached.Hash != hex.EncodeToString(hash[:]) {
    changedFiles = append(changedFiles, path)
}
```

---

#### 问题 4: 无自动触发机制

**现状**:
- 用户修改代码 → 索引不更新 → `vic ask` 返回旧结果
- 用户需要记住手动运行 `vic deps sync`
- 没有 git hook 集成
- 没有 IDE 保存监听

**影响**: AI 在过时的代码上下文中工作，可能产生错误的代码或建议。

---

## 3. 竞品分析

### 3.1 Cursor 的索引策略

```
┌─────────────────────────────────────────────────────────┐
│                    Cursor Indexing                       │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. 启动时全量索引                                       │
│     └─ 扫描工作区，建立符号表和嵌入索引                   │
│                                                         │
│  2. 文件保存时增量更新                                   │
│     └─ VSCode onDidSaveTextDocument 事件监听            │
│     └─ 只更新变更的文件                                  │
│                                                         │
│  3. 智能防抖                                            │
│     └─ 等待 300ms 无新变更后再更新                       │
│     └─ 避免频繁重建                                     │
│                                                         │
│  4. 后台异步处理                                        │
│     └─ 不阻塞编辑器                                     │
│     └─ 进度显示在状态栏                                  │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 3.2 Sourcegraph 的索引策略

```
┌─────────────────────────────────────────────────────────┐
│                  Sourcegraph Indexing                    │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Git 集成                                            │
│     └─ 监听 git diff，只索引变更的提交                   │
│     └─ 使用 blob ID 作为缓存键                          │
│                                                         │
│  2. 符号精确索引                                        │
│     └─ 使用 tree-sitter 精确解析                        │
│     └─ 函数、类、变量级别的索引                          │
│                                                         │
│  3. 分布式索引                                          │
│     └─ 大仓库分片索引                                   │
│     └─ 并行处理加速                                     │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### 3.3 可借鉴的设计模式

| 模式 | 描述 | 适用场景 |
|------|------|----------|
| **Write-Ahead Log** | 先记录变更，后台批量处理 | 大型代码库 |
| **Content-Addressed Storage** | 用内容哈希作为索引键 | 避免重复索引 |
| **Event-Driven Update** | 文件系统事件触发更新 | 实时性要求高 |
| **Lazy Indexing** | 首次访问时才索引 | 节省资源 |

---

## 4. 改进方案

### 4.1 短期修复（立即）

#### 修复 1: ask.go 调用增量同步

```go
// cmd/vic-go/internal/commands/ask.go
func runAsk(cfg *config.Config, query string) error {
    // ... Ollama 检查 ...

    // 新增: 自动增量同步
    sync := embedding.NewSync(cfg.ProjectDir, cfg.EmbeddingDir, cfg.EmbeddingIndexFile)
    if added, updated, removed, err := sync.IncrementalSync(); err == nil {
        if added+updated+removed > 0 {
            fmt.Printf("🔄 Auto-synced: +%d ~%d -%d chunks\n", added, updated, removed)
        }
    }

    // ... 继续搜索 ...
}
```

#### 修复 2: sync.go 删除检测

```go
// cmd/vic-go/internal/embedding/sync.go
// 新增方法
func (s *Store) GetAllIndexedFiles() ([]string, error) {
    rows, err := s.db.Query("SELECT DISTINCT file_path FROM chunks")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var files []string
    for rows.Next() {
        var f string
        rows.Scan(&f)
        files = append(files, f)
    }
    return files, nil
}

// 修改 IncrementalSync 中的删除检测
func (s *Sync) IncrementalSync() (added, updated, removed int, err error) {
    // ... 现有代码 ...

    // 新增: 检测已删除的文件
    indexedFiles, _ := store.GetAllIndexedFiles()
    for _, filePath := range indexedFiles {
        if !seenFiles[filePath] {
            del, _ := store.DeleteChunksByFile(filePath)
            removed += int(del)
        }
    }

    // ...
}
```

---

### 4.2 中期改进（1-2周）

#### Git Hook 集成

```bash
# .git/hooks/pre-commit
#!/bin/bash

echo "🔄 Syncing embedding index before commit..."
vic deps sync --incremental

echo "✅ Running gate checks..."
vic gate check --blocking

# 如果失败，阻止提交
if [ $? -ne 0 ]; then
    echo "❌ Gate check failed, commit blocked"
    exit 1
fi
```

#### 自动安装 Hook 的命令

```go
// cmd/vic-go/internal/commands/hooks.go
func NewHooksInstallCmd(cfg *config.Config) *cobra.Command {
    return &cobra.Command{
        Use:   "hooks install",
        Short: "Install git hooks for auto-sync",
        RunE: func(cmd *cobra.Command, args []string) error {
            hookContent := `#!/bin/bash
vic deps sync --incremental 2>/dev/null
vic gate check --blocking
`
            hookPath := ".git/hooks/pre-commit"
            os.WriteFile(hookPath, []byte(hookContent), 0755)
            fmt.Println("✅ Git hooks installed")
            return nil
        },
    }
}
```

#### 新增 vic 命令

```bash
# 安装 git hooks
vic hooks install

# 手动触发同步（已有）
vic deps sync

# 检查索引状态
vic deps status
# 输出:
# 📊 Embedding Index Status
# ────────────────────────
# Total Chunks:    1,234
# Last Full Sync:  2026-03-29 10:30:00
# Last Sync:       2026-03-29 14:22:15
# Model:           all-minilm-l6-v2
# Files Indexed:   45
# Pending Changes: 3 files modified
```

---

### 4.3 长期方案（1-3月）

#### 文件系统监听 (fsnotify)

```go
// cmd/vic-go/internal/embedding/watcher.go
package embedding

import (
    "github.com/fsnotify/fsnotify"
    "time"
)

type SmartIndexer struct {
    watcher   *fsnotify.Watcher
    store     *Store
    embedder  *Embedder
    pending   map[string]time.Time
    debounce  time.Duration
}

func NewSmartIndexer(projectDir, indexFile string) (*SmartIndexer, error) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, err
    }

    return &SmartIndexer{
        watcher:  watcher,
        pending:  make(map[string]time.Time),
        debounce: 500 * time.Millisecond,
    }, nil
}

func (s *SmartIndexer) Start() {
    // 后台处理队列
    go s.processQueue()

    // 监听事件
    for {
        select {
        case event := <-s.watcher.Events:
            if event.Op&fsnotify.Write == fsnotify.Write {
                s.pending[event.Name] = time.Now()
            }
            if event.Op&fsnotify.Remove == fsnotify.Remove {
                s.store.DeleteChunksByFile(event.Name)
            }
        }
    }
}

func (s *SmartIndexer) processQueue() {
    ticker := time.NewTicker(s.debounce)
    for range ticker.C {
        now := time.Now()
        for file, modTime := range s.pending {
            if now.Sub(modTime) >= s.debounce {
                // 文件稳定，执行增量索引
                go s.indexFile(file)
                delete(s.pending, file)
            }
        }
    }
}
```

#### 内容哈希缓存

```go
// cmd/vic-go/internal/embedding/hash_cache.go
package embedding

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "os"
)

type HashCache struct {
    file    string
    entries map[string]CacheEntry
}

type CacheEntry struct {
    Hash     string `json:"hash"`
    Mtime    int64  `json:"mtime"`
    ChunkIDs []int  `json:"chunk_ids"` // 关联的chunk ID
}

func (c *HashCache) IsChanged(filePath string, content []byte, mtime int64) bool {
    entry, exists := c.entries[filePath]
    if !exists {
        return true
    }

    // mtime 未变，跳过
    if mtime == entry.Mtime {
        return false
    }

    // mtime 变了，检查内容哈希
    hash := sha256.Sum256(content)
    return hex.EncodeToString(hash[:]) != entry.Hash
}
```

---

## 5. 实施路线图

### Phase 1: 关键修复（本周）

| 任务 | 优先级 | 工作量 |
|------|--------|--------|
| 修复 ask.go sync 调用 | P0 | 1h |
| 修复删除检测逻辑 | P0 | 2h |
| 添加内容哈希检测 | P1 | 4h |

### Phase 2: 自动化集成（1-2周）

| 任务 | 优先级 | 工作量 |
|------|--------|--------|
| 实现 `vic hooks install` | P1 | 4h |
| 实现 `vic deps status` | P2 | 4h |
| 集成 pre-commit hook | P1 | 2h |

### Phase 3: 实时监听（1-3月）

| 任务 | 优先级 | 工作量 |
|------|--------|--------|
| 实现 fsnotify 监听 | P2 | 8h |
| 实现智能防抖 | P2 | 4h |
| 后台异步索引 | P2 | 8h |
| IDE 扩展集成 | P3 | 40h |

---

## 6. 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| Ollama 不可用 | 索引失败 | 降级到关键词搜索 |
| 大文件索引慢 | 用户体验差 | 异步处理 + 进度提示 |
| 频繁变更导致索引抖动 | 资源消耗 | 防抖 + 批量处理 |
| 索引文件损坏 | 搜索不可用 | 自动检测 + 全量重建 |

---

## 7. 验收标准

### 7.1 功能验收

- [ ] 修改代码后立即运行 `vic ask`，返回最新结果
- [ ] 删除文件后索引自动清理
- [ ] `git commit` 前自动同步索引
- [ ] 索引状态可查询

### 7.2 性能验收

| 指标 | 目标 |
|------|------|
| 增量同步延迟 | < 2s (100个文件变更) |
| ask 响应时间 | < 500ms (含同步) |
| 内存占用 | < 100MB (10万chunk) |

### 7.3 可靠性验收

- [ ] 索引损坏后自动恢复
- [ ] Ollama 不可用时优雅降级
- [ ] 并发修改不导致索引不一致

---

## 8. 附录

### A. 相关文件清单

```
需要修改的文件:
├── cmd/vic-go/internal/commands/ask.go       # 调用同步
├── cmd/vic-go/internal/embedding/sync.go     # 修复删除检测
├── cmd/vic-go/internal/embedding/store.go    # 新增GetAllIndexedFiles
└── cmd/vic-go/internal/embedding/hash_cache.go  # 新增内容哈希

需要新增的文件:
├── cmd/vic-go/internal/commands/hooks.go     # git hook管理
├── cmd/vic-go/internal/embedding/watcher.go  # 文件监听
└── cmd/vic-go/internal/embedding/hash_cache.go
```

### B. 参考资料

- [fsnotify - File System Notifications for Go](https://github.com/fsnotify/fsnotify)
- [Cursor Indexing Architecture](https://cursor.sh/docs/codebase-indexing)
- [Sourcegraph Code Intelligence](https://docs.sourcegraph.com/code_intelligence)
