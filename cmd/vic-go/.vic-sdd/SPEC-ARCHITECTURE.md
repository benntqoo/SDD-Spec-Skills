# SPEC-ARCHITECTURE.md

## Architecture Overview

> Generated: 2026-03-20

> High-level architecture description

## System Design

### Components

- VIC CLI Tool - Command-line interface for VIBE-SDD
- Core Libraries - Go libraries for SDD operations
- Configuration Management - YAML and file handling
- Gate Checking - Quality gate validation system
- Embedding System - Semantic search with embeddings

### Data Flow

> User input → VIC CLI → Core libraries → Configuration → Gate checks → Report generation

## Data Model

### Core Entities

- **Project**: Contains all project metadata and configuration
- **Phase**: SDD phase with gates and status tracking
- **Gate**: Individual gate check with results
- **TechRecord**: Technology decision records
- **RiskRecord**: Risk assessment records
- **Event**: Audit trail for all actions

### Entity Relationships

```
Project 1---N Phase
Phase 1---2 Gate
Project 1---N TechRecord
Project 1---N RiskRecord
Project 1---N Event
```

## Technology Stack

| Layer | Technology | Rationale |
|-------|------------|----------|
| CLI Tool | Go | Compiled language for performance and cross-platform deployment |
| Configuration | YAML | Human-readable format for complex configurations |
| Database | SQLite | Lightweight, embedded database for project state storage |
| Embeddings | Sentence Transformers | High-quality semantic embeddings for search |
| Git Integration | Git CLI | Native git support for version control |
| Output Formats | JSON/YAML/Plain | Flexible output for different use cases |

## Decision Rationale

### Go for CLI Development
- **Why**: Compiled performance, single binary deployment, strong typing
- **Alternative**: Python (slower startup), Rust (steeper learning curve)
- **Impact**: Faster execution, easier distribution

### SQLite for State Management
- **Why**: Zero configuration, embedded, good for local development
- **Alternative**: PostgreSQL (requires server), JSON files (no querying)
- **Impact**: Reliable state storage with full SQL capabilities

### YAML Configuration
- **Why**: Hierarchical structure, comments support, standard format
- **Alternative**: JSON (no comments), TOML (less adoption)
- **Impact**: Human-readable configurations with documentation

### Sentence Transformers
- **Why**: State-of-the-art semantic embeddings, multilingual support
- **Alternative**: TF-IDF (no semantic understanding), custom models (training cost)
- **Impact**: Accurate code search and similarity matching

## API Design

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/resource | List resources |
| POST | /api/resource | Create resource |

## Security

- Authentication:
- Authorization:
- Data Protection:

## Open Questions

- Question 1
- Question 2
