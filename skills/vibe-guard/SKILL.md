#BQ|---
#XS|name: vibe-guard
#JK|description: AI completion integrity checker - detects false completion claims, TODO placeholders, fake tests, and ensures code is actually executable
#ZM|---
#BT|

#MJ|# Vibe Guard

#HN|

#MS|## Overview

#HZ|Vibe Guard is an AI completion integrity checker designed to prevent false completion claims from AI coding assistants. It detects common AI hallucination patterns including TODO placeholders, empty functions, fake tests, and ensures the code actually works before allowing progression to the next phase.

#KV|

#KV|**Design Philosophy:**
#KV|- **Single entry point** - One skill to rule all AI integrity checks
#KV|- **Configurable strictness** - From vibe (minimal) to strict (full) modes
#KV|- **Non-blocking by default** - Warnings don't stop flow, only critical errors block
#KV|- **Evidence-based** - Requires actual proof of work, not AI self-assessment

#SP|

#KV|## When to Use

#TH|**Use when:**
#HT|- AI claims completion but you're not sure it's actually done
#ZY|- You want to prevent AI from leaving TODO placeholders in production code
#XZ|- You need to verify tests are real and actually passing
#YJ|- You want to ensure code builds and type-checks before proceeding
#WB|- You're transitioning from POC to production and need a sanity check

#KT|

#NJ|**When NOT to use:**
#HQ|- When you're still exploring and TODOs are acceptable
#YQ|- When you explicitly want rapid prototyping without checks
#QW|- For documentation-only changes

#BT|

#MH|## Invocation Alignment

#KW|

#YT|- If `sdd-orchestrator` is present, it can invoke vibe-guard automatically
#QP|- Can be invoked directly for ad-hoc integrity checks
#RV|- Multiple invocations within grace period are skipped (avoids redundant work)
#WH|- Results are cached and can be queried without re-running

#RT|

#BP|## Modes

#NT|Vibe Guard supports three modes that control check strictness:

#XZ|

#HT|| Mode | Use Case | Blocking Conditions |
#XB||------|----------|---------------------|
#HZ|| `vibe` | Rapid prototyping, POC | Build failure, critical security |
#KB|| `standard` | SME projects, team development | Build + security + core tests |
#JJ|| `strict` | Enterprise, production | All checks fail |

#RB|

#BK|### Mode Selection

#YQ|Mode can be set via:
#YQ|1. **Config file**: `.sdd-spec/vibe-guard.config.json`
#YQ|2. **Environment variable**: `VIBE_GUARD_MODE=vibe|standard|strict`
#YQ|3. **Auto-detection**: Orchestrator determines based on project state

#RT|

#BP|## Check Categories

#NT|

#HT|### A. Completeness Checks

#XB|| Check | vibe | standard | strict | Description |
#QM||---------|------|----------|--------|-------------|
#PQ|| TODO/FIXME | ⚠️ warn | ⚠️ warn | ❌ error | Scans for incomplete implementations |
#YJ|| Empty functions | ⚠️ warn | ⚠️ warn | ❌ error | Detects stub/placeholder functions |
#RX|| Hardcoded fake data | ⚠️ warn | ❌ error | ❌ error | Finds mock data left in code |
#QP|| Unimplemented imports | ⚠️ warn | ⚠️ warn | ❌ error | Imports that aren't used |

#RN|

#HT|### B. Test Authenticity Checks

#XB|| Check | vibe | standard | strict | Description |
#QM||---------|------|----------|--------|-------------|
#RT|| Tests run | optional | required | required | Must actually execute tests |
#HQ|| Fake tests | ⚠️ warn | ⚠️ warn | ❌ error | Empty tests, always-pass assertions |
#NM|| Test coverage | - | >60% | >80% | Minimum coverage threshold |
#RR|| Error path tests | ⚠️ warn | ⚠️ warn | ❌ error | Must test error scenarios |

#RN|

#HT|### C. Executability Checks (Always Blocking)

#XB|| Check | vibe | standard | strict | Description |
#QM||---------|------|----------|--------|-------------|
#QM|| Build success | ❌ error | ❌ error | ❌ error | `npm run build` or equivalent |
#RT|| Type check | - | ❌ error | ❌ error | TypeScript/mypy errors |
#RX|| Lint | - | ⚠️ warn | ❌ error | Linting violations |

#RN|

#HT|### D. Security Checks (Always Blocking in standard+)

#XB|| Check | vibe | standard | strict | Description |
#QM||---------|------|----------|--------|-------------|
#QP|| Hardcoded secrets | ❌ error | ❌ error | ❌ error | API keys, passwords in code |
#NM|| SQL injection risk | ⚠️ warn | ❌ error | ❌ error | Raw SQL with string concat |
#RN|| XSS vulnerability | ⚠️ warn | ❌ error | ❌ error | Unsanitized user input |
#TH|| Empty catch blocks | ⚠️ warn | ⚠️ warn | ❌ error | Swallowed exceptions |

#RN|

#BP|## Trigger Conditions

#NT|

#YT|Vibe Guard can be triggered automatically or manually:

#RB|

#HT|### Automatic Triggers
#XB|1. **Phrase detection**: AI uses completion words (done, ready, complete, finished, 完成)
#XZ|2. **Phase transition**: Moving from one SDD state to another
#XZ|3. **PR/MR creation**: When AI creates a pull/merge request
#XZ|4. **Periodic**: Configurable interval (default: disabled)

#RB|

#HT|### Manual Triggers
#XB|```bash
#XZ|# Direct invocation
#XZ|python validate-vibe-guard.py --check

#XZ|# Force strict mode
#XZ|python validate-vibe-guard.py --mode strict

#XZ|# Specific category only
#XZ|python validate-vibe-guard.py --category security
#```

#RB|

#BP|## Grace Period & Deduplication

#NT|

#HB|To avoid redundant checks:
#HB|- **Grace period**: 10 minutes (configurable) - skip if checked recently
#HB|- **Skip if no changes**: Don't re-check if no files changed since last run
#HB|- **Cache results**: Query cached results without re-running

#RB|

#BP|## Required Outputs

#NT|

#HV|`.sdd-spec/vibe-guard-report.json` - Detailed check results:
#HQ|
#HQ|```json
#HQ|{
#HQ|  "timestamp": "2026-03-12T10:30:00Z",
#HQ|  "mode": "standard",
#HQ|  "trigger": "phase_transition",
#HQ|  "duration_ms": 3200,
#HQ|  "results": {
#HQ|    "completeness": {
#HQ|      "status": "warning",
#HQ|      "findings": [
#HQ|        {
#HQ|          "type": "TODO",
#HQ|          "file": "src/payment.ts",
#HQ|          "line": 42,
#HQ|          "content": "// TODO: implement refund"
#HQ|        }
#HQ|      ]
#HQ|    },
#HQ|    "test_authenticity": {
#HQ|      "status": "pass",
#HQ|      "tests_run": {
#HQ|        "command": "npm test",
#HQ|        "exit_code": 0,
#HQ|        "passed": 42,
#HQ|        "failed": 0,
#HQ|        "skipped": 3
#HQ|      }
#HQ|    },
#HQ|    "executability": {
#HQ|      "status": "pass",
#HQ|      "build": { "status": "pass" },
#HQ|      "type_check": { "status": "pass" }
#HQ|    },
#HQ|    "security": {
#HQ|      "status": "pass",
#HQ|      "secrets_found": []
#HQ|    }
#HQ|  },
#HQ|  "gate_decision": "pass_with_warning",
#HQ|  "blocking_issues": [],
#HQ|  "warnings": ["TODO found in src/payment.ts:42"]
#HQ|}
#```

#RB|

#BP|## Gate Decision Logic

#NT|

#KV|```
#KV|                    ┌──────────────────┐
#KV|                    │  Run all checks   │
#KV|                    └────────┬─────────┘
#KV|                             │
#KV|              ┌──────────────┼──────────────┐
#KV|              ▼              ▼              ▼
#KV|       ┌──────────┐   ┌──────────┐   ┌──────────┐
#KV|       │ 0 errors │   │ 1+ errors│   │ warnings │
#KV|       │ (any warn)│   │ (blocking│   │ only     │
#KV|       └────┬─────┘   └────┬─────┘   └────┬─────┘
#KV|            │              │              │
#KV|            ▼              ▼              ▼
#KV|     pass_with_    BLOCK      pass_with_
#KV|       warning                   warning
#KV|
#KV|Note: In 'vibe' mode, warnings don't block
#KV|      In 'strict' mode, warnings block
#KV```

#RB|

#BP|## Configuration

#NT|

#KV|Default config file: `.sdd-spec/vibe-guard.config.json`

#KV|```json
#KV|{
#KV|  "mode": "standard",
#KV|  "auto_trigger": true,
#KV|  "trigger_phrases": ["done", "ready", "complete", "finished", "完成", "完成了"],
#KV|  "grace_period_minutes": 10,
#KV|  "skip_if_no_changes": true,
#KV|  "checks": {
#KV|    "completeness": {
#KV|      "scan_todo": true,
#KV|      "scan_stub": true,
#KV|      "scan_fake_data": true,
#KV|      "severity_by_mode": {
#KV|        "vibe": "warning",
#KV|        "standard": "warning",
#KV|        "strict": "error"
#KV|      }
#KV|    },
#KV|    "test_authenticity": {
#KV|      "require_execution": {
#KV|        "vibe": false,
#KV|        "standard": true,
#KV|        "strict": true
#KV|      },
#KV|      "detect_fake_tests": true,
#KV|      "min_coverage": {
#KV|        "vibe": 0,
#KV|        "standard": 60,
#KV|        "strict": 80
#KV|      }
#KV|    },
#KV|    "executability": {
#KV|      "build": "error",
#KV|      "type_check": {
#KV|        "vibe": "skip",
#KV|        "standard": "error",
#KV|        "strict": "error"
#KV|      },
#KV|      "lint": {
#KV|        "vibe": "skip",
#KV|        "standard": "warning",
#KV|        "strict": "error"
#KV|      }
#KV|    },
#KV|    "security": {
#KV|      "scan_secrets": true,
#KV|      "scan_hardcoded": true,
#KV|      "severity": "error"
#KV|    }
#KV|  },
#KV|  "file_patterns": {
#KV|    "scan": ["**/*.ts", "**/*.js", "**/*.jsx", "**/*.tsx", "**/*.py", "**/*.java"],
#KV|    "exclude": ["node_modules/**", "dist/**", "build/**", ".git/**"]
#KV|  }
#KV|}
#```

#RB|

#BP|## Common Detection Patterns

#NT|

#HT|### TODO/FIXME Detection
#XZ|```typescript
#XZ|// ❌ BLOCK - Incomplete implementation
#XZ|function processPayment() {
#XZ|  // TODO: implement
#XZ|  return;
#XZ|}
#XZ|
#XZ|// ⚠️ WARN - Explanation is OK but should be tracked
#XZ|// TODO(username): Research Stripe vs PayPal for EU compliance
#XZ|// See: https://docs.company.com/payment-options
#XZ|
#XZ|// ✅ OK - Completed TODO with explanation
#XZ|// TODO: Migrated to Stripe API v3 (2026-01-15)
#XZ|const stripe = new Stripe(config);
#```

#RB|

#HT|### Fake Test Detection
#XZ|```typescript
#XZ|// ❌ BLOCK - Fake test that always passes
#XZ|it('should work', () => {
#XZ|  expect(true).toBe(true);
#XZ|});
#XZ|
#XZ|// ❌ BLOCK - Empty test
#XZ|it('handles error', () => {
#XZ|  // TODO: add tests
#XZ|});
#XZ|
#XZ|// ❌ BLOCK - Skipped test without reason
#XZ|it.skip('important test', () => { ... });
#XZ|
#XZ|// ✅ OK - Real assertion
#XZ|it('validates email', () => {
#XZ|  expect(validateEmail('test@test.com')).toBe(true);
#XZ|  expect(validateEmail('invalid')).toBe(false);
#XZ|});
#```

#RB|

#HT|### Empty Function Detection
#XZ|```typescript
#XZ|// ❌ BLOCK - Empty function
#XZ|function calculateDiscount(price: number): number {
#XZ|  // TODO: implement
#XZ|}
#XZ|
#XZ|// ❌ BLOCK - Only logs without action
#XZ|async function sendEmail(to: string, subject: string) {
#XZ|  console.log('Would send email to', to);
#XZ|  // TODO: implement
#XZ|}
#XZ|
#XZ|// ✅ OK - Real implementation
#XZ|function calculateDiscount(price: number): number {
#XZ|  if (price > 100) return price * 0.1;
#XZ|  return 0;
#XZ|}
#```

#RB|

#HT|### Secret Detection
#XZ|```typescript
#XZ|// ❌ BLOCK - Hardcoded secrets
#XZ|const apiKey = 'sk_live_1234567890abcdef';
#XZ|const password = 'admin123';
#XZ|
#XZ|// ✅ OK - Environment variables
#XZ|const apiKey = process.env.STRIPE_API_KEY;
#XZ|```

#RB|

#BP|## Integration with Orchestrator

#NT|

#HT|When invoked via `sdd-orchestrator`:
#HV|1. Orchestrator detects completion phrase or phase transition
#HV|2. Vibe Guard runs with current mode (vibe/standard/strict)
#HV|3. Results determine if transition is allowed:
#HV|   - `pass` or `pass_with_warning` → allow transition
#HV|   - `block` → return to previous state with feedback

#RB|

#BP|## Quick Reference

#NT|

#HT|| Category | vibe | standard | strict |
#HT||-----------|------|----------|--------|
#HT|| TODO/FIXME | ⚠️ | ⚠️ | ❌ |
#HT|| Empty functions | ⚠️ | ⚠️ | ❌ |
#HT|| Fake data | ⚠️ | ❌ | ❌ |
#HT|| Tests run | ○ | ✓ | ✓ |
#HT|| Fake tests | ⚠️ | ⚠️ | ❌ |
#HT|| Coverage | - | >60% | >80% |
#HT|| Build | ❌ | ❌ | ❌ |
#HT|| Type check | - | ❌ | ❌ |
#HT|| Lint | - | ⚠️ | ❌ |
#HT|| Secrets | ❌ | ❌ | ❌ |
#HT|| SQL injection | ⚠️ | ❌ | ❌ |
#HT|| XSS | ⚠️ | ❌ | ❌ |
#HT|| Empty catch | ⚠️ | ⚠️ | ❌ |

#RB|

#BP|## Examples

#NT|

#HT|### Example 1: Vibe Mode (Quick Check)
#XZ|```
#XZ|> User: "Create a simple todo app"
#XZ|> AI: (writes code, claims done)
#XZ|> Vibe Guard: scans, finds TODO but only warns
#XZ|> Result: pass_with_warning, allows continuation
#XZ|```

#HT|### Example 2: Standard Mode (Team Project)
#XZ|```
#XZ|> User: "Feature complete, ready to deploy"
#XZ|> AI: (implements feature, runs tests)
#XZ|> Vibe Guard: 
#XZ|>   - Finds TODO → warning
#XZ|>   - Tests ran and passed → pass
#XZ|>   - Build success → pass
#XZ|>   - Type check success → pass
#XZ|>   - No secrets → pass
#XZ|> Result: pass_with_warning (TODO should be addressed before prod)
#XZ|```

#HT|### Example 3: Strict Mode (Enterprise)
#XZ|```
#XZ|> AI: (claims feature complete)
#XZ|> Vibe Guard:
#XZ|>   - Finds TODO → BLOCK
#XZ|>   - Coverage 50% < 80% → BLOCK
#XZ|>   - Lint warning → BLOCK
#XZ|> Result: BLOCK, returns with blocking_issues list
#XZ|```

#BP|## Related Skills

#NT|

#HV|- `sdd-orchestrator` - Can invoke vibe-guard automatically
#HV|- `spec-driven-test` - Tests vibe-guard verifies
#HV|- `sdd-release-guard` - Final gate before production

#RB|

#XH|## Machine Contracts

#BN|Report structure conforms to `skills/sdd-orchestrator/sdd-machine-schema.json`.
