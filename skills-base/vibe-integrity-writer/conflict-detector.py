#!/usr/bin/env python3
"""
Conflict Detector for Vibe Integrity
=====================================

Detects potential conflicts in .vibe-integrity/ YAML files:
- Duplicate record IDs
- Similar decisions
- Concurrent modifications
"""

import os
import sys
import yaml
import difflib
from pathlib import Path
from typing import Dict, List, Tuple, Set
from dataclasses import dataclass
from datetime import datetime

@dataclass
class Conflict:
    type: str
    file: str
    description: str
    severity: str  # low, medium, high
    details: Dict

class ConflictDetector:
    def __init__(self, root_path: str = None):
        self.root_path = Path(root_path) if root_path else Path.cwd()
        self.vibe_dir = self.root_path / '.vibe-integrity'
        self.conflicts: List[Conflict] = []
    
    def load_yaml(self, filepath: Path) -> Dict:
        """Load YAML file with error handling"""
        if not filepath.exists():
            return {}
        
        try:
            with open(filepath, 'r') as f:
                return yaml.safe_load(f) or {}
        except Exception as e:
            print(f"Warning: Could not load {filepath}: {e}")
            return {}
    
    def detect_duplicate_ids(self) -> List[Conflict]:
        """Detect duplicate record IDs across all YAML files"""
        conflicts = []
        id_to_files: Dict[str, List[str]] = {}
        
        # Check tech-records.yaml
        tech_records = self.load_yaml(self.vibe_dir / 'tech-records.yaml')
        if 'records' in tech_records:
            for record in tech_records['records']:
                if 'id' in record:
                    record_id = record['id']
                    if record_id not in id_to_files:
                        id_to_files[record_id] = []
                    id_to_files[record_id].append('tech-records.yaml')
        
        # Check risk-zones.yaml
        risk_zones = self.load_yaml(self.vibe_dir / 'risk-zones.yaml')
        if 'zones' in risk_zones:
            for zone_name, zone_data in risk_zones['zones'].items():
                # Zone names should be unique
                if zone_name not in id_to_files:
                    id_to_files[zone_name] = []
                id_to_files[zone_name].append('risk-zones.yaml')
        
        # Find duplicates
        for record_id, files in id_to_files.items():
            if len(files) > 1:
                conflicts.append(Conflict(
                    type='duplicate_id',
                    file=', '.join(files),
                    description=f"Duplicate ID '{record_id}' found in multiple files",
                    severity='high',
                    details={'id': record_id, 'files': files}
                ))
        
        return conflicts
    
    def detect_similar_decisions(self, threshold: float = 0.8) -> List[Conflict]:
        """Detect similar decisions that might be duplicates"""
        conflicts = []
        
        tech_records = self.load_yaml(self.vibe_dir / 'tech-records.yaml')
        if 'records' not in tech_records:
            return conflicts
        
        records = tech_records['records']
        
        # Compare each pair of records
        for i, record1 in enumerate(records):
            for j, record2 in enumerate(records[i+1:], i+1):
                title1 = record1.get('title', '').lower()
                title2 = record2.get('title', '').lower()
                
                # Calculate similarity
                similarity = difflib.SequenceMatcher(None, title1, title2).ratio()
                
                if similarity > threshold:
                    conflicts.append(Conflict(
                        type='similar_decisions',
                        file='tech-records.yaml',
                        description=f"Similar decisions detected: '{record1.get('title')}' and '{record2.get('title')}'",
                        severity='medium',
                        details={
                            'similarity': similarity,
                            'id1': record1.get('id'),
                            'id2': record2.get('id'),
                            'title1': record1.get('title'),
                            'title2': record2.get('title')
                        }
                    ))
        
        return conflicts
    
    def detect_concurrent_modifications(self) -> List[Conflict]:
        """Detect files modified by multiple agents recently"""
        conflicts = []
        
        # Check for recent modifications in backups
        backup_dir = self.vibe_dir / 'backups'
        if backup_dir.exists():
            backup_files = list(backup_dir.glob('*.yaml.*'))
            
            # Group by original filename
            file_groups = {}
            for backup in backup_files:
                # Extract original filename and timestamp
                parts = backup.name.split('.')
                if len(parts) >= 3:
                    original = '.'.join(parts[:-1])
                    if original not in file_groups:
                        file_groups[original] = []
                    file_groups[original].append(backup)
            
            # Check for files with multiple recent backups
            for filename, backups in file_groups.items():
                if len(backups) > 2:  # More than 2 backups suggests multiple modifications
                    conflicts.append(Conflict(
                        type='frequent_modifications',
                        file=filename,
                        description=f"File '{filename}' has {len(backups)} recent backups, suggesting frequent modifications",
                        severity='low',
                        details={'backup_count': len(backups)}
                    ))
        
        return conflicts
    
    def detect_missing_metadata(self) -> List[Conflict]:
        """Detect YAML files missing agent metadata"""
        conflicts = []
        
        yaml_files = [
            'tech-records.yaml',
            'risk-zones.yaml',
            'dependency-graph.yaml',
            'module-map.yaml',
            'schema-evolution.yaml'
        ]
        
        for filename in yaml_files:
            filepath = self.vibe_dir / filename
            if filepath.exists():
                data = self.load_yaml(filepath)
                if 'metadata' not in data:
                    conflicts.append(Conflict(
                        type='missing_metadata',
                        file=filename,
                        description=f"File '{filename}' is missing metadata section",
                        severity='medium',
                        details={'filename': filename}
                    ))
        
        return conflicts
    
    def run_detection(self) -> List[Conflict]:
        """Run all conflict detection checks"""
        print("Running conflict detection...")
        
        self.conflicts = []
        
        # Run all detection methods
        self.conflicts.extend(self.detect_duplicate_ids())
        self.conflicts.extend(self.detect_similar_decisions())
        self.conflicts.extend(self.detect_concurrent_modifications())
        self.conflicts.extend(self.detect_missing_metadata())
        
        return self.conflicts
    
    def print_report(self):
        """Print conflict detection report"""
        if not self.conflicts:
            print("✓ No conflicts detected")
            return
        
        print(f"\n⚠ Found {len(self.conflicts)} potential conflicts:\n")
        
        severity_order = {'high': 0, 'medium': 1, 'low': 2}
        sorted_conflicts = sorted(
            self.conflicts,
            key=lambda c: severity_order.get(c.severity, 3)
        )
        
        for i, conflict in enumerate(sorted_conflicts, 1):
            severity_icon = {
                'high': '🔴',
                'medium': '🟡',
                'low': '🔵'
            }.get(conflict.severity, '⚪')
            
            print(f"{severity_icon} [{conflict.severity.upper()}] {conflict.type}")
            print(f"   File: {conflict.file}")
            print(f"   Description: {conflict.description}")
            if conflict.details:
                print(f"   Details: {conflict.details}")
            print()
    
    def get_conflicts_by_severity(self, severity: str) -> List[Conflict]:
        """Get conflicts filtered by severity"""
        return [c for c in self.conflicts if c.severity == severity]
    
    def has_high_severity_conflicts(self) -> bool:
        """Check if there are any high-severity conflicts"""
        return any(c.severity == 'high' for c in self.conflicts)

def main():
    """Command-line interface"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Detect conflicts in Vibe Integrity YAML files')
    parser.add_argument('--root', help='Root path of project')
    parser.add_argument('--severity', choices=['low', 'medium', 'high'], 
                       help='Filter conflicts by severity')
    parser.add_argument('--json', action='store_true', 
                       help='Output in JSON format')
    
    args = parser.parse_args()
    
    detector = ConflictDetector(args.root)
    conflicts = detector.run_detection()
    
    if args.json:
        import json
        output = {
            'conflicts': [
                {
                    'type': c.type,
                    'file': c.file,
                    'description': c.description,
                    'severity': c.severity,
                    'details': c.details
                }
                for c in conflicts
            ]
        }
        print(json.dumps(output, indent=2))
    else:
        if args.severity:
            filtered = detector.get_conflicts_by_severity(args.severity)
            if not filtered:
                print(f"✓ No {args.severity} severity conflicts detected")
                return
            
            print(f"Found {len(filtered)} {args.severity} severity conflicts:")
            for conflict in filtered:
                print(f"  - {conflict.description}")
        else:
            detector.print_report()
    
    # Exit with error code if high severity conflicts found
    if detector.has_high_severity_conflicts():
        sys.exit(1)

if __name__ == '__main__':
    main()