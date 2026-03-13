#!/usr/bin/env python3
"""
Agent Registry for Vibe Integrity
==================================

Manages agent registration and tracking for multi-agent collaboration.
"""

import os
import sys
import yaml
import uuid
from pathlib import Path
from datetime import datetime
from typing import Dict, List, Optional
from dataclasses import dataclass, asdict

@dataclass
class Agent:
    """Agent registration information"""
    agent_id: str
    name: str
    session_id: str
    created_at: str
    last_seen: str
    branch: str
    status: str  # active, idle, completed
    metadata: Dict[str, str]
    
    @classmethod
    def create(cls, name: str = None, branch: str = None) -> 'Agent':
        """Create a new agent instance"""
        agent_id = f"agent-{uuid.uuid4().hex[:8]}"
        session_id = f"ses-{uuid.uuid4().hex[:12]}"
        
        return cls(
            agent_id=agent_id,
            name=name or f"Agent-{agent_id[-4:]}",
            session_id=session_id,
            created_at=datetime.now().isoformat(),
            last_seen=datetime.now().isoformat(),
            branch=branch or get_current_git_branch(),
            status='active',
            metadata={}
        )

class AgentRegistry:
    """Manages agent registration and tracking"""
    
    def __init__(self, root_path: Optional[str] = None):
        self.root_path = Path(root_path) if root_path else Path.cwd()
        self.registry_file = self.root_path / '.vibe-integrity' / 'agents.yaml'
        
        # Ensure directory exists
        self.registry_file.parent.mkdir(parents=True, exist_ok=True)
    
    def load_registry(self) -> Dict:
        """Load agent registry from file"""
        if not self.registry_file.exists():
            return {'agents': [], 'version': '1.0'}
        
        try:
            with open(self.registry_file, 'r') as f:
                return yaml.safe_load(f) or {'agents': [], 'version': '1.0'}
        except Exception as e:
            print(f"Warning: Could not load agent registry: {e}")
            return {'agents': [], 'version': '1.0'}
    
    def save_registry(self, registry: Dict):
        """Save agent registry to file"""
        try:
            with open(self.registry_file, 'w') as f:
                yaml.dump(registry, f, default_flow_style=False, allow_unicode=True)
        except Exception as e:
            print(f"Error: Could not save agent registry: {e}")
    
    def register_agent(self, agent: Agent) -> Agent:
        """Register a new agent"""
        registry = self.load_registry()
        
        # Check if agent already exists
        for existing_agent in registry.get('agents', []):
            if existing_agent['agent_id'] == agent.agent_id:
                # Update existing agent
                existing_agent.update(asdict(agent))
                self.save_registry(registry)
                return agent
        
        # Add new agent
        if 'agents' not in registry:
            registry['agents'] = []
        registry['agents'].append(asdict(agent))
        
        self.save_registry(registry)
        return agent
    
    def update_agent_status(self, agent_id: str, status: str):
        """Update agent status"""
        registry = self.load_registry()
        
        for agent in registry.get('agents', []):
            if agent['agent_id'] == agent_id:
                agent['status'] = status
                agent['last_seen'] = datetime.now().isoformat()
                break
        
        self.save_registry(registry)
    
    def get_active_agents(self) -> List[Agent]:
        """Get all active agents"""
        registry = self.load_registry()
        active_agents = []
        
        for agent_data in registry.get('agents', []):
            if agent_data.get('status') == 'active':
                active_agents.append(Agent(**agent_data))
        
        return active_agents
    
    def get_agent_by_session(self, session_id: str) -> Optional[Agent]:
        """Get agent by session ID"""
        registry = self.load_registry()
        
        for agent_data in registry.get('agents', []):
            if agent_data.get('session_id') == session_id:
                return Agent(**agent_data)
        
        return None
    
    def cleanup_stale_agents(self, max_age_hours: int = 24):
        """Mark inactive agents as idle"""
        registry = self.load_registry()
        now = datetime.now()
        
        for agent_data in registry.get('agents', []):
            if agent_data.get('status') == 'active':
                try:
                    last_seen = datetime.fromisoformat(agent_data['last_seen'])
                    age_hours = (now - last_seen).total_seconds() / 3600
                    
                    if age_hours > max_age_hours:
                        agent_data['status'] = 'idle'
                except:
                    pass
        
        self.save_registry(registry)

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

def main():
    """Command-line interface"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Agent Registry Manager')
    parser.add_argument('--register', action='store_true', help='Register a new agent')
    parser.add_argument('--name', help='Agent name')
    parser.add_argument('--list-active', action='store_true', help='List active agents')
    parser.add_argument('--update-status', help='Update agent status')
    parser.add_argument('--agent-id', help='Agent ID for operations')
    parser.add_argument('--cleanup', action='store_true', help='Cleanup stale agents')
    
    args = parser.parse_args()
    
    registry = AgentRegistry()
    
    if args.register:
        agent = Agent.create(name=args.name)
        registry.register_agent(agent)
        print(f"Registered agent: {agent.agent_id}")
        print(f"Session ID: {agent.session_id}")
        print(f"Branch: {agent.branch}")
    
    elif args.list_active:
        active_agents = registry.get_active_agents()
        if not active_agents:
            print("No active agents")
        else:
            print(f"Active agents ({len(active_agents)}):")
            for agent in active_agents:
                print(f"  - {agent.name} ({agent.agent_id})")
                print(f"    Session: {agent.session_id}")
                print(f"    Branch: {agent.branch}")
                print(f"    Last seen: {agent.last_seen}")
    
    elif args.update_status and args.agent_id:
        registry.update_agent_status(args.agent_id, args.update_status)
        print(f"Updated status of {args.agent_id} to {args.update_status}")
    
    elif args.cleanup:
        registry.cleanup_stale_agents()
        print("Cleaned up stale agents")

if __name__ == '__main__':
    main()