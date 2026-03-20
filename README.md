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
```

## Commands

| Command | Description |
|---------|-------------|
| `vic init` | Initialize .vic-sdd/ |
| `vic spec init` | Initialize SPEC documents |
| `vic spec status` | Show SPEC status |
| `vic spec gate [0-3\|1.5]` | Run Gate checks (validation) |
| `vic spec watch` | Monitor SPEC changes and auto-run drift detection |
| `vic spec diff` | Detect SPEC changes since last check |
| `vic spec changes` | Show SPEC change history |
| `vic phase advance` | Advance phase (auto-validates gates) |
| `vic gate check --blocking` | Pre-commit hook check |
| `vic rt` | Record technical decision |
| `vic rr` | Record risk |
| `vic check` | Check code alignment |
| `vic validate` | Full validation |

See [cmd/vic-go/README.md](./cmd/vic-go/README.md) for full documentation.

## Development Workflow

```
定图纸 (Requirements)     打地基 (Architecture)    立规矩 (Implementation)
        │                          │                         │
   requirements             architecture             sdd-orchestrator
        │                          │                         │
        ▼                          ▼                         ▼
SPEC-REQUIREMENTS.md  ──▶  SPEC-ARCHITECTURE.md  ──▶  Implementation
        │                          │                         │
        ▼                          ▼                         ▼
   Gate 0                    Gate 1                  Gate 2 + 3
(Requirements)          (Architecture)           (Code + Tests)
```

## Directory Structure

```
project/
├── cmd/vic-go/                 # Go CLI (compiled, fast)
│   ├── internal/
│   │   └── commands/          # Gate implementations
│   │       ├── gate0.go       # Requirements validation
│   │       ├── gate1.go       # Architecture validation
│   │       ├── gate2.go       # Code alignment check
│   │       ├── gate3.go       # Test coverage check
│   │       └── ...
│   └── README.md
│
├── skills/                     # 10 core skills (simplified from 19)
│   ├── context-tracker/      # Self-awareness (4→1)
│   ├── requirements/          # Requirements (2→1)
│   ├── architecture/          # Tech architecture
│   ├── design-review/         # Design system
│   ├── debugging/            # Debug (2→1)
│   ├── qa/                    # Testing (3→1)
│   ├── sdd-orchestrator/      # SDD pipeline
│   ├── spec-architect/        # Spec contracts
│   ├── spec-contract-diff/     # Drift detection
│   └── spec-traceability/     # Traceability
│
├── docs/                      # Documentation
├── .vic-sdd/                  # Project memory
│   ├── SPEC-REQUIREMENTS.md    # Requirements
│   ├── SPEC-ARCHITECTURE.md    # Architecture
│   ├── PROJECT.md              # Status
│   ├── agent-prompt.md        # AI workflow prompt
│   ├── context.yaml            # Unified context
│   └── status/
│       ├── gate-status.yaml    # Gate check status
│       ├── spec-hash.json      # SPEC file hashes for change detection
│       └── change-log.yaml     # SPEC change history
└── .pre-commit-config.yaml    # Gate enforcement
```

## Core Concepts

### 定图纸 (Requirements)
- Define user stories and acceptance criteria
- Plan development phases
- Create SPEC-REQUIREMENTS.md

### 打地基 (Architecture)
- Evaluate technology options
- Design system architecture
- Create SPEC-ARCHITECTURE.md

### 立规矩 (Implementation)
- Small iteration cycles
- Gate checks before progression
- Merge to PRD/ARCH/PROJECT

### 自我认知 (Self-Awareness)
VIBE-SDD gives AI "self-awareness" through unified context tracking:
- **Context Tracker** — knows/infers/assumes/unknown + signals + confidence

## AI Quick Start

When AI starts on this project, read in order:

```
1. .vic-sdd/agent-prompt.md    → Workflow overview (displayed at session start)
2. .vic-sdd/PROJECT.md         → Project status, milestones
3. .vic-sdd/SPEC-REQUIREMENTS.md → Requirements, acceptance criteria
4. .vic-sdd/SPEC-ARCHITECTURE.md → Architecture, tech stack
```

**Result**: AI understands project context in ~15 seconds.

## Skills Reference (10 Core Skills)

| Category | Skill | Purpose |
|----------|-------|---------|
| Self-Awareness | `context-tracker` | Unified: known/inferred/assumed/unknown + signals |
| Vibe | `requirements` | User stories, acceptance criteria |
| Vibe | `architecture` | Tech selection, system design |
| Vibe | `design-review` | Design system, AI slop detection |
| Vibe | `debugging` | Root cause analysis (SURVEY→PATTERN→HYPOTHESIS→IMPLEMENT) |
| QA | `qa` | TDD, test coverage, E2E |
| SDD | `sdd-orchestrator` | State machine, gate enforcement |
| SDD | `spec-architect` | Freeze requirements into contracts |
| SDD | `spec-contract-diff` | Detect spec drift |
| SDD | `spec-traceability` | Story→contract→code→test mapping |

## Gate Enforcement

### Automatic Gate Checks

```bash
# Run before claiming "done"
vic spec gate 0   # Validates SPEC-REQUIREMENTS.md structure
vic spec gate 1   # Validates SPEC-ARCHITECTURE.md structure
vic spec gate 1.5 # Validates DESIGN.md completeness (optional)
vic spec gate 2   # Checks code vs SPEC alignment
vic spec gate 3   # Validates test coverage
```

### SPEC Change Detection

```bash
# Start monitoring SPEC files for changes
vic spec watch     # Detects changes + auto-runs Gate 2 drift check

# Compare current SPEC with last known state
vic spec diff      # Show what changed

# View change history
vic spec changes   # Show all recorded changes
```

Changes are logged to `.vic-sdd/status/change-log.yaml` with:
- Change type (tech_stack, api, module, security)
- Impact level (high, medium, low)
- Review status tracking

### Pre-commit Hook

`.pre-commit-config.yaml` includes `vic gate check --blocking` to prevent commits until gates pass.

### Phase Advance

```bash
vic phase advance --to 1  # Auto-runs all required gates first
```

## Workflow

| Scenario | Command |
|----------|---------|
| Start new project | `vic init` |
| Check requirements | `vic spec gate 0` |
| Check architecture | `vic spec gate 1` |
| Monitor SPEC changes | `vic spec watch` |
| Check code alignment | `vic spec gate 2` |
| Check test coverage | `vic spec gate 3` |
| Advance phase | `vic phase advance --to N` |
| Pre-commit check | `vic gate check --blocking` |

## Installation

```bash
# Build from source
cd cmd/vic-go
make build

# Install to PATH
sudo ln -s $(pwd)/vic /usr/local/bin/vic

# Or use Go
go install github.com/vic-sdd/vic@latest
```

## License

MIT License. See [LICENSE](./LICENSE).
