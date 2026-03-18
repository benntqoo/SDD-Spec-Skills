"""
SDD Content Quality Validator

Validates that SDD artifacts meet quality standards beyond just structure.
"""

from __future__ import annotations

import json
import re
from pathlib import Path
from typing import Any


class ValidationResult:
    def __init__(self):
        self.passed: list[str] = []
        self.warnings: list[str] = []
        self.errors: list[str] = []
    
    @property
    def is_valid(self) -> bool:
        return len(self.errors) == 0
    
    def add_pass(self, message: str) -> None:
        self.passed.append(message)
    
    def add_warning(self, message: str) -> None:
        self.warnings.append(message)
    
    def add_error(self, message: str) -> None:
        self.errors.append(message)
    
    def summary(self) -> str:
        lines = []
        if self.passed:
            lines.append(f"Passed ({len(self.passed)}):")
            for p in self.passed:
                lines.append(f"  ✓ {p}")
        if self.warnings:
            lines.append(f"Warnings ({len(self.warnings)}):")
            for w in self.warnings:
                lines.append(f"  ⚠ {w}")
        if self.errors:
            lines.append(f"Errors ({len(self.errors)}):")
            for e in self.errors:
                lines.append(f"  ✗ {e}")
        return "\n".join(lines)


def validate_json_file(path: Path) -> tuple[bool, str | None]:
    """Validate that a file is valid JSON."""
    try:
        with open(path, 'r', encoding='utf-8') as f:
            json.load(f)
        return True, None
    except json.JSONDecodeError as e:
        return False, f"Invalid JSON: {e}"
    except Exception as e:
        return False, f"Error reading file: {e}"


def validate_spec_file(path: Path) -> ValidationResult:
    """Validate a spec markdown file."""
    result = ValidationResult()
    
    if not path.exists():
        result.add_error(f"Spec file not found: {path}")
        return result
    
    try:
        content = path.read_text(encoding='utf-8')
    except Exception as e:
        result.add_error(f"Cannot read spec file: {e}")
        return result
    
    # Check required sections
    required_sections = [
        "Objective",
        "User Stories",
        "Acceptance Criteria",
        "Domain Models",
    ]
    
    for section in required_sections:
        if section.lower() not in content.lower():
            result.add_warning(f"Missing section: {section}")
    
    # Check for TBD/TODO
    if re.search(r'\bTBD\b', content, re.IGNORECASE):
        result.add_error("Contains unresolved TBD markers")
    
    if re.search(r'\bTODO\b', content, re.IGNORECASE):
        result.add_warning("Contains TODO markers")
    
    # Check for unnamed acceptance criteria
    if re.search(r'^\s*-\s+Acceptance', content, re.MULTILINE):
        if not re.search(r'\[.*?\]', content):
            result.add_warning("Acceptance criteria may lack IDs")
    
    if result.is_valid:
        result.add_pass(f"Spec file valid: {path.name}")
    
    return result


def validate_contract_file(path: Path) -> ValidationResult:
    """Validate a contract JSON file."""
    result = ValidationResult()
    
    if not path.exists():
        result.add_error(f"Contract file not found: {path}")
        return result
    
    # Validate JSON syntax
    valid, error = validate_json_file(path)
    if not valid:
        result.add_error(f"Invalid JSON: {error}")
        return result
    
    try:
        with open(path, 'r', encoding='utf-8') as f:
            contract = json.load(f)
    except Exception as e:
        result.add_error(f"Cannot parse contract: {e}")
        return result
    
    # Validate required fields
    required_fields = ["version", "operations"]
    for field in required_fields:
        if field not in contract:
            result.add_error(f"Missing required field: {field}")
    
    # Validate operations structure
    if "operations" in contract:
        if not isinstance(contract["operations"], list):
            result.add_error("'operations' must be an array")
        elif len(contract["operations"]) == 0:
            result.add_warning("No operations defined")
        else:
            for i, op in enumerate(contract["operations"]):
                if not isinstance(op, dict):
                    result.add_error(f"Operation {i} is not an object")
                    continue
                
                op_id = op.get("id", f"index_{i}")
                if not op.get("input_schema"):
                    result.add_warning(f"Operation {op_id} missing input_schema")
                if not op.get("output_schema"):
                    result.add_warning(f"Operation {op_id} missing output_schema")
    
    # Check compatibility mode
    valid_modes = ["backward", "forward", "strict"]
    if "compatibility_mode" in contract:
        if contract["compatibility_mode"] not in valid_modes:
            result.add_error(f"Invalid compatibility_mode: {contract['compatibility_mode']}")
    
    if result.is_valid:
        result.add_pass(f"Contract file valid: {path.name}")
    
    return result


def validate_traceability_file(path: Path) -> ValidationResult:
    """Validate a traceability YAML/JSON file."""
    result = ValidationResult()
    
    if not path.exists():
        result.add_error(f"Traceability file not found: {path}")
        return result
    
    # Try to parse as JSON first
    valid, error = validate_json_file(path)
    
    if not valid:
        # Could be YAML, just warn
        result.add_warning(f"Cannot parse as JSON (may be YAML): {error}")
    
    if valid:
        try:
            with open(path, 'r', encoding='utf-8') as f:
                trace = json.load(f)
            
            if not isinstance(trace, list):
                result.add_error("Traceability must be an array")
            else:
                required_fields = ["story_id", "acceptance_id", "contract_operation_id"]
                
                for i, row in enumerate(trace):
                    if not isinstance(row, dict):
                        result.add_error(f"Row {i} is not an object")
                        continue
                    
                    for field in required_fields:
                        if field not in row or not row[field]:
                            result.add_error(f"Row {i} missing {field}")
                    
                    # Check for orphan entries
                    if row.get("status") == "verified":
                        if not row.get("test_case_ids"):
                            result.add_warning(f"Row {i} verified but no test cases")
            
            if result.is_valid:
                result.add_pass(f"Traceability file valid: {path.name}")
                
        except Exception as e:
            result.add_error(f"Error validating traceability: {e}")
    
    return result


def validate_test_file(path: Path) -> ValidationResult:
    """Validate a test file exists and has basic structure."""
    result = ValidationResult()
    
    if not path.exists():
        result.add_error(f"Test file not found: {path}")
        return result
    
    # Check file extension
    valid_extensions = [".spec.ts", ".spec.js", ".spec.py", "_test.go", ".test.java"]
    if not any(path.name.endswith(ext) for ext in valid_extensions):
        result.add_warning(f"Unusual test file extension: {path.suffix}")
    
    try:
        content = path.read_text(encoding='utf-8')
        
        # Check for basic test structure
        if "describe" not in content and "it(" not in content and "def test_" not in content:
            result.add_warning("File may not contain valid test cases")
        
        # Check for empty tests
        if re.search(r'it\([\'"][^\'"]+[\'"]\s*,\s*\)\s*{', content):
            result.add_warning("Found empty test case")
        
        if result.is_valid:
            result.add_pass(f"Test file valid: {path.name}")
            
    except Exception as e:
        result.add_error(f"Error reading test file: {e}")
    
    return result


def validate_feature_root(root: Path, feature: str) -> ValidationResult:
    """
    Validate all artifacts for a feature.
    
    Args:
        root: Project root path
        feature: Feature name
    
    Returns:
        ValidationResult with all findings
    """
    result = ValidationResult()
    
    # Support both legacy (docs/specs) and new (.sdd-spec) locations
    sdd_specs = root / ".sdd-spec"
    docs_specs = root / "docs" / "specs"
    tests_specs = root / ".sdd-spec" / "tests"
    
    # Check both locations for specs
    specs_root = sdd_specs if sdd_specs.exists() else docs_specs
    if not specs_root.exists():
        result.add_error("Neither .sdd-spec nor docs/specs directory found")
        return result
    
    # Expected files
    spec_file = specs_root / f"{feature}.md"
    contract_file = specs_root / f"{feature}.contract.json"
    traceability_file = specs_root / f"{feature}.traceability.yaml"
    test_report = specs_root / f"{feature}.test.report.json"
    
    docs_specs = root / "docs" / "specs"
    tests_specs = root / "tests" / "specs"
    
    # Expected files
    spec_file = docs_specs / f"{feature}.md"
    contract_file = docs_specs / f"{feature}.contract.json"
    traceability_file = docs_specs / f"{feature}.traceability.yaml"
    test_report = docs_specs / f"{feature}.test.report.json"
    
    # Validate spec
    spec_result = validate_spec_file(spec_file)
    result.warnings.extend(spec_result.warnings)
    result.errors.extend(spec_result.errors)
    result.passed.extend(spec_result.passed)
    
    # Validate contract
    contract_result = validate_contract_file(contract_file)
    result.warnings.extend(contract_result.warnings)
    result.errors.extend(contract_result.errors)
    result.passed.extend(contract_result.passed)
    
    # Validate traceability
    trace_result = validate_traceability_file(traceability_file)
    result.warnings.extend(trace_result.warnings)
    result.errors.extend(trace_result.errors)
    result.passed.extend(trace_result.passed)
    
    return result


def validate_project(root: Path) -> ValidationResult:
    """
    Validate entire project for SDD compliance.
    
    Args:
        root: Project root path
    
    Returns:
        ValidationResult with all findings
    """
    result = ValidationResult()
    
    # Support both legacy (docs/specs) and new (.sdd-spec) locations
    sdd_specs = root / ".sdd-spec"
    docs_specs = root / "docs" / "specs"
    
    # Use .sdd-spec if exists, otherwise fall back to docs/specs
    specs_root = sdd_specs if sdd_specs.exists() else docs_specs
    
    if not specs_root.exists():
        result.add_error("Neither .sdd-spec nor docs/specs directory found")
        return result
    
    # Find all features
    features = set()
    for f in specs_root.glob("*.md"):
        # Extract feature name (remove .md extension)
        feature = f.stem
        # Skip special files
        if feature in ["README", "INDEX"]:
            continue
        features.add(feature)
    
    if not features:
        location = ".sdd-spec" if sdd_specs.exists() else "docs/specs"
        result.add_warning(f"No feature specs found in {location}/")
        return result
    
    location = ".sdd-spec" if sdd_specs.exists() else "docs/specs"
    result.add_pass(f"Found {len(features)} features in {location}/: {', '.join(sorted(features))}")
    """
    Validate entire project for SDD compliance.
    
    Args:
        root: Project root path
    
    Returns:
        ValidationResult with all findings
    """
    result = ValidationResult()
    
    docs_specs = root / "docs" / "specs"
    
    if not docs_specs.exists():
        result.add_error("docs/specs directory not found")
        return result
    
    # Find all features
    features = set()
    for f in docs_specs.glob("*.md"):
        # Extract feature name (remove .md extension)
        feature = f.stem
        # Skip special files
        if feature in ["README", "INDEX"]:
            continue
        features.add(feature)
    
    if not features:
        HJ|        result.add_warning("No feature specs found in .sdd-spec/specs/")
        return result
    
    result.add_pass(f"Found {len(features)} features: {', '.join(sorted(features))}")
    
    # Validate each feature
    for feature in sorted(features):
        feature_result = validate_feature_root(root, feature)
        result.warnings.extend(feature_result.warnings)
        result.errors.extend(feature_result.errors)
        result.passed.extend(feature_result.passed)
    
    return result


if __name__ == "__main__":
    import sys
    
    if len(sys.argv) < 2:
        print("Usage: python validate-content.py <project_root>")
        sys.exit(1)
    
    root = Path(sys.argv[1]).resolve()
    result = validate_project(root)
    
    print(result.summary())
    sys.exit(0 if result.is_valid else 1)
