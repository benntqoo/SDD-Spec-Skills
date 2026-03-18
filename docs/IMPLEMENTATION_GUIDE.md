# Multi-Agent Collaboration Implementation Guide

## Quick Start: Choose Your Approach

### For Most Teams: Enhanced Git Workflow (Recommended)

**Use this if:** You have 1-5 agents, sequential work patterns, or want minimal complexity

**Setup (5 minutes):**
1. Add agent identity tracking to your records
2. Install pre-commit hooks for validation
3. Configure `.gitattributes` for union merge

**Example Implementation:**
```python
# vibe-integrity-writer: Add agent tracking
def write_tech_record(record, agent_id, session_id):
    record['agent'] = agent_id
    record['session'] = session_id
    record['timestamp'] = get_iso_timestamp()
    # Write to YAML file
```

**Pros:**
- ✅ No new infrastructure needed
- ✅ Familiar Git workflow
- ✅ Automatic audit trail
- ✅ Easy to understand

**Cons:**
- ❌ No real-time coordination
- ❌ Manual conflict resolution
- ❌ Long feedback loops

---

### For Teams Needing Real-Time: Redis Coordination

**Use this if:** Multiple agents need to coordinate simultaneously, or you want automatic conflict detection

**Setup (30 minutes):**
1. Install Redis server
2. Add Redis client to vibe-integrity-writer
3. Implement lock manager

**Example Implementation:**
```python
# Redis-based lock manager
class RedisLockManager:
    def acquire_lock(self, file_path, agent_id, timeout=30):
        lock_key = f"lock:{file_path}"
        # Use SETNX for atomic lock acquisition
        acquired = self.redis.setnx(lock_key, f"{agent_id}:{time.time()}")
        if acquired:
            self.redis.expire(lock_key, timeout)
        return acquired
    
    def release_lock(self, file_path, agent_id):
        lock_key = f"lock:{file_path}"
        current_owner = self.redis.get(lock_key)
        if current_owner and current_owner.startswith(agent_id):
            self.redis.delete(lock_key)
```

**Pros:**
- ✅ Real-time coordination
- ✅ Prevents race conditions
- ✅ Central management
- ✅ Pub/sub for notifications

**Cons:**
- ❌ Requires Redis infrastructure
- ❌ Single point of failure
- ❌ Network dependencies

---

### For Large Teams: Coordination Service

**Use this if:** 10+ agents, complex coordination needs, enterprise deployment

**Setup (1-2 days):**
1. Deploy coordination service (Redis, RabbitMQ, or custom)
2. Implement agent registry
3. Create decision queue system
4. Add conflict resolution logic

**Architecture:**
```
┌─────────────────────────────────────────────┐
│         Coordination Service                 │
│  ┌──────────────┐  ┌──────────────────┐     │
│  │  Lock Manager│  │  Decision Queue  │     │
│  │  - Redis     │  │  - RabbitMQ      │     │
│  └──────────────┘  └──────────────────┘     │
└─────────────────────────────────────────────┘
         │                           │
    ┌────┴────┐               ┌─────┴────┐
    │ Agent A │               │ Agent B  │
    └─────────┘               └──────────┘
```

**Pros:**
- ✅ Highly scalable
- ✅ Automatic coordination
- ✅ Enterprise features
- ✅ Central monitoring

**Cons:**
- ❌ Complex infrastructure
- ❌ Higher operational cost
- ❌ Learning curve

## Step-by-Step Implementation

### Step 1: Add Agent Identity Tracking

**File:** `vibe-integrity-writer/SKILL.md`

Update the writer to include agent information:

```python
def update_yaml_file(file_path, data, agent_id, session_id):
    # Add metadata
    data['metadata'] = {
        'agent_id': agent_id,
        'session_id': session_id,
        'timestamp': get_iso_timestamp(),
        'branch': get_current_git_branch()
    }
    
    # Create backup
    backup_file = create_backup(file_path)
    
    # Write with lock
    with file_lock(file_path):
        # Read current content
        current = read_yaml(file_path)
        # Merge changes
        updated = merge_changes(current, data)
        # Write back
        write_yaml(file_path, updated)
    
    # Validate
    validate_yaml(file_path)
```

### Step 2: Implement File Locking

**Option A: Simple File Locks (No External Dependencies)**

```python
import fcntl
import os
import time

class FileLock:
    def __init__(self, file_path, timeout=30):
        self.lock_file = f"{file_path}.lock"
        self.timeout = timeout
        self.fd = None
        
    def acquire(self):
        start_time = time.time()
        while time.time() - start_time < self.timeout:
            try:
                self.fd = open(self.lock_file, 'w')
                fcntl.flock(self.fd, fcntl.LOCK_EX | fcntl.LOCK_NB)
                return True
            except IOError:
                time.sleep(0.1)
        return False
    
    def release(self):
        if self.fd:
            fcntl.flock(self.fd, fcntl.LOCK_UN)
            self.fd.close()
            os.remove(self.lock_file)
    
    def __enter__(self):
        self.acquire()
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.release()
```

**Usage:**
```python
with FileLock(".vic-sdd/tech-records.yaml"):
    # Safe to edit
    writer.update_tech_record(...)
```

**Option B: Redis-based Locks**

```python
import redis
import time

class RedisLock:
    def __init__(self, redis_url, key, timeout=30):
        self.redis = redis.from_url(redis_url)
        self.key = f"lock:{key}"
        self.timeout = timeout
        self.agent_id = None
        
    def acquire(self, agent_id):
        self.agent_id = agent_id
        lock_value = f"{agent_id}:{time.time()}"
        
        # SETNX: Set if Not eXists (atomic)
        acquired = self.redis.setnx(self.key, lock_value)
        
        if acquired:
            self.redis.expire(self.key, self.timeout)
            return True
        return False
    
    def release(self):
        if self.agent_id:
            current = self.redis.get(self.key)
            if current and current.startswith(self.agent_id):
                self.redis.delete(self.key)
```

### Step 3: Create Decision Registry

**File:** `.vic-sdd/decisions.yaml`

```yaml
decisions:
  - id: DEC-001
    title: "Use PostgreSQL for main database"
    agent: design-agent-001
    session: ses_abc123
    timestamp: "2026-03-13T10:30:00Z"
    status: approved
    related_records:
      - DB-001
      - DB-002
    
  - id: DEC-002
    title: "Implement microservices architecture"
    agent: design-agent-002
    session: ses_def456
    timestamp: "2026-03-13T10:35:00Z"
    status: proposed
    votes:
      - agent: design-agent-001
        vote: approve
        reason: "Better scalability"
```

### Step 4: Add Conflict Detection

```python
def detect_conflicts(new_record, existing_records):
    conflicts = []
    
    # Check for duplicate IDs
    if new_record['id'] in [r['id'] for r in existing_records]:
        conflicts.append({
            'type': 'duplicate_id',
            'record_id': new_record['id']
        })
    
    # Check for similar decisions
    for existing in existing_records:
        if is_similar_decision(new_record, existing):
            conflicts.append({
                'type': 'similar_decision',
                'existing': existing['id'],
                'new': new_record['id']
            })
    
    return conflicts

def is_similar_decision(record1, record2):
    # Simple similarity check based on title keywords
    title1 = record1.get('title', '').lower()
    title2 = record2.get('title', '').lower()
    
    # Check for common keywords
    keywords = ['postgresql', 'postgres', 'database', 'sql']
    for keyword in keywords:
        if keyword in title1 and keyword in title2:
            return True
    
    return False
```

### Step 5: Implement Consensus Building

```python
class ConsensusManager:
    def __init__(self, redis_url=None):
        self.redis = redis.from_url(redis_url) if redis_url else None
        
    def request_consensus(self, question, options, agents):
        request_id = f"cons-{uuid.uuid4().hex[:8]}"
        
        # Store in Redis or file
        consensus_request = {
            'id': request_id,
            'question': question,
            'options': options,
            'agents': agents,
            'votes': {},
            'status': 'pending'
        }
        
        if self.redis:
            self.redis.set(f"consensus:{request_id}", json.dumps(consensus_request))
        else:
            # Store in file
            write_consensus_request(consensus_request)
        
        # Notify agents
        self.notify_agents(agents, {
            'type': 'consensus_request',
            'request': consensus_request
        })
        
        return request_id
    
    def vote(self, request_id, agent_id, vote, reason=None):
        # Get consensus request
        request = self.get_consensus_request(request_id)
        
        # Record vote
        request['votes'][agent_id] = {
            'vote': vote,
            'reason': reason,
            'timestamp': get_iso_timestamp()
        }
        
        # Check if all agents voted
        if len(request['votes']) >= len(request['agents']):
            request['status'] = 'completed'
            request['decision'] = self.determine_decision(request)
            
            # Notify agents of result
            self.notify_agents(request['agents'], {
                'type': 'consensus_result',
                'request': request
            })
        
        # Save updated request
        self.save_consensus_request(request)
        
        return request['decision'] if request['status'] == 'completed' else None
    
    def determine_decision(self, request):
        # Simple majority voting
        vote_counts = {}
        for agent_id, vote_data in request['votes'].items():
            vote = vote_data['vote']
            vote_counts[vote] = vote_counts.get(vote, 0) + 1
        
        # Return option with most votes
        return max(vote_counts.items(), key=lambda x: x[1])[0]
```

## Testing Multi-Agent Scenarios

### Test 1: Concurrent Write Test

```python
def test_concurrent_writes():
    """Test two agents writing to same file simultaneously"""
    import threading
    
    def write_record(agent_id, record):
        writer.update_tech_record(record, agent_id, f"session-{agent_id}")
    
    # Start two threads
    thread1 = threading.Thread(
        target=write_record,
        args=("agent1", {"id": "DB-001", "title": "PostgreSQL"})
    )
    thread2 = threading.Thread(
        target=write_record,
        args=("agent2", {"id": "DB-002", "title": "MongoDB"})
    )
    
    thread1.start()
    thread2.start()
    thread1.join()
    thread2.join()
    
    # Verify both records exist
    records = read_tech_records()
    assert len(records) >= 2
    assert "DB-001" in [r['id'] for r in records]
    assert "DB-002" in [r['id'] for r in records]
```

### Test 2: Conflict Detection Test

```python
def test_conflict_detection():
    """Test that conflicts are detected and reported"""
    
    # Add first record
    writer.update_tech_record(
        {"id": "DB-001", "title": "PostgreSQL"},
        "agent1", "session1"
    )
    
    # Try to add duplicate
    conflicts = detect_conflicts(
        {"id": "DB-001", "title": "MongoDB"},
        read_tech_records()
    )
    
    assert len(conflicts) > 0
    assert any(c['type'] == 'duplicate_id' for c in conflicts)
```

### Test 3: Consensus Building Test

```python
def test_consensus_building():
    """Test that agents can reach consensus on decisions"""
    
    consensus_mgr = ConsensusManager()
    
    # Request consensus
    request_id = consensus_mgr.request_consensus(
        question="Use microservices or monolith?",
        options=["microservices", "monolith"],
        agents=["agent1", "agent2"]
    )
    
    # Agents vote
    consensus_mgr.vote(request_id, "agent1", "microservices", "Better scalability")
    decision = consensus_mgr.vote(request_id, "agent2", "microservices", "Easier deployment")
    
    assert decision == "microservices"
```

## Deployment Options

### Option 1: Single Machine (Development)

```
┌─────────────────────────────────────┐
│         Development Machine          │
│  ┌──────────────────────────────┐   │
│  │  Vibe Integrity System       │   │
│  │  - Git repository            │   │
│  │  - Local file locks          │   │
│  │  - Agent registry in YAML    │   │
│  └──────────────────────────────┘   │
└─────────────────────────────────────┘
```

**Best for:** Individual developers, small teams, testing

### Option 2: Multi-Machine (Team)

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Agent A    │    │  Redis      │    │  Agent B    │
│  (Design)   │◄──►│  Coordination│◄──►│  (Debug)    │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    ┌─────────────┐
                    │  Git Repo   │
                    │  (Shared)   │
                    └─────────────┘
```

**Best for:** Distributed teams, medium concurrency

### Option 3: Enterprise (Cloud)

```
┌─────────────────────────────────────────────────┐
│                Cloud Infrastructure              │
│  ┌──────────────┐  ┌──────────────┐            │
│  │  Coordination│  │  Message     │            │
│  │  Service     │  │  Queue       │            │
│  │  (Redis/     │  │  (RabbitMQ)  │            │
│  │   Consul)    │  │              │            │
│  └──────────────┘  └──────────────┘            │
│         │                    │                  │
│    ┌────┴────┐          ┌────┴────┐            │
│    │ Agent A │          │ Agent B │            │
│    └─────────┘          └─────────┘            │
│         │                    │                  │
│         └────────────────────┘                  │
│                    │                            │
│              ┌─────────────┐                   │
│              │  Git Repo   │                   │
│              │  (GitHub/   │                   │
│              │   GitLab)   │                   │
│              └─────────────┘                   │
└─────────────────────────────────────────────────┘
```

**Best for:** Large enterprises, high availability requirements

## Monitoring and Debugging

### Key Metrics to Track

1. **Lock Wait Time**: How long agents wait for locks
2. **Conflict Rate**: How often conflicts occur
3. **Consensus Time**: Time to reach decisions
4. **Agent Activity**: Which agents are active
5. **File Access Patterns**: Which files are most contested

### Debugging Commands

```bash
# Check active locks
redis-cli keys "lock:*"

# View agent registry
cat .vic-sdd/agents.yaml

# Check for conflicts (use vic validate instead)
vic validate

# View decision history
cat .vic-sdd/decisions.yaml
```

## Troubleshooting

### Problem: Agents waiting indefinitely for locks

**Solution:** Implement lock timeout and retry logic
```python
def acquire_lock_with_retry(file_path, agent_id, max_retries=3):
    for attempt in range(max_retries):
        if lock.acquire(file_path, agent_id):
            return True
        time.sleep(1)
    raise Exception(f"Failed to acquire lock after {max_retries} attempts")
```

### Problem: Conflicting decisions not being resolved

**Solution:** Implement automatic conflict resolution
```python
def resolve_conflict(conflict):
    if conflict['type'] == 'duplicate_id':
        # Rename duplicate ID
        return rename_duplicate(conflict)
    elif conflict['type'] == 'similar_decision':
        # Merge similar decisions
        return merge_decisions(conflict)
```

### Problem: Agents not aware of each other's work

**Solution:** Implement broadcast notifications
```python
def broadcast_decision(decision):
    # Send to all registered agents
    for agent in get_registered_agents():
        send_notification(agent['id'], {
            'type': 'decision_made',
            'decision': decision
        })
```

## Conclusion

The best multi-agent collaboration approach depends on your specific needs:

1. **Start with Enhanced Git Workflow** (easiest, no new infrastructure)
2. **Add File Locks** (prevents race conditions, minimal complexity)
3. **Add Redis Coordination** (real-time features, requires Redis)
4. **Build Coordination Service** (enterprise features, complex)

For most teams, steps 1-2 provide sufficient coordination. Only move to step 3+ if you have specific requirements for real-time collaboration or have many concurrent agents.