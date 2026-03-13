# Vibe Integrity Team Adoption Guide

## Collaboration Best Practices

### 1. Branching Strategy
- Use feature branches for all changes
- Keep main branch stable and always deployable
- Use descriptive branch names: `feature/`, `fix/`, `refactor/`
- Squash and merge pull requests to keep history clean

### 2. .vibe-integrity File Management
- Treat .vibe-integrity files as regular code - review them in PRs
- Avoid editing the same .vibe-integrity files simultaneously when possible
- If simultaneous edits are necessary, communicate with team members
- Use the validation script regularly to catch inconsistencies early

### 3. Regular Maintenance
- Run the validation script as part of your pre-commit hook
- Schedule weekly grooming sessions to archive old records
- Monitor index file sizes and split if they become too large
- Update documentation when significant architectural decisions are made

### 4. Conflict Resolution
- When merge conflicts occur in .vibe-integrity files:
  1. Understand the intent of both changes
  2. Preserve all valid records from both versions
  3. Ensure no duplicate IDs are introduced
  4. Run validation script to verify consistency
  5. Test that the application still works with the merged memory

### 5. Record Keeping Standards
- Use clear, descriptive titles for all records
- Include context: why a decision was made, not just what
- Link related records using IDs when appropriate
- Tag records with relevant modules or components
- Update records when decisions change (don't delete, mark as superseded)

## Team Adoption Roadmap

### Phase 1: Foundation (Weeks 1-2)
- [x] Enhanced validation script with consistency checking and auto-indexing
- [x] Team training on Vibe Integrity concepts and usage
- [x] Set up pre-commit hooks to run validation script
- [x] Establish basic naming conventions and record formats

### Phase 2: Collaboration Optimization (Weeks 3-4)
- [ ] Implement entry-level storage for tech-records (one file per record)
- [ ] Set up automatic index updates via validation script
- [ ] Create conflict resolution guidelines document
- [ ] Pilot with a small team on a non-critical feature

### Phase 3: Scale and Maintain (Weeks 5-6)
- [ ] Extend entry-level storage to risk-zones and other growing files
- [ ] Implement maintenance workflows (groom, stats, query commands)
- [ ] Set up automated archiving based on relevance scoring
- [ ] Establish monthly review process for Vibe Integrity health

### Phase 4: Advanced Features (Ongoing)
- [ ] Implement value-based relevance scoring with decay
- [ ] Create visualization tools for dependency graphs and risk zones
- [ ] Build query interface for searching project memory
- [ ] Integrate with CI/CD pipelines for validation gates

## Recommended Tooling

### Pre-commit Hook Example
```bash
#!/bin/bash
# .git/hooks/pre-commit
python skills-base/vibe-integrity/validate-vibe-integrity.py
```

### Git Attributes for .vibe-integrity Files
Create `.gitattributes` with:
```
.vibe-integrity/*.yaml merge=union
.vibe-integrity/index/*.yaml merge=union
```
This helps reduce merge conflicts by taking both versions when conflicts occur.

## Troubleshooting Common Issues

### Duplicate ID Errors
- Run validation script to identify conflicting IDs
- Decide which record to keep or merge information
- Update IDs to be unique (use UUIDs or timestamp-based IDs)

### Large Index Files
- If index files become too large (>1000 lines), consider:
  - Splitting by module or component
  - Implementing tiered indexing (primary/secondary)
  - Adding search functionality instead of loading entire index

### Merge Conflicts in YAML
- Use a YAML-aware merge tool
- Remember that YAML is sensitive to indentation
- When in doubt, preserve both versions and run validation

## Measuring Success

Track these metrics to gauge adoption and effectiveness:
- Reduction in merge conflicts involving .vibe-integrity files
- Frequency of validation script runs (should be pre-commit)
- Number of records archived per grooming cycle
- Team satisfaction with project memory usability
- Time saved when onboarding new team members

## References
- Entire.io design principles: https://entire.io/
- Vibe Integrity documentation: skills-base/vibe-integrity/README.md
- Validation script: skills-base/vibe-integrity/validate-vibe-integrity.py