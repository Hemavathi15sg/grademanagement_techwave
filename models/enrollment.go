package models

import (
	"encoding/json"
	"errors"
	"time"
)

// EnrollmentStatus represents the current state of an enrollment
type EnrollmentStatus string

const (
	// StatusPending indicates enrollment is awaiting approval or processing
	StatusPending EnrollmentStatus = "pending"
	// StatusActive indicates student is actively enrolled in the course
	StatusActive EnrollmentStatus = "active"
	// StatusCompleted indicates student has completed the course
	StatusCompleted EnrollmentStatus = "completed"
)

// Enrollment represents a student's enrollment in a course
type Enrollment struct {
	ID             string           `json:"id"`
	StudentID      string           `json:"student_id"`
	CourseID       string           `json:"course_id"`
	EnrollmentDate time.Time        `json:"enrollment_date"`
	Status         EnrollmentStatus `json:"status"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// Validate checks if the enrollment has all required fields and valid status
func (e *Enrollment) Validate() error {
	if e.StudentID == "" {
		return errors.New("student_id is required")
	}
	if e.CourseID == "" {
		return errors.New("course_id is required")
	}
	if e.Status == "" {
		return errors.New("status is required")
	}
	
	// Validate status is one of the allowed values
	if e.Status != StatusPending && e.Status != StatusActive && e.Status != StatusCompleted {
		return errors.New("status must be one of: pending, active, completed")
	}
	
	return nil
}

// ValidateStatusTransition checks if transitioning from current status to new status is allowed
// Valid transitions:
//   - pending -> active or completed
//   - active -> completed
//   - completed -> no transitions (final state)
func (e *Enrollment) ValidateStatusTransition(newStatus EnrollmentStatus) error {
	// Validate the new status is valid
	if newStatus != StatusPending && newStatus != StatusActive && newStatus != StatusCompleted {
		return errors.New("new status must be one of: pending, active, completed")
	}
	
	// If status hasn't changed, allow it
	if e.Status == newStatus {
		return nil
	}
	
	switch e.Status {
	case StatusPending:
		// Can transition to active or completed
		if newStatus == StatusActive || newStatus == StatusCompleted {
			return nil
		}
		return errors.New("pending status can only transition to active or completed")
		
	case StatusActive:
		// Can only transition to completed
		if newStatus == StatusCompleted {
			return nil
		}
		return errors.New("active status can only transition to completed")
		
	case StatusCompleted:
		// Cannot transition from completed (final state)
		return errors.New("completed is a final status and cannot be changed")
		
	default:
		return errors.New("unknown current status")
	}
}

// ToJSON serializes the enrollment to JSON bytes for Redis storage
func (e *Enrollment) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON deserializes enrollment from JSON bytes retrieved from Redis
func (e *Enrollment) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}
