package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"techwave/models"

	"github.com/gorilla/mux"
)

// Global in-memory storage - MESSY! Mixed with HTTP logic
var (
	grades   = make(map[int]models.Grade)
	gradesMu sync.RWMutex
	nextGradeID = 1
)

type GradeHandler struct {
	// No repository or cache - direct data access
}

func (h *GradeHandler) CreateGrade(w http.ResponseWriter, r *http.Request) {
	var grade models.Grade
	if err := json.NewDecoder(r.Body).Decode(&grade); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Direct database logic mixed with HTTP handler - BAD!
	gradesMu.Lock()
	grade.ID = nextGradeID
	nextGradeID++
	grades[grade.ID] = grade
	gradesMu.Unlock()
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(grade)
}

func (h *GradeHandler) GetGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Direct data access in handler - BAD!
	gradesMu.RLock()
	grade, ok := grades[id]
	gradesMu.RUnlock()
	
	if !ok {
		http.Error(w, "Grade not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(grade)
}

func (h *GradeHandler) ListGrades(w http.ResponseWriter, r *http.Request) {
	// Direct data access in handler - BAD!
	gradesMu.RLock()
	list := make([]models.Grade, 0, len(grades))
	for _, g := range grades {
		list = append(list, g)
	}
	gradesMu.RUnlock()
	
	json.NewEncoder(w).Encode(list)
}

func (h *GradeHandler) UpdateGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var grade models.Grade
	if err := json.NewDecoder(r.Body).Decode(&grade); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Direct data access in handler - BAD!
	gradesMu.Lock()
	_, ok := grades[id]
	if !ok {
		gradesMu.Unlock()
		http.Error(w, "Grade not found", http.StatusNotFound)
		return
	}
	grade.ID = id
	grades[id] = grade
	gradesMu.Unlock()

	json.NewEncoder(w).Encode(grade)
}

func (h *GradeHandler) DeleteGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	// Direct data access in handler - BAD!
	gradesMu.Lock()
	_, ok := grades[id]
	if !ok {
		gradesMu.Unlock()
		http.Error(w, "Grade not found", http.StatusNotFound)
		return
	}
	delete(grades, id)
	gradesMu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
