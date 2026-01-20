package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"techwave/internal/services"
	"techwave/models"

	"github.com/gorilla/mux"
)

// GradeHandler handles HTTP requests for grade operations
type GradeHandler struct {
	service services.GradeService
}

// NewGradeHandler creates a new GradeHandler with dependency injection
func NewGradeHandler(service services.GradeService) *GradeHandler {
	return &GradeHandler{
		service: service,
	}
}

func (h *GradeHandler) CreateGrade(w http.ResponseWriter, r *http.Request) {
	var grade models.Grade
	if err := json.NewDecoder(r.Body).Decode(&grade); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateGrade(r.Context(), &grade); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	grade, err := h.service.GetGrade(r.Context(), id)
	if err != nil {
		http.Error(w, "Grade not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(grade)
}

func (h *GradeHandler) ListGrades(w http.ResponseWriter, r *http.Request) {
	grades, err := h.service.ListGrades(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(grades)
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

	if err := h.service.UpdateGrade(r.Context(), id, &grade); err != nil {
		http.Error(w, "Grade not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(grade)
}

func (h *GradeHandler) DeleteGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteGrade(r.Context(), id); err != nil {
		http.Error(w, "Grade not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
