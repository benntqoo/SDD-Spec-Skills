# SPEC-ARCHITECTURE.md

## Architecture Overview

> Generated: 2026-04-01

> VIBE-SDD is a CLI-based Spec-Driven Development system with a layered architecture. The core CLI (vic) manages project state, SPEC documents, and gate checks. It uses a modular design with separate packages for configuration, output formatting, code checking, and utilities.

## System Design

### Components

**CLI Layer (`cmd/vic-go`)**
- Main entry point and command registration
- Command handlers for all user-facing operations
- Cobra-based command structure

**Core Layer (`internal/`)**
- `commands/` - Command implementations (init, spec, gate, hooks, etc.)
- `config/` - Configuration management with Viper
- `output/` - Unified output formatting (plain, JSON, YAML)
- `checker/` - Code alignment checking and validation
- `embedding/` - Semantic search with SQLite vector store

**Data Layer**
- YAML files in `.vic-sdd/` directory
- SQLite database for vector embeddings
- Git hooks for integration

**Skills System**
- Modular skill definitions in `skills/` directory
- Skill validation and execution
- Traceability tracking

### Data Flow

1. **User Input** → CLI commands
2. **Command Processing** → Appropriate handler invoked
3. **State Management** → YAML/SQLite updates
4. **Gate Validation** → Checks run and results returned
5. **Output Formatting** → Results formatted and displayed
6. **Git Integration** → Pre-commit hooks enforce quality

## Technology Stack

| Layer | Technology | Rationale |
|-------|------------|------------|
| CLI Framework | Cobra | De facto standard for Go CLIs, excellent subcommand support |
| Configuration | Viper | Flexible configuration from files, env vars, flags |
| Output | Custom `output` package | Unified formatting for all commands |
| Persistence | YAML, SQLite | Human-readable configs + efficient vector search |
| Code Scanning | Regexp, filepath | Built-in Go packages for pattern matching |
| Vector Store | SQLite + embeddings | Fast local semantic search without external deps |
| Git Integration | Bash hooks, git commands | Platform-agnostic approach |
| Testing | Go testing framework | Built-in, well-maintained test runner |

## Data Model

### Project State (`.vic-sdd/project.yaml`)
```yaml
name: string
description: string
tech_stack: []string
created_at: string
current_phase: int
```

### Phase Status (`.vic-sdd/status/phase.yaml`)
```yaml
current_phase: int
phases:
  - id: int
    name: string
    status: string
    gates:
      - id: string
        name: string
        status: string  # pending/passed/failed
        checked_at: string
        checked_by: string
        notes: string
```

### SPEC Documents
- `.vic-sdd/SPEC-REQUIREMENTS.md` - Requirements specification
- `.vic-sdd/SPEC-ARCHITECTURE.md` - Architecture specification

### Dependency Graph (`.vic-sdd/dependency-graph.yaml`)
```yaml
nodes:
  - id: string
    type: string  # spec/feature/gate
    name: string
    file: string
edges:
  - from: string
    to: string
    type: string  # depends_on/validates/tests
```

### Events (`.vic-sdd/status/events.yaml`)
```yaml
events:
  - timestamp: string
    type: string  # state_change/gate_result/decision/risk
    phase: int
    description: string
    metadata: map[string]interface{}
```

## API Design

### Command Structure

**Root Commands**
- `vic init` - Initialize new project
- `vic status` - Show project status
- `vic spec` - SPEC management subcommands
- `vic gate` - Gate management subcommands
- `vic hooks` - Git hooks management

**Spec Subcommands**
- `vic spec init` - Initialize SPEC files
- `vic spec status` - Show SPEC status
- `vic spec gate <n>` - Run specific gate check
- `vic spec hash` - Check SPEC hashes

**Gate Subcommands**
- `vic gate status` - Show all gate status
- `vic gate check` - Run gate check
- `vic gate pass --gate <n>` - Manually pass gate
- `vic gate smart` - Smart gate selection

### Internal APIs

**Gate Report API**
```go
type GateReport struct {
    GateNumber     int
    Checks         []CheckResult
    PassedChecks    int
    TotalChecks    int
}

type CheckResult struct {
    ID          string
    Name        string
    Passed      bool
    Message     string
    Details     string
}
```

**Output Formatter API**
```go
type OutputFormatter interface {
    Plain(data interface{}) string
    JSON(data interface{}) string
    YAML(data interface{}) string
}
```

## Security

**Authentication:**
- CLI does not require authentication (local-only tool)
- Git hooks run in user's git context

**Authorization:**
- File system permissions control access
- Respect user's git repository permissions

**Data Protection:**
- No sensitive data stored in default configuration
- Optional encryption for embeddings can be added
- Secrets not logged or displayed

**Input Validation:**
- All file paths validated before access
- Gate check inputs sanitized
- Configuration values validated on load

## Decision Rationale

### Why Go?
- **Performance**: Native compilation, fast execution for CLI tools
- **Cross-platform**: Single binary for Windows, macOS, Linux
- **Tooling**: Excellent build tooling (go build, make)
- **Dependencies**: Rich standard library reduces external dependencies
- **Community**: Strong Go CLI ecosystem (Cobra, Viper)

### Why YAML for Configuration?
- **Human-readable**: Easy to edit manually when needed
- **Hierarchical**: Supports nested configurations naturally
- **Standard**: Widely used in DevOps and tooling
- **Validation**: Can be validated against schemas

### Why SQLite for Embeddings?
- **Local-first**: No external service dependencies
- **Fast**: Efficient for vector similarity searches
- **Portable**: Single file database, easy to backup
- **Mature**: Well-tested, stable, supported everywhere

### Why Pre-commit Hooks over CI?
- **Immediate Feedback**: Developers see gate failures immediately
- **No Setup**: Works in any git repository
- **Flexible**: Can be bypassed when needed (--no-verify)
- **Complementary**: Works alongside CI for multiple layers of checks

### Why Gate-Based Validation?
- **Incremental Quality**: Enforce standards at each development phase
- **Clear Milestones**: Define what "done" means for each phase
- **Adaptable**: Gates can be customized per project
- **Traceable**: Each gate check is audited and tracked

## Open Questions

None - architecture is well-defined and implementation is in progress.

## Evolution Roadmap

**Phase 1** (Complete)
- Basic CLI structure
- SPEC file templates
- Gate 0-3 checks
- Git hooks integration

**Phase 2** (Planned)
- Enhanced output formatting
- Better error messages
- Configuration file support

**Phase 3** (Future)
- GUI interface option
- Team collaboration features
- Integration with popular CI/CD platforms
