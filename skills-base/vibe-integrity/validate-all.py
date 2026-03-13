#!/usr/bin/env python3
"""
Vibe Integrity Framework - Full Validation Script
=================================================
Runs both vibe-guard validation and vibe-integrity validation.
Provides a single command to ensure project completeness and integrity.
"""

import subprocess
import sys
import os
from pathlib import Path

def run_command(cmd, cwd=None):
    """Run a command and return (success, output)"""
    try:
        result = subprocess.run(
            cmd, 
            shell=True, 
            capture_output=True, 
            text=True,
            cwd=cwd
        )
        return result.returncode == 0, result.stdout + result.stderr
    except Exception as e:
        return False, str(e)

def main():
    print("=" * 60)
    print("Vibe Integrity Framework - Full Validation")
    print("=" * 60)
    print()
    
    # Get the script directory to find relative paths
    script_dir = Path(__file__).parent
    project_root = script_dir.parent.parent  # skills/vibe-integrity -> project root
    
    all_passed = True
    
    # 1. Run vibe-guard validation
    print("1. Running vibe-guard validation...")
    print("-" * 40)
    success, output = run_command(
        "python skills/vibe-guard/validate-vibe-guard.py ",
        cwd=project_root
    )
    print(output)
    if success:
        print("✓ vibe-guard validation PASSED")
    else:
        print("✗ vibe-guard validation FAILED")
        all_passed = False
    print()
    
    # 2. Run vibe-integrity validation
    print("2. Running vibe-integrity validation...")
    print("-" * 40)
    success, output = run_command(
        "python skills/vibe-integrity/validate-vibe-integrity.py",
        cwd=project_root
    )
    print(output)
    if success:
        print("✓ vibe-integrity validation PASSED")
    else:
        print("✗ vibe-integrity validation FAILED")
        all_passed = False
    print()
    
    # Summary
    print("=" * 60)
    if all_passed:
        print("🎉 ALL VALIDATIONS PASSED")
        print("Project is complete and structurally sound.")
        return 0
    else:
        print("❌ SOME VALIDATIONS FAILED")
        print("Please fix the issues above before proceeding.")
        return 1

if __name__ == '__main__':
    sys.exit(main())
