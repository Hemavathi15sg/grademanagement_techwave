# Enrollment Testing Guide

This guide demonstrates how to use the gomock interface, mock implementations, and test data factory for enrollment testing.

## Overview

The enrollment testing infrastructure includes:
- **EnrollmentRepositoryInterface**: Interface for mocking repository operations
- **MockEnrollmentRepositoryInterface**: Generated mock implementation using gomock
- **EnrollmentFactory**: Builder pattern for creating test enrollment data
- **Comprehensive test examples**: Demonstrating various testing scenarios

## Installation

The gomock dependency is already included in the project:

```bash
go get go.uber.org/mock
```

## Using the Test Data Factory

### Basic Usage

Create a default enrollment with valid values:

```go
enrollment := NewEnrollmentFactory().Build()
```

### Builder Pattern

Chain methods to customize enrollment properties:

```go
enrollment := NewEnrollmentFactory().
    WithID(42).
    WithStudentID(999).
    WithCourseID(888).
    WithStatus(models.StatusActive).
    Build()
```

### Scenario Helpers

Quick creation methods for common scenarios:

```go
// Create enrollments with specific statuses
pending := NewPendingEnrollment()
active := NewActiveEnrollment()
completed := NewCompletedEnrollment()

// Create enrollment for specific student or course
enrollment := NewEnrollmentForStudent(123)
enrollment := NewEnrollmentForCourse(456)
enrollment := NewEnrollmentForStudentAndCourse(123, 456)

// Create multiple enrollments
enrollments := NewEnrollmentBatch(10)
mixedStatus := NewMixedStatusEnrollments()
```

### Edge Case Helpers

Test validation with invalid data:

```go
// Invalid status scenarios
invalid := NewEnrollmentWithInvalidStatus()
empty := NewEnrollmentWithEmptyStatus()

// Missing required fields
noStudent := NewEnrollmentWithMissingStudentID()
noCourse := NewEnrollmentWithMissingCourseID()

// Edge cases
zeroDate := NewEnrollmentWithZeroEnrollmentDate()
```

## Using Mocks

### Setup

Import the required packages:

```go
import (
    "testing"
    "techwave/mocks"
    "techwave/models"
    "techwave/repository"
    "go.uber.org/mock/gomock"
)
```

### Basic Mock Usage

```go
func TestWithMock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)
    enrollment := NewPendingEnrollment()

    // Setup expectation
    mockRepo.EXPECT().
        Create(gomock.Any()).
        Return(enrollment, nil)

    // Call the mock
    created, err := mockRepo.Create(enrollment)

    // Verify
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}
```

### Mocking with Specific Arguments

```go
mockRepo.EXPECT().
    Get(1).
    Return(enrollment, true)

result, found := mockRepo.Get(1)
```

### Mocking Error Cases

```go
expectedErr := errors.New("student_id is required")

mockRepo.EXPECT().
    Create(enrollment).
    Return(models.Enrollment{}, expectedErr)

_, err := mockRepo.Create(enrollment)
if err == nil {
    t.Error("expected error")
}
```

### Ordered Expectations

Test scenarios with multiple sequential operations:

```go
gomock.InOrder(
    mockRepo.EXPECT().Create(enrollment).Return(created, nil),
    mockRepo.EXPECT().Get(1).Return(created, true),
    mockRepo.EXPECT().Update(1, gomock.Any()).Return(updated, nil),
)

// Calls must happen in this order
mockRepo.Create(enrollment)
mockRepo.Get(1)
mockRepo.Update(1, updates)
```

### Mocking Collection Returns

```go
enrollments := NewEnrollmentBatch(3)

mockRepo.EXPECT().
    List().
    Return(enrollments)

result := mockRepo.List()
```

### Mocking Statistics

```go
stats := repository.EnrollmentStats{
    TotalEnrollments: 10,
    ByStatus: map[string]int{
        "pending": 3,
        "active": 5,
        "completed": 2,
    },
}

mockRepo.EXPECT().
    GetEnrollmentStats().
    Return(stats)

result := mockRepo.GetEnrollmentStats()
```

## Example Test Patterns

### Testing Validation

```go
func TestValidation(t *testing.T) {
    enrollment := NewEnrollmentWithMissingStudentID()
    err := enrollment.Validate()
    if err == nil {
        t.Error("expected validation error")
    }
}
```

### Testing Status Transitions

```go
func TestStatusTransition(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)
    
    // Mock invalid transition from completed
    mockRepo.EXPECT().
        Update(1, gomock.Any()).
        Return(models.Enrollment{}, errors.New("completed is a final status"))

    _, err := mockRepo.Update(1, models.Enrollment{Status: models.StatusActive})
    if err == nil {
        t.Error("expected error for invalid transition")
    }
}
```

### Testing Complex Scenarios

```go
func TestCompleteWorkflow(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)
    enrollment := NewPendingEnrollment()

    // Setup workflow: Create -> Get -> Update -> Get
    gomock.InOrder(
        mockRepo.EXPECT().Create(enrollment).Return(created, nil),
        mockRepo.EXPECT().Get(1).Return(created, true),
        mockRepo.EXPECT().Update(1, gomock.Any()).Return(updated, nil),
        mockRepo.EXPECT().Get(1).Return(updated, true),
    )

    // Execute workflow and verify each step
    // ... implementation ...
}
```

## Regenerating Mocks

If the EnrollmentRepositoryInterface changes, regenerate the mock:

```bash
go generate ./repository
```

Or manually from the repository root:

```bash
mockgen -destination=./mocks/mock_enrollment_repository.go -package=mocks techwave/repository EnrollmentRepositoryInterface
```

## Best Practices

1. **Use the factory for all test data creation** - Ensures consistency and reduces boilerplate
2. **Mock only what you need** - Don't mock the entire repository if you only need one method
3. **Test edge cases** - Use the edge case helpers to test validation and error handling
4. **Use gomock.InOrder for workflows** - When testing multi-step scenarios
5. **Clean up with defer ctrl.Finish()** - Ensures all expectations are verified
6. **Combine factory and mocks** - Use factory to create test data, mocks to control behavior

## Running Tests

Run all enrollment tests:

```bash
go test -v ./tests -run "Enrollment"
```

Run only factory tests:

```bash
go test -v ./tests -run "TestEnrollmentFactory"
```

Run only mock tests:

```bash
go test -v ./tests -run "TestMock"
```

## See Also

- `tests/enrollment_mock_test.go` - Comprehensive examples of mock usage
- `tests/enrollment_factory.go` - Factory implementation
- `repository/enrollment_repository_interface.go` - Interface definition
