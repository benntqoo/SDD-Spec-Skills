# VIBE-SDD 作为 LLM Tool Call 工具 - 质量评估报告

> 评估日期: 2026-03-29
> 评估视角: AI First, Human Verify
> 评估原则: 不增加功能，只修复现有功能使其可靠

---

## 1. 核心定位

### 1.1 工具定位

```
┌─────────────────────────────────────────────────────────────┐
│                    vic 工具定位                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   AI 的角色              人的角色                           │
│   ┌─────────────┐       ┌─────────────┐                    │
│   │   写入      │       │   查询      │                    │
│   │   编辑      │  ───▶ │   验证      │                    │
│   │   执行      │       │   确认      │                    │
│   └─────────────┘       └─────────────┘                    │
│                                                             │
│   vic 是 AI 的工具，不是人的工具                             │
│   人只负责验证 AI 做的事情对不对                             │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 设计原则

| 原则 | 描述 |
|------|------|
| **可靠性** | 每个命令必须正确执行，不能有静默失败 |
| **可预测性** | 相同输入必须产生相同输出（幂等性） |
| **原子性** | 每个命令要么完全成功，要么完全失败 |
| **清晰反馈** | 错误信息必须让AI知道下一步该怎么做 |
| **结构化输出** | 输出应该是JSON或易于解析的格式 |

---

## 2. 现有命令质量评估

### 2.1 评估标准

| 维度 | 权重 | 评分标准 |
|------|------|----------|
| 正确性 | 30% | 是否真正执行了预期的功能？ |
| 可靠性 | 25% | 是否有静默失败或未处理的错误？ |
| 幂等性 | 15% | 多次调用是否产生相同结果？ |
| 输出质量 | 15% | 输出是否清晰、结构化？ |
| 错误处理 | 15% | 错误信息是否有指导意义？ |

### 2.2 命令评估矩阵

#### 查询类命令（Human主要使用）

| 命令 | 正确性 | 可靠性 | 幂等性 | 输出质量 | 错误处理 | 总分 |
|------|--------|--------|--------|----------|----------|------|
| `vic status` | ✅ 9 | ✅ 9 | ✅ 10 | ⚠️ 7 | ⚠️ 6 | **8.2** |
| `vic spec status` | ✅ 9 | ✅ 9 | ✅ 10 | ✅ 8 | ⚠️ 6 | **8.4** |
| `vic spec hash` | ✅ 9 | ✅ 9 | ✅ 10 | ✅ 8 | ⚠️ 6 | **8.4** |
| `vic spec diff` | ⚠️ 7 | ⚠️ 7 | ✅ 10 | ⚠️ 6 | ⚠️ 5 | **7.0** |
| `vic ask` | ❌ 4 | ❌ 4 | ✅ 10 | ✅ 8 | ⚠️ 6 | **5.8** |
| `vic gate check` | ⚠️ 6 | ⚠️ 6 | ✅ 10 | ⚠️ 6 | ⚠️ 5 | **6.4** |
| `vic history` | ✅ 8 | ✅ 8 | ✅ 10 | ⚠️ 6 | ⚠️ 5 | **7.4** |
| `vic search` | ✅ 8 | ✅ 8 | ✅ 10 | ⚠️ 6 | ⚠️ 5 | **7.4** |

#### 写入类命令（AI主要使用）

| 命令 | 正确性 | 可靠性 | 幂等性 | 输出质量 | 错误处理 | 总分 |
|------|--------|--------|--------|----------|----------|------|
| `vic init` | ✅ 9 | ✅ 9 | ✅ 10 | ✅ 8 | ⚠️ 6 | **8.4** |
| `vic spec init` | ✅ 9 | ✅ 9 | ⚠️ 7 | ✅ 8 | ⚠️ 6 | **7.9** |
| `vic rt` (决策记录) | ✅ 8 | ✅ 8 | ⚠️ 6 | ⚠️ 6 | ⚠️ 5 | **6.7** |
| `vic rr` (风险记录) | ✅ 8 | ✅ 8 | ⚠️ 6 | ⚠️ 6 | ⚠️ 5 | **6.7** |
| `vic deps sync` | ❌ 4 | ❌ 4 | ⚠️ 6 | ✅ 8 | ⚠️ 5 | **5.2** |
| `vic phase advance` | ⚠️ 5 | ⚠️ 5 | ❌ 4 | ⚠️ 5 | ❌ 4 | **4.6** |
| `vic gate 0-3` | ⚠️ 5 | ⚠️ 5 | ✅ 8 | ⚠️ 5 | ❌ 4 | **5.4** |

### 2.3 关键问题清单

#### 🔴 P0 - 必须立即修复（工具不可用）

| # | 问题 | 位置 | 影响 |
|---|------|------|------|
| 1 | `vic ask` 中 sync 创建但未调用 | ask.go:67-69 | 返回过时结果 |
| 2 | `vic deps sync` 删除检测逻辑错误 | sync.go:240-248 | 索引不一致 |
| 3 | Gate 检查只验证文件存在 | gate2.go | 无法检测真正的问题 |
| 4 | 多处 `if err != nil { return }` 静默失败 | 多个文件 | AI无法知道发生了什么 |

#### 🟡 P1 - 应该尽快修复（影响可靠性）

| # | 问题 | 位置 | 影响 |
|---|------|------|------|
| 5 | 无测试覆盖 | 整个项目 | 无法保证修改不引入bug |
| 6 | 输出格式不统一 | 多个命令 | AI解析困难 |
| 7 | 错误码未定义 | 多个命令 | AI无法程序化处理错误 |
| 8 | `vic phase advance` 状态验证缺失 | phase.go | 状态转换不可靠 |

---

## 3. 系统性缺陷分析

### 3.1 架构层面

```
问题: 缺乏统一的输出/错误处理层

