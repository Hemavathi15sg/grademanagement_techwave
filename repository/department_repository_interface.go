package repository

import "techwave/models"

// DepartmentRepositoryInterface defines the contract for department data operations
// This interface enables mocking for testing purposes
//
//go:generate mockgen -destination=../mocks/mock_department_repository.go -package=mocks techwave/repository DepartmentRepositoryInterface
type DepartmentRepositoryInterface interface {
	// Create stores a new department with auto-generated ID and timestamps
	Create(department models.Department) (models.Department, error)

	// Get retrieves a department by its ID
	Get(id int) (models.Department, bool)

	// List retrieves all departments
	List() []models.Department

	// GetByCode retrieves a department by its unique code
	GetByCode(code string) (models.Department, bool)

	// Update updates an existing department with validation
	Update(id int, updates models.Department) (models.Department, error)

	// Delete removes a department
	Delete(id int) bool
}
