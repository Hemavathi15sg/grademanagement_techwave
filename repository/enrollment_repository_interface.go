package repository

import "techwave/models"

// EnrollmentRepositoryInterface defines the contract for enrollment data operations
// This interface enables mocking for testing purposes
//
//go:generate mockgen -destination=../mocks/mock_enrollment_repository.go -package=mocks techwave/repository EnrollmentRepositoryInterface
type EnrollmentRepositoryInterface interface {
	// Create stores a new enrollment with auto-generated ID and timestamps
	Create(enrollment models.Enrollment) (models.Enrollment, error)

	// Get retrieves an enrollment by its ID
	Get(id int) (models.Enrollment, bool)

	// List retrieves all enrollments
	List() []models.Enrollment

	// GetByStudentID retrieves all enrollments for a specific student
	GetByStudentID(studentID int) []models.Enrollment

	// GetByCourseID retrieves all enrollments for a specific course
	GetByCourseID(courseID int) []models.Enrollment

	// GetByStudentAndCourse retrieves a specific enrollment for a student in a course
	GetByStudentAndCourse(studentID, courseID int) (models.Enrollment, bool)

	// Update updates an existing enrollment with validation
	Update(id int, updates models.Enrollment) (models.Enrollment, error)

	// Delete removes an enrollment
	Delete(id int) bool

	// GetEnrollmentStats returns statistics about enrollments grouped by status
	GetEnrollmentStats() EnrollmentStats
}