当前状态:
┌─────────────┐
│  Command A  │──▶ fmt.Println("✅ Success")
└─────────────┘
┌─────────────┐
│  Command B  │──▶ fmt.Printf("Error: %v\n", err)
└─────────────┘
┌─────────────┐
│  Command C  │──▶ return err  // 静默失败
└─────────────┘

应该的状态:
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Command A  │──▶  │   Output    │──▶  │  JSON/Text  │
└─────────────┘     │   Layer     │     │  统一格式   │
┌─────────────┐     │             │     └─────────────┘
│  Command B  │──▶  │  - Success  │
└─────────────┘     │  - Error    │
┌─────────────┐     │  - Warning  │
│  Command C  │──▶  │  - Progress │
└─────────────┘     └─────────────┘
```

### 3.2 数据层面

```
问题: 缺乏事务保证

场景: vic deps sync 失败中途
结果: 索引部分更新，数据不一致

应该:
1. 开始事务
2. 删除旧索引
3. 插入新索引
4. 提交事务 (或回滚)
```

### 3.3 状态层面

```
问题: 状态转换缺乏验证

当前:
vic phase advance --to 3
直接跳到阶段3，不检查前置条件

应该:
Phase 1 ─▶ Phase 2 ─▶ Phase 3
   │          │          │
   └──Gate0───┴──Gate1───┴──Gate2
              │
              ▼
         必须通过才能推进
```

---

## 4. 修复方案（按优先级）

### 4.1 Phase 1: 让工具可用（1周）

#### 修复 1: ask.go 调用 sync

```go
// 修复前 (ask.go:67-69)
sync := embedding.NewSync(...)
_ = sync  // 被忽略

// 修复后
sync := embedding.NewSync(cfg.ProjectDir, cfg.EmbeddingDir, cfg.EmbeddingIndexFile)
if added, updated, removed, err := sync.IncrementalSync(); err != nil {
    fmt.Printf("⚠️ Sync failed: %v (results may be stale)\n", err)
} else if added+updated+removed > 0 {
    fmt.Printf("🔄 Auto-synced: +%d ~%d -%d\n", added, updated, removed)
}
```

#### 修复 2: sync.go 删除检测

```go
// 修复后 (sync.go)
func (s *Sync) IncrementalSync() (added, updated, removed int, err error) {
    // ... 现有代码 ...

    // 修复: 从索引中获取所有文件，检测已删除的
    indexedFiles, _ := store.GetAllIndexedFiles()
    for _, filePath := range indexedFiles {
        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            // 文件已被删除
            del, _ := store.DeleteChunksByFile(filePath)
            removed += int(del)
        }
    }
    // ...
}

// 新增方法 (store.go)
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
```

#### 修复 3: 消除静默失败

```go
// 修复前
if err != nil {
    return  // 静默失败
}

// 修复后
if err != nil {
    return fmt.Errorf("failed to X: %w", err)  // 带上下文
}
```

### 4.2 Phase 2: 统一输出层（1周）

#### 定义输出结构

```go
// internal/output/output.go
package output

type Result struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
    Warnings []string   `json:"warnings,omitempty"`
}

type ErrorInfo struct {
    Code    string `json:"code"`     // MACHINE_READABLE_CODE
    Message string `json:"message"`  // Human readable
    Hint    string `json:"hint"`     // AI action hint
}

// 预定义错误码
const (
    ErrSpecNotFound     = "SPEC_NOT_FOUND"
    ErrOllamaUnavailable = "OLLAMA_UNAVAILABLE"
    ErrIndexCorrupted   = "INDEX_CORRUPTED"
    ErrGateFailed       = "GATE_FAILED"
)
```

#### 统一命令输出

```go
// 修复前
fmt.Println("✅ SPEC unchanged since last check")

