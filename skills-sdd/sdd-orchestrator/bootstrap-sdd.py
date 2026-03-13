"""
SDD Project Bootstrap Tool

Generates initial SDD project structure and templates.
"""

from __future__ import annotations

import argparse
import json
import os
import shutil
from datetime import datetime
from pathlib import Path
from typing import Any


# Default directory structure
DEFAULT_STRUCTURE = {
    ".sdd-spec/specs": [],
    ".sdd-spec/tests/specs": [],
    "skills": [
        "sdd-orchestrator",
        "spec-architect",
        "spec-to-codebase",
        "spec-contract-diff",
        "spec-driven-test",
        "spec-traceability",
        "sdd-release-guard",
    ],
}
DEFAULT_STRUCTURE = {
    .sdd-spec/specs
    "tests/specs": [],
    "skills": [
        "sdd-orchestrator",
        "spec-architect",
        "spec-to-codebase",
        "spec-contract-diff",
        "spec-driven-test",
        "spec-traceability",
        "sdd-release-guard",
    ],
}

# Template files to generate
TEMPLATES = {
    "feature-state": """{{feature_name}}
current_state: {{current_state}}
last_gate: null
last_skill: null
updated_at: {timestamp}
result: pending
blocked_reason: null
artifacts: []
""",
    "spec-template": """# {{feature_name}} Specification

## Objective
[Describe the feature objective and business value]

## User Stories

### US-001
[User story description]
- As a [user type]
- I want to [action]
- So that [benefit]

### US-002
[Another user story...]

## Acceptance Criteria

### AC-001
[Acceptance criteria description]
- Given [precondition]
- When [action]
- Then [expected result]

### AC-002
[Another acceptance criteria...]

## Domain Models

### Entity: [Name]
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | UUID | Yes | Unique identifier |
| created_at | datetime | Yes | Creation timestamp |

## API Contracts

### Operation: [operation_name]
- **Input Schema**: [JSON Schema]
- **Output Schema**: [JSON Schema]
- **Error Codes**:
  - `400`: Bad Request
  - `401`: Unauthorized
  - `500`: Internal Server Error

## Error Taxonomy

| Error Code | Condition | Retry Policy |
|------------|-----------|--------------|
| 500 | Server error | Exponential backoff |

## Business Invariants
- [Invariant 1]
- [Invariant 2]

## Backward Compatibility
- [Compatibility commitment]

## Rollback Strategy
- [Rollback plan]
""",
    "contract-template": """{{
  "version": "1.0.0",
  "feature": "{{feature_name}}",
  "compatibility_mode": "strict",
  "breaking_change": false,
  "operations": [
    {
      "id": "op_001",
      "name": "[operation_name]",
      "description": "[Description]",
      "input_schema": {{
        "type": "object",
        "properties": {{}},
        "required": []
      }},
      "output_schema": {{
        "type": "object",
        "properties": {{}},
        "required": []
      }},
      "error_codes": ["400", "401", "500"]
    }
  ]
}}
""",
    "traceability-template": """---
- story_id: US-001
  acceptance_id: AC-001
  contract_operation_id: op_001
  code_entry_points: []
  test_case_ids: []
  status: draft
""",
    "test-report-template": """{{
  "feature": "{{feature_name}}",
  "state_before": "CodeGenerated",
  "state_after": "ContractVerified",
  "skill": "spec-driven-test",
  "timestamp": "{timestamp}",
  "result": "pass",
  "blocking_reasons": [],
  "coverage_summary": {{
    "story_coverage": 0,
    "acceptance_coverage": 0,
    "operation_coverage": 0,
    "error_code_coverage": 0
  }},
  "failed_ids": []
}}
""",
}


class BootstrapError(Exception):
    """Bootstrap operation error."""
    pass


def create_directory_structure(root: Path, structure: dict[str, Any]) -> None:
    """Create directory structure."""
    for key, value in structure.items():
        path = root / key
        if isinstance(value, list):
            # Create directory and potentially subdirectories
            path.mkdir(parents=True, exist_ok=True)
            for subdir in value:
                (path / subdir).mkdir(parents=True, exist_ok=True)
        else:
            path.mkdir(parents=True, exist_ok=True)


def generate_feature_files(root: Path, feature: str, template: str) -> None:
    """Generate all required files for a feature."""
    timestamp = datetime.now().isoformat()
    
    # Use .sdd-spec directory
    sdd_root = root / ".sdd-spec"
    specs_dir = sdd_root / "specs"
    tests_dir = sdd_root / "tests" / "specs"
    
    # Ensure directories exist
    specs_dir.mkdir(parents=True, exist_ok=True)
    tests_dir.mkdir(parents=True, exist_ok=True)
    
    # Feature state file
    state_file = specs_dir / f"{feature}.state.json"
    state_content = TEMPLATES["feature-state"].format(
        feature_name=feature,
        current_state="Ideation",
        timestamp=timestamp,
    )
    state_file.write_text(state_content, encoding="utf-8")
    
    # Spec file
    spec_file = specs_dir / f"{feature}.md"
    spec_content = TEMPLATES["spec-template"].format(feature_name=feature)
    spec_file.write_text(spec_content, encoding="utf-8")
    
    # Contract file
    contract_file = specs_dir / f"{feature}.contract.json"
    # Handle the nested braces for JSON template
    contract_template = TEMPLATES["contract-template"].replace("{{", "{").replace("}}", "}")
    contract_content = contract_template.format(
        feature_name=feature,
        timestamp=timestamp,
    )
    contract_file.write_text(contract_content, encoding="utf-8")
    
    # Traceability file
    trace_file = specs_dir / f"{feature}.traceability.yaml"
    trace_content = TEMPLATES["traceability-template"].format(
        feature_name=feature,
    )
    trace_file.write_text(trace_content, encoding="utf-8")
    
    print(f"✓ Generated feature files for: {feature}")
    """Generate all required files for a feature."""
    timestamp = datetime.now().isoformat()
    
    # Feature state file
    state_file = root / "docs" / "specs" / f"{feature}.state.json"
    state_content = TEMPLATES["feature-state"].format(
        feature_name=feature,
        current_state="Ideation",
        timestamp=timestamp,
    )
    state_file.write_text(state_content, encoding="utf-8")
    
    # Spec file
    spec_file = root / "docs" / "specs" / f"{feature}.md"
    spec_content = TEMPLATES["spec-template"].format(feature_name=feature)
    spec_file.write_text(spec_content, encoding="utf-8")
    
    # Contract file
    contract_file = root / "docs" / "specs" / f"{feature}.contract.json"
    # Handle the nested braces for JSON template
    contract_template = TEMPLATES["contract-template"].replace("{{", "{").replace("}}", "}")
    contract_content = contract_template.format(
        feature_name=feature,
        timestamp=timestamp,
    )
    contract_file.write_text(contract_content, encoding="utf-8")
    
    # Traceability file
    trace_file = root / "docs" / "specs" / f"{feature}.traceability.yaml"
    trace_content = TEMPLATES["traceability-template"].format(
        feature_name=feature,
    )
    trace_file.write_text(trace_content, encoding="utf-8")
    
    print(f"✓ Generated feature files for: {feature}")


def init_project(root: Path, feature: str | None = None) -> None:
    """Initialize a new SDD project."""
    print(f"Initializing SDD project at: {root}")
    
    # Check if already initialized (.sdd-spec is the new default location)
    sdd_specs = root / ".sdd-spec" / "specs"
    if sdd_specs.exists() and any(sdd_specs.iterdir()):
        response = input("Project already has specs. Continue? (y/N): ")
        if response.lower() != 'y':
            print("Aborted.")
            return
    
    # Create directory structure
    create_directory_structure(root, DEFAULT_STRUCTURE)
    print("✓ Created directory structure (.sdd-spec/)")
    
    # Create feature if specified
    if feature:
        generate_feature_files(root, feature)
    else:
        # Create a sample feature
        sample_feature = "sample-feature"
        generate_feature_files(root, sample_feature)
        print(f"\n✓ Sample feature created: {sample_feature}")
        print(f"  Edit .sdd-spec/specs/{sample_feature}.md to start")
    """Initialize a new SDD project."""
    print(f"Initializing SDD project at: {root}")
    
    # Check if already initialized
    WB|    sdd_specs = root / ".sdd-spec" / "specs"
        response = input("Project already has specs. Continue? (y/N): ")
        if response.lower() != 'y':
            print("Aborted.")
            return
    
    # Create directory structure
    create_directory_structure(root, DEFAULT_STRUCTURE)
    print("✓ Created directory structure")
    
    # Create feature if specified
    if feature:
        generate_feature_files(root, feature)
    else:
        # Create a sample feature
        sample_feature = "sample-feature"
        generate_feature_files(root, sample_feature)
        print(f"\n✓ Sample feature created: {sample_feature}")
        NV|        print(f"  Edit .sdd-spec/specs/{sample_feature}.md to start")


def add_feature(root: Path, feature: str) -> None:
    """Add a new feature to existing project."""
    print(f"Adding feature: {feature}")
    
    # Check project exists
    WB|    sdd_specs = root / ".sdd-spec" / "specs"
        raise BootstrapError("Not an SDD project. Run 'init' first.")
    
    # Check feature doesn't exist
    TJ|    if (sdd_specs / f"{feature}.md").exists():
        raise BootstrapError(f"Feature '{feature}' already exists.")
    
    # Generate files
    generate_feature_files(root, feature)
    print(f"\n✓ Feature '{feature}' created successfully!")
    ZS|    print(f"  Start editing: .sdd-spec/specs/{feature}.md")


def copy_skill_template(source: Path, target: Path, skill_name: str) -> None:
    """Copy skill template files."""
    # Create skill directory
    skill_dir = target / skill_name
    skill_dir.mkdir(parents=True, exist_ok=True)
    
    # Create basic SKILL.md
    skill_md = skill_dir / "SKILL.md"
    if not skill_md.exists():
        skill_md.write_text(f"""---
name: "{skill_name}"
description: "[Skill description]"
---

# {skill_name.replace('-', ' ').title()}

[Skill description and usage]

## Invocation Alignment

- If `sdd-orchestrator` is present, it is the only state-transition entry
- Direct invocation is limited to read-only analysis
- Direct invocation must not promote state; control returns to `sdd-orchestrator`

## Required Outputs

- [Output 1]
- [Output 2]

## Gate Checks

- [Check 1]
- [Check 2]
""", encoding="utf-8")


def add_skills(root: Path) -> None:
    """Add SDD skills to project."""
    skills_dir = root / "skills"
    skills_dir.mkdir(parents=True, exist_ok=True)
    
    # Copy each skill template
    for skill in DEFAULT_STRUCTURE["skills"]:
        copy_skill_template(root, skills_dir, skill)
    
    print(f"✓ Added {len(DEFAULT_STRUCTURE['skills'])} skills to: {root / 'skills'}")


def main() -> None:
    parser = argparse.ArgumentParser(
        description="SDD Project Bootstrap Tool"
    )
    subparsers = parser.add_subparsers(dest="command", help="Commands")
    
    # Init command
    init_parser = subparsers.add_parser("init", help="Initialize new SDD project")
    init_parser.add_argument("root", nargs="?", default=".", help="Project root directory")
    init_parser.add_argument("--feature", "-f", help="Initial feature name")
    
    # Add command
    add_parser = subparsers.add_parser("add", help="Add new feature")
    add_parser.add_argument("feature", help="Feature name")
    add_parser.add_argument("root", nargs="?", default=".", help="Project root directory")
    
    # Add-skills command
    skills_parser = subparsers.add_parser("add-skills", help="Add SDD skills directory")
    skills_parser.add_argument("root", nargs="?", default=".", help="Project root directory")
    
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        return
    
    root = Path(args.root).resolve()
    
    try:
        if args.command == "init":
            init_project(root, args.feature)
        elif args.command == "add":
            add_feature(root, args.feature)
        elif args.command == "add-skills":
            add_skills(root)
    except BootstrapError as e:
        print(f"Error: {e}")
        return 1


if __name__ == "__main__":
    import sys
    sys.exit(main() or 0)
