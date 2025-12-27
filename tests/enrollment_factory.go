package tests

import (
	"techwave/models"
	"time"
)

// EnrollmentFactory provides a builder pattern for creating test enrollment data
type EnrollmentFactory struct {
	enrollment models.Enrollment
}

// NewEnrollmentFactory creates a new factory with default valid values
func NewEnrollmentFactory() *EnrollmentFactory {
	return &EnrollmentFactory{
		enrollment: models.Enrollment{
			ID:             1,
			StudentID:      100,
			CourseID:       200,
			EnrollmentDate: time.Now(),
			Status:         models.StatusPending,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}
}

// WithID sets the enrollment ID
func (f *EnrollmentFactory) WithID(id int) *EnrollmentFactory {
	f.enrollment.ID = id
	return f
}

// WithStudentID sets the student ID
func (f *EnrollmentFactory) WithStudentID(studentID int) *EnrollmentFactory {
	f.enrollment.StudentID = studentID
	return f
}

// WithCourseID sets the course ID
func (f *EnrollmentFactory) WithCourseID(courseID int) *EnrollmentFactory {
	f.enrollment.CourseID = courseID
	return f
}

// WithEnrollmentDate sets the enrollment date
func (f *EnrollmentFactory) WithEnrollmentDate(date time.Time) *EnrollmentFactory {
	f.enrollment.EnrollmentDate = date
	return f
}

// WithStatus sets the enrollment status
func (f *EnrollmentFactory) WithStatus(status models.EnrollmentStatus) *EnrollmentFactory {
	f.enrollment.Status = status
	return f
}

// WithCreatedAt sets the created timestamp
func (f *EnrollmentFactory) WithCreatedAt(t time.Time) *EnrollmentFactory {
	f.enrollment.CreatedAt = t
	return f
}

// WithUpdatedAt sets the updated timestamp
func (f *EnrollmentFactory) WithUpdatedAt(t time.Time) *EnrollmentFactory {
	f.enrollment.UpdatedAt = t
	return f
}

// Build returns the constructed enrollment
func (f *EnrollmentFactory) Build() models.Enrollment {
	return f.enrollment
}

// Common scenario helpers

// NewPendingEnrollment creates an enrollment with pending status
func NewPendingEnrollment() models.Enrollment {
	return NewEnrollmentFactory().
		WithStatus(models.StatusPending).
		Build()
}

// NewActiveEnrollment creates an enrollment with active status
func NewActiveEnrollment() models.Enrollment {
	return NewEnrollmentFactory().
		WithStatus(models.StatusActive).
		Build()
}

// NewCompletedEnrollment creates an enrollment with completed status
func NewCompletedEnrollment() models.Enrollment {
	return NewEnrollmentFactory().
		WithStatus(models.StatusCompleted).
		Build()
}

// Edge case helpers for invalid status

// NewEnrollmentWithInvalidStatus creates an enrollment with an invalid status
// This is useful for testing validation logic
func NewEnrollmentWithInvalidStatus() models.Enrollment {
	enrollment := NewEnrollmentFactory().Build()
	// Use type conversion to set an invalid status that bypasses compile-time type safety
	enrollment.Status = models.EnrollmentStatus("invalid_status")
	return enrollment
}

// NewEnrollmentWithEmptyStatus creates an enrollment with empty status
func NewEnrollmentWithEmptyStatus() models.Enrollment {
	return NewEnrollmentFactory().
		WithStatus("").
		Build()
}

// NewEnrollmentWithMissingStudentID creates an enrollment without student ID
func NewEnrollmentWithMissingStudentID() models.Enrollment {
	return NewEnrollmentFactory().
		WithStudentID(0).
		Build()
}

// NewEnrollmentWithMissingCourseID creates an enrollment without course ID
func NewEnrollmentWithMissingCourseID() models.Enrollment {
	return NewEnrollmentFactory().
		WithCourseID(0).
		Build()
}

// NewEnrollmentWithZeroEnrollmentDate creates an enrollment with zero enrollment date
func NewEnrollmentWithZeroEnrollmentDate() models.Enrollment {
	return NewEnrollmentFactory().
		WithEnrollmentDate(time.Time{}).
		Build()
}

// Batch creation helpers

// NewEnrollmentBatch creates multiple enrollments with different IDs
func NewEnrollmentBatch(count int) []models.Enrollment {
	enrollments := make([]models.Enrollment, count)
	for i := 0; i < count; i++ {
		enrollments[i] = NewEnrollmentFactory().
			WithID(i + 1).
			WithStudentID(100 + i).
			WithCourseID(200 + i).
			Build()
	}
	return enrollments
}

// NewMixedStatusEnrollments creates enrollments with different statuses
func NewMixedStatusEnrollments() []models.Enrollment {
	return []models.Enrollment{
		NewPendingEnrollment(),
		NewActiveEnrollment(),
		NewCompletedEnrollment(),
	}
}

// Scenario helpers for specific test cases

// NewEnrollmentForStudent creates an enrollment for a specific student
func NewEnrollmentForStudent(studentID int) models.Enrollment {
	return NewEnrollmentFactory().
		WithStudentID(studentID).
		Build()
}

// NewEnrollmentForCourse creates an enrollment for a specific course
func NewEnrollmentForCourse(courseID int) models.Enrollment {
	return NewEnrollmentFactory().
		WithCourseID(courseID).
		Build()
}

// NewEnrollmentForStudentAndCourse creates an enrollment for a specific student and course
func NewEnrollmentForStudentAndCourse(studentID, courseID int) models.Enrollment {
	return NewEnrollmentFactory().
		WithStudentID(studentID).
		WithCourseID(courseID).
		Build()
}
