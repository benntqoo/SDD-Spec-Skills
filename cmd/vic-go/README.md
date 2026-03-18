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

## License

MIT
