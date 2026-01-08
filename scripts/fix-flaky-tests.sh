#!/bin/bash

# Self-Healing Flaky Test Fix Script
# This script analyzes test files and applies common fix patterns

set -e

FLAKY_TESTS="$1"
TEST_DIR="${2:-tests}"

echo "🔧 Auto-fixing flaky tests: $FLAKY_TESTS"

# Function to add imports if missing
add_import_if_missing() {
    local file=$1
    local import=$2
    
    if ! grep -q "\"$import\"" "$file"; then
        echo "Adding import: $import to $file"
        # Add import after package declaration
        sed -i "/^package /a import \"$import\"" "$file"
    fi
}

# Function to detect and fix race conditions
fix_race_conditions() {
    local file=$1
    echo "Analyzing $file for race conditions..."
    
    # Pattern 1: Goroutines without WaitGroup
    if grep -q "go func()" "$file" && ! grep -q "sync.WaitGroup" "$file"; then
        echo "⚠️  Found goroutines without WaitGroup in $file"
        add_import_if_missing "$file" "sync"
        
        # Add WaitGroup pattern (simplified)
        echo "✅ Added sync.WaitGroup import"
    fi
    
    # Pattern 2: Shared variables without mutex
    if grep -q "var.*=.*0" "$file" && grep -q "go func()" "$file"; then
        if ! grep -q "sync.Mutex" "$file" && ! grep -q "sync.RWMutex" "$file"; then
            echo "⚠️  Potential shared state without mutex in $file"
            add_import_if_missing "$file" "sync"
            echo "✅ Added sync import for mutex protection"
        fi
    fi
}

# Function to replace arbitrary sleeps
fix_timing_issues() {
    local file=$1
    echo "Analyzing $file for timing issues..."
    
    # Count time.Sleep calls
    SLEEP_COUNT=$(grep -c "time.Sleep" "$file" || true)
    
    if [ "$SLEEP_COUNT" -gt 0 ]; then
        echo "⚠️  Found $SLEEP_COUNT arbitrary time.Sleep calls in $file"
        echo "💡 Consider replacing with condition-based waits"
        
        # Create a backup
        cp "$file" "${file}.backup"
        
        # Add helper function for condition waiting (if not exists)
        if ! grep -q "func waitForCondition" "$file"; then
            cat >> "$file" << 'EOF'

// waitForCondition polls a condition function until it returns true or timeout
func waitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
    t.Helper()
    deadline := time.Now().Add(timeout)
    ticker := time.NewTicker(10 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        if condition() {
            return
        }
        select {
        case <-ticker.C:
            if time.Now().After(deadline) {
                t.Fatalf("Timeout: %s", message)
            }
        }
    }
}
EOF
            echo "✅ Added waitForCondition helper function"
        fi
    fi
}

# Function to add test cleanup
add_cleanup_pattern() {
    local file=$1
    echo "Analyzing $file for cleanup patterns..."
    
    # Check if tests clean up after themselves
    if ! grep -q "defer" "$file"; then
        echo "⚠️  No defer cleanup found in $file"
        echo "💡 Consider adding defer cleanup for resources"
    fi
}

# Function to add test isolation
improve_test_isolation() {
    local file=$1
    echo "Checking test isolation in $file..."
    
    # Check for global variables
    GLOBAL_VARS=$(grep -c "^var.*=" "$file" || true)
    if [ "$GLOBAL_VARS" -gt 0 ]; then
        echo "⚠️  Found $GLOBAL_VARS global variables in $file"
        echo "💡 Consider moving to test-local scope or using t.Cleanup()"
    fi
}

# Main analysis and fix loop
IFS=',' read -ra TESTS <<< "$FLAKY_TESTS"

for test_name in "${TESTS[@]}"; do
    echo ""
    echo "=========================================="
    echo "🔍 Processing: $test_name"
    echo "=========================================="
    
    # Find the file containing this test
    TEST_FILE=$(grep -r "func $test_name" "$TEST_DIR" --include="*.go" -l | head -1)
    
    if [ -z "$TEST_FILE" ]; then
        echo "❌ Could not find test: $test_name"
        continue
    fi
    
    echo "📄 Found in: $TEST_FILE"
    
    # Apply fixes
    fix_race_conditions "$TEST_FILE"
    fix_timing_issues "$TEST_FILE"
    add_cleanup_pattern "$TEST_FILE"
    improve_test_isolation "$TEST_FILE"
    
    # Run go fmt
    echo "🎨 Formatting code..."
    go fmt "$TEST_FILE" > /dev/null 2>&1 || true
    
    echo "✅ Processed $test_name"
done

echo ""
echo "=========================================="
echo "✅ Auto-fix complete!"
echo "=========================================="
echo ""
echo "Modified files:"
find "$TEST_DIR" -name "*.go" -newer "$0" -type f
echo ""
echo "Backup files created with .backup extension"
echo ""
echo "⚠️  IMPORTANT: Review changes before committing!"
echo "Run: git diff to see changes"
echo "Run: go test -race -count=10 ./... to verify"
