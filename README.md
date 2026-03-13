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
- `superpowers/test-driven-development` - TDD workflow (optional)
- `sdd-orchestrator` - SDD workflow (optional)

**Note**: Vibe Integrity works with ANY development approach. You can use Vibe Integrity alone, or combine it with SDD, TDD, Agile, or any other methodology. The SDD and TDD skills listed above are optional add-ons for teams that wish to follow those specific methodologies while still benefiting from Vibe Integrity's completion guards and project memory.

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