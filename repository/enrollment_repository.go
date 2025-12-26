package repository

import (
	"context"
	"fmt"
	"techwave/models"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// EnrollmentRepository manages enrollment data in Redis
type EnrollmentRepository struct {
	client *redis.Client
	ctx    context.Context
}

// NewEnrollmentRepository creates a new enrollment repository with Redis client
func NewEnrollmentRepository(client *redis.Client) *EnrollmentRepository {
	return &EnrollmentRepository{
		client: client,
		ctx:    context.Background(),
	}
}

// Create stores a new enrollment in Redis with generated UUID and timestamps
// Uses Redis pipeline for atomic operations across multiple keys
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) (*models.Enrollment, error) {
	// Validate enrollment before creating
	if err := enrollment.Validate(); err != nil {
		return nil, err
	}
	
	// Generate UUID for new enrollment
	enrollment.ID = uuid.New().String()
	
	// Set timestamps
	now := time.Now()
	enrollment.CreatedAt = now
	enrollment.UpdatedAt = now
	
	// If enrollment date not set, use current time
	if enrollment.EnrollmentDate.IsZero() {
		enrollment.EnrollmentDate = now
	}
	
	// Serialize to JSON
	data, err := enrollment.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize enrollment: %w", err)
	}
	
	// Use pipeline for atomic operations
	pipe := r.client.Pipeline()
	
	// Store enrollment data
	enrollmentKey := fmt.Sprintf("enrollment:%s", enrollment.ID)
	pipe.Set(r.ctx, enrollmentKey, data, 0)
	
	// Add to all enrollments set
	pipe.SAdd(r.ctx, "enrollments:all", enrollment.ID)
	
	// Add to student index
	studentKey := fmt.Sprintf("student:enrollments:%s", enrollment.StudentID)
	pipe.SAdd(r.ctx, studentKey, enrollment.ID)
	
	// Add to course index
	courseKey := fmt.Sprintf("course:enrollments:%s", enrollment.CourseID)
	pipe.SAdd(r.ctx, courseKey, enrollment.ID)
	
	// Execute pipeline
	if _, err := pipe.Exec(r.ctx); err != nil {
		return nil, fmt.Errorf("failed to create enrollment: %w", err)
	}
	
	return enrollment, nil
}

// GetByID retrieves an enrollment by its ID
func (r *EnrollmentRepository) GetByID(id string) (*models.Enrollment, error) {
	enrollmentKey := fmt.Sprintf("enrollment:%s", id)
	data, err := r.client.Get(r.ctx, enrollmentKey).Bytes()
	if err == redis.Nil {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}
	
	var enrollment models.Enrollment
	if err := enrollment.FromJSON(data); err != nil {
		return nil, fmt.Errorf("failed to deserialize enrollment: %w", err)
	}
	
	return &enrollment, nil
}

// GetAll retrieves all enrollments
func (r *EnrollmentRepository) GetAll() ([]*models.Enrollment, error) {
	// Get all enrollment IDs
	ids, err := r.client.SMembers(r.ctx, "enrollments:all").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment IDs: %w", err)
	}
	
	return r.getEnrollmentsByIDs(ids)
}

// GetByStudentID retrieves all enrollments for a specific student
func (r *EnrollmentRepository) GetByStudentID(studentID string) ([]*models.Enrollment, error) {
	studentKey := fmt.Sprintf("student:enrollments:%s", studentID)
	ids, err := r.client.SMembers(r.ctx, studentKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get student enrollments: %w", err)
	}
	
	return r.getEnrollmentsByIDs(ids)
}

// GetByCourseID retrieves all enrollments for a specific course
func (r *EnrollmentRepository) GetByCourseID(courseID string) ([]*models.Enrollment, error) {
	courseKey := fmt.Sprintf("course:enrollments:%s", courseID)
	ids, err := r.client.SMembers(r.ctx, courseKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get course enrollments: %w", err)
	}
	
	return r.getEnrollmentsByIDs(ids)
}

// GetByStudentAndCourse retrieves a specific enrollment for a student in a course
func (r *EnrollmentRepository) GetByStudentAndCourse(studentID, courseID string) (*models.Enrollment, error) {
	// Get intersection of student and course enrollments
	studentKey := fmt.Sprintf("student:enrollments:%s", studentID)
	courseKey := fmt.Sprintf("course:enrollments:%s", courseID)
	
	ids, err := r.client.SInter(r.ctx, studentKey, courseKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to find enrollment: %w", err)
	}
	
	if len(ids) == 0 {
		return nil, nil // Not found
	}
	
	// Return the first match (should only be one)
	enrollments, err := r.getEnrollmentsByIDs(ids[:1])
	if err != nil {
		return nil, err
	}
	
	if len(enrollments) > 0 {
		return enrollments[0], nil
	}
	
	return nil, nil
}

// Update updates an existing enrollment with validation
func (r *EnrollmentRepository) Update(id string, updates *models.Enrollment) (*models.Enrollment, error) {
	// Get existing enrollment
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil // Not found
	}
	
	// Validate status transition if status is changing
	if updates.Status != "" && updates.Status != existing.Status {
		if err := existing.ValidateStatusTransition(updates.Status); err != nil {
			return nil, err
		}
		existing.Status = updates.Status
	}
	
	// Note: StudentID and CourseID cannot be updated to maintain index consistency
	// If these need to change, the enrollment should be deleted and recreated
	
	// Update enrollment date if provided
	if !updates.EnrollmentDate.IsZero() {
		existing.EnrollmentDate = updates.EnrollmentDate
	}
	
	// Validate updated enrollment
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	
	// Update timestamp
	existing.UpdatedAt = time.Now()
	
	// Serialize to JSON
	data, err := existing.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize enrollment: %w", err)
	}
	
	// Store updated enrollment
	enrollmentKey := fmt.Sprintf("enrollment:%s", id)
	if err := r.client.Set(r.ctx, enrollmentKey, data, 0).Err(); err != nil {
		return nil, fmt.Errorf("failed to update enrollment: %w", err)
	}
	
	return existing, nil
}

// Delete removes an enrollment and all its index entries atomically
func (r *EnrollmentRepository) Delete(id string) error {
	// Get enrollment to find student and course IDs for index cleanup
	enrollment, err := r.GetByID(id)
	if err != nil {
		return err
	}
	if enrollment == nil {
		return nil // Already deleted or not found
	}
	
	// Use pipeline for atomic deletion
	pipe := r.client.Pipeline()
	
	// Delete enrollment data
	enrollmentKey := fmt.Sprintf("enrollment:%s", id)
	pipe.Del(r.ctx, enrollmentKey)
	
	// Remove from all enrollments set
	pipe.SRem(r.ctx, "enrollments:all", id)
	
	// Remove from student index
	studentKey := fmt.Sprintf("student:enrollments:%s", enrollment.StudentID)
	pipe.SRem(r.ctx, studentKey, id)
	
	// Remove from course index
	courseKey := fmt.Sprintf("course:enrollments:%s", enrollment.CourseID)
	pipe.SRem(r.ctx, courseKey, id)
	
	// Execute pipeline
	if _, err := pipe.Exec(r.ctx); err != nil {
		return fmt.Errorf("failed to delete enrollment: %w", err)
	}
	
	return nil
}

// EnrollmentStats represents enrollment statistics by status
type EnrollmentStats struct {
	TotalEnrollments int            `json:"total_enrollments"`
	ByStatus         map[string]int `json:"by_status"`
}

// GetEnrollmentStats returns statistics about enrollments grouped by status
func (r *EnrollmentRepository) GetEnrollmentStats() (*EnrollmentStats, error) {
	enrollments, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	
	stats := &EnrollmentStats{
		TotalEnrollments: len(enrollments),
		ByStatus:         make(map[string]int),
	}
	
	// Initialize status counts
	stats.ByStatus["pending"] = 0
	stats.ByStatus["active"] = 0
	stats.ByStatus["completed"] = 0
	
	// Count by status
	for _, enrollment := range enrollments {
		stats.ByStatus[string(enrollment.Status)]++
	}
	
	return stats, nil
}

// getEnrollmentsByIDs is a helper to retrieve multiple enrollments by their IDs
// Uses Redis MGET for efficient batch retrieval instead of N separate GET operations
func (r *EnrollmentRepository) getEnrollmentsByIDs(ids []string) ([]*models.Enrollment, error) {
	if len(ids) == 0 {
		return []*models.Enrollment{}, nil
	}
	
	// Build keys for MGET
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = fmt.Sprintf("enrollment:%s", id)
	}
	
	// Use MGET to retrieve all enrollments in a single operation
	values, err := r.client.MGet(r.ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollments: %w", err)
	}
	
	enrollments := make([]*models.Enrollment, 0, len(ids))
	
	for _, val := range values {
		if val == nil {
			// Enrollment was deleted or doesn't exist, skip it
			continue
		}
		
		// Type assert to string
		data, ok := val.(string)
		if !ok {
			continue
		}
		
		var enrollment models.Enrollment
		if err := enrollment.FromJSON([]byte(data)); err != nil {
			// Skip malformed data but don't fail entire operation
			continue
		}
		enrollments = append(enrollments, &enrollment)
	}
	
	return enrollments, nil
}
