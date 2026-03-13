# Vibe Integrity

[中文说明](./README.zh-CN.md)

Vibe Integrity is an **AI Project Memory & Safety System** designed specifically for AI-assisted development (vibe coding). It prevents false completion claims from AI coding assistants and provides structured project knowledge for rapid AI understanding.

## Overview

Vibe Integrity solves two critical problems in AI-assisted development:

1. **Completion Guard** - Detects when AI falsely claims work is complete (TODO/FIXME placeholders, empty functions, fake tests, etc.)
2. **Architecture Memory** - Provides structured project knowledge so AI can quickly understand project state without reading hundreds of files

Unlike traditional development methodologies, Vibe Integrity is **methodology agnostic** - it works with TDD, SDD, Agile, or pure vibe coding approaches.

## Core Concepts

### Two Pillars

#### Pillar 1: Completion Guard
Detection and validation to ensure AI actually completed the work.

| Skill | Purpose |
|-------|---------|
| `vibe-guard` | Detects TODO, empty functions, fake tests |
| `cascade-check` | Prevents cascading errors after fixes |
| `integration-check` | Validates component integration |

#### Pillar 2: Architecture Memory
Structured project knowledge base for AI quick understanding.

| File | Purpose |
|------|---------|
| `project.yaml` | Project metadata, tech stack |
| `dependency-graph.yaml` | Module dependencies |
| `module-map.yaml` | Directory structure |
| `risk-zones.yaml` | High-risk areas |
| `tech-records.yaml` | Technical decisions |
| `schema-evolution.yaml` | Data model changes |

## AI Quick Start

When AI starts work on this project, read in this order:

```
1. .vibe-integrity/project.yaml
   → Understand project status, tech stack

2. .vibe-integrity/risk-zones.yaml  
   → Know what areas are high-risk

3. .vibe-integrity/dependency-graph.yaml
   → Understand module relationships

4. .vibe-integrity/module-map.yaml
   → Find where files are located

5. .vibe-integrity/tech-records.yaml
   → Understand why system is designed this way
```

**Result**: AI understands project in ~15 seconds instead of 3 minutes.


## Base Workflow

Vibe Integrity provides a complete AI-assisted development workflow:

```
[User提出需求] → vibe-design (需求澄清/苏格拉底提问)
                                    ↓
                        [做出架构决策] → vibe-integrity-writer (自动更新 tech-records.yaml)
                                    ↓
                        [完成设计] → 生成 .vibe-integrity/ 更新
                                    ↓
                [实现过程] → vibe-integrity-debug (发现问题时进行根因分析)
                                    ↓
                        [发现新风险/决策] → vibe-integrity-writer (自动更新 risk-zones.yaml/tech-records.yaml)
                                    ↓
                        [实现完成] → vibe-guard (验证完整性)
                                    ↓
                        [验证通过] → 开发完成
```

### Workflow Stages:

| Stage | Primary Skill | Secondary Skills | Purpose |
|-------|---------------|------------------|---------|
| **Clarification** | `vibe-design` | `vibe-integrity-writer` | Understand requirements, capture decisions |
| **Implementation** | Developer/AI | `vibe-integrity-debug` | Build features, debug issues |
| **Verification** | `vibe-guard` | `validate-vibe-integrity.py` | Ensure completion and integrity |
| **Discovery** | `vibe-integrity-debug` | `vibe-integrity-writer` | Root cause analysis, update project memory |

## YAML File Responsibilities

Each YAML file in `.vibe-integrity/` has a specific responsibility and is maintained by different skills:

| File | Responsibility | Maintained By | Trigger Events |
|------|----------------|---------------|----------------|
| `project.yaml` | Project metadata, tech stack, status | `vibe-design` | Project creation, scope changes, tech stack updates |
| `tech-records.yaml` | Technical decisions and their rationale | `vibe-design` → `vibe-integrity-writer` | Architecture decisions, technology choices, implementation patterns |
| `dependency-graph.yaml` | Module/service dependencies and relationships | `vibe-design` → `vibe-integrity-writer` | Adding/removing modules, changing dependencies, service integration |
| `module-map.yaml` | Directory structure and file organization | `vibe-design` → `vibe-integrity-writer` | Reorganization, new directories, file relocation |
| `risk-zones.yaml` | Identified risks and high-risk areas | `vibe-integrity-debug` → `vibe-integrity-writer` | Bug discovery, security issues, performance bottlenecks, architectural flaws |
| `schema-evolution.yaml` | Data model changes and migrations | `vibe-design` → `vibe-integrity-writer` | Database schema changes, API model updates, data structure evolution |

### Automatic Updates by Workflow:

| Workflow Stage | Skills Involved | YAML Files Updated | Example |
|----------------|-----------------|-------------------|---------|
| **Requirement Clarification** | `vibe-design` → `vibe-integrity-writer` | `tech-records.yaml`, `project.yaml` | User confirms PostgreSQL over MongoDB → Record decision |
| **Architecture Discussion** | `vibe-design` → `vibe-integrity-writer` | `tech-records.yaml`, `dependency-graph.yaml` | Discussing module boundaries → Update dependency graph |
| **Risk Identification** | `vibe-integrity-debug` → `vibe-integrity-writer` | `risk-zones.yaml`, `tech-records.yaml` | Finding tight coupling → Record as risk area |
| **Schema Changes** | `vibe-design` → `vibe-integrity-writer` | `schema-evolution.yaml`, `tech-records.yaml` | Adding new table → Record schema evolution |
| **Debugging Insights** | `vibe-integrity-debug` → `vibe-integrity-writer` | `risk-zones.yaml`, `tech-records.yaml` | Discovering architectural issue → Record insight |

### Key Workflow Principles:

1. **AI Maintains Memory**: No manual YAML editing required - AI automatically updates project memory during normal workflow
2. **Decisions Recorded in Real-time**: Decisions are captured as they're made, not after the fact
3. **Safe Updates**: All YAML modifications go through `vibe-integrity-writer` with backup and validation
4. **Audit Trail**: Every change is tracked and reversible
5. **Cross-Referenced**: YAML files reference each other to maintain consistency

## Usage

### For AI: Before Making Changes

```bash
# 1. Check risk zone
cat .vibe-integrity/risk-zones.yaml

# 2. Check dependencies
cat .vibe-integrity/dependency-graph.yaml

# 3. Check schema
cat .vibe-integrity/schema-evolution.yaml
```

### For AI: After "Completing"

```bash
# Run vibe-guard
python skills/vibe-guard/validate-vibe-guard.py --check
```

### For Humans: After Significant Changes

```bash
# Update tech-records
python skills/vibe-integrity/validate-vibe-integrity.py  # First check integrity

# Add new decision to .vibe-integrity/tech-records.yaml
# Add new version to .vibe-integrity/schema-evolution.yaml  
# Reflect new module relationships in .vibe-integrity/dependency-graph.yaml
```

## Directory Structure


## Directory Structure with Skills

```
.vibe-integrity/
├── project.yaml              # Project metadata
├── dependency-graph.yaml     # Module dependencies
├── module-map.yaml          # Directory structure
├── risk-zones.yaml          # Risk areas
├── tech-records.yaml        # Technical decisions
└── schema-evolution.yaml   # Data model changes

skills/
├── vibe-guard/              # Completion detection (Pillar 1)
├── vibe-integrity/          # Validation framework
│   ├── SKILL.md
│   ├── validate-vibe-integrity.py
│   ├── validate-all.py
│   └── template/           # Schema templates
│       ├── project.schema.json
│       ├── dependency-graph.schema.json
│       ├── module-map.schema.json
│       ├── risk-zones.schema.json
│       ├── tech-records.schema.json
│       └── schema-evolution.schema.json
├── vibe-design/             # Requirement clarification (Pillar 2)
│   ├── SKILL.md
│   └── Uses vibe-integrity-writer for updates
├── vibe-integrity-debug/    # Systematic debugging
│   ├── SKILL.md
│   └── Uses vibe-integrity-writer for insights
└── vibe-integrity-writer/   # YAML file writer (Pillar 2)
    ├── SKILL.md
    └── Handles all .vibe-integrity/ updates
```

### Skill Responsibilities:

| Skill | Pillar | Purpose | Outputs |
|-------|--------|---------|----------|
| `vibe-guard` | Completion Guard | Detects false completion claims | validation reports |
| `vibe-design` | Architecture Memory | Clarifies requirements, captures decisions | updated .vibe-integrity/ files |
| `vibe-integrity-debug` | Architecture Memory | Root cause analysis, risk identification | updated .vibe-integrity/ files |
| `vibe-integrity-writer` | Architecture Memory | Safe YAML updates with validation | updated .vibe-integrity/ files |
| `vibe-integrity` | Both | Validates structure and integrity | validation reports |

```
.vibe-integrity/
├── project.yaml              # Project metadata
├── dependency-graph.yaml     # Module dependencies
├── module-map.yaml          # Directory structure
├── risk-zones.yaml          # Risk areas
├── tech-records.yaml        # Technical decisions
└── schema-evolution.yaml   # Data model changes

skills/
├── vibe-guard/             # Completion detection
└── vibe-integrity/         # This skill
    ├── SKILL.md
    ├── validate-vibe-integrity.py
    ├── validate-all.py
    └── template/           # Schema templates
        ├── project.schema.json
        ├── dependency-graph.schema.json
        ├── module-map.schema.json
        ├── risk-zones.schema.json
        ├── tech-records.schema.json
        └── schema-evolution.schema.json
```

## Validation

Run validation to ensure integrity:

```bash
python skills/vibe-integrity/validate-vibe-integrity.py  # checks .vibe-integrity/ files
python skills/vibe-integrity/validate-all.py             # runs both vibe-guard and vibe-integrity validations
python skills/vibe-guard/validate-vibe-guard.py --check  # AI completion check
```

## Related Skills

- `vibe-guard` - Completion detection
- `vibe-design` - Requirement clarification and design helper that uses vibe-integrity-writer to update project memory
- `vibe-integrity-debug` - Systematic debugging helper ensuring root cause analysis before fixes
- `vibe-integrity-writer` - Specialized skill for safely updating .vibe-integrity/ YAML files (called by other skills)
- `superpowers/test-driven-development` - TDD workflow (optional)
- `sdd-orchestrator` - SDD workflow (optional)

**Note**: Vibe Integrity works with ANY development approach. You can use Vibe Integrity alone, or combine it with SDD, TDD, Agile, or any other methodology. The SDD and TDD skills listed above are optional add-ons for teams that wish to follow those specific methodologies while still benefiting from Vibe Integrity's completion guards and project memory.

**Workflow Summary**:
- **vibe-design**: Clarifies requirements, makes decisions, calls writer to update YAML
- **vibe-integrity-debug**: Performs root cause analysis, identifies risks, calls writer to update YAML  
- **vibe-integrity-writer**: Safely updates .vibe-integrity/ YAML files with backup and validation
- **vibe-guard**: Verifies AI completion after implementation
- **vibe-integrity**: Validates .vibe-integrity/ directory structure

## Multi-Agent Collaboration

Vibe Integrity is designed for both single-agent and multi-agent use, with specific mechanisms to support collaboration:

### Current Collaboration Capabilities

| Mechanism | Description | Use Case |
|-----------|-------------|----------|
| **Git Branching** | Separate branches for different agents/sessions | Multi-agent parallel development |
| **Union Merge** | `.vibe-integrity/*.yaml merge=union` in `.gitattributes` | Prevents data loss on conflicts |
| **Validation Script** | Detects duplicate IDs and inconsistencies | Early conflict detection |
| **Grace Periods** | 10-minute deduplication window | Prevents redundant writes |
| **Session Isolation** | Different git worktrees for different workflows | Avoiding workflow confusion |

### File Locking & Concurrency

**Current Implementation:**
- ✅ Backup creation before modifications
- ✅ Atomic batch operations
- ✅ Post-update validation
- ✅ **File-level locking** (implemented in vibe-integrity-writer)
- ✅ **Agent identity tracking** (all updates include agent metadata)
- ✅ **Conflict detection** (detects duplicate IDs and similar decisions)
- ❌ No distributed locking (for multi-instance deployments)

**Recommendation:** For multi-agent scenarios, use the new file locking and conflict detection features. Each agent should still use separate branches when possible, but the system now supports concurrent writes with proper locking.

**Current Implementation:**
- ✅ Backup creation before modifications
- ✅ Atomic batch operations
- ✅ Post-update validation
- ❌ No file-level locking (recommended for sequential use)
- ❌ No distributed locking (for multi-instance deployments)

**Recommendation:** For multi-agent scenarios, use separate branches and coordinate via Git workflow rather than simultaneous writes to the same files.

### Conflict Resolution

When multiple agents modify the same .vibe-integrity/ files:

1. **Git will detect conflicts** during merge/pull request
2. **Union merge** preserves both versions (may create duplicates)
3. **Validation script** identifies duplicate IDs
4. **Manual resolution** required to merge similar decisions

### Best Practices for Multi-Agent Use

1. **Use Separate Branches**: Each agent gets its own feature branch
2. **Coordinate Architecture Decisions**: Document major decisions in shared records
3. **Regular Validation**: Run validation script frequently
4. **Pull Request Reviews**: Review .vibe-integrity/ changes before merging
5. **Conflict Resolution**: Use validation script to detect and fix duplicates

### Multi-Instance Considerations

For deployments with multiple agent instances:

- **File Locking**: Implement file-level locks in vibe-integrity-writer
- **Distributed Locking**: Use Redis or similar for multi-server deployments
- **Lease Mechanism**: Implement TTL for locks to prevent deadlocks
- **Fallback Strategy**: Graceful degradation when locks unavailable

See [MULTI_AGENT_COLLABORATION.md](skills-base/vibe-integrity/MULTI_AGENT_COLLABORATION.md) for detailed analysis.

**Workflow Separation Guidance**: To avoid confusion when using multiple methodologies:
- Use `vibe-design` during the clarification phase to understand requirements and record decisions
- Use `vibe-guard` after implementation claims to verify completion
- Use `vibe-integrity-debug` for investigating issues with root cause analysis
- When following TDD or SDD, consider using separate work sessions or explicitly declaring the workflow at the start of each session
- Vibe Design automatically updates `.vibe-integrity/` files as decisions are made, reducing manual documentation burden

**Important**: Do not mix `vibe-design` with active TDD/SDD implementation sessions without clear context switching to prevent AI confusion about which workflow is being followed.

## Quick Start

1) Run default validation (scans `<root>/skills`):

```bash
python skills/vibe-integrity/validate-all.py
```

2) Initialize Vibe Integrity in your project:

```bash
# Create .vibe-integrity directory with template files
python skills/vibe-integrity/validate-vibe-integrity.py --init

# Or manually copy template files:
cp -r skills/vibe-integrity/template/* .vibe-integrity/
```

3) Customize the files for your project:
   - Edit `.vibe-integrity/project.yaml` with your project details
   - Update `.vibe-integrity/tech-records.yaml` with your technical decisions
   - Customize `.vibe-integrity/risk-zones.yaml` for your project's risk areas

## Example Output

A successful validation run looks like this:

```text
Vibe Integrity validation passed
Root: D:\Code\aaa
Files checked:
- .vibe-integrity/project.yaml ✓
- .vibe-integrity/dependency-graph.yaml ✓
- .vibe-integrity/module-map.yaml ✓
- .vibe-integrity/risk-zones.yaml ✓
- .vibe-integrity/tech-records.yaml ✓
- .vibe-integrity/schema-evolution.yaml ✓

Vibe Guard validation:
- TODO/FIXME check: PASSED
- Empty functions check: PASSED
- Fake tests check: PASSED
- Build success: PASSED
- Type check: PASSED
- Lint check: PASSED
- Security check: PASSED
- Test authenticity: PASSED

All validations PASSED
```

If `Vibe Integrity validation passed` is shown, all files are present and structurally valid.

## Configuration

Vibe Integrity uses YAML files in the `.vibe-integrity/` directory for configuration.

### project.yaml
```yaml
name: my-project
version: 0.1.0
status: mvp
description: "My amazing project"
created_at: 2026-01-15
last_updated: 2026-03-12
tech_stack:
  frontend: [Vue, Vite]
  backend: [Express, Node]
  database: [SQLite]
```

### tech-records.yaml
```yaml
records:
  - id: DB-001
    date: "2026-01-15"
    category: database
    title: "Choose SQLite for MVP"
    decision: "Use SQLite for fast iteration"
    reason: "MVP phase prioritizes speed over scalability"
    impact: low
    status: completed
```

## Common Operations

### Initialize new project structure
```bash
python skills/vibe-integrity/validate-vibe-integrity.py --init
```

### Validate integrity
```bash
python skills/vibe-integrity/validate-all.py
```

### AI completion check
```bash
python skills/vibe-guard/validate-vibe-guard.py --check
```

## License

This project is licensed under MIT. See [LICENSE](./LICENSE).