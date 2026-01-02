package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"techwave/models"
	"techwave/repository"

	"github.com/gorilla/mux"
)

// DepartmentHandler handles HTTP requests for department operations
type DepartmentHandler struct {
	Repo repository.DepartmentRepositoryInterface
}

// CreateDepartment handles POST requests to create a new department
// Returns 201 Created on success, 400 Bad Request on validation errors
func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var department models.Department
	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.Repo.Create(department)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetDepartment handles GET requests to retrieve a specific department by ID
// Returns 200 OK with department, 404 Not Found if doesn't exist
func (h *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	department, ok := h.Repo.Get(id)
	if !ok {
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(department)
}

// GetAllDepartments handles GET requests to retrieve departments
// Supports query parameter:
//   - code: filter by department code
//
// Returns 200 OK with array of departments or single department
func (h *DepartmentHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code != "" {
		// Get by code
		department, ok := h.Repo.GetByCode(code)
		if !ok {
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(department)
		return
	}

	// Get all departments
	departments := h.Repo.List()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departments)
}

// UpdateDepartment handles PUT requests to update an existing department
// Returns 200 OK with updated department, 400 on validation errors, 404 if not found
func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updates models.Department
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
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteDepartment handles DELETE requests to remove a department
// Returns 204 No Content on success, 404 if not found
func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ok := h.Repo.Delete(id)
	if !ok {
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
