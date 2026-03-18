# Agent Coordination Protocol

## Protocol Overview

A lightweight coordination protocol for Vibe Integrity agents to work together without conflicts.

## Protocol Messages

### 1. Agent Registration
```json
{
  "type": "register",
  "agent_id": "design-agent-001",
  "session_id": "ses_abc123",
  "timestamp": "2026-03-13T10:30:00Z",
  "capabilities": ["design", "architecture"],
  "branch": "feature/auth-module"
}
```

### 2. File Interest Declaration
```json
{
  "type": "interest",
  "agent_id": "design-agent-001",
  "files": [
    ".vic-sdd/tech-records.yaml",
    ".vic-sdd/dependency-graph.yaml"
  ],
  "expected_duration": "5m"
}
```

### 3. Lock Request
```json
{
  "type": "lock_request",
  "agent_id": "design-agent-001",
  "file": ".vic-sdd/tech-records.yaml",
  "operation": "append",
  "timeout": 30
}
```

### 4. Lock Grant
```json
{
  "type": "lock_grant",
  "file": ".vic-sdd/tech-records.yaml",
  "granted_to": "design-agent-001",
  "expires_at": "2026-03-13T10:30:30Z",
  "token": "lock-token-xyz"
}
```

### 5. Decision Broadcast
```json
{
  "type": "decision",
  "agent_id": "design-agent-001",
  "decision_type": "tech_choice",
  "content": {
    "title": "Use PostgreSQL for main database",
    "reason": "Need ACID transactions for financial module"
  },
  "target_files": ["tech-records.yaml"],
  "timestamp": "2026-03-13T10:30:00Z"
}
```

### 6. Conflict Notification
```json
{
  "type": "conflict",
  "conflict_id": "conf-001",
  "agents": ["design-agent-001", "debug-agent-002"],
  "file": ".vic-sdd/tech-records.yaml",
  "conflicting_changes": [
    {"agent": "design-agent-001", "change": "Add DB choice"},
    {"agent": "debug-agent-002", "change": "Add risk about DB"}
  ],
  "suggested_resolution": "merge_similar"
}
```

### 7. Consensus Request
```json
{
  "type": "consensus_request",
  "request_id": "cons-001",
  "question": "Should we use microservices or monolith?",
  "options": ["microservices", "monolith"],
  "agents": ["design-agent-001", "debug-agent-002"],
  "timeout": 60
}
```

### 8. Consensus Vote
```json
{
  "type": "consensus_vote",
  "request_id": "cons-001",
  "agent_id": "design-agent-001",
  "vote": "microservices",
  "reason": "Better scalability for future growth"
}
```

### 9. Consensus Result
```json
{
  "type": "consensus_result",
  "request_id": "cons-001",
  "decision": "microservices",
  "vote_count": {"microservices": 2, "monolith": 0},
  "timestamp": "2026-03-13T10:35:00Z"
}
```

## Protocol Flow Examples

### Example 1: Coordinated Design Session

```
Agent A (design)          Coordinator           Agent B (debug)
     |                         |                      |
     |---register------------->|                      |
     |                         |                      |
     |---interest------------->|                      |
     |   [tech-records.yaml]   |                      |
     |                         |                      |
     |---lock_request--------->|                      |
     |   [tech-records.yaml]   |                      |
     |                         |                      |
     |<--lock_grant------------|                      |
     |                         |                      |
     |---decision------------->|                      |
     |   [PostgreSQL choice]   |                      |
     |                         |                      |
     |<--release_lock----------|                      |
     |                         |                      |
                                |<--register-----------|
                                |<--interest-----------|
                                |   [tech-records.yaml]|
                                |                      |
                                |<--lock_request-------|
                                |   [tech-records.yaml]|
                                |                      |
                                |---lock_grant-------->|
                                |                      |
                                |<--decision-----------|
                                |   [Risk added]       |
                                |                      |
                                |---release_lock------>|
```

### Example 2: Conflict Resolution

```
Agent A                    Coordinator                    Agent B
  |                            |                            |
  |---decision-------------->|                            |
  |   [Use PostgreSQL]        |                            |
  |                            |                            |
  |                            |<--decision-----------------|
  |                            |   [Use MongoDB]           |
  |                            |                            |
  |<--conflict----------------|                            |
  |   [Conflicting DB choice] |                            |
  |                            |                            |
  |---consensus_request------>|                            |
  |   [DB choice question]    |                            |
  |                            |                            |
  |                            |<--consensus_vote----------|
  |                            |   [Agent A: PostgreSQL]   |
  |                            |<--consensus_vote----------|
  |                            |   [Agent B: MongoDB]      |
  |                            |                            |
  |<--consensus_result---------|                            |
  |   [PostgreSQL chosen]     |                            |
  |                            |                            |
```

## Implementation Approaches

### Option 1: File-based Coordination

```python
# Simple implementation using lock files
class FileCoordination:
    def lock(self, file_path, timeout=30):
        lock_file = f"{file_path}.lock"
        # Create lock file with agent ID and timestamp
        # Wait for lock with timeout
        
    def release(self, file_path):
        # Remove lock file
        
    def get_lock_info(self, file_path):
        # Read lock file to see who has lock
```

**Pros:**
- Simple to implement
- No external dependencies
- Works across processes

**Cons:**
- Doesn't work across machines
- No central coordination
- Race conditions possible

### Option 2: Redis-based Coordination

```python
# Redis-backed coordination service
class RedisCoordination:
    def __init__(self, redis_url):
        self.redis = redis.from_url(redis_url)
        
    def lock(self, key, agent_id, timeout=30):
        lock_key = f"lock:{key}"
        # Use Redis SETNX for atomic lock acquisition
        # Set expiration time
        
    def release(self, key, agent_id):
        # Remove lock if owned by agent
        
    def subscribe(self, pattern, callback):
        # Use Redis pub/sub for notifications
```

**Pros:**
- Works across machines
- Central coordination
- Pub/sub for notifications
- Atomic operations

**Cons:**
- Requires Redis infrastructure
- Single point of failure
- Network dependencies

### Option 3: Message Queue (RabbitMQ/Kafka)

```python
# Message queue-based coordination
class QueueCoordination:
    def __init__(self, queue_url):
        self.connection = pika.BlockingConnection(queue_url)
        
    def publish_decision(self, decision):
        # Publish to decision exchange
        # Route to relevant agents
        
    def subscribe_decisions(self, agent_id, callback):
        # Create queue for agent
        # Bind to decision exchange
```

**Pros:**
- Reliable message delivery
- Decoupled agents
- Scalable
- Persistent messages

**Cons:**
- Complex setup
- Overhead for simple coordination
- Message ordering challenges

### Option 4: Peer-to-Peer (Gossip Protocol)

```python
# Distributed gossip protocol
class GossipCoordination:
    def __init__(self, node_id, peers):
        self.node_id = node_id
        self.peers = peers
        
    def broadcast(self, message):
        # Send to random subset of peers
        # Peers forward to their peers
        # Eventually consistent
        
    def sync_state(self):
        # Exchange state with peers
        # Resolve conflicts via last-write-wins
```

**Pros:**
- No central coordinator
- Highly available
- Works offline
- Scalable

**Cons:**
- Eventual consistency
- Complex conflict resolution
- Network overhead
- Hard to debug

## Recommended Implementation

### For Vibe Integrity: Hybrid Approach

1. **Git-based primary workflow** (current)
2. **File locks for concurrent access** (short-term)
3. **Redis for coordination** (medium-term, optional)
4. **Agent registry in .vic-sdd/** (always)

### Agent Registry Format

```yaml
# .vic-sdd/agents.yaml
agents:
  - id: design-agent-001
    session: ses_abc123
    branch: feature/auth
    last_active: "2026-03-13T10:30:00Z"
    capabilities: [design, architecture]
    
  - id: debug-agent-002
    session: ses_def456
    branch: feature/database
    last_active: "2026-03-13T10:31:00Z"
    capabilities: [debug, analysis]
```

### Coordination Workflow

1. **Agent starts**: Register in agents.yaml
2. **Before editing**: Check if other agents interested in same file
3. **Make decision**: Broadcast to other agents
4. **Finish work**: Update last_active timestamp
5. **Agent ends**: Remove from registry (or timeout)

## Tools and Libraries

### Python Libraries
- `redis-py`: For Redis coordination
- `filelock`: For file-based locking
- `pika`: For RabbitMQ
- `kafka-python`: For Kafka

### Existing Solutions
- **Consul**: Service discovery and coordination
- **ZooKeeper**: Distributed coordination
- **etcd**: Key-value store for coordination
- **Redis**: In-memory data store with pub/sub

## Conclusion

The best approach for Vibe Integrity depends on your deployment scenario:

- **Single machine, few agents**: File-based locks + agent registry
- **Multiple machines, few agents**: Redis coordination
- **Multiple machines, many agents**: Message queue or CRDT
- **Enterprise deployment**: Consul/etcd coordination service

For most teams, starting with file-based locks and agent registry in Git will provide sufficient coordination without added complexity.