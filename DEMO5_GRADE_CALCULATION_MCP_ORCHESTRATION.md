# Demo 5: Grade Calculation Engine - Cross-Tool MCP Orchestration
## Complete SDLC Automation with GitHub, Grafana, Figma, and Jira

**Duration:** 16-18 minutes  
**Objective:** Demonstrate end-to-end SDLC automation by orchestrating multiple MCP servers (GitHub, Grafana, Figma, Jira) in a unified workflow.

---

## 🎯 Demo Overview

### What's Different in Demo 5?

**Demos 1-4 Showed:**
- Individual tool capabilities (Copilot, Jira MCP, BDD)
- Single workflow per demo
- Linear progression

**Demo 5 Shows:**
- **Orchestration of 4 MCP servers working together**
- **Cross-tool automation** (GitHub → Grafana → Figma → Jira)
- **Single source of truth** approach
- **Real SDLC workflow** from issue to production metrics

### New Feature: Grade Calculation Engine

**Business Value:** Automate final grade calculations for courses based on weighted assignments, attendance bonuses, and grade curves.

**Technical Showcase:**
- GitHub MCP: Automated branch/PR creation from issue assignments
- Grafana MCP: Performance baseline queries for test threshold generation
- Figma MCP: Design token extraction for grade display validation
- Jira MCP: Central orchestration hub collecting data from all sources

---

## 📋 Pre-Demo Checklist

### Required MCP Extensions (Install Before Demo)

#### 1. GitHub MCP
```bash
# Check if installed
code --list-extensions | grep -i github

# If not installed
code --install-extension GitHub.copilot-mcp
```

**Test Connection:**
```
@github get repository info
```

#### 2. Grafana MCP (if available, otherwise simulate)
```bash
# Install Grafana MCP extension
code --install-extension grafana.grafana-mcp
```

**Test Connection:**
```
@grafana list dashboards
```

**Alternative (if MCP not available):** Show manual Grafana query, then explain how MCP would automate it.

#### 3. Figma MCP (Already have this)
**Test Connection:**
```
@figma get current file
```

#### 4. Jira MCP (Already configured from Demo 4)
**Test Connection:**
```
@atlassian get cloudId
```

---

### Setup Steps (Complete Before Demo)

#### 1. Create Figma Design Mockup

**Option A:** Use existing Figma file and add grade display component

**Option B:** Create simple text mockup showing design tokens:
```
GRADE DISPLAY DESIGN TOKENS
─────────────────────────────
Letter Grades:
- A+  (97-100)  Color: #22C55E (Green)
- A   (93-96)   Color: #22C55E
- A-  (90-92)   Color: #86EFAC
- B+  (87-89)   Color: #3B82F6 (Blue)
- B   (83-86)   Color: #3B82F6
- B-  (80-82)   Color: #93C5FD
- C+  (77-79)   Color: #F59E0B (Orange)
- C   (73-76)   Color: #F59E0B
- C-  (70-72)   Color: #FCD34D
- D   (60-69)   Color: #EF4444 (Red)
- F   (0-59)    Color: #DC2626 (Dark Red)

Font: Inter, 16px, Semi-bold
Display: Badge with rounded corners (8px radius)
```

#### 2. Setup Grafana Dashboard (or Simulation)

**Option A - Real Grafana:** Create dashboard with API response time metrics

**Option B - Simulation:** Prepare JSON response simulating Grafana query:
```json
{
  "dashboard": "grade-api-performance",
  "metrics": {
    "avg_response_time": 45,
    "p95_response_time": 120,
    "p99_response_time": 250,
    "requests_per_second": 150
  },
  "baseline_period": "last_30_days"
}
```

#### 3. Create GitHub Issue Template

Create issue TEC-15 manually or via API:
```markdown
Title: Grade Calculation Engine

Description:
Implement automated grade calculation system for course final grades.

Requirements:
- Calculate weighted average from assignments, exams, attendance
- Support grade curves (add percentage to all grades)
- Convert numeric grades to letter grades (A+, A, B+, etc.)
- API endpoint: POST /api/grades/calculate
- Performance: Response time < 200ms for 100 students

Acceptance Criteria:
1. Input: Student ID, Course ID, weighting rules
2. Output: Numeric grade, letter grade, percentile rank
3. Support grade curves (linear addition)
4. Validate against Figma design tokens for letter grades
5. Performance tests using Grafana baseline thresholds
6. All tests passing
```

#### 4. Update Custom Instructions

Add to `.github/instructions/copilot.instructions.md`:
```markdown
## Grade Calculation Standards
- Grade calculations use float64 for precision
- Letter grade constants: GradeAPlus, GradeA, GradeAMinus, etc.
- Grade thresholds configurable via struct
- Performance requirement: < 200ms response time
- Validate letter grade colors match Figma design tokens
```

#### 5. Prepare Branch
```powershell
git checkout -b demo/grade-calculation-orchestration
```

---

## 🎬 Demo Script

---

### **INTRODUCTION (1 min)**

**NARRATION:**
> "In the previous four demos, we showed individual tools and capabilities:
> - Demo 1: Development delegation with Coding Agent
> - Demo 2: Test automation with mocks and factories
> - Demo 3: Security automation with GHAS
> - Demo 4: Product owner workflow with Jira MCP and BDD
> 
> Today, Demo 5 is different. We're orchestrating **four MCP servers** to automate the complete SDLC: GitHub for workflow automation, Grafana for performance baselines, Figma for design validation, and Jira as the central hub.
>
> The feature: Grade Calculation Engine - automatically calculate final course grades with curves and convert to letter grades."

---

### **PHASE 1: GitHub MCP - Automated Branch Creation (2 min)**

**NARRATION:**
> "In traditional workflows, developer sees Jira assignment, manually creates branch, sets up environment. Let's automate it using GitHub MCP."

#### Step 1.1: Assign Issue via Jira

**ACTION:** In Jira, assign issue TEC-15 to yourself (or show it's assigned)

**TALKING POINTS:**
- "Product owner assigned me this story in Jira"
- "Normally I'd copy the issue key, switch to terminal, create branch manually"
- "With GitHub MCP, one command does everything"

#### Step 1.2: Auto-Create Branch from Issue

**ACTION:** In Copilot Chat:
```
@github create a new branch for Jira issue TEC-15:

Branch name: feature/grade-calculation-engine
Base: main
Description: Automated grade calculation with curves and letter grades

After creating the branch, checkout locally.
```

**EXPECTED OUTPUT:**
```
✓ Connecting to GitHub repository: Hemavathi15sg/grademanagement_techwave
✓ Creating branch: feature/grade-calculation-engine from main
✓ Branch created successfully
✓ Checking out branch locally...

Branch: feature/grade-calculation-engine
Ready for development.

Would you like me to:
- Create a draft PR now? (recommended for early visibility)
- Set up branch protection rules?
- Link this branch to TEC-15 in Jira?
```

**ACTION:** Respond:
```
Yes, create a draft PR and link to TEC-15
```

**EXPECTED OUTPUT:**
```
✓ Creating draft pull request...
✓ PR #10 created: Grade Calculation Engine
✓ Linked to Jira issue: TEC-15
✓ Status: Draft (not ready for review)

PR: https://github.com/Hemavathi15sg/grademanagement_techwave/pull/10
Jira: https://ecanarys-team-y31whl7q.atlassian.net/browse/TEC-15
```

**TALKING POINTS:**
- "GitHub MCP created branch, checked it out, created draft PR"
- "PR automatically linked to Jira issue"
- "Developer gets straight to coding, zero manual setup"

**CLIENT VALUE:**
> "Developer onboarding task: 10 minutes → 30 seconds. New developer gets assigned issue, one command sets up entire environment. No wiki pages to read about branch naming conventions."

---

### **PHASE 2: Figma MCP - Extract Design Tokens (3 min)**

**NARRATION:**
> "Designers created mockups for how grades should display. Instead of manually reading specs and hoping we implement correctly, let's extract design tokens programmatically."

#### Step 2.1: Show Figma Design

**ACTION:** Open Figma mockup (or show text version)

**TALKING POINTS:**
- Point out grade display badges: A+, A, B+, etc.
- Point out color palette: Green for A's, Blue for B's, Orange for C's, Red for D/F
- Point out font specifications
- "This is what the designer created. Now let's use Figma MCP to extract these as code."

#### Step 2.2: Extract Design Tokens via Figma MCP

**ACTION:** In Copilot Chat:
```
@figma extract design tokens from the Grade Display component:

Component name: GradeDisplayBadge
Extract:
- Letter grade values and numeric ranges (A+ = 97-100, etc.)
- Color codes for each grade category
- Font family, size, and weight
- Border radius and spacing values

Generate Go constants that I can use in validation tests.
```

**EXPECTED OUTPUT:**
```go
// Design tokens extracted from Figma: GradeDisplayBadge component
// File: design_tokens.go

package constants

// Grade letter values and thresholds
const (
    GradeAPlus  = "A+"  // 97-100
    GradeA      = "A"   // 93-96
    GradeAMinus = "A-"  // 90-92
    GradeBPlus  = "B+"  // 87-89
    GradeB      = "B"   // 83-86
    GradeBMinus = "B-"  // 80-82
    GradeCPlus  = "C+"  // 77-79
    GradeC      = "C"   // 73-76
    GradeCMinus = "C-"  // 70-72
    GradeD      = "D"   // 60-69
    GradeF      = "F"   // 0-59
)

// Grade color codes (hex)
const (
    ColorGradeA = "#22C55E" // Green for A's
    ColorGradeB = "#3B82F6" // Blue for B's
    ColorGradeC = "#F59E0B" // Orange for C's
    ColorGradeD = "#EF4444" // Red for D
    ColorGradeF = "#DC2626" // Dark Red for F
)

// Typography
const (
    FontFamily = "Inter"
    FontSize   = 16 // pixels
    FontWeight = "Semi-bold"
)

// Layout
const (
    BorderRadius = 8 // pixels
)

// Grade thresholds
var GradeThresholds = map[string]struct{ Min, Max float64 }{
    GradeAPlus:  {97.0, 100.0},
    GradeA:      {93.0, 96.99},
    GradeAMinus: {90.0, 92.99},
    GradeBPlus:  {87.0, 89.99},
    GradeB:      {83.0, 86.99},
    GradeBMinus: {80.0, 82.99},
    GradeCPlus:  {77.0, 79.99},
    GradeC:      {73.0, 76.99},
    GradeCMinus: {70.0, 72.99},
    GradeD:      {60.0, 69.99},
    GradeF:      {0.0, 59.99},
}
```

**TALKING POINTS:**
- "Figma MCP read the design file and generated Go constants"
- "Grade thresholds exactly match designer's specs"
- "Colors extracted as hex codes for validation"
- "No manual copy-paste, no translation errors"

#### Step 2.3: Generate Validation Tests from Design Tokens

**ACTION:** In Copilot Chat:
```
Using these design tokens, generate validation tests that verify:
1. Grade calculation logic returns correct letter grades
2. Letter grade thresholds match Figma specifications
3. Any grade display in API must use approved color codes

Create test file: tests/grade_design_validation_test.go
```

**EXPECTED OUTPUT:**
Test file created with tests like:
```go
func TestGradeConversion_MatchesFigmaTokens(t *testing.T) {
    tests := []struct {
        numericGrade float64
        expectedLetter string
    }{
        {98.5, constants.GradeAPlus},
        {94.0, constants.GradeA},
        {90.5, constants.GradeAMinus},
        // ... more cases
    }
    
    for _, tt := range tests {
        result := CalculateLetterGrade(tt.numericGrade)
        assert.Equal(t, tt.expectedLetter, result,
            "Grade %v should be %s per Figma design tokens",
            tt.numericGrade, tt.expectedLetter)
    }
}

func TestGradeColors_MatchFigmaTokens(t *testing.T) {
    // Verify any grade display uses approved colors
    gradeDisplay := GetGradeDisplay(95.0) // Should be A
    assert.Equal(t, constants.ColorGradeA, gradeDisplay.Color,
        "Grade A must use Figma-approved color")
}
```

**TALKING POINTS:**
- "Tests validate code matches design specifications"
- "Designer changes Figma → Tests update → Code must follow"
- "Design-dev consistency enforced automatically"

**CLIENT VALUE:**
> "Design handoff nightmares: gone. Designer updates Figma, Copilot extracts new tokens, tests fail if implementation drifts. Design and code stay in sync without meetings."

---

### **PHASE 3: Grafana MCP - Performance Baseline & Test Thresholds (3 min)**

**NARRATION:**
> "Product owner said: response time must be under 200ms. But how do we validate that? Let's query Grafana to understand current API performance and generate intelligent test thresholds."

#### Step 3.1: Query Grafana for Production Metrics

**ACTION:** In Copilot Chat:
```
@grafana query the grade-api-performance dashboard:

Metrics needed:
- Average API response time (last 30 days)
- 95th percentile response time
- 99th percentile response time
- Requests per second (peak load)

Time range: Last 30 days
Endpoint: /api/enrollments (similar feature for baseline)
```

**EXPECTED OUTPUT (if Grafana MCP working):**
```json
{
  "dashboard": "grade-api-performance",
  "query_time": "2026-01-02T10:30:00Z",
  "time_range": "last_30_days",
  "endpoint": "/api/enrollments",
  
  "metrics": {
    "avg_response_time_ms": 45,
    "p50_response_time_ms": 38,
    "p95_response_time_ms": 120,
    "p99_response_time_ms": 250,
    "max_response_time_ms": 450,
    "requests_per_second_avg": 150,
    "requests_per_second_peak": 420
  },
  
  "recommendation": {
    "new_endpoint_target": "< 200ms for p95",
    "performance_budget": "150ms avg, 200ms p95, 300ms p99"
  }
}
```

**ALTERNATIVE (if Grafana MCP unavailable):** Show pre-prepared JSON and explain:
> "In production, this query would be automated. Grafana MCP returns real metrics from our monitoring stack."

**TALKING POINTS:**
- "Grafana MCP queried production metrics for similar endpoint"
- "Average response: 45ms, p95: 120ms, p99: 250ms"
- "This gives us intelligent baseline for new feature"
- "We're not guessing performance targets anymore"

#### Step 3.2: Generate Performance Test with Grafana Baselines

**ACTION:** In Copilot Chat:
```
Based on these Grafana metrics, generate a performance test for the new grade calculation endpoint:

Test requirements:
- Simulate 100 concurrent requests
- Calculate grades for 100 students
- Target: Average < 150ms, p95 < 200ms (product owner requirement)
- Fail if average exceeds 150ms or p95 exceeds 200ms
- Use Grafana baselines to validate we're not regressing

Create: tests/grade_performance_test.go
```

**EXPECTED OUTPUT:**
```go
package tests

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
)

// Performance thresholds based on Grafana production metrics
// Baseline from /api/enrollments (similar complexity):
//   - Avg: 45ms, P95: 120ms, P99: 250ms
// Target for /api/grades/calculate:
//   - Avg: < 150ms (3x baseline, accounts for calculation overhead)
//   - P95: < 200ms (product owner requirement)
//   - P99: < 300ms (acceptable outlier threshold)

const (
    TargetAvgResponseTime = 150 * time.Millisecond
    TargetP95ResponseTime = 200 * time.Millisecond
    TargetP99ResponseTime = 300 * time.Millisecond
)

func TestGradeCalculation_PerformanceBaseline(t *testing.T) {
    // Setup: 100 student grade calculations
    students := GenerateTestStudents(100)
    
    // Measure response times
    var responseTimes []time.Duration
    
    for i := 0; i < 100; i++ {
        start := time.Now()
        _, err := CalculateGrade(students[i])
        elapsed := time.Since(start)
        
        assert.NoError(t, err)
        responseTimes = append(responseTimes, elapsed)
    }
    
    // Calculate percentiles
    avg := calculateAverage(responseTimes)
    p95 := calculatePercentile(responseTimes, 0.95)
    p99 := calculatePercentile(responseTimes, 0.99)
    
    // Validate against Grafana-derived thresholds
    assert.Less(t, avg, TargetAvgResponseTime,
        "Average response time %v exceeds target %v (Grafana baseline: 45ms)",
        avg, TargetAvgResponseTime)
    
    assert.Less(t, p95, TargetP95ResponseTime,
        "P95 response time %v exceeds target %v (Product requirement)",
        p95, TargetP95ResponseTime)
    
    assert.Less(t, p99, TargetP99ResponseTime,
        "P99 response time %v exceeds target %v (Acceptable outlier)",
        p99, TargetP99ResponseTime)
    
    // Log metrics for trend analysis
    t.Logf("Performance metrics:")
    t.Logf("  Average: %v (target: %v)", avg, TargetAvgResponseTime)
    t.Logf("  P95: %v (target: %v)", p95, TargetP95ResponseTime)
    t.Logf("  P99: %v (target: %v)", p99, TargetP99ResponseTime)
}

func TestGradeCalculation_ConcurrentLoad(t *testing.T) {
    // Simulate real-world concurrent load based on Grafana RPS data
    // Peak load: 420 requests/second
    // Test: 100 concurrent requests (simulates 10 seconds of peak load)
    
    concurrentRequests := 100
    results := make(chan time.Duration, concurrentRequests)
    
    for i := 0; i < concurrentRequests; i++ {
        go func() {
            student := GenerateTestStudent()
            start := time.Now()
            CalculateGrade(student)
            results <- time.Since(start)
        }()
    }
    
    // Collect results
    var responseTimes []time.Duration
    for i := 0; i < concurrentRequests; i++ {
        responseTimes = append(responseTimes, <-results)
    }
    
    p95 := calculatePercentile(responseTimes, 0.95)
    assert.Less(t, p95, TargetP95ResponseTime,
        "Under concurrent load, P95 %v exceeds target %v", p95, TargetP95ResponseTime)
}
```

**TALKING POINTS:**
- "Performance test uses real Grafana metrics as baseline"
- "Test thresholds: 150ms avg, 200ms p95 (not arbitrary numbers)"
- "Comments explain why these numbers: 3x baseline + product requirement"
- "Concurrent load test simulates peak traffic from Grafana data"

**CLIENT VALUE:**
> "Performance requirements used to be guesses: 'make it fast.' Now we query production metrics, set intelligent thresholds, validate with data. Performance regressions caught before production."

---

### **PHASE 4: Implement Feature with Copilot (3 min)**

**NARRATION:**
> "We have design tokens from Figma, performance thresholds from Grafana, tests ready. Now let's implement the feature."

#### Step 4.1: Generate Grade Calculation Logic

**ACTION:** In Copilot Chat:
```
@atlassian get issue TEC-15

Based on the acceptance criteria, design tokens from Figma, and performance requirements from Grafana, generate:

1. models/grade_calculation.go
   - GradeCalculation struct with numeric and letter grade
   - Weighted average calculation logic
   - Grade curve support (add percentage to all grades)
   - Convert numeric to letter using Figma design tokens

2. handlers/grade_calculation_handler.go
   - POST /api/grades/calculate
   - Input: StudentID, CourseID, Curve percentage (optional)
   - Output: Numeric grade, Letter grade, Percentile
   - Validate performance < 200ms per Grafana baseline

3. Use repository pattern following existing conventions

@github #github-pull-request_copilot-coding-agent implement this in background.
```

**EXPECTED OUTPUT:**
```
✓ Reading Jira acceptance criteria...
✓ Loading Figma design tokens...
✓ Loading Grafana performance baselines...
✓ Background Coding Agent started

Generating:
- models/grade_calculation.go
- handlers/grade_calculation_handler.go  
- repository/grade_repository.go
- tests/grade_calculation_test.go (using design validation)
- tests/grade_performance_test.go (using Grafana thresholds)

Estimated time: 6-8 minutes
PR will be updated at: https://github.com/.../pull/10
```

**TALKING POINTS:**
- "Agent is reading requirements from Jira"
- "Using design tokens from Figma for letter grade logic"
- "Using Grafana baselines for performance validation"
- "All three MCP sources feeding into code generation"

**NARRATION:**
> "While the agent works, let me show you what happens behind the scenes..."

#### Step 4.2: Show Cross-Tool Data Flow

**ACTION:** Draw or display diagram:
```
TEC-15 (Jira)
  ├─ Acceptance Criteria → Business logic requirements
  │
  ├─ Figma Design Tokens → Grade thresholds, colors, display rules
  │
  ├─ Grafana Metrics → Performance targets, test thresholds
  │
  └─→ Copilot Agent
       ├─ models/grade_calculation.go (uses Figma tokens)
       ├─ handlers/grade_calculation_handler.go (validates Grafana perf)
       ├─ tests/grade_design_validation_test.go (Figma validation)
       └─ tests/grade_performance_test.go (Grafana thresholds)
```

**TALKING POINTS:**
- "Single story (TEC-15) pulls data from three MCP sources"
- "Design tokens ensure UI consistency"
- "Performance baselines ensure production-ready code"
- "Jira orchestrates the entire workflow"

---

### **PHASE 5: Validate & Update Status (4 min)**

**NARRATION:**
> "Background agent should be done. Let's validate the implementation meets all requirements from all three tools."

#### Step 5.1: Review Generated Code

**ACTION:** Check PR #10, show key files:

**1. models/grade_calculation.go:**
```go
// Grade letter constants from Figma design tokens
const (
    GradeAPlus  = "A+"  // 97-100, Color: #22C55E
    GradeA      = "A"   // 93-96,  Color: #22C55E
    // ... (uses Figma tokens)
)

func ConvertToLetterGrade(numeric float64) string {
    // Uses Figma-extracted thresholds
    switch {
    case numeric >= 97.0:
        return GradeAPlus
    case numeric >= 93.0:
        return GradeA
    // ...
    }
}
```

**2. tests/grade_performance_test.go:**
```go
// Performance thresholds from Grafana production metrics
const (
    TargetAvgResponseTime = 150 * time.Millisecond // 3x baseline
    TargetP95ResponseTime = 200 * time.Millisecond // Product requirement
)
```

**TALKING POINTS:**
- "Code uses exact Figma design tokens for grade thresholds"
- "Performance test uses Grafana-derived targets"
- "All comments reference source of truth (Figma, Grafana, Jira)"

#### Step 5.2: Run All Tests

**ACTION:** In terminal:
```powershell
# Run design validation tests
go test ./tests/grade_design_validation_test.go -v

# Run performance tests
go test ./tests/grade_performance_test.go -v

# Run full suite
go test ./tests/grade_* -v -cover
```

**EXPECTED OUTPUT:**
```
=== RUN   TestGradeConversion_MatchesFigmaTokens
--- PASS: TestGradeConversion_MatchesFigmaTokens (0.00s)

=== RUN   TestGradeColors_MatchFigmaTokens
--- PASS: TestGradeColors_MatchFigmaTokens (0.00s)

=== RUN   TestGradeCalculation_PerformanceBaseline
    Performance metrics:
      Average: 42ms (target: 150ms) ✓
      P95: 105ms (target: 200ms) ✓
      P99: 210ms (target: 300ms) ✓
--- PASS: TestGradeCalculation_PerformanceBaseline (2.34s)

=== RUN   TestGradeCalculation_ConcurrentLoad
    Concurrent P95: 185ms (target: 200ms) ✓
--- PASS: TestGradeCalculation_ConcurrentLoad (1.52s)

PASS
coverage: 91.2% of statements
ok      techwave/tests  4.123s
```

**TALKING POINTS:**
- "Design validation: All letter grades match Figma tokens ✓"
- "Performance: Average 42ms, well under 150ms target ✓"
- "P95 response: 105ms, beats 200ms requirement ✓"
- "Concurrent load: 185ms, within threshold ✓"

#### Step 5.3: Cross-Tool Status Update

**ACTION:** In Copilot Chat:
```
Update all tools with completion status:

@github Update PR #10:
- Mark as "Ready for Review"
- Add test results summary
- Request review from team

@grafana Log this metric for trending:
- Endpoint: /api/grades/calculate
- Average response: 42ms
- P95 response: 105ms
- Tag: new-feature-baseline

@atlassian Update TEC-15:
- Status: Ready for QA
- Comment with cross-tool validation results:
  * Figma design validation: All tests passing
  * Grafana performance validation: Exceeds targets
  * GitHub: PR #10 ready for review
  * Coverage: 91.2%
- Set custom fields: Copilot Delegated=Yes, Time Saved=5 hours
```

**EXPECTED OUTPUT:**
```
✓ GitHub PR #10 updated:
  - Status: Ready for Review
  - Test results added to description
  - Review requested from @team

✓ Grafana metric logged:
  - Dashboard: grade-api-performance
  - New baseline recorded for /api/grades/calculate
  - Available for future comparisons

✓ Jira TEC-15 updated:
  - Status: To Do → Ready for QA
  - Comment added with validation from all sources
  - Custom fields updated

Cross-tool synchronization complete ✓
```

**TALKING POINTS:**
- "One command updated three tools"
- "GitHub knows PR is ready"
- "Grafana has new baseline for future features"
- "Jira tracks that design and performance validated"
- "Product owner sees complete picture without asking"

---

### **PHASE 6: Unified Dashboard View (2 min)**

**NARRATION:**
> "Product owner wants to see: Is TEC-15 done? Let's query all MCP sources for a unified status."

#### Step 6.1: Generate Cross-Tool Status Report

**ACTION:** In Copilot Chat:
```
Generate a unified status report for TEC-15 by querying all MCP sources:

@atlassian Get TEC-15 status and custom fields
@github Get PR #10 status and checks
@grafana Get latest performance metrics for /api/grades/calculate
@figma Verify design token file last updated date

Consolidate into single report showing:
- Overall status (In Progress / Ready for QA / Done)
- Code review status
- Test results
- Performance validation
- Design compliance
- Time saved metrics
```

**EXPECTED OUTPUT:**
```markdown
# Feature Status Report: TEC-15 Grade Calculation Engine
Generated: 2026-01-02 10:45:00

## Overall Status: ✅ READY FOR QA

---

## Code Status (GitHub)
📍 PR #10: Grade Calculation Engine
- Status: ✅ Ready for Review
- Branch: feature/grade-calculation-engine
- Files changed: 8 files, 1,234 lines added
- CI/CD: ✅ All checks passed
- Tests: ✅ 28 tests passing
- Coverage: 91.2%
- Review: Pending (requested 5 minutes ago)
- Link: https://github.com/Hemavathi15sg/.../pull/10

---

## Design Validation (Figma)
🎨 Design Tokens: grade-display-component
- Token extraction: ✅ Complete
- Last updated: 2026-01-02 09:30:00
- Letter grades validated: ✅ 11/11 passing
- Color compliance: ✅ All grades use approved colors
- Font specifications: ✅ Match design system
- Status: ✅ DESIGN COMPLIANT

Design validation test results:
  ✓ Grade thresholds match Figma specs (A+=97-100, A=93-96, etc.)
  ✓ Color codes validated (#22C55E for A, #3B82F6 for B, etc.)
  ✓ Display format matches typography requirements

---

## Performance Validation (Grafana)
📊 Metrics: /api/grades/calculate
- Baseline comparison: /api/enrollments (similar complexity)
- Test results:
  * Average response: 42ms ✅ (target: <150ms, baseline: 45ms)
  * P95 response: 105ms ✅ (target: <200ms, baseline: 120ms)
  * P99 response: 210ms ✅ (target: <300ms, baseline: 250ms)
  * Concurrent load (100 req): 185ms P95 ✅
- Status: ✅ EXCEEDS PERFORMANCE TARGETS
- Improvement: 7% faster than baseline feature

Performance trend:
  📈 New feature performs better than existing similar endpoints
  📉 No performance regression
  🎯 Product requirement (<200ms) exceeded by 48%

---

## Project Management (Jira)
📋 TEC-15: Grade Calculation Engine
- Status: Ready for QA
- Sprint: Sprint 1
- Story Points: 5
- Assignee: You
- Created: 2025-12-30
- Updated: 2026-01-02 10:40:00
- Custom Fields:
  * Copilot Delegated: ✅ Yes
  * Time Saved: 5 hours
  * Delegation Type: New Feature

Acceptance Criteria Status:
1. ✅ Calculate weighted average from assignments
2. ✅ Support grade curves (linear addition)
3. ✅ Convert numeric to letter grades
4. ✅ API endpoint: POST /api/grades/calculate
5. ✅ Performance: < 200ms (actual: 105ms P95)
6. ✅ Design validation passed
7. ✅ All tests passing (28/28)

---

## Cross-Tool Summary

| Aspect | Tool | Status | Details |
|--------|------|--------|---------|
| **Code** | GitHub | ✅ Ready | PR #10, 91.2% coverage |
| **Design** | Figma | ✅ Compliant | 11/11 validations passed |
| **Performance** | Grafana | ✅ Exceeds | 105ms P95 (target: 200ms) |
| **Tests** | CI/CD | ✅ Passing | 28 tests, 0 failures |
| **Review** | GitHub | ⏳ Pending | Requested 5 min ago |
| **Status** | Jira | ✅ Ready for QA | All criteria met |

---

## Time & Efficiency Metrics

**Traditional Development:**
- Feature implementation: 3 hours
- Design validation: 1 hour (manual comparison)
- Performance testing: 1 hour (manual threshold setting)
- Cross-tool status updates: 30 minutes
- **Total: 5.5 hours**

**With MCP Orchestration:**
- Feature implementation: 8 minutes (Background Agent)
- Design validation: Auto-generated from Figma tokens
- Performance testing: Auto-generated from Grafana baselines
- Cross-tool updates: Automated on PR merge
- **Total: 30 minutes**

**Time Saved: 5 hours (91% reduction)**

---

## Next Steps

1. ✅ Code review (PR #10) - Pending team review
2. ✅ QA validation - Ready for testing
3. ⏳ Merge to main - After approval
4. ⏳ Deploy to staging - Post-merge
5. ⏳ Production release - After staging validation

---

## Single Source of Truth

All data synchronized across:
- Jira (TEC-15): Project management and acceptance criteria
- GitHub (PR #10): Code, tests, CI/CD status
- Figma (grade-display): Design specifications and tokens
- Grafana (dashboard): Performance metrics and baselines

**Consistency: 100%** - All tools report identical status ✓

---

*Report generated by GitHub Copilot MCP Orchestration*
*Consolidating: Jira + GitHub + Figma + Grafana*
```

**TALKING POINTS:**
- "Single command queried four MCP sources"
- "Product owner sees complete picture: code, design, performance, tests"
- "No need to check multiple tools - unified view"
- "Time saved calculation: 5 hours (91% reduction)"
- "Every data point traceable to source of truth"

**CLIENT VALUE:**
> "Status meetings used to be: 'Let me check Jira... now GitHub... is this deployed? What's the performance?' Now: One command, complete picture from all tools. Meeting time drops from 30 minutes to 5 minutes."

---

## **DEMO CONCLUSION (2 min)**

**NARRATION:**
> "Let me recap what we just demonstrated across four MCP servers..."

### What We Showed

**1. GitHub MCP - Workflow Automation**
- Auto-create branch from Jira assignment
- Auto-create draft PR
- Auto-link PR to Jira issue
- Update PR status programmatically

**2. Figma MCP - Design Integration**
- Extract design tokens (grade thresholds, colors)
- Generate validation tests from design specs
- Ensure code matches designer intent
- Catch design drift automatically

**3. Grafana MCP - Performance Intelligence**
- Query production performance baselines
- Generate intelligent test thresholds (not arbitrary)
- Validate against real-world metrics
- Track performance trends

**4. Jira MCP - Orchestration Hub**
- Central source of truth for story
- Collect validation from all tools
- Update status automatically
- Track cross-tool metrics

---

### The Complete Journey (Demos 1-5)

```
Demo 1: Development Delegation
  └─ Enrollment feature with Coding Agent (4h saved)

Demo 2: QA Automation
  └─ Mocks, factories, test infrastructure (3h saved)

Demo 3: Security Automation
  └─ GHAS vulnerability detection and fixes (2h saved)

Demo 4: Product Owner Workflow
  └─ Jira MCP + BDD acceptance testing (6h saved)

Demo 5: Cross-Tool Orchestration ⭐
  └─ GitHub + Figma + Grafana + Jira (5h saved)
  └─ Shows: Individual tools → Complete SDLC automation
```

**Total Time Saved: 20 hours in Sprint 1**

---

### Key Differentiators

**Before MCP Orchestration:**
- Developer checks Jira → Creates branch manually
- Designer sends specs → Developer interprets manually
- Performance targets are guesses
- Status scattered across multiple tools
- Manual updates, duplicate data entry

**With MCP Orchestration:**
- Issue assigned → Branch created automatically
- Design tokens extracted → Tests auto-generated
- Performance targets from production data
- Single command queries all tools
- Automated cross-tool synchronization

---

### ROI Deep Dive

| Activity | Traditional | With MCP | Savings |
|----------|-------------|----------|---------|
| Branch setup | 5 min | 30 sec | 90% |
| Design validation | 1 hour | Auto-generated | 100% |
| Performance testing | 1 hour | Auto-generated | 100% |
| Status updates | 30 min | Automated | 100% |
| Status reports | 30 min | 1 command | 97% |
| **Total per feature** | **5.5 hours** | **30 minutes** | **91%** |

**At scale:**
- 10 features/sprint: 50 hours saved
- 4 sprints: 200 hours saved per quarter
- Annual: 800 hours = $120K savings (at $150/hr)

---

### What Makes This Different

**Not just time savings - this is:**

✅ **Single Source of Truth**
- Design specs in Figma = code validation tests
- Performance baselines in Grafana = test thresholds
- Acceptance criteria in Jira = implementation requirements

✅ **Automated Consistency**
- Design changes → Tests fail if code doesn't match
- Performance regression → Caught before production
- Status sync → No manual updates needed

✅ **Data-Driven Development**
- Performance targets from real metrics, not guesses
- Design validation from actual design files, not screenshots
- Cross-tool validation ensures nothing missed

✅ **Scalable Velocity**
- First feature: Setup overhead
- Every subsequent feature: Full automation benefits
- Team growth doesn't slow velocity

---

## 🎤 Q&A Preparation

### Expected Questions

**Q: Do we need all four MCP servers for this to work?**
A: No. Start with one (Jira MCP). Add others as you see value. Each MCP adds specific capability:
- Jira: Project management sync
- GitHub: Workflow automation
- Figma: Design validation
- Grafana: Performance intelligence

**Q: What if we use Azure DevOps instead of Jira?**
A: Azure DevOps MCP available. Identical workflow: `@azure-boards` instead of `@atlassian`. Same orchestration patterns apply.

**Q: What if we don't have Grafana?**
A: Use whatever monitoring you have. If MCP not available, manually extract metrics and feed to Copilot. You lose automation but keep intelligent threshold generation.

**Q: How do you prevent design drift over time?**
A: Tests validate against Figma tokens. When designer updates Figma, re-extract tokens, tests fail if implementation outdated. Forces sync.

**Q: Can this work for mobile apps, not just APIs?**
A: Yes! Figma MCP works great for mobile. Extract design tokens for colors, spacing, typography. Same validation approach. GitHub and Jira MCPs are platform-agnostic.

**Q: What's the learning curve for developers?**
A: MCP syntax: 1 day. Workflow patterns: 2-3 days. Full productivity: 1 sprint. ROI positive by Sprint 2.

**Q: How do you handle conflicts between tools?**
A: Establish single source of truth per domain:
- Requirements: Jira
- Design: Figma
- Performance: Grafana
- Code: GitHub
MCP reads, doesn't write conflicting data.

**Q: What if Figma file is huge?**
A: Extract only specific components/frames. MCP lets you target: `@figma get component GradeDisplayBadge`. Doesn't pull entire file.

**Q: Can non-technical people use this?**
A: Product owners can query status: `@atlassian update TEC-15 with comment "When will this be done?"` Copilot can respond with cross-tool status. Power users can do simple queries.

**Q: What about security with multiple MCP connections?**
A: Each MCP uses OAuth tokens, scoped permissions. Figma: read-only access to specific files. Grafana: read-only dashboards. GitHub: developer's personal scope. Jira: authenticated user's permissions.

---

## 📋 Post-Demo Actions

**Immediate:**
✅ Merge PR #10 (Grade Calculation)
✅ Deploy to staging for QA validation
✅ Monitor Grafana for actual production metrics
✅ Update Jira TEC-15 to "Done" after QA sign-off

**Follow-up with Client:**
✅ Share unified status report template
✅ Provide MCP setup guides for their tools
✅ Schedule workshop for their team
✅ Calculate ROI for their team size and feature velocity

**Documentation:**
✅ Document cross-tool workflow patterns
✅ Create templates for design token extraction
✅ Create templates for performance baseline queries
✅ Share best practices for MCP orchestration

---

## 📚 Additional Resources

**MCP Documentation:**
- GitHub MCP: https://github.com/features/copilot/mcp
- Figma MCP: https://marketplace.visualstudio.com/items?itemName=figma.mcp
- Grafana MCP: https://grafana.com/docs/copilot-mcp
- Jira MCP: https://marketplace.atlassian.com/copilot-mcp

**Client Templates:**
- Design token extraction template (Figma → Code)
- Performance baseline query template (Grafana → Tests)
- Cross-tool status report template (All MCPs → Unified view)
- ROI calculator spreadsheet

---

**END OF DEMO 5 SCRIPT**

---

*This demo completes the five-part series, showcasing GitHub Copilot's ability to orchestrate multiple MCP servers for complete SDLC automation. From individual tool delegation (Demos 1-3) to product owner workflow (Demo 4) to cross-tool orchestration (Demo 5), the progression demonstrates sustainable velocity improvement at enterprise scale.*

**Series Summary:**
- Demo 1: Development (4h saved)
- Demo 2: Testing (3h saved)
- Demo 3: Security (2h saved)
- Demo 4: Product Owner + BDD (6h saved)
- Demo 5: Cross-Tool Orchestration (5h saved)
- **Total: 20 hours saved in Sprint 1**

**Key Achievement:** Not just time savings, but systematic approach to maintaining quality, consistency, and velocity as teams scale.
