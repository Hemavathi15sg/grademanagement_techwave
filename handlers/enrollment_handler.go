package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"techwave/models"
	"techwave/repository"

	"github.com/gorilla/mux"
)

// EnrollmentHandler handles HTTP requests for enrollment operations
type EnrollmentHandler struct {
	Repo *repository.EnrollmentRepository
}

// CreateEnrollment handles POST requests to create a new enrollment
// Returns 201 Created on success, 400 Bad Request on validation errors
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var enrollment models.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	created, err := h.Repo.Create(enrollment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetEnrollment handles GET requests to retrieve a specific enrollment by ID
// Returns 200 OK with enrollment, 404 Not Found if doesn't exist
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	enrollment, ok := h.Repo.Get(id)
	if !ok {
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
// Returns 200 OK with array of enrollments
func (h *EnrollmentHandler) GetAllEnrollments(w http.ResponseWriter, r *http.Request) {
	studentIDStr := r.URL.Query().Get("student_id")
	courseIDStr := r.URL.Query().Get("course_id")
	
	var enrollments []models.Enrollment
	
	// If both student_id and course_id are provided, get specific enrollment
	if studentIDStr != "" && courseIDStr != "" {
		studentID, err1 := strconv.Atoi(studentIDStr)
		courseID, err2 := strconv.Atoi(courseIDStr)
		if err1 != nil || err2 != nil {
			http.Error(w, "Invalid student_id or course_id", http.StatusBadRequest)
			return
		}
		
		enrollment, ok := h.Repo.GetByStudentAndCourse(studentID, courseID)
		if ok {
			enrollments = []models.Enrollment{enrollment}
		} else {
			enrollments = []models.Enrollment{}
		}
	} else if studentIDStr != "" {
		// Filter by student
		studentID, err := strconv.Atoi(studentIDStr)
		if err != nil {
			http.Error(w, "Invalid student_id", http.StatusBadRequest)
			return
		}
		enrollments = h.Repo.GetByStudentID(studentID)
	} else if courseIDStr != "" {
		// Filter by course
		courseID, err := strconv.Atoi(courseIDStr)
		if err != nil {
			http.Error(w, "Invalid course_id", http.StatusBadRequest)
			return
		}
		enrollments = h.Repo.GetByCourseID(courseID)
	} else {
		// Get all enrollments
		enrollments = h.Repo.List()
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrollments)
}

// UpdateEnrollment handles PUT requests to update an existing enrollment
// Returns 200 OK with updated enrollment, 400 on validation errors, 404 if not found
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	var updates models.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	updated, err := h.Repo.Update(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if updated.ID == 0 {
		http.Error(w, "Enrollment not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteEnrollment handles DELETE requests to remove an enrollment
// Returns 204 No Content on success, 404 if not found
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	ok := h.Repo.Delete(id)
	if !ok {
		http.Error(w, "Enrollment not found", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// GetEnrollmentStats handles GET requests to retrieve enrollment statistics
// Returns 200 OK with statistics grouped by status
func (h *EnrollmentHandler) GetEnrollmentStats(w http.ResponseWriter, r *http.Request) {
	stats := h.Repo.GetEnrollmentStats()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