// 修复后
output.Print(output.Result{
    Success: true,
    Data: map[string]interface{}{
        "spec_status": "unchanged",
        "last_check": record.LastCheck,
    },
})
// --format json 时输出:
// {"success":true,"data":{"spec_status":"unchanged","last_check":"2026-03-29T10:00:00Z"}}
```

### 4.3 Phase 3: 增强可靠性（2周）

#### 添加测试

```go
// internal/embedding/sync_test.go
func TestIncrementalSync_DeletedFile(t *testing.T) {
    // 创建测试文件
    // 建立索引
    // 删除文件
    // 运行 IncrementalSync
    // 验证索引中不再有该文件
}

func TestIncrementalSync_ModifiedFile(t *testing.T) {
    // 创建测试文件
    // 建立索引
    // 修改文件内容
    // 运行 IncrementalSync
    // 验证索引已更新
}
```

#### Gate 检查增强

```go
// gate2.go 增强
func runGate2(cfg *config.Config) error {
    result := &output.Result{Success: true}

    // 1. 检查 SPEC 文件存在
    spec, err := loadSpec(cfg)
    if err != nil {
        return output.Error(output.ErrSpecNotFound, err.Error())
    }

    // 2. 检查 SPEC 中定义的模块是否存在
    for _, module := range spec.Modules {
        if !dirExists(module.Path) {
            result.Warnings = append(result.Warnings,
                fmt.Sprintf("Module %s defined but not found at %s", module.Name, module.Path))
        }
    }

    // 3. 检查代码中的 TODO/FIXME
    todos, _ := scanTODOs(cfg.ProjectDir)
    if len(todos) > 0 {
        result.Warnings = append(result.Warnings,
            fmt.Sprintf("Found %d TODO/FIXME comments", len(todos)))
    }

    return result
}
```

---

## 5. 不做的事（明确边界）

### 5.1 不添加新命令

```
❌ 不做 vic hooks install
❌ 不做 vic watcher
❌ 不做 vic auto
❌ 不做 vic doctor
❌ 不做 vic dashboard
```

**理由**: 现有命令还没做好，不应该扩展功能范围

### 5.2 不添加新语言支持

```
❌ 不做 Java chunker
❌ 不做 Rust chunker
❌ 不做 C/C++ chunker
```

**理由**: 现有 chunker 的质量还没验证

### 5.3 不做复杂功能

```
❌ 不做分布式索引
❌ 不做实时文件监听
❌ 不做云同步
```

**理由**: 超出当前工具的定位

---

## 6. 验收标准

### 6.1 功能验收

```
✅ vic ask 返回的永远是最新代码
✅ vic deps sync 正确处理新增/修改/删除
✅ 所有命令都有有意义的错误信息
✅ 所有命令支持 --format json
✅ Gate 检查能检测出真正的问题
```

### 6.2 质量验收

```
✅ 核心模块测试覆盖率 > 70%
✅ 无 P0 级别的已知 bug
✅ 所有命令有文档
✅ 错误码有定义
```

### 6.3 AI 友好性验收

```
✅ JSON 输出格式稳定（AI可以解析）
✅ 错误码有对应的修复建议（AI知道下一步做什么）
✅ 命令幂等（AI可以安全重试）
```

---

## 7. 实施计划

### Week 1: P0 修复

| 天 | 任务 |
|----|------|
| 1-2 | 修复 ask.go, sync.go, store.go |
| 3 | 消除静默失败，添加错误上下文 |
| 4-5 | 添加核心模块测试 |

### Week 2: 输出层统一

| 天 | 任务 |
|----|------|
| 1-2 | 实现 output 包 |
| 2-3 | 改造所有命令使用统一输出 |
| 4-5 | 增强 Gate 检查 |

### Week 3: 可靠性增强

| 天 | 任务 |
|----|------|
| 1-3 | 完善测试覆盖 |
| 4-5 | 文档更新，验收测试 |

---

## 8. 文件修改清单

### 必须修改

```
cmd/vic-go/internal/
├── commands/
│   ├── ask.go          # 调用 sync
│   ├── deps_sync.go    # 错误处理
│   ├── gate2.go        # 增强检查
│   ├── gate0.go        # 增强检查
│   ├── gate1.go        # 增强检查
│   ├── gate3.go        # 增强检查
│   └── *.go            # 统一输出格式
├── embedding/
│   ├── sync.go         # 修复删除检测
│   └── store.go        # 新增 GetAllIndexedFiles
└── output/             # 新增包
    └── output.go
```

### 必须新增

```
cmd/vic-go/internal/
├── embedding/
│   ├── sync_test.go
│   └── store_test.go
├── commands/
│   ├── ask_test.go
│   └── gate_test.go
└── output/
    └── output_test.go
```

---

## 9. 总结

### 核心问题

**vic 工具目前不是一个可靠的 LLM Tool Call 工具**

- 静默失败多
- 输出不统一
- 缺乏测试
- 错误处理弱

### 解决方向

**不是做更多，而是做更好**

1. 修复现有 bug
2. 统一输出格式
3. 添加测试
4. 增强错误处理

### 预期结果

3周后，vic 成为：
- AI 可以可靠调用的工具
- 人可以清楚验证的工具
- 有质量保证的工具
