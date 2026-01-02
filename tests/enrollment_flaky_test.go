package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEnrollmentConcurrency_Flaky demonstrates a flaky test with race condition
// This test intentionally has timing issues to demonstrate auto-fix capability
//
// FIXED: Replaced arbitrary timeout with sync.WaitGroup for proper synchronization
func TestEnrollmentConcurrency_Flaky(t *testing.T) {
	// Simulate concurrent enrollment operations
	results := make(chan bool, 5)
	var wg sync.WaitGroup

	// ✅ STABLE: Use WaitGroup to track goroutine completion
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(id int) {
			defer wg.Done() // ✅ STABLE: Ensure WaitGroup is decremented

			// ✅ STABLE: Simulate variable timing (no longer causes flakiness)
			time.Sleep(time.Duration(id*10) * time.Millisecond)

			// Simulate enrollment operation
			enrollment := createTestEnrollment(id)

			// ✅ STABLE: Send result to buffered channel
			if enrollment.ID > 0 {
				results <- true
			} else {
				results <- false
			}
		}(i)
	}

	// ✅ STABLE: Wait for all goroutines to complete
	wg.Wait()
	close(results) // ✅ STABLE: Close channel after all senders are done

	// ✅ STABLE: Collect all results without timeout racing
	successCount := 0
	for success := range results {
		if success {
			successCount++
		}
	}

	// ✅ STABLE: Assertion now always has all 5 results
	assert.Equal(t, 5, successCount, "Expected all 5 enrollments to succeed")
}

func createTestEnrollment(id int) struct{ ID int } {
	// Simulate database operation with variable timing
	time.Sleep(time.Duration(id*15) * time.Millisecond)
	return struct{ ID int }{ID: id + 1}
}
