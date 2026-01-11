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
	students      = make(map[int]models.Student)
	studentsMu    sync.RWMutex
	nextStudentID = 1
)

type StudentHandler struct {
	// No repository - direct data access
}

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Direct database logic mixed with HTTP handler - BAD!
	studentsMu.Lock()
	student.ID = nextStudentID
	nextStudentID++
	students[student.ID] = student
	studentsMu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func (h *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Direct data access in handler - BAD!
	studentsMu.RLock()
	student, ok := students[id]
	studentsMu.RUnlock()

	if !ok {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(student)
}

func (h *StudentHandler) ListStudents(w http.ResponseWriter, r *http.Request) {
	// Direct data access in handler - BAD!
	studentsMu.RLock()
	list := make([]models.Student, 0, len(students))
	for _, s := range students {
		list = append(list, s)
	}
	studentsMu.RUnlock()

	json.NewEncoder(w).Encode(list)
}

func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Direct data access in handler - BAD!
	studentsMu.Lock()
	_, ok := students[id]
	if !ok {
		studentsMu.Unlock()
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	student.ID = id
	students[id] = student
	studentsMu.Unlock()

	json.NewEncoder(w).Encode(student)
}

func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Direct data access in handler - BAD!
	studentsMu.Lock()
	_, ok := students[id]
	if !ok {
		studentsMu.Unlock()
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	delete(students, id)
	studentsMu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
