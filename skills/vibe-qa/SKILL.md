---
name: vibe-qa
description: Use when need to run end-to-end tests, browser automation, visual regression testing, or verify UI functionality.
---

# Vibe QA - End-to-End Testing

**Test the whole user journey, not just individual components.**

---

## When to Use

**Use when:**
- Running end-to-end tests
- Verifying user flows work
- Browser automation tasks
- Visual regression testing
- Smoke testing before release

**NOT use when:**
- Unit testing (use `spec-driven-test`)
- Code review (use `vibe-debug`)
- Design discussions (use `vibe-design`)

---

## QA Modes

| Mode | Purpose | Time |
|------|---------|------|
| `quick` | Smoke test critical paths | ~30s |
| `diff-aware` | Test changed features | 5-10 min |
| `full` | Complete application test | 5-15 min |
| `regression` | Compare against baseline | Varies |

---

## Test Types

```
┌─────────────────────────────────────────────────────────┐
│                    QA Pyramid                             │
├─────────────────────────────────────────────────────────┤
│                                                          │
│                     ┌─────┐                              │
│                    │ E2E │    ← Vibe QA focus           │
│                   ┌┴─────┴┐                             │
│                  │Integr. │                              │
│                 ┌┴───────┴┐                             │
│                │  Unit   │                              │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 1. Unit Tests
- Test individual functions
- Fast, isolated
- High coverage
- Managed by `spec-driven-test`

### 2. Integration Tests
- Test component interactions
- Moderate speed
- Critical paths covered

### 3. E2E Tests (Vibe QA)
- Test complete user flows
- Real browser automation
- Playwright-based
- Slowest but most realistic

---

## Browser Automation (Playwright)

### Core Capabilities

| Capability | Command | Purpose |
|------------|---------|---------|
| Navigate | `page.goto(url)` | Open URL |
| Click | `page.click(selector)` | Interact |
| Type | `page.fill(selector, text)` | Input text |
| Snapshot | `page.snapshot()` | Get page structure |
| Screenshot | `page.screenshot()` | Visual capture |

### Example E2E Test

```javascript
// test/user-flow.spec.ts
import { test, expect } from '@playwright/test';

test('user can login and view dashboard', async ({ page }) => {
  // Navigate
  await page.goto('/login');
  
  // Fill form
  await page.fill('[name="email"]', 'user@example.com');
  await page.fill('[name="password"]', 'password123');
  
  // Submit
  await page.click('button[type="submit"]');
  
  // Verify
  await expect(page).toHaveURL('/dashboard');
  await expect(page.locator('h1')).toContainText('Dashboard');
});
```

---

## Running Tests

### Quick Mode (Smoke Test)

```bash
# Test only critical paths
vic qa --mode quick
```

### Diff-Aware Mode

```bash
# Test only changed features
vic qa --mode diff-aware
```

### Full Mode

```bash
# Complete application test
vic qa --mode full
```

### Regression Mode

```bash
# Compare against baseline
vic qa --mode regression
```

---

## Visual Regression

### Before/After Comparison

```bash
# Capture baseline
vic qa screenshot --name login-page

# After changes
vic qa screenshot --name login-page --compare
```

### Diff Detection

```bash
# Run visual diff
vic qa visual-diff --baseline ./baseline --current ./current
```

---

## CI/CD Integration

### GitHub Actions

```yaml
name: E2E Tests
on: [push, pull_request]

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run QA
        run: vic qa --mode quick
```

---

## QA Commands

| Command | Purpose |
|---------|---------|
| `vic qa init` | Initialize QA setup |
| `vic qa test` | Run E2E tests |
| `vic qa screenshot` | Capture screenshot |
| `vic qa visual-diff` | Compare visuals |
| `vic qa report` | Generate QA report |

---

## Test Report

```bash
$ vic qa report

═══════════════════════════════════════════════════════════
  QA Report
═══════════════════════════════════════════════════════════

  Mode: full
  Duration: 8m 32s
  
  Results:
    Passed:  47
    Failed:   2
    Skipped:  3
    
  Critical Paths:
    ✅ Login flow
    ✅ Checkout process
    ✅ User profile
    ❌ Payment method
    ❌ Order history
    
  Screenshots: 12 captured
  
═══════════════════════════════════════════════════════════
```

---

## Integration with VIBE-SDD

```
实现完成 → Gate 3 (测试覆盖)
    ↓
vibe-qa → Gate 4 (E2E 验证)
    ↓
生成回归测试
    ↓
vic check → 准备发布
```

---

## Playwright Setup

### Installation

```bash
# Install Playwright
npm install -D @playwright/test

# Install browsers
npx playwright install
```

### Configuration

```javascript
// playwright.config.ts
import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  use: {
    baseURL: 'http://localhost:3000',
    headless: true,
  },
  projects: [
    { name: 'chromium', use: { browserName: 'chromium' } },
    { name: 'firefox', use: { browserName: 'firefox' } },
    { name: 'webkit', use: { browserName: 'webkit' } },
  ],
});
```

---

## Common QA Issues

| Issue | Solution |
|-------|----------|
| Flaky tests | Add proper waits, retry logic |
| Slow tests | Parallel execution, test isolation |
| Missing coverage | Add critical paths |
| Visual noise | Use viewport consistency |

---

## Quick Checklist

Before declaring QA complete:
- [ ] All critical paths tested
- [ ] No failing tests
- [ ] Screenshots captured
- [ ] Report generated
- [ ] Regression tests added

---

## Related Skills

| Skill | Relationship |
|-------|--------------|
| `spec-driven-test` | Unit/integration tests + TDD |
| `vibe-debug` | Debug failed tests |
| `vibe-design` | Visual consistency |
| `sdd-release-guard` | Final QA gate |

---

**Remember: Trust, but verify. Test the complete user journey.**
