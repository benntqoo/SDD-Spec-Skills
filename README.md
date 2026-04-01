# VIBE-SDD

[中文说明](./README.zh-CN.md)

VIBE-SDD is a **Vibe-Driven Software Development System** combining structured SDD (Spec-Driven Development) with flexible Vibe Coding. It provides a complete workflow for AI-assisted development with proper gates and documentation.

## Overview

VIBE-SDD solves three critical problems in AI-assisted development:

1. **Specification** - Structured requirements and architecture documentation
2. **Gates** - Quality checkpoints before progression
3. **Memory** - Project knowledge for AI quick understanding

## Quick Start

```bash
# Initialize project
vic init --name "My Project" --tech "React,Node,PostgreSQL"

# Run Gate checks (blocks commit until passed)
vic spec gate 0  # Requirements completeness
vic spec gate 1  # Architecture completeness
vic spec gate 2  # Code alignment
vic spec gate 3  # Test coverage

# Advance phase (auto-runs gate checks)
vic phase advance --to 1

# Record decisions
vic rt --id DB-001 --title "Use PostgreSQL" --decision "Primary database"
vic rr --id RISK-001 --area auth --desc "JWT not validated"

# Semantic search (requires embedding)
vic ask "how does authentication work?"
```

## Commands

| Command | Description |
|---------|-------------|
| `vic init` | Initialize .vic-sdd/ |
| `vic status` | Show project status |
| `vic spec init` | Initialize SPEC documents |
| `vic spec status` | Show SPEC status |
| `vic spec gate [0-3]` | Run Gate checks (validation) |
| `vic spec hash` | Check SPEC hashes and detect changes |
| `vic spec diff` | Detect SPEC changes since last check |
| `vic phase advance` | Advance phase (auto-validates gates) |
| `vic gate check --blocking` | Pre-commit hook check |
| `vic rt` / `vic record tech` | Record technical decision |
| `vic rr` / `vic record risk` | Record risk |
| `vic check` | Check code alignment |
| `vic validate` | Full validation |
| `vic ask <query>` | Semantic search about codebase |

See [docs/VIC-CLI-GUIDE.md](./docs/VIC-CLI-GUIDE.md) for full documentation.

## Quality Gates

VIBE-SDD includes a comprehensive gate system to ensure quality at each phase:

- **[Gate 0: Requirements](./docs/GATE-GUIDE.md#gate-0-需求完整性)** - Ensure requirements are complete and testable
- **[Gate 1: Architecture](./docs/GATE-GUIDE.md#gate-1-架构完整性)** - Validate architecture completeness
- **[Gate 2: Code Alignment](./docs/GATE-GUIDE.md#gate-2-代码对齐)** - Check code-SPEC consistency and TODOs
- **[Gate 3: Test Coverage](./docs/GATE-GUIDE.md#gate-3-测试覆盖率)** - Verify adequate test coverage

### Quick Reference
- **[Gate Cheat Sheet](./docs/GATE-CHEAT-SHEET.md)** - Quick reference for gate checks
- **[Git Hooks Guide](./docs/FIX-LOG.md#2026-04-01---phase-5-git-hooks-集成完成)** - Automated quality checks

## Development Workflow

```
Ideation → Explore → SpecCheckpoint → Build → Verify → ReleaseReady → Released
    │         │            │             │        │          │
    └─────────┴────────────┘             └────────┴──────────┘
         spec-workflow                        implementation
                                              unified-workflow
```

### 5 Core Skills

| Skill | Auto-Activate | Responsibility |
|-------|---------------|----------------|
| `context-tracker` | ✅ Yes | AI self-awareness, confidence tracking |
| `spec-workflow` | No | Requirements → Architecture → SPEC freezing |
| `implementation` | No | Code/debugging/testing/SPEC alignment |
| `unified-workflow` | No | SDD orchestration/Constitution/traceability |
| `quick` | No | Simple single-file changes |

## Directory Structure

```
project/
├── cmd/vic-go/                 # Go CLI (compiled, fast)
│   ├── main.go
│   └── internal/
│       ├── commands/           # CLI command implementations
│       │   ├── root.go
│       │   ├── spec.go
│       │   ├── gate.go
│       │   ├── gate0-3.go      # Gate implementations
│       │   └── ...
│       ├── config/             # Configuration (Viper)
│       ├── checker/            # Code alignment checking
│       ├── types/              # Type definitions
│       └── embedding/          # Semantic search with embeddings
│           ├── store.go        # SQLite vector store
│           ├── embedder.go     # Embedding generation
│           └── chunker/        # Code chunking by language
│
├── skills/                     # 5 core skills
│   ├── context-tracker/        # Self-awareness (auto-activate)
│   ├── spec-workflow/          # Requirements/Architecture/SPEC
│   ├── implementation/         # Code/Debug/Test
│   ├── unified-workflow/       # SDD orchestration
│   └── quick/                  # Simple changes
│
├── docs/                       # Documentation
│   ├── VIC-CLI-GUIDE.md        # CLI reference
│   ├── SDD-PROCESS-CN.md       # SDD process (Chinese)
│   ├── INTELLIGENT_JUDGMENT_DESIGN.md
│   └── ...
│
├── .vic-sdd/                   # Project memory
│   ├── SPEC-REQUIREMENTS.md    # Requirements
│   ├── SPEC-ARCHITECTURE.md    # Architecture
│   ├── PROJECT.md              # Status
│   ├── constitution.yaml       # Unbreakable rules
│   ├── context.yaml            # AI self-awareness state
│   ├── agent-prompt.md         # AI workflow prompt
│   └── status/
│       ├── spec-hash.json      # SPEC file hashes
│       ├── gate-status.yaml    # Gate check status
│       └── state.yaml          # System state
│
└── .pre-commit-config.yaml     # Gate enforcement
```

## Core Concepts

### SDD State Machine

```
Ideation → Explore → SpecCheckpoint → Build → Verify → ReleaseReady → Released
              │              │              │        │
        spec-workflow     implementation    unified-workflow
```

### 4 Gates

| Gate | Name | Checks |
|------|------|--------|
| Gate 0 | Requirements | SPEC-REQUIREMENTS.md completeness |
| Gate 1 | Architecture | SPEC-ARCHITECTURE.md completeness |
| Gate 2 | Code | Code alignment with SPEC |
| Gate 3 | Tests | Test coverage verification |

### Constitution Rules

Defined in `.vic-sdd/constitution.yaml`:

| Rule | Description |
|------|-------------|
| `SPEC-FIRST` | Update SPEC before changing code |
| `SPEC-ALIGNED` | Code must match SPEC |
| `GATE-BEFORE-COMMIT` | All gates must pass before commit |
| `NO-TODO-IN-CODE` | No TODO/FIXME comments |
| `NO-CONSOLE-IN-PROD` | No console.log in production |
| `TESTS-REQUIRED` | New features must have tests |

## AI Quick Start

When AI starts on this project, read in order:

```
1. AGENTS.md                  → Entry point, skills overview
2. .vic-sdd/PROJECT.md        → Project status, milestones
3. .vic-sdd/SPEC-REQUIREMENTS.md → Requirements, acceptance criteria
4. .vic-sdd/SPEC-ARCHITECTURE.md → Architecture, tech stack
```

**Result**: AI understands project context in ~15 seconds.

## Pre-commit Hooks

Pre-commit is configured in `.pre-commit-config.yaml`:

```bash
pre-commit install
pre-commit run --all-files
```

Hooks include:
- `vic-gate-check`: Blocks commits until gates pass
- `vic-spec-drift`: Detects code drift from SPEC

## Installation

```bash
# Build from source
cd cmd/vic-go
make build

# Install to PATH
make install

# Or run directly
make run ARGS="--help"
```

## Build Commands

```bash
cd cmd/vic-go

# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Install to PATH
make install

# Run locally with arguments
make run ARGS="--help"
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `VIC_DIR` | `.vic-sdd` | VIC directory name |
| `VIC_PROJECT_DIR` | current dir | Project directory |
| `VIC_OUTPUT` | `plain` | Output format (json/yaml/plain) |
| `VIC_VERBOSE` | `false` | Verbose output |

## License

MIT License. See [LICENSE](./LICENSE).
