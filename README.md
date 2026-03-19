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

# Initialize SPEC documents
vic spec init --name "My Project"

# Record a technical decision
vic rt --id DB-001 --title "Use PostgreSQL" --decision "Primary database" --reason "Need ACID"

# Check SPEC status
vic spec status

# Run Gate checks
vic spec gate 0  # Requirements
vic spec gate 1  # Architecture

# Validate
vic validate
```

## Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `vic init` | - | Initialize .vic-sdd/ |
| `vic spec init` | - | Initialize SPEC documents |
| `vic spec status` | - | Show SPEC status |
| `vic spec gate [0-3]` | - | Run Gate checks |
| `vic rt` | `record-tech` | Record technical decision |
| `vic rr` | `record-risk` | Record risk |
| `vic rd` | `record-dep` | Record dependency |
| `vic check` | - | Check code alignment |
| `vic validate` | - | Full validation |
| `vic status` | - | Show project status |
| `vic search` | - | Search records |
| `vic history` | - | Show event history |
| `vic export` | - | Export data |
| `vic import` | - | Import data |

See [cmd/vic/README.md](./cmd/vic/README.md) for full documentation.

## Development Workflow

```
定图纸 (Requirements)     打地基 (Architecture)    立规矩 (Implementation)
        │                          │                         │
   vibe-think              vibe-architect            sdd-orchestrator
        │                          │                         │
        ▼                          ▼                         ▼
SPEC-REQUIREMENTS.md  ──▶  SPEC-ARCHITECTURE.md  ──▶  Implementation
        │                          │                         │
        ▼                          ▼                         ▼
   Gate 0                    Gate 1                  Gate 2 + 3
(Requirements)          (Architecture)           (Code + Tests)
                                                        │
                                                        ▼
                                              Merge to PRD/ARCH/PROJECT
```

## Directory Structure

```
project/
├── cmd/
│   └── vic/                    # CLI tool
│       ├── vic                  # Main CLI
│       ├── README.md            # English docs
│       └── *.py                 # Scripts
│
├── skills/                         # 18 skills total
│   │
│   ├── Self-Awareness (4):        # AI self-awareness mechanisms
│   │   ├── knowledge-boundary/ # Knows/infers/assumes/unknown
│   │   ├── pre-decision-check/ # Gate check before decisions
│   │   ├── signal-register/    # Evidence-based progress
│   │   └── exploration-journal/ # Exploration memory
│   │
│   ├── Vibe Exploration (7):     # Flexible exploration
│   │   ├── vibe-think/         # Requirements clarification
│   │   ├── vibe-architect/     # Tech selection + architecture
│   │   ├── vibe-redesign/      # Product redesign
│   │   ├── vibe-design/        # Design system
│   │   ├── vibe-debug/         # Systematic debugging
│   │   ├── vibe-qa/           # Quality assurance
│   │   └── adaptive-planning/  # Adaptive replanning
│   │
│   └── SDD Core (7):              # Strict spec-driven delivery
│       ├── sdd-orchestrator/    # State machine + gate enforcement
│       ├── spec-architect/      # Freeze requirements into contracts
│       ├── spec-to-codebase/    # Generate implementation
│       ├── spec-contract-diff/  # Detect spec drift
│       ├── spec-driven-test/    # Contract + TDD tests
│       ├── spec-traceability/   # Story-to-code traceability
│       └── sdd-release-guard/  # Final release gates
│
├── docs/                      # Design docs
│   ├── VIC-CLI-GUIDE.md      # CLI操作指南
│   └── *.md
│
└── .vic-sdd/                  # Project memory & specs
    ├── SPEC-REQUIREMENTS.md    # Requirements spec
    ├── SPEC-ARCHITECTURE.md    # Architecture spec
    ├── PROJECT.md              # Project status
    ├── knowledge-boundary.yaml # AI 认知地图
    ├── decision-guardrails.yaml # Decision constraints
    ├── signal-register.yaml    # Evidence-based progress
    ├── exploration-journal.yaml # Exploration memory
    ├── status/
    │   ├── events.yaml         # Event history
    │   └── state.yaml         # Current state
    ├── tech/
    │   └── tech-records.yaml  # Technical decisions
    ├── risk-zones.yaml        # Risk records
    ├── project.yaml           # AI quick reference
    └── dependency-graph.yaml  # Module dependencies
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
VIBE-SDD gives AI "self-awareness" through 4 mechanisms:
- **Knowledge Boundary** — Knows what it knows, infers, assumes, or doesn't know
- **Pre-Decision Check** — Gates before major decisions (scope, quality, signals)
- **Signal Register** — Evidence-based progress instead of "60% done"
- **Exploration Journal** — Remembers exploration process to avoid repeating failures

## AI Quick Start

When AI starts on this project, read in order:

```
1. .vic-sdd/PROJECT.md          → Project status, milestones
2. .vic-sdd/SPEC-REQUIREMENTS.md → Requirements, acceptance criteria
3. .vic-sdd/SPEC-ARCHITECTURE.md → Architecture, tech stack
4. .vic-sdd/risk-zones.yaml    → High-risk areas
```

**Result**: AI understands project context in ~15 seconds.

## Workflow

| Scenario | Command |
|----------|---------|
| Start new project | `vic init` |
| Initialize SPEC | `vic spec init` |
| Made a decision | `vic rt` |
| Found a risk | `vic rr` |
| Before progression | `vic phase advance` |
| Check phase | `vic phase status` |
| Pass gate | `vic gate pass --gate N` |
| AI claims "done" | `vic check` |
| Before commit | `vic validate` |
| Backup memory | `vic export` |

## Related Skills

| Skill | Purpose |
|-------|---------|
| `vibe-think` | Requirements clarification, user story discovery |
| `vibe-architect` | Tech selection, architecture design |
| `vibe-redesign` | Product redesign, scope re-evaluation |
| `vibe-design` | Design system, UI/UX specifications |
| `vibe-debug` | Systematic debugging with root cause analysis |
| `vibe-qa` | Quality assurance, verification against specs |
| `adaptive-planning` | Adaptive replanning when scope changes |
| `sdd-orchestrator` | Entry point for SDD phase, state machine |
| `spec-architect` | Freeze requirements into contracts |
| `spec-to-codebase` | Generate implementation from contracts |
| `spec-contract-diff` | Detect drift between spec and code |
| `spec-driven-test` | Build and enforce test gates |
| `spec-traceability` | Story-to-contract-to-code-to-test mapping |
| `sdd-release-guard` | Final SDD release gates |
| `knowledge-boundary` | Self-awareness: knows/infers/assumes/unknown |
| `pre-decision-check` | Self-awareness: gate check before decisions |
| `signal-register` | Self-awareness: evidence-based progress tracking |
| `exploration-journal` | Self-awareness: exploration memory |

> 注：CLI命令详细用法见 [VIC-CLI-GUIDE.md](./docs/VIC-CLI-GUIDE.md)

## Installation

```bash
# Dependencies
pip install pyyaml

# Linux/macOS
chmod +x cmd/vic/vic
sudo ln -s $(pwd)/cmd/vic/vic /usr/local/bin/vic

# Windows PowerShell
Set-Alias vic "python D:\path\to\cmd\vic\vic"
```

## License

MIT License. See [LICENSE](./LICENSE).
