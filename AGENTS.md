# Agent Collaboration Guide

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