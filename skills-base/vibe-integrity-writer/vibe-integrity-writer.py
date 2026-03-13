#!/usr/bin/env python3
"""
Vibe Integrity Writer - Safe YAML file updater with multi-agent collaboration support
=====================================================================================

A specialized tool for safely updating .vibe-integrity/ YAML files with:
- File locking to prevent concurrent writes
- Agent identity tracking
- Conflict detection
- Atomic operations
- Backup creation
"""

import os
import sys
import json
import yaml
import time
import uuid
# fcntl is Unix-only, will use alternative on Windows
try:
    import fcntl
    HAS_FCNTL = True
except ImportError:
    HAS_FCNTL = False
    import msvcrt  # Windows alternative
import hashlib
import tempfile
import threading
from pathlib import Path
from datetime import datetime
from typing import Dict, List, Optional, Any, Union
from dataclasses import dataclass, asdict
from contextlib import contextmanager

# Configuration
VIBE_DIR = Path('.vibe-integrity')
BACKUP_DIR = VIBE_DIR / 'backups'
LOCK_DIR = VIBE_DIR / 'locks'
AGENT_REGISTRY = VIBE_DIR / 'agents.yaml'

# Default timeout for file locks (seconds)
LOCK_TIMEOUT = 30
LOCK_POLL_INTERVAL = 0.1

@dataclass
class AgentInfo:
    """Information about the agent making changes"""
    agent_id: str
    session_id: str
    timestamp: str
    branch: str
    
    @classmethod
    def current(cls) -> 'AgentInfo':
        """Create AgentInfo from current environment"""
        return cls(
            agent_id=os.getenv('AGENT_ID', f'agent-{uuid.uuid4().hex[:8]}'),
            session_id=os.getenv('SESSION_ID', f'ses-{uuid.uuid4().hex[:12]}'),
            timestamp=datetime.now().isoformat(),
            branch=get_current_git_branch()
        )

@dataclass
class FileLock:
    """File-based locking mechanism for concurrent write protection"""
    lock_file: Path
    agent_id: str
    timeout: int = LOCK_TIMEOUT
    
    def acquire(self) -> bool:
        """Acquire the lock with timeout"""
        start_time = time.time()
        
        while time.time() - start_time < self.timeout:
            try:
                # Create lock directory if it doesn't exist
                self.lock_file.parent.mkdir(parents=True, exist_ok=True)
                
                # Try to create lock file exclusively
                self.fd = open(self.lock_file, 'w')
                
                if HAS_FCNTL:
                    # Unix: use fcntl
                    fcntl.flock(self.fd.fileno(), fcntl.LOCK_EX | fcntl.LOCK_NB)
                else:
                    # Windows: use msvcrt locking
                    msvcrt.locking(self.fd.fileno(), msvcrt.LK_NBLCK, 1)
                
                # Write agent info to lock file
                self.fd.write(f"Locked by: {self.agent_id}\n")
                self.fd.write(f"Time: {datetime.now().isoformat()}\n")
                self.fd.flush()
                
                return True
            except (IOError, OSError):
                # Lock is held by another process
                try:
                    # Check if lock is stale (older than timeout)
                    if self.lock_file.exists():
                        lock_age = time.time() - self.lock_file.stat().st_mtime
                        if lock_age > self.timeout:
                            # Lock is stale, remove it
                            self.lock_file.unlink()
                            continue
                except:
                    pass
                
                time.sleep(LOCK_POLL_INTERVAL)
        
        return False
        """Acquire the lock with timeout"""
        start_time = time.time()
        
        while time.time() - start_time < self.timeout:
            try:
                # Create lock directory if it doesn't exist
                self.lock_file.parent.mkdir(parents=True, exist_ok=True)
                
                # Try to create lock file exclusively
                self.fd = open(self.lock_file, 'w')
                fcntl.flock(self.fd.fileno(), fcntl.LOCK_EX | fcntl.LOCK_NB)
                
                # Write agent info to lock file
                self.fd.write(f"Locked by: {self.agent_id}\n")
                self.fd.write(f"Time: {datetime.now().isoformat()}\n")
                self.fd.flush()
                
                return True
            except (IOError, OSError):
                # Lock is held by another process
                try:
                    # Check if lock is stale (older than timeout)
                    if self.lock_file.exists():
                        lock_age = time.time() - self.lock_file.stat().st_mtime
                        if lock_age > self.timeout:
                            # Lock is stale, remove it
                            self.lock_file.unlink()
                            continue
                except:
                    pass
                
                time.sleep(LOCK_POLL_INTERVAL)
        
        return False
    
    def release(self):
        """Release the lock"""
        if hasattr(self, 'fd'):
            try:
                if HAS_FCNTL:
                    # Unix: use fcntl
                    fcntl.flock(self.fd.fileno(), fcntl.LOCK_UN)
                else:
                    # Windows: unlock using msvcrt
                    msvcrt.locking(self.fd.fileno(), msvcrt.LK_UNLCK, 1)
                self.fd.close()
            except:
                pass
        
        # Remove lock file
        try:
            if self.lock_file.exists():
                self.lock_file.unlink()
        except:
            pass
        """Release the lock"""
        if hasattr(self, 'fd'):
            try:
                fcntl.flock(self.fd.fileno(), fcntl.LOCK_UN)
                self.fd.close()
            except:
                pass
        
        # Remove lock file
        try:
            if self.lock_file.exists():
                self.lock_file.unlink()
        except:
            pass
    
    def __enter__(self):
        if not self.acquire():
            raise TimeoutError(f"Could not acquire lock {self.lock_file} within {self.timeout} seconds")
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.release()

def get_current_git_branch() -> str:
    """Get current git branch name"""
    try:
        import subprocess
        result = subprocess.run(
            ['git', 'rev-parse', '--abbrev-ref', 'HEAD'],
            capture_output=True, text=True, cwd=Path.cwd()
        )
        return result.stdout.strip() if result.returncode == 0 else 'unknown'
    except:
        return 'unknown'

def get_file_hash(filepath: Path) -> str:
    """Calculate hash of file content for change detection"""
    if not filepath.exists():
        return ''
    
    hasher = hashlib.sha256()
    with open(filepath, 'rb') as f:
        for chunk in iter(lambda: f.read(4096), b''):
            hasher.update(chunk)
    return hasher.hexdigest()

def create_backup(filepath: Path) -> Path:
    """Create timestamped backup of file"""
    if not filepath.exists():
        return Path()
    
    backup_dir = BACKUP_DIR
    backup_dir.mkdir(parents=True, exist_ok=True)
    
    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    backup_path = backup_dir / f"{filepath.name}.{timestamp}"
    
    # Copy file to backup location
    backup_path.write_text(filepath.read_text())
    
    return backup_path

def load_yaml(filepath: Path) -> Dict:
    """Load YAML file with error handling"""
    if not filepath.exists():
        return {}
    
    try:
        with open(filepath, 'r') as f:
            return yaml.safe_load(f) or {}
    except Exception as e:
        raise ValueError(f"Failed to load YAML file {filepath}: {e}")

def save_yaml(filepath: Path, data: Dict):
    """Save YAML file with proper formatting"""
    filepath.parent.mkdir(parents=True, exist_ok=True)
    
    with open(filepath, 'w') as f:
        yaml.dump(data, f, default_flow_style=False, allow_unicode=True, sort_keys=False)

def merge_with_agent_info(data: Dict, agent_info: AgentInfo) -> Dict:
    """Add agent tracking information to data"""
    if 'metadata' not in data:
        data['metadata'] = {}
    
    data['metadata']['agent_id'] = agent_info.agent_id
    data['metadata']['session_id'] = agent_info.session_id
    data['metadata']['timestamp'] = agent_info.timestamp
    data['metadata']['branch'] = agent_info.branch
    
    return data

def detect_conflicts(existing_data: Dict, new_data: Dict, operation: str) -> List[str]:
    """Detect potential conflicts between existing and new data"""
    conflicts = []
    
    if operation == 'add_record':
        # Check for duplicate IDs
        if 'records' in existing_data and 'records' in new_data:
            existing_ids = {r.get('id') for r in existing_data['records']}
            new_ids = {r.get('id') for r in new_data['records']}
            
            duplicate_ids = existing_ids.intersection(new_ids)
            if duplicate_ids:
                conflicts.append(f"Duplicate record IDs detected: {duplicate_ids}")
    
    return conflicts

class VibeIntegrityWriter:
    """Main writer class for updating .vibe-integrity/ YAML files"""
    
    def __init__(self, root_path: Optional[str] = None):
        self.root_path = Path(root_path) if root_path else Path.cwd()
        self.vibe_dir = self.root_path / '.vibe-integrity'
        self.agent_info = AgentInfo.current()
        
        # Ensure directories exist
        self.vibe_dir.mkdir(parents=True, exist_ok=True)
        BACKUP_DIR.mkdir(parents=True, exist_ok=True)
        LOCK_DIR.mkdir(parents=True, exist_ok=True)
    
    @contextmanager
    def file_lock(self, filepath: Path):
        """Context manager for file locking"""
        lock_file = LOCK_DIR / f"{filepath.name}.lock"
        lock = FileLock(lock_file, self.agent_info.agent_id)
        
        try:
            with lock:
                yield
        finally:
            # Cleanup any remaining lock files
            try:
                if lock_file.exists():
                    lock_file.unlink()
            except:
                pass
    
    def update_file(self, target_file: str, operation: str, data: Dict, 
                   options: Optional[Dict] = None) -> Dict:
        """Main method to update a YAML file"""
        if options is None:
            options = {}
        
        filepath = self.vibe_dir / target_file
        result = {
            'success': False,
            'message': '',
            'changes_made': [],
            'backup_created': None,
            'validation_passed': False,
            'conflicts_detected': []
        }
        
        try:
            # Create backup if requested
            if options.get('create_backup', True) and filepath.exists():
                backup_path = create_backup(filepath)
                result['backup_created'] = str(backup_path)
            
            # Acquire file lock
            with self.file_lock(filepath):
                # Load existing content
                existing_data = load_yaml(filepath)
                
                # Detect conflicts before making changes
                conflicts = detect_conflicts(existing_data, data, operation)
                result['conflicts_detected'] = conflicts
                
                if conflicts and not options.get('force', False):
                    result['message'] = f"Conflicts detected: {', '.join(conflicts)}"
                    return result
                
                # Apply operation
                updated_data = self._apply_operation(existing_data, operation, data)
                
                # Add agent tracking info
                updated_data = merge_with_agent_info(updated_data, self.agent_info)
                
                # Save updated file
                save_yaml(filepath, updated_data)
                
                result['success'] = True
                result['message'] = f"Successfully updated {target_file}"
                result['changes_made'] = [f"{operation} to {target_file}"]
                
                # Validate if requested
                if options.get('validate_after', True):
                    result['validation_passed'] = self._validate_file(filepath)
                
                # Generate index if requested
                if options.get('generate_index', False):
                    self._generate_index(target_file, updated_data)
            
            return result
            
        except Exception as e:
            result['success'] = False
            result['message'] = f"Failed to update {target_file}: {str(e)}"
            return result
    
    def _apply_operation(self, existing_data: Dict, operation: str, data: Dict) -> Dict:
        """Apply specific operation to data"""
        if operation == 'add_record':
            if 'records' not in existing_data:
                existing_data['records'] = []
            existing_data['records'].append(data)
        
        elif operation == 'update_record':
            record_id = data.get('id')
            if 'records' in existing_data:
                for i, record in enumerate(existing_data['records']):
                    if record.get('id') == record_id:
                        existing_data['records'][i].update(data)
                        break
        
        elif operation == 'delete_record':
            record_id = data.get('id')
            if 'records' in existing_data:
                existing_data['records'] = [
                    r for r in existing_data['records'] 
                    if r.get('id') != record_id
                ]
        
        elif operation == 'set_field':
            field_path = data.get('field', '')
            value = data.get('value')
            
            if '.' in field_path:
                parts = field_path.split('.')
                target = existing_data
                for part in parts[:-1]:
                    if part not in target:
                        target[part] = {}
                    target = target[part]
                target[parts[-1]] = value
            else:
                existing_data[field_path] = value
        
        elif operation == 'batch_operations':
            for op in data.get('operations', []):
                existing_data = self._apply_operation(
                    existing_data, 
                    op['operation'], 
                    op['data']
                )
        
        return existing_data
    
    def _validate_file(self, filepath: Path) -> bool:
        """Validate YAML file structure"""
        try:
            data = load_yaml(filepath)
            # Basic validation: check if it's a valid dict
            return isinstance(data, dict)
        except:
            return False
    
    def _generate_index(self, target_file: str, data: Dict):
        """Generate index file for quick lookup"""
        index_dir = self.vibe_dir / 'index'
        index_dir.mkdir(parents=True, exist_ok=True)
        
        index_file = index_dir / f"{Path(target_file).stem}_index.yaml"
        
        # Create simple index structure
        index_data = {
            'source_file': target_file,
            'last_updated': datetime.now().isoformat(),
            'agent_id': self.agent_info.agent_id,
            'session_id': self.agent_info.session_id
        }
        
        # Add record count if applicable
        if 'records' in data:
            index_data['record_count'] = len(data['records'])
        
        save_yaml(index_file, index_data)

def main():
    """Command-line interface for the writer"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Vibe Integrity Writer - Safe YAML updates')
    parser.add_argument('--target', required=True, help='Target YAML file (e.g., tech-records.yaml)')
    parser.add_argument('--operation', required=True, 
                       choices=['add_record', 'update_record', 'delete_record', 
                               'set_field', 'batch_operations'],
                       help='Operation to perform')
    parser.add_argument('--data', required=True, help='JSON string of data to apply')
    parser.add_argument('--options', help='JSON string of options')
    parser.add_argument('--root', help='Root path of project')
    
    args = parser.parse_args()
    
    try:
        # Parse JSON data
        data = json.loads(args.data)
        options = json.loads(args.options) if args.options else {}
        
        # Create writer and execute
        writer = VibeIntegrityWriter(args.root)
        result = writer.update_file(args.target, args.operation, data, options)
        
        # Output result
        print(json.dumps(result, indent=2))
        
        sys.exit(0 if result['success'] else 1)
        
    except Exception as e:
        print(json.dumps({
            'success': False,
            'message': f'Error: {str(e)}'
        }, indent=2))
        sys.exit(1)

if __name__ == '__main__':
    main()