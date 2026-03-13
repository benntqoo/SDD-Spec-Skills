# Agent Collaboration Guide

**Updated with Multi-Agent Collaboration Features**

## Current System Status

✅ **Multi-Agent Supported**: Yes, via Git branching workflow
✅ **Multi-User Supported**: Yes, via Git collaboration
❌ **Real-Time Collaboration**: No, requires Git merge workflow
❌ **File Locking**: No (use separate branches)

## When Using Multiple Agents

### Scenario 1: Sequential Agents (Recommended)
```
Agent A (design): Completes work → Pushes branch → Creates PR
Agent B (review): Reviews PR → Merges → Continues work
```

### Scenario 2: Parallel Agents (Use Caution)
```
Agent A: Working on branch feature/auth
Agent B: Working on branch feature/database
Both: Use separate branches, merge independently
```

### Scenario 3: Same Branch (Avoid if Possible)
```
⚠️ Risk: Merge conflicts in .vibe-integrity/ files
⚠️ Solution: Coordinate via PR reviews, use union merge
```

## Conflict Resolution Workflow

When multiple agents modify the same YAML files:

1. **Git detects conflict** during merge/pull request
2. **Union merge** preserves both versions (may create duplicates)
3. **Run validation script** to detect duplicate IDs
4. **Manual resolution** required to merge similar decisions
5. **Verify** application still works with merged memory

## Best Practices

1. **Use Separate Branches**: Each agent gets own branch
2. **Document Decisions**: Use tech-records.yaml for major choices
3. **Regular Validation**: Run validation script frequently
4. **PR Reviews**: Review .vibe-integrity/ changes before merging
5. **Conflict Detection**: Use validation script to find duplicates

## Known Limitations

- ❌ No file locking mechanism (use branch isolation)
- ❌ No real-time sync (requires Git workflow)
- ❌ No automatic decision merging (manual resolution needed)
- ❌ No agent identity tracking (can be added to records)

## Future Enhancements

1. File-level locking in vibe-integrity-writer
2. Automatic duplicate ID resolution
3. Custom Git merge driver for YAML files
4. Agent coordination protocol
5. Real-time collaboration interface


## New Multi-Agent Collaboration Features

### File Locking
- ✅ Implemented in `vibe-integrity-writer`
- Lock files created in `.vibe-integrity/locks/`
- 30-second timeout with stale lock detection
- Prevents concurrent writes to same file

### Agent Identity Tracking
- ✅ All YAML updates include agent metadata:
  - `agent_id`: Unique agent identifier
  - `session_id`: Session identifier
  - `timestamp`: ISO format timestamp
  - `branch`: Git branch name

### Conflict Detection
- ✅ `conflict-detector.py` script for:
  - Duplicate ID detection
  - Similar decision detection
  - Concurrent modification detection
  - Missing metadata detection

### Agent Registry
- ✅ `agent-registry.py` for:
  - Agent registration and tracking
  - Status management (active/idle/completed)
  - Session tracking
  - Stale agent cleanup

## Updated Workflow

### Using Multi-Agent Features

1. **Register Agent** (optional):
   ```bash
   python skills-base/vibe-integrity-writer/agent-registry.py --register --name "My Agent"
   ```

2. **Check for Conflicts**:
   ```bash
   python skills-base/vibe-integrity-writer/conflict-detector.py
   ```

3. **Use vibe-integrity-writer** (with automatic agent tracking):
   ```bash
   python skills-base/vibe-integrity-writer/vibe-integrity-writer.py \
     --target tech-records.yaml \
     --operation add_record \
     --data '{"id": "DB-001", "title": "Use PostgreSQL"}'
   ```

4. **Monitor Active Agents**:
   ```bash
   python skills-base/vibe-integrity-writer/agent-registry.py --list-active
   ```

## Best Practices for Multi-Agent Collaboration

1. **Use Separate Branches**: Each agent should work on its own branch
2. **Check for Conflicts**: Run conflict detector before and after changes
3. **Review Agent Activity**: Monitor active agents to avoid conflicts
4. **Clean Up Stale Agents**: Regularly run cleanup to mark inactive agents
5. **Validate After Updates**: Always validate YAML structure after modifications

## Example Multi-Agent Scenario

```bash
# Agent 1: Branch feature/auth
# Agent 2: Branch feature/database

# Both agents register
python skills-base/vibe-integrity-writer/agent-registry.py --register --name "Auth Agent"
python skills-base/vibe-integrity-writer/agent-registry.py --register --name "Database Agent"

# Agent 1 adds authentication decision
python skills-base/vibe-integrity-writer/vibe-integrity-writer.py \
  --target tech-records.yaml \
  --operation add_record \
  --data '{"id": "AUTH-001", "title": "Use JWT for authentication"}'

# Agent 2 adds database decision
python skills-base/vibe-integrity-writer/vibe-integrity-writer.py \
  --target tech-records.yaml \
  --operation add_record \
  --data '{"id": "DB-001", "title": "Use PostgreSQL for main database"}'

# Check for conflicts
python skills-base/vibe-integrity-writer/conflict-detector.py

# List active agents
python skills-base/vibe-integrity-writer/agent-registry.py --list-active
```