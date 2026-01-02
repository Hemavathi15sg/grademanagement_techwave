---
applyTo: '**'
---
# Flaky Test Remediation Instructions

When analyzing flaky tests:

## Analysis Steps
1. Review test code for common flaky patterns:
   - Race conditions (goroutines without synchronization)
   - Timing dependencies (sleep, timeouts)
   - Shared state without locks
   - Network/external dependencies
   - Random data without seeds

2. Check test logs for patterns:
   - "timeout" messages
   - Variable execution times
   - "race detected" warnings
   - Inconsistent assertion failures

3. Identify root cause category:
   - **Timing:** Insufficient waits, arbitrary timeouts
   - **Concurrency:** Missing sync primitives (WaitGroup, Mutex)
   - **State:** Shared mutable state, test pollution
   - **External:** Network calls, file system, random data

## Fix Patterns

### For Timing Issues
- Replace fixed timeouts with retry loops + max duration
- Use proper wait conditions (e.g., wait for specific state)
- Add exponential backoff for retries

### For Concurrency Issues
- Use sync.WaitGroup for goroutine completion
- Use sync.Mutex for shared state
- Use channels with proper buffering and closing
- Avoid sleeping for synchronization

### For State Issues
- Isolate test data (unique IDs, separate databases)
- Clean up after tests (defer cleanup functions)
- Use test fixtures, not shared globals

## Fix Template
```go
// BEFORE: Flaky test
func TestFlaky(t *testing.T) {
    go doSomething()
    time.Sleep(50 * time.Millisecond) // ❌ FLAKY: arbitrary sleep
    assert.True(t, done)
}

// AFTER: Stable test
func TestStable(t *testing.T) {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        doSomething()
    }()
    wg.Wait() // ✅ STABLE: explicit synchronization
    assert.True(t, done)
}
```

## Validation
- Run test 10+ times to verify stability
- Use `go test -race` to detect race conditions
- Check execution time variance (should be <10%)