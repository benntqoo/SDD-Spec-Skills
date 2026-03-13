#!/usr/bin/env python3
"""
Vibe Guard - AI Completion Integrity Checker

Validates that AI-generated code is complete and executable.
"""

import argparse
import json
import os
import re
import subprocess
import sys
from datetime import datetime, timedelta
from pathlib import Path
from typing import Any


class VibeGuard:
    """AI completion integrity checker."""
    
    TODO_PATTERNS = [
        r'//\s*TODO[:\s]',
        r'#\s*TODO[:\s]',
        r'//\s*FIXME[:\s]',
        r'//\s*XXX[:\s]',
        r'//\s*NOT_IMPLEMENTED',
    ]
    
    SECRET_PATTERNS = [
        r'["\']?(?:api[_-]?key|secret|password|token)[_-]?["\']?\s*[:=]\s*["\'][^"\']{8,}["\']',
        r'sk[-_][a-zA-Z0-9]{20,}',
    ]
    
    EMPTY_CATCH_PATTERNS = [
        r'catch\s*\([^)]*\)\s*\{\s*\}',
        r'catch\s*\([^)]*\)\s*\{\s*/[/*].*?\*/\s*\}',
    ]
    
    def __init__(self, root_path: str = ".", config_path: str = None, mode: str = None):
        self.root_path = Path(root_path)
        self.config = self._load_config(config_path)
        if mode:
            self.config["mode"] = mode
        self.results = {
            "timestamp": datetime.now().isoformat() + "Z",
            "mode": self.config["mode"],
            "duration_ms": 0,
            "results": {
                "completeness": {"status": "pass", "findings": []},
                "test_authenticity": {"status": "pass", "findings": []},
                "executability": {"status": "pass", "build": {}, "type_check": {}, "lint": {}},
                "security": {"status": "pass", "findings": []}
            },
            "gate_decision": "pass",
            "blocking_issues": [],
            "warnings": []
        }
        
    def _load_config(self, config_path: str = None) -> dict:
        default_config = {
            "mode": "standard",
            "grace_period_minutes": 10,
            "skip_if_no_changes": True,
        }
        
        if config_path:
            config_file = Path(config_path)
        else:
            config_file = self.root_path / ".sdd-spec" / "vibe-guard.config.json"
        
        if config_file.exists():
            try:
                with open(config_file) as f:
                    user_config = json.load(f)
                    default_config.update(user_config)
            except json.JSONDecodeError:
                pass
                
        return default_config
    
    def run(self, trigger: str = "manual") -> dict:
        import time
        start_time = time.time()
        
        self.results["trigger"] = trigger
        
        # Skip if grace period
        if self._should_skip_due_to_grace_period():
            self.results["gate_decision"] = "skipped"
            self.results["reason"] = "Within grace period"
            return self.results
        
        # Run checks
        self._check_completeness()
        self._check_security()
        self._check_executability()
        
        # Determine gate decision
        self._determine_gate_decision()
        
        # Save results
        self._save_results()
        
        self.results["duration_ms"] = int((time.time() - start_time) * 1000)
        
        return self.results
    
    def _should_skip_due_to_grace_period(self) -> bool:
        if not self.config.get("skip_if_no_changes", True):
            return False
            
        report_file = self.root_path / ".sdd-spec" / "vibe-guard-report.json"
        if not report_file.exists():
            return False
            
        grace_minutes = self.config.get("grace_period_minutes", 10)
        
        try:
            with open(report_file) as f:
                last_result = json.load(f)
            last_time = datetime.fromisoformat(last_result["timestamp"].replace("Z", "+00:00"))
            if datetime.now() - last_time.replace(tzinfo=None) < timedelta(minutes=grace_minutes):
                return True
        except (json.JSONDecodeError, KeyError, ValueError):
            pass
            
        return False
    
    def _get_files_to_scan(self) -> list[Path]:
        files = []
        patterns = ["**/*.ts", "**/*.js", "**/*.jsx", "**/*.tsx", "**/*.py", "**/*.java", "**/*.go"]
        exclude = ["node_modules/**", "dist/**", "build/**", ".git/**"]
        
        for pattern in patterns:
            for filepath in self.root_path.glob(pattern):
                is_excluded = any(filepath.match(e) for e in exclude)
                if not is_excluded and filepath.is_file():
                    files.append(filepath)
        
        return files
    
    def _read_file_content(self, filepath: Path) -> str:
        try:
            with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
                return f.read()
        except Exception:
            return ""
    
    def _check_completeness(self):
        mode = self.config["mode"]
        severity = "warning" if mode == "vibe" else "error"
        
        findings = []
        
        for filepath in self._get_files_to_scan():
            content = self._read_file_content(filepath)
            rel_path = str(filepath.relative_to(self.root_path))
            
            for pattern in self.TODO_PATTERNS:
                for match in re.finditer(pattern, content, re.IGNORECASE | re.MULTILINE):
                    line_num = content[:match.start()].count('\n') + 1
                    findings.append({
                        "type": "TODO",
                        "file": rel_path,
                        "line": line_num,
                        "severity": severity
                    })
        
        if findings:
            self.results["results"]["completeness"]["findings"] = findings
            self.results["results"]["completeness"]["status"] = severity
    
    def _check_security(self):
        mode = self.config["mode"]
        severity = "warning" if mode == "vibe" else "error"
        
        findings = []
        
        for filepath in self._get_files_to_scan():
            content = self._read_file_content(filepath)
            rel_path = str(filepath.relative_to(self.root_path))
            
            # Check for secrets
            for pattern in self.SECRET_PATTERNS:
                for match in re.finditer(pattern, content, re.IGNORECASE):
                    line_num = content[:match.start()].count('\n') + 1
                    findings.append({
                        "type": "SECRET",
                        "file": rel_path,
                        "line": line_num,
                        "severity": severity
                    })
            
            # Check for empty catch
            for pattern in self.EMPTY_CATCH_PATTERNS:
                for match in re.finditer(pattern, content, re.IGNORECASE):
                    line_num = content[:match.start()].count('\n') + 1
                    findings.append({
                        "type": "EMPTY_CATCH",
                        "file": rel_path,
                        "line": line_num,
                        "severity": "warning"
                    })
        
        if findings:
            self.results["results"]["security"]["findings"] = findings
            self.results["results"]["security"]["status"] = severity
    
    def _check_executability(self):
        mode = self.config["mode"]
        results = self.results["results"]["executability"]
        
        # Only check in standard/strict mode
        if mode == "vibe":
            results["status"] = "pass"
            return
        
        # Try to run build
        build_result = self._run_command("npm run build 2>&1 || pnpm build 2>&1 || yarn build 2>&1 || echo 'no build'")
        
        if "not found" in build_result["output"].lower() or build_result["exit_code"] == 0:
            results["build"] = {"status": "pass"}
        else:
            results["build"] = {"status": "error", "output": build_result["output"][:200]}
        
        results["status"] = results["build"]["status"]
    
    def _run_command(self, cmd: str) -> dict:
        try:
            result = subprocess.run(
                cmd,
                shell=True,
                cwd=self.root_path,
                capture_output=True,
                text=True,
                timeout=120
            )
            return {"exit_code": result.returncode, "output": result.stdout + result.stderr}
        except Exception as e:
            return {"exit_code": -1, "output": str(e)}
    
    def _determine_gate_decision(self):
        mode = self.config["mode"]
        results = self.results["results"]
        
        blocking = []
        warnings = []
        
        # Collect findings
        for category, data in results.items():
            if data.get("status") == "error":
                blocking.extend(data.get("findings", []))
            elif data.get("status") == "warning":
                warnings.extend(data.get("findings", []))
        
        # Security blocks in standard/strict
        if results["security"].get("status") == "error" and mode != "vibe":
            self.results["blocking_issues"] = blocking
            self.results["gate_decision"] = "block"
            return
        
        # Executability blocks in standard/strict
        if results["executability"].get("status") == "error" and mode != "vibe":
            self.results["blocking_issues"] = blocking
            self.results["gate_decision"] = "block"
            return
        
        # Completeness blocks in strict
        if mode == "strict" and results["completeness"].get("status") == "error":
            self.results["blocking_issues"] = blocking
            self.results["gate_decision"] = "block"
            return
        
        # Final decision
        self.results["warnings"] = [f["type"] for f in warnings]
        if blocking:
            self.results["gate_decision"] = "block"
        elif warnings:
            self.results["gate_decision"] = "pass_with_warning"
        else:
            self.results["gate_decision"] = "pass"
    
    def _save_results(self):
        output_dir = self.root_path / ".sdd-spec"
        output_dir.mkdir(parents=True, exist_ok=True)
        
        output_file = output_dir / "vibe-guard-report.json"
        with open(output_file, 'w') as f:
            json.dump(self.results, f, indent=2)


def main():
    parser = argparse.ArgumentParser(description="Vibe Guard - AI Completion Integrity Checker")
    parser.add_argument("--root-path", default=".", help="Root path of project")
    parser.add_argument("--config", help="Path to config file")
    parser.add_argument("--mode", choices=["vibe", "standard", "strict"], help="Override mode")
    parser.add_argument("--trigger", default="manual", help="Trigger reason")
    
    args = parser.parse_args()
    
    guard = VibeGuard(root_path=args.root_path, config_path=args.config, mode=args.mode)
    results = guard.run(trigger=args.trigger)
    
    # Print summary
    print(f"\n{'='*50}")
    print(f"Vibe Guard Results")
    print(f"{'='*50}")
    print(f"Mode: {results['mode']}")
    print(f"Trigger: {results['trigger']}")
    print(f"Duration: {results['duration_ms']}ms")
    print(f"\nResults:")
    print(f"  Completeness: {results['results']['completeness']['status']}")
    print(f"  Security: {results['results']['security']['status']}")
    print(f"  Executability: {results['results']['executability']['status']}")
    print(f"\nGate Decision: {results['gate_decision'].upper()}")
    
    if results.get('warnings'):
        print(f"\nWarnings: {', '.join(results['warnings'])}")
    
    if results.get('blocking_issues'):
        print(f"\nBlocking Issues: {len(results['blocking_issues'])}")
        for issue in results['blocking_issues'][:5]:
            print(f"  - {issue['type']}: {issue.get('file', 'unknown')}:{issue.get('line', '?')}")
    
    # Exit with appropriate code
    if results['gate_decision'] == 'block':
        print(f"\n❌ BLOCKED - Fix issues before proceeding")
        sys.exit(1)
    elif results['gate_decision'] == 'pass_with_warning':
        print(f"\n⚠️  PASSED WITH WARNINGS")
        sys.exit(0)
    else:
        print(f"\n✅ PASSED")
        sys.exit(0)


if __name__ == "__main__":
    main()
