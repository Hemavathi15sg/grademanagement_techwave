# 🤖 Self-Healing Flaky Tests - Quick Start Guide

## What This Does

The self-healing workflow automatically:
1. ✅ **Detects** flaky tests (tests that sometimes pass, sometimes fail)
2. 🔧 **Fixes** them using common patterns
3. ✅ **Verifies** the fixes work (10 consecutive runs)
4. 📝 **Creates PR** automatically if tests pass
5. 🚨 **Creates issue** if auto-fix fails

## Files Created

### 1. Main Workflow
**File:** `.github/workflows/self-healing-flaky-tests.yml`
- Runs on: Push, PR, Daily schedule, Manual trigger
- Detects flaky tests by running 5 times
- Auto-applies fixes
- Creates PR with fixes

### 2. Fix Script
**File:** `scripts/fix-flaky-tests.sh`
- Bash script for applying fixes
- Detects common patterns:
  - Race conditions (missing WaitGroup)
  - Timing issues (arbitrary sleeps)
  - Missing cleanup
  - Poor test isolation

### 3. Documentation
**File:** `.github/workflows/SELF_HEALING_README.md`
- Complete reference guide
- Fix patterns
- Best practices
- Troubleshooting

## How to Use

### Option 1: Manual Trigger
```bash
# Go to GitHub Actions
# Select "Self-Healing Flaky Tests"
# Click "Run workflow"
```

### Option 2: Let It Run Automatically
- Workflow runs daily at 2 AM UTC
- Runs on every push to feature branches
- Runs on PRs to main

### Option 3: Test Locally First
```bash
# Make script executable
chmod +x scripts/fix-flaky-tests.sh

# Run on specific test
./scripts/fix-flaky-tests.sh "TestEnrollmentConcurrency_Flaky"

# Review changes
git diff

# Test fixes
go test -race -count=10 ./tests
```

## What Happens When Flaky Test Detected

### Success Path ✅
```
1. Detect flaky test
   ↓
2. Create branch: auto-fix/flaky-tests-TIMESTAMP
   ↓
3. Apply fixes (WaitGroup, cleanup, etc.)
   ↓
4. Run test 10 times → All pass
   ↓
5. Run full test suite → Pass
   ↓
6. Create PR automatically
   ↓
7. Assign to you for review
   ↓
8. You review & merge
```

### Failure Path 🚨
```
1. Detect flaky test
   ↓
2. Create branch & apply fixes
   ↓
3. Run test 10 times → Some fail
   ↓
4. Create GitHub Issue
   ↓
5. Issue assigned to you
   ↓
6. Manual fix required
```

## Example PR Created

**Title:** 🤖 Auto-fix: Resolve flaky tests detected in CI

**Body:**
- List of flaky tests
- Fixes applied
- Verification results (10/10 passed)
- Review checklist
- Detailed analysis

**Labels:**
- `automated-fix`
- `flaky-tests`
- `priority-high`
- `testing`

**Assignee:** You (the committer)

## Reviewing Auto-Generated PRs

### ✅ What to Check

1. **Review Code Changes**
   ```bash
   git checkout auto-fix/flaky-tests-TIMESTAMP
   git diff main
   ```

2. **Run Tests Locally**
   ```bash
   # With race detector
   go test -race -count=20 ./tests
   
   # Specific test
   go test -race -count=50 -run TestEnrollmentConcurrency_Flaky ./tests
   ```

3. **Check for Side Effects**
   - Does the fix introduce new issues?
   - Are there performance implications?
   - Is the fix addressing root cause?

4. **Verify Full Suite**
   ```bash
   go test ./... -cover
   ```

### ✅ Merge Checklist

- [ ] Code changes make sense
- [ ] No race conditions (`-race` flag)
- [ ] Tests pass consistently (20+ runs)
- [ ] No regressions in other tests
- [ ] Fix addresses root cause (not just symptom)

## Common Fixes Applied

### 1. Race Condition Fix
**Problem:** Goroutines without synchronization
```go
// Before (Flaky)
go updateCounter()
time.Sleep(100 * time.Millisecond)
assert.Equal(t, 1, counter)
```

**Solution:** Add WaitGroup
```go
// After (Stable)
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    updateCounter()
}()
wg.Wait()
assert.Equal(t, 1, counter)
```

### 2. Timing Fix
**Problem:** Arbitrary sleep durations
```go
// Before (Flaky)
time.Sleep(50 * time.Millisecond)
assert.True(t, isReady)
```

**Solution:** Condition-based wait
```go
// After (Stable)
waitForCondition(t, func() bool {
    return isReady
}, 5*time.Second, "waiting for ready state")
```

### 3. Cleanup Fix
**Problem:** No cleanup between tests
```go
// Before (Flaky)
func TestA(t *testing.T) {
    globalState = "test"
    // test logic
}
```

**Solution:** Proper cleanup
```go
// After (Stable)
func TestA(t *testing.T) {
    old := globalState
    t.Cleanup(func() {
        globalState = old
    })
    globalState = "test"
    // test logic
}
```

## Monitoring Results

### GitHub Actions
- Check workflow runs: Actions → Self-Healing Flaky Tests
- View test results artifacts
- Check PR/Issue creation

### Metrics to Track
- Number of flaky tests detected
- Auto-fix success rate
- Time to fix
- PR merge rate

## Customization

### Adjust Detection Sensitivity
In workflow file, change:
```yaml
# Run tests 5 times → Change to 10 for more accurate detection
for i in {1..5}; do
```

### Add More Fix Patterns
Edit `scripts/fix-flaky-tests.sh`:
```bash
# Add custom fix function
fix_custom_pattern() {
    local file=$1
    # Your custom fix logic
}
```

### Change Schedule
```yaml
schedule:
  - cron: '0 2 * * *'  # Daily 2 AM UTC
  # Change to: '0 */6 * * *' for every 6 hours
```

## Troubleshooting

### Workflow Not Running
- Check branch name matches trigger pattern
- Verify workflow file syntax (use GitHub Actions validator)
- Check repository permissions

### Auto-Fix Not Working
- Review workflow logs for errors
- Check if test file was found
- Verify Go environment setup
- Test script locally first

### PRs Not Created
- Check `GITHUB_TOKEN` permissions
- Verify PR creation step logs
- Check if verification passed
- May need `peter-evans/create-pull-request@v6` action permissions

## Integration with Other Tools

### Grafana/Prometheus
Track flaky test metrics:
```promql
# Flaky test detection rate
rate(flaky_tests_detected_total[1d])

# Auto-fix success rate
rate(auto_fix_success_total[1d]) / rate(auto_fix_attempts_total[1d])
```

### Jira Integration
Workflow can create Jira tickets:
```yaml
- name: Create Jira Issue
  if: failure()
  uses: atlassian/gajira-create@v3
  with:
    project: TECH
    issuetype: Bug
    summary: "Flaky test auto-fix failed: ${{ needs.detect-flaky-tests.outputs.flaky_tests }}"
```

## Best Practices

1. **Review Before Merge**: Always review auto-generated fixes
2. **Test Thoroughly**: Run with `-race` and `-count=20`
3. **Root Cause**: Ensure fix addresses root cause
4. **Monitor CI**: Watch for recurring patterns
5. **Document**: Add comments explaining why fix works

## Support

For issues or questions:
1. Check workflow logs in GitHub Actions
2. Review this guide and main README
3. Check `.github/instructions/flaky-test-fix.instructions.md`
4. Create issue with `flaky-tests` label

---

🤖 **Generated by**: Self-Healing Flaky Tests Workflow
📅 **Last Updated**: January 8, 2026
