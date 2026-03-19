# vic-go

VIBE-SDD CLI written in Go.

## Features

- Single binary, no dependencies required
- Fast startup time
- Cross-platform (Linux, macOS, Windows)
- Full support for all vic commands

## Installation

### From Source

```bash
# Clone and build
cd cmd/vic-go
make build

# Install to PATH
sudo ln -s $(pwd)/vic /usr/local/bin/vic

# Or use make install
make install
```

### Pre-built Binaries

Download from [Releases](https://github.com/vic-sdd/vic/releases)

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `VIC_DIR` | `.vic-sdd` | Override VIC directory name |
| `VIC_PROJECT_DIR` | (current dir) | Override project directory |
| `VIC_OUTPUT` | `plain` | Output format (json/yaml/plain) |
| `VIC_VERBOSE` | `false` | Verbose output |

### Examples

```bash
# Use custom VIC directory
VIC_DIR=.my-vic vic init

# Use custom project directory
VIC_PROJECT_DIR=/path/to/project vic status

# JSON output
VIC_OUTPUT=json vic status
```

## Usage

```bash
# Initialize project
vic init --name "My Project" --tech "Go,PostgreSQL"

# Record technical decision
vic record tech --id DB-001 --title "Use PostgreSQL" --decision "Primary DB"

# Record risk
vic record risk --id RISK-001 --area auth --desc "JWT not validated"

# Check code alignment
vic check

# Full validation
vic validate

# Show status
vic status

# Search records
vic search postgres

# SPEC management
vic spec init
vic spec gate 0
```

## Development

```bash
# Build
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run locally
make run ARGS="--help"
```

## Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `init` | - | Initialize .vic-sdd/ |
| `record tech` | `rt` | Record technical decision |
| `record risk` | `rr` | Record risk |
| `record dep` | `rd` | Record dependency |
| `check` | - | Check code alignment |
| `validate` | - | Full validation |
| `fold` | - | Fold events to state |
| `status` | - | Show project status |
| `search` | - | Search records |
| `history` | - | Show event history |
| `export` | - | Export data |
| `import` | - | Import data |
| `spec init` | - | Initialize SPEC |
| `spec status` | - | Show SPEC status |
| `spec gate` | - | Run SPEC gate |
| `auto start` | - | Start autonomous mode |
| `auto status` | - | Show auto mode status |
| `auto pause` | - | Pause autonomous mode |
| `auto resume` | - | Resume autonomous mode |
| `auto stop` | - | Stop autonomous mode |
| `cost init` | - | Initialize cost tracking |
| `cost status` | - | Show cost tracking status |
| `cost set-budget` | - | Set budget ceiling |
| `cost add` | - | Add cost record |
| `product record` | - | Record product redesign |
| `product list` | - | List product decisions |
| `product modes` | - | Show four modes |
| `replan trigger` | - | Trigger adaptive replan |
| `replan list` | - | List replan history |
| `replan show` | - | Show replan details |
| `slop scan` | - | Scan for AI slop patterns |
| `slop report` | - | Show last scan report |
| `slop list` | - | List configured patterns |
| `slop fix` | - | Auto-fix AI slop patterns |
| `tdd start` | - | Start TDD session |
| `tdd red` | - | RED phase - write failing test |
| `tdd green` | - | GREEN phase - make it pass |
| `tdd refactor` | - | REFACTOR phase |
| `tdd status` | - | Show TDD status |
| `tdd checkpoint` | - | Save TDD checkpoint |
| `tdd history` | - | Show TDD history |
| `debug start` | - | Start debug session |
| `debug survey` | - | Gather evidence |
| `debug pattern` | - | Find similar issues |
| `debug hypothesis` | - | Form and test hypothesis |
| `debug implement` | - | Implement fix |
| `debug status` | - | Show debug status |
| `debug report` | - | Generate debug report |
| `qa init` | - | Initialize QA setup |
| `qa quick` | - | Quick smoke test |
| `qa full` | - | Full application test |
| `qa screenshot` | - | Capture screenshot |
| `qa report` | - | Show QA report |

## Auto Mode (Phase 1 Enhancement)

VIC-SDD Phase 1 introduces autonomous execution mode with state persistence and crash recovery:

```bash
# Start autonomous mode
vic auto start

# Monitor progress
vic auto status

# Pause to take a break
vic auto pause

# Resume later
vic auto resume

# Stop when done
vic auto stop
```

Auto mode persists state to `.vic-sdd/status/auto.yaml` for crash recovery.

## Cost Tracking (Phase 1 Enhancement)

Track token usage and costs with budget management:

```bash
# Initialize cost tracking
vic cost init

# Show current costs
vic cost status

# Set budget ceiling ($50 default)
vic cost set-budget 100

# Record token usage
vic cost add --input 1000 --output 500 --cost 0.50

# Budget alerts and auto-pause
# When 80% of budget is reached, a warning is issued
# When budget is exceeded, auto mode will pause
```

## Product Redesign (Phase 2 Enhancement)

Inspired by gstack's /plan-ceo-review, implement the four modes of product discovery:

```bash
# Show the four modes
vic product modes

# Record a product redesign decision
vic product record \
  --original "Photo upload" \
  --real "Help sellers create sellable listings" \
  --mode expansion

# List all product decisions
vic product list
```

**The Four Modes:**

| Mode | Icon | When to Use |
|------|------|-------------|
| EXPANSION | 🚀 | Explore ambitious possibilities |
| SELECTIVE | ⚖️ | Neutral presentation, let user choose |
| HOLD | 🔒 | Stay focused, no expansions |
| REDUCTION | ✂️ | Find minimum viable version |

## Adaptive Planning (Phase 2 Enhancement)

Reassess and adapt plans when new information is discovered:

```bash
# Trigger a replan
vic replan trigger --finding "Found better library"

# List past replans
vic replan list

# Show replan details
vic replan show REPLAN-001
```

Replans are recorded to `.vic-sdd/status/replan-log.yaml` for audit trail.

## Integration with Skills

Phase 1, 2 & 3 features integrate with VIC-SDD skills:

```bash
# Activate vibe-redesign skill
skill vibe-redesign

# Activate adaptive-planning skill
skill adaptive-planning

# Activate vibe-design skill
skill vibe-design
```

See `skills/` directory for skill definitions.

## AI Slop Detection (Phase 3 Enhancement)

Detect and fix AI-generated code and design patterns that hurt quality:

```bash
# List all detection patterns
vic slop list

# Scan current directory
vic slop scan

# Scan specific directory
vic slop scan ./src/components

# Scan for design patterns only
vic slop scan --type design

# Show last scan report
vic slop report

# Auto-fix detected patterns
vic slop fix
```

**Detection Categories:**

| Category | Patterns | Description |
|----------|----------|-------------|
| Design | 6 | Gradient backgrounds, dark mode only, floating cards, etc. |
| Code | 5 | console.log, TODO comments, generic naming, etc. |
| Text | 4 | AI-style phrasing, excessive emojis, tutorial patterns, etc. |

**Scoring:**

| Score | Meaning | Action |
|-------|---------|--------|
| A | Clean | ✅ Excellent |
| B | Minor issues | 👍 Good |
| C | Some issues | ⚠️ Review |
| D | Major issues | ❌ Needs fixing |

**Auto-fix** replaces detected patterns with better alternatives. Always review changes before committing.

## TDD Mode (Phase 4 Enhancement)

Test-Driven Development with Red-Green-Refactor cycle:

```bash
# Start TDD session
vic tdd start --feature "user login"

# RED phase - write failing test first
vic tdd red --test "should validate email"

# GREEN phase - write minimal code to pass
vic tdd green --test "should validate email" --passed

# REFACTOR phase - improve code structure
vic tdd refactor

# Show current status
vic tdd status

# Save checkpoint
vic tdd checkpoint --note "Email validation complete"

# Show history
vic tdd history
```

**TDD Cycle:**

```
🔴 RED → 🟢 GREEN → 🔵 REFACTOR → (repeat)
  Write   Make it    Improve
  failing  pass      code
  test
```

## Systematic Debugging (Phase 4 Enhancement)

4-phase root cause analysis methodology:

```bash
# Start debug session
vic debug start --problem "Login fails in production"

# Phase 1: SURVEY - gather evidence
vic debug survey

# Phase 2: PATTERN - find similar issues
vic debug pattern

# Phase 3: HYPOTHESIS - form and test hypothesis
vic debug hypothesis --explain "Token expired because server time is off"

# Phase 4: IMPLEMENT - fix root cause
vic debug implement --fix "Added token refresh logic" --root-cause "Token expires without refresh"

# Show status
vic debug status

# Generate report
vic debug report

# Show history
vic debug history
```

**Debug Phase Flow:**

```
1️⃣ SURVEY → 2️⃣ PATTERN → 3️⃣ HYPOTHESIS → 4️⃣ IMPLEMENT
Gather       Find         Form &        Fix root
evidence     patterns    test          cause
```

⚠️ **STOP after 3 failed attempts** - question the architecture itself.

## E2E Testing (Phase 4 Enhancement)

Browser automation with Playwright:

```bash
# Initialize QA setup
vic qa init

# Quick smoke test (~30s)
vic qa quick

# Full application test (5-15 min)
vic qa full

# Capture screenshot
vic qa screenshot --name "login-page"

# Show last report
vic qa report

# Show history
vic qa history
```

**QA Modes:**

| Mode | Purpose | Time |
|------|---------|------|
| `quick` | Smoke test critical paths | ~30s |
| `full` | Complete application test | 5-15 min |

**Test Pyramid:**

```
         E2E      ← Vibe QA
      Integration
        Unit Tests  ← spec-driven-test
```

## Integration with Skills

All phases integrate with VIC-SDD skills:

```bash
# Phase 1: Autonomous execution
skill subagent-driven-development
skill executing-plans

# Phase 2: Product thinking
skill vibe-redesign
skill adaptive-planning

# Phase 3: Design system
skill vibe-design

# Phase 4: Testing & QA
skill spec-driven-test # TDD enforcement
skill vibe-debug     # Systematic debugging
skill vibe-qa        # E2E testing
```

See `skills/` directory for full skill definitions.

## License

MIT
