package handlers

import (
	"encoding/json"
	"net/http"
	"techwave/models"
	"techwave/repository"

	"github.com/gorilla/mux"
)

// EnrollmentHandler handles HTTP requests for enrollment operations
type EnrollmentHandler struct {
	Repo *repository.EnrollmentRepository
}

// CreateEnrollment handles POST requests to create a new enrollment
// Returns 201 Created on success, 400 Bad Request on validation errors, 500 on server errors
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var enrollment models.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	created, err := h.Repo.Create(&enrollment)
	if err != nil {
		// Check if it's a validation error (return 400) or server error (return 500)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetEnrollment handles GET requests to retrieve a specific enrollment by ID
// Returns 200 OK with enrollment, 404 Not Found if doesn't exist, 500 on server errors
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	if id == "" {
		http.Error(w, "Enrollment ID is required", http.StatusBadRequest)
		return
	}
	
	enrollment, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if enrollment == nil {
		http.Error(w, "Enrollment not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollment)
}

// GetAllEnrollments handles GET requests to retrieve enrollments
// Supports query parameters:
//   - student_id: filter by student
//   - course_id: filter by course
//   - both: get specific enrollment for student in course
// Returns 200 OK with array of enrollments, 500 on server errors
func (h *EnrollmentHandler) GetAllEnrollments(w http.ResponseWriter, r *http.Request) {
	studentID := r.URL.Query().Get("student_id")
	courseID := r.URL.Query().Get("course_id")
	
	var enrollments []*models.Enrollment
	var err error
	
	// If both student_id and course_id are provided, get specific enrollment
	if studentID != "" && courseID != "" {
		enrollment, err := h.Repo.GetByStudentAndCourse(studentID, courseID)
		if err != nil {
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if enrollment != nil {
			enrollments = []*models.Enrollment{enrollment}
		} else {
			enrollments = []*models.Enrollment{}
		}
	} else if studentID != "" {
		// Filter by student
		enrollments, err = h.Repo.GetByStudentID(studentID)
		if err != nil {
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if courseID != "" {
		// Filter by course
		enrollments, err = h.Repo.GetByCourseID(courseID)
		if err != nil {
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Get all enrollments
		enrollments, err = h.Repo.GetAll()
		if err != nil {
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollments)
}

// UpdateEnrollment handles PUT requests to update an existing enrollment
// Returns 200 OK with updated enrollment, 400 on validation errors, 404 if not found, 500 on server errors
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	if id == "" {
		http.Error(w, "Enrollment ID is required", http.StatusBadRequest)
		return
	}
	
	var updates models.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	updated, err := h.Repo.Update(id, &updates)
	if err != nil {
		// Check if it's a validation error (return 400) or server error (return 500)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if updated == nil {
		http.Error(w, "Enrollment not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteEnrollment handles DELETE requests to remove an enrollment
// Returns 204 No Content on success, 404 if not found, 500 on server errors
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	if id == "" {
		http.Error(w, "Enrollment ID is required", http.StatusBadRequest)
		return
	}
	
	// Check if enrollment exists before deletion
	enrollment, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if enrollment == nil {
		http.Error(w, "Enrollment not found", http.StatusNotFound)
		return
	}
	
	if err := h.Repo.Delete(id); err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// GetEnrollmentStats handles GET requests to retrieve enrollment statistics
// Returns 200 OK with statistics grouped by status
func (h *EnrollmentHandler) GetEnrollmentStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Repo.GetEnrollmentStats()
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
