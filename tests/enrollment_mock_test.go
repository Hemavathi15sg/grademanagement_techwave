package tests

import (
	"errors"
	"techwave/mocks"
	"techwave/models"
	"techwave/repository"
	"testing"

	"go.uber.org/mock/gomock"
)

// TestEnrollmentFactory_DefaultValues tests the factory creates valid default enrollments
func TestEnrollmentFactory_DefaultValues(t *testing.T) {
	enrollment := NewEnrollmentFactory().Build()

	if enrollment.ID == 0 {
		t.Error("expected ID to be set")
	}
	if enrollment.StudentID == 0 {
		t.Error("expected StudentID to be set")
	}
	if enrollment.CourseID == 0 {
		t.Error("expected CourseID to be set")
	}
	if enrollment.Status == "" {
		t.Error("expected Status to be set")
	}
}

// TestEnrollmentFactory_BuilderPattern tests the builder pattern chaining
func TestEnrollmentFactory_BuilderPattern(t *testing.T) {
	enrollment := NewEnrollmentFactory().
		WithID(42).
		WithStudentID(999).
		WithCourseID(888).
		WithStatus(models.StatusActive).
		Build()

	if enrollment.ID != 42 {
		t.Errorf("expected ID 42, got %d", enrollment.ID)
	}
	if enrollment.StudentID != 999 {
		t.Errorf("expected StudentID 999, got %d", enrollment.StudentID)
	}
	if enrollment.CourseID != 888 {
		t.Errorf("expected CourseID 888, got %d", enrollment.CourseID)
	}
	if enrollment.Status != models.StatusActive {
		t.Errorf("expected Status active, got %s", enrollment.Status)
	}
}

// TestEnrollmentFactory_ScenarioHelpers tests the scenario helper methods
func TestEnrollmentFactory_ScenarioHelpers(t *testing.T) {
	tests := []struct {
		name           string
		enrollment     models.Enrollment
		expectedStatus models.EnrollmentStatus
	}{
		{"Pending", NewPendingEnrollment(), models.StatusPending},
		{"Active", NewActiveEnrollment(), models.StatusActive},
		{"Completed", NewCompletedEnrollment(), models.StatusCompleted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.enrollment.Status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, tt.enrollment.Status)
			}
		})
	}
}

// TestEnrollmentFactory_EdgeCaseHelpers tests the edge case helper methods
func TestEnrollmentFactory_EdgeCaseHelpers(t *testing.T) {
	t.Run("InvalidStatus", func(t *testing.T) {
		enrollment := NewEnrollmentWithInvalidStatus()
		err := enrollment.Validate()
		if err == nil {
			t.Error("expected validation error for invalid status")
		}
	})

	t.Run("EmptyStatus", func(t *testing.T) {
		enrollment := NewEnrollmentWithEmptyStatus()
		err := enrollment.Validate()
		if err == nil {
			t.Error("expected validation error for empty status")
		}
	})

	t.Run("MissingStudentID", func(t *testing.T) {
		enrollment := NewEnrollmentWithMissingStudentID()
		err := enrollment.Validate()
		if err == nil {
			t.Error("expected validation error for missing student ID")
		}
	})

	t.Run("MissingCourseID", func(t *testing.T) {
		enrollment := NewEnrollmentWithMissingCourseID()
		err := enrollment.Validate()
		if err == nil {
			t.Error("expected validation error for missing course ID")
		}
	})
}

// TestMockEnrollmentRepository_Create demonstrates mocking Create method
func TestMockEnrollmentRepository_Create(t *testing.T) {
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
	if created.Status != models.StatusPending {
		t.Errorf("expected pending status, got %s", created.Status)
	}
}

// TestMockEnrollmentRepository_CreateWithError demonstrates mocking error cases
func TestMockEnrollmentRepository_CreateWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)
	enrollment := NewEnrollmentWithMissingStudentID()
	expectedErr := errors.New("student_id is required")

	// Setup expectation for error
	mockRepo.EXPECT().
		Create(enrollment).
		Return(models.Enrollment{}, expectedErr)

	// Call the mock
	_, err := mockRepo.Create(enrollment)

	// Verify error
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error '%s', got '%s'", expectedErr.Error(), err.Error())
	}
}

// TestMockEnrollmentRepository_Get demonstrates mocking Get method
func TestMockEnrollmentRepository_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)
	enrollment := NewEnrollmentFactory().WithID(1).Build()

	// Setup expectation
	mockRepo.EXPECT().
		Get(1).
		Return(enrollment, true)

	// Call the mock
	result, found := mockRepo.Get(1)

	// Verify
	if !found {
		t.Error("expected enrollment to be found")
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

// TestMockEnrollmentRepository_GetNotFound demonstrates mocking not found case
func TestMockEnrollmentRepository_GetNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)

	// Setup expectation for not found
	mockRepo.EXPECT().
		Get(999).
		Return(models.Enrollment{}, false)

	// Call the mock
	_, found := mockRepo.Get(999)

	// Verify
	if found {
		t.Error("expected enrollment not to be found")
	}
}

// TestMockEnrollmentRepository_List demonstrates mocking List method
func TestMockEnrollmentRepository_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)
	enrollments := NewEnrollmentBatch(3)

	// Setup expectation
	mockRepo.EXPECT().
		List().
		Return(enrollments)

	// Call the mock
	result := mockRepo.List()

	// Verify
	if len(result) != 3 {
		t.Errorf("expected 3 enrollments, got %d", len(result))
	}
}

// TestMockEnrollmentRepository_GetByStudentID demonstrates filtering by student
func TestMockEnrollmentRepository_GetByStudentID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)
	studentID := 123
	enrollments := []models.Enrollment{
		NewEnrollmentForStudent(studentID),
		NewEnrollmentForStudent(studentID),
	}

	// Setup expectation
	mockRepo.EXPECT().
		GetByStudentID(studentID).
		Return(enrollments)

	// Call the mock
	result := mockRepo.GetByStudentID(studentID)

	// Verify
	if len(result) != 2 {
		t.Errorf("expected 2 enrollments, got %d", len(result))
	}
	for _, e := range result {
		if e.StudentID != studentID {
			t.Errorf("expected StudentID %d, got %d", studentID, e.StudentID)
		}
	}
}

// TestMockEnrollmentRepository_Update demonstrates mocking Update method
func TestMockEnrollmentRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)
	enrollment := NewEnrollmentFactory().
		WithID(1).
		WithStatus(models.StatusActive).
		Build()

	// Setup expectation
	mockRepo.EXPECT().
		Update(1, gomock.Any()).
		Return(enrollment, nil)

	// Call the mock
	updated, err := mockRepo.Update(1, enrollment)

	// Verify
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if updated.Status != models.StatusActive {
		t.Errorf("expected active status, got %s", updated.Status)
	}
}

// TestMockEnrollmentRepository_UpdateInvalidTransition demonstrates error on invalid status transition
func TestMockEnrollmentRepository_UpdateInvalidTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)
	expectedErr := errors.New("completed is a final status and cannot be changed")

	// Setup expectation for invalid transition
	mockRepo.EXPECT().
		Update(1, gomock.Any()).
		Return(models.Enrollment{}, expectedErr)

	// Call the mock
	_, err := mockRepo.Update(1, models.Enrollment{Status: models.StatusActive})

	// Verify error
	if err == nil {
		t.Error("expected error for invalid status transition")
	}
}

// TestMockEnrollmentRepository_Delete demonstrates mocking Delete method
func TestMockEnrollmentRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)

	// Setup expectation
	mockRepo.EXPECT().
		Delete(1).
		Return(true)

	// Call the mock
	deleted := mockRepo.Delete(1)

	// Verify
	if !deleted {
		t.Error("expected enrollment to be deleted")
	}
}

// TestMockEnrollmentRepository_GetEnrollmentStats demonstrates mocking stats method
func TestMockEnrollmentRepository_GetEnrollmentStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)
	stats := repository.EnrollmentStats{
		TotalEnrollments: 10,
		ByStatus: map[string]int{
			"pending":   3,
			"active":    5,
			"completed": 2,
		},
	}

	// Setup expectation
	mockRepo.EXPECT().
		GetEnrollmentStats().
		Return(stats)

	// Call the mock
	result := mockRepo.GetEnrollmentStats()

	// Verify
	if result.TotalEnrollments != 10 {
		t.Errorf("expected 10 total enrollments, got %d", result.TotalEnrollments)
	}
	if result.ByStatus["active"] != 5 {
		t.Errorf("expected 5 active enrollments, got %d", result.ByStatus["active"])
	}
}

// TestMockEnrollmentRepository_ComplexScenario demonstrates a complex test scenario
func TestMockEnrollmentRepository_ComplexScenario(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepositoryInterface(ctrl)

	// Scenario: Create enrollment, update status, verify
	enrollment := NewPendingEnrollment()
	createdEnrollment := NewEnrollmentFactory().
		WithID(1).
		WithStatus(models.StatusPending).
		Build()
	updatedEnrollment := NewEnrollmentFactory().
		WithID(1).
		WithStatus(models.StatusActive).
		Build()

	// Setup expectations in order
	gomock.InOrder(
		mockRepo.EXPECT().Create(enrollment).Return(createdEnrollment, nil),
		mockRepo.EXPECT().Get(1).Return(createdEnrollment, true),
		mockRepo.EXPECT().Update(1, gomock.Any()).Return(updatedEnrollment, nil),
		mockRepo.EXPECT().Get(1).Return(updatedEnrollment, true),
	)

	// Execute scenario
	created, err := mockRepo.Create(enrollment)
	if err != nil {
		t.Fatalf("unexpected error creating: %v", err)
	}

	retrieved, found := mockRepo.Get(1)
	if !found {
		t.Fatal("enrollment not found after creation")
	}
	if retrieved.Status != models.StatusPending {
		t.Errorf("expected pending status, got %s", retrieved.Status)
	}

	updated, err := mockRepo.Update(created.ID, models.Enrollment{Status: models.StatusActive})
	if err != nil {
		t.Fatalf("unexpected error updating: %v", err)
	}

	final, found := mockRepo.Get(1)
	if !found {
		t.Fatal("enrollment not found after update")
	}
	if final.Status != models.StatusActive {
		t.Errorf("expected active status after update, got %s", final.Status)
	}
	if updated.Status != final.Status {
		t.Error("updated and final status should match")
	}
}
