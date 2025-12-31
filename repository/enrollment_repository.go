package repository

import (
	"sync"
	"techwave/models"
	"time"
)

// EnrollmentRepository manages enrollment data in memory
type EnrollmentRepository struct {
	data   map[int]models.Enrollment
	mu     sync.RWMutex
	nextID int
}

// NewEnrollmentRepository creates a new enrollment repository with in-memory storage
func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{
		data:   make(map[int]models.Enrollment),
		nextID: 1,
	}
}

// Create stores a new enrollment with auto-generated ID and timestamps
func (r *EnrollmentRepository) Create(enrollment models.Enrollment) (models.Enrollment, error) {
	// Validate enrollment before creating
	if err := enrollment.Validate(); err != nil {
		return models.Enrollment{}, err
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Generate ID
	enrollment.ID = r.nextID
	r.nextID++
	
	// Set timestamps
	now := time.Now()
	enrollment.CreatedAt = now
	enrollment.UpdatedAt = now
	
	// If enrollment date not set, use current time
	if enrollment.EnrollmentDate.IsZero() {
		enrollment.EnrollmentDate = now
	}
	
	r.data[enrollment.ID] = enrollment
	return enrollment, nil
}

// Get retrieves an enrollment by its ID
func (r *EnrollmentRepository) Get(id int) (models.Enrollment, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	enrollment, ok := r.data[id]
	return enrollment, ok
}

// List retrieves all enrollments
func (r *EnrollmentRepository) List() []models.Enrollment {
	r.mu.RLock()
	defer r.mu.RUnlock()
	enrollments := make([]models.Enrollment, 0, len(r.data))
	for _, e := range r.data {
		enrollments = append(enrollments, e)
	}
	return enrollments
}

// GetByStudentID retrieves all enrollments for a specific student
func (r *EnrollmentRepository) GetByStudentID(studentID int) []models.Enrollment {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var enrollments []models.Enrollment
	for _, e := range r.data {
		if e.StudentID == studentID {
			enrollments = append(enrollments, e)
		}
	}
	return enrollments
}

// GetByCourseID retrieves all enrollments for a specific course
func (r *EnrollmentRepository) GetByCourseID(courseID int) []models.Enrollment {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var enrollments []models.Enrollment
	for _, e := range r.data {
		if e.CourseID == courseID {
			enrollments = append(enrollments, e)
		}
	}
	return enrollments
}

// GetByStudentAndCourse retrieves a specific enrollment for a student in a course
func (r *EnrollmentRepository) GetByStudentAndCourse(studentID, courseID int) (models.Enrollment, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, e := range r.data {
		if e.StudentID == studentID && e.CourseID == courseID {
			return e, true
		}
	}
	return models.Enrollment{}, false
}

// Update updates an existing enrollment with validation
func (r *EnrollmentRepository) Update(id int, updates models.Enrollment) (models.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Get existing enrollment
	existing, exists := r.data[id]
	if !exists {
		return models.Enrollment{}, nil
	}
	
	// Validate status transition if status is changing
	if updates.Status != "" && updates.Status != existing.Status {
		if err := existing.ValidateStatusTransition(updates.Status); err != nil {
			return models.Enrollment{}, err
		}
		existing.Status = updates.Status
	}
	
	// Update enrollment date if provided
	if !updates.EnrollmentDate.IsZero() {
		existing.EnrollmentDate = updates.EnrollmentDate
	}
	
	// Validate updated enrollment
	if err := existing.Validate(); err != nil {
		return models.Enrollment{}, err
	}
	
	// Update timestamp
	existing.UpdatedAt = time.Now()
	
	r.data[id] = existing
	return existing, nil
}

// Delete removes an enrollment
func (r *EnrollmentRepository) Delete(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[id]; exists {
		delete(r.data, id)
		return true
	}
	return false
}

// EnrollmentStats represents enrollment statistics by status
type EnrollmentStats struct {
	TotalEnrollments int            `json:"total_enrollments"`
	ByStatus         map[string]int `json:"by_status"`
}

// GetEnrollmentStats returns statistics about enrollments grouped by status
func (r *EnrollmentRepository) GetEnrollmentStats() EnrollmentStats {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	stats := EnrollmentStats{
		TotalEnrollments: len(r.data),
		ByStatus:         make(map[string]int),
	}
	
	// Initialize status counts
	stats.ByStatus["pending"] = 0
	stats.ByStatus["active"] = 0
	stats.ByStatus["completed"] = 0
	
	// Count by status
	for _, enrollment := range r.data {
		stats.ByStatus[string(enrollment.Status)]++
	}
	
	return stats
}
