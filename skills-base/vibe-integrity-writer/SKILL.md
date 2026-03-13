---
name: vibe-integrity-writer
description: Specialized skill for safely updating .vibe-integrity/ YAML files with multi-agent collaboration support, file locking, and conflict detection
implementation: skills-base/vibe-integrity-writer/vibe-integrity-writer.py
---
name: vibe-integrity-writer
description: Specialized skill for safely updating .vibe-integrity/ YAML files, designed to be called by other skills like vibe-design and vibe-integrity-debug
---

# Vibe Integrity Writer

## Overview

Vibe Integrity Writer is a specialized skill whose sole purpose is to safely update the YAML files in the `.vibe-integrity/` directory. Unlike other skills that focus on clarification, debugging, or validation, this skill handles the mechanical work of modifying project architecture memory files.

This skill is designed to be **called by other skills** (such as `vibe-design` and `vibe-integrity-debug`) when they have insights or decisions to record, ensuring that:
1. Updates are made safely without breaking file structure
2. Consistent formatting is maintained
3. Associated index files are regenerated when appropriate
4. The process is auditable and reversible
5. Other skills can focus on their core competencies without worrying about YAML mechanics

## Core Philosophy

**Separation of Concerns**: 
- `vibe-design` focuses on requirement clarification and decision capture
- `vibe-integrity-debug` focuses on root cause analysis  
- `vibe-integrity-writer` focuses solely on safe YAML file updates
- `vibe-integrity` validation skills focus on checking integrity

This allows each skill to do one thing well while working together to maintain a self-updating project memory system.

## When to Use

This skill should **almost never be invoked directly by humans**. Instead, it is designed to be:
- Called by `vibe-design` when recording architectural decisions
- Called by `vibe-integrity-debug` when documenting insights from root cause analysis
- Used in automated workflows where YAML updates are needed
- Invoked by custom scripts or orchestrator processes

**Direct human use cases** (rare):
- Recovering from failed automated updates
- Performing bulk imports of historical data
- Executing predefined maintenance tasks
- Testing the writer skill itself

## How It Works

### Input Format
The writer skill expects structured input specifying:
- **Target file**: Which `.vibe-integrity/` YAML file to update
- **Operation**: What type of modification to perform
- **Data**: The information to add, modify, or remove
- **Options**: Additional controls (index generation, validation, etc.)

### Supported Operations
| Operation | Description | Typical Use Case |
|-----------|-------------|------------------|
| `add_record` | Add new entry to a list (e.g., new tech record) | `vibe-design` records a new architecture decision |
| `update_record` | Modify existing entry by ID | Correcting information in a tech record |
| `delete_record` | Remove entry by ID | Removing outdated or incorrect information |
| `append_list` | Add items to a list field | Adding new dependencies to dependency graph |
| `set_field` | Set/replace a scalar value | Updating project version or last_updated timestamp |
| `merge_object` | Deep merge into an object | Updating complex configuration sections |
| `batch_operations` | Perform multiple operations in sequence | Complex updates affecting multiple files |

### Safety Features
1. **Backup Creation**: Automatically creates timestamped backups before modifications
2. **Schema Validation**: Validates against known schemas when available
3. **Atomic Operations**: Either all changes in a batch succeed or none are applied
4. **Post-Update Validation**: Runs structural validation after updates
5. **Index Regeneration**: Automatically updates associated index files when needed
6. **Change Tracking**: Records what was changed for audit trails

## Integration with Other Skills

### Typical Usage by vibe-design
When `vibe-design` captures a decision during clarification:
```python
# After user confirms: "We'll use PostgreSQL for the main database"
writer_input = {
    "target_file": "tech-records.yaml",
    "operation": "add_record",
    "data": {
        "id": "DB-002",
        "date": "2026-03-13",
        "category": "database",
        "title": "Select PostgreSQL for primary data storage",
        "decision": "Use PostgreSQL as primary database",
        "reason": "Need ACID transactions and complex queries for financial module",
        "impact": "medium",
        "status": "completed"
    },
    "options": {
        "generate_index": True,
        "validate_after": True,
        "create_backup": True
    }
}
# Call vibe-integrity-writer with this input
```

### Typical Usage by vibe-integrity-debug
When debugging reveals an architectural insight:
```python
# After discovering: "The auth service actually calls user-service directly, not through API gateway"
writer_input = {
    "target_file": "risk-zones.yaml", 
    "operation": "add_record",
    "data": {
        "id": "RISK-005",
        "date": "2026-03-13",
        "category": "architecture",
        "title": "Tight coupling between auth and user services",
        "description": "Auth service makes direct database calls to user-service bypassing API gateway",
        "impact": "high",
        "status": "identified"
    },
    "options": {
        "generate_index": True,
        "validate_after": True
    }
}
```

## Supported Files
The writer skill can safely update all standard `.vibe-integrity/` files:
- `project.yaml`
- `tech-records.yaml` 
- `dependency-graph.yaml`
- `module-map.yaml`
- `risk-zones.yaml`
- `schema-evolution.yaml`

## Error Handling
The skill provides detailed error reporting for:
- File not found or access denied
- Invalid YAML format in target file
- Schema validation failures
- Operation-specific errors (e.g., trying to update non-existent record)
- Post-update validation failures
- Index generation failures

All errors include context about what operation was attempted and what data was involved.

## Output Format
On success, returns:
```yaml
success: true
message: "Successfully updated tech-records.yaml"
changes_made: ["added record DB-002"]
backup_created: ".vibe-integrity/backups/tech-records.yaml.20260313_143022"
validation_passed: true
index_regenerated: true
```

On failure, returns:
```yaml
success: false
message: "Failed to update tech-records.yaml: Schema validation failed"
error_details: { ... }
changes_attempted: [...]
```

## Best Practices for Caller Skills

1. **Validate Input**: Ensure data conforms to expected structure before calling writer
2. **Batch Related Updates**: When making multiple related changes, use batch_operations
3. **Handle Errors Gracefully**: Have fallback plans if writer fails (e.g., manual edit instructions)
4. **Log Calls**: Record when and why the writer was invoked for audit trails
5. **Respect Backups**: Don't delete backup files unless implementing a cleanup policy

## Example Implementation Snippets

### For vibe-design recording a decision:
```python
def record_tech_decision(title, decision, reason, category="architecture"):
    writer_input = {
        "target_file": "tech-records.yaml",
        "operation": "add_record", 
        "data": {
            "id": generate_uuid(),
            "date": today_iso(),
            "category": category,
            "title": title,
            "decision": decision,
            "reason": reason,
            "impact": assess_impact(decision),
            "status": "completed"
        },
        "options": {
            "generate_index": True,
            "validate_after": True
        }
    }
    return call_skill("vibe-integrity-writer", writer_input)
```

### For vibe-integrity-debug documenting a risk:
```python
def record_discovered_risk(title, description, impact="medium"):
    writer_input = {
        "target_file": "risk-zones.yaml",
        "operation": "add_record",
        "data": {
            "id": generate_uuid(),
            "date": today_iso(),
            "category": "architecture",
            "title": title,
            "description": description,
            "impact": impact,
            "status": "identified"
        },
        "options": {
            "generate_index": True,
            "validate_after": True
        }
    }
    return call_skill("vibe-integrity-writer", writer_input)
```

## Relationship to Existing Validation

This skill complements rather than replaces the existing validation skills:
- Run `vibe-integrity-writer` to make updates
- Run `validate-vibe-integrity.py` to check structural integrity
- Run `validate-all.py` for full system validation
- The writer skill can optionally trigger validation after updates

## Machine Interface

The skill is designed to be called programmatically with JSON/YAML input specifying the update to perform. It returns structured output indicating success/failure and details about what was done.

This makes it suitable for:
- Direct skill invocation from other Opencode skills
- Use in custom automation scripts
- Integration with orchestrator systems
- Testing in isolation

## Implementation Details

### File Locking
The writer implements file-based locking to prevent concurrent writes:
- Lock files are created in `.vibe-integrity/locks/`
- Timeout: 30 seconds (configurable)
- Stale lock detection and cleanup

### Agent Tracking
All updates include agent metadata:
- `agent_id`: Unique identifier for the agent
- `session_id`: Session identifier for tracking
- `timestamp`: ISO format timestamp
- `branch`: Current git branch

### Conflict Detection
The writer detects potential conflicts:
- Duplicate record IDs
- Simultaneous edits to same file
- Provides conflict information in result

### Command-Line Usage
```bash
# Add a new tech record
python skills-base/vibe-integrity-writer/vibe-integrity-writer.py \
  --target tech-records.yaml \
  --operation add_record \
  --data '{"id": "DB-001", "title": "Use PostgreSQL"}' \
  --options '{"validate_after": true, "generate_index": true}'

# Update existing record
python skills-base/vibe-integrity-writer/vibe-integrity-writer.py \
  --target tech-records.yaml \
  --operation update_record \
  --data '{"id": "DB-001", "status": "completed"}'
```