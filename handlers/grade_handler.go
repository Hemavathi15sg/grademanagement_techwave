package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "techwave/models"
    "techwave/repository"

    "github.com/gorilla/mux"
)

type GradeHandler struct {
    Repo *repository.GradeRepository
}

func (h *GradeHandler) CreateGrade(w http.ResponseWriter, r *http.Request) {
    var grade models.Grade
    if err := json.NewDecoder(r.Body).Decode(&grade); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    created := h.Repo.Create(grade)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(created)
}

func (h *GradeHandler) GetGrade(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    grade, ok := h.Repo.Get(id)
    if !ok {
        http.Error(w, "Grade not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(grade)
}

func (h *GradeHandler) ListGrades(w http.ResponseWriter, r *http.Request) {
    grades := h.Repo.List()
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
    updated, ok := h.Repo.Update(id, grade)
    if !ok {
        http.Error(w, "Grade not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(updated)
}

func (h *GradeHandler) DeleteGrade(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    ok := h.Repo.Delete(id)
    if !ok {
        http.Error(w, "Grade not found", http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}