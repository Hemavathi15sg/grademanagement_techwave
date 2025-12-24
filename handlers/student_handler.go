package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "techwave/models"
    "techwave/repository"

    "github.com/gorilla/mux"
)

type StudentHandler struct {
    Repo *repository.StudentRepository
}

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
    var student models.Student
    if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    created := h.Repo.Create(student)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(created)
}

func (h *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    student, ok := h.Repo.Get(id)
    if !ok {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(student)
}

func (h *StudentHandler) ListStudents(w http.ResponseWriter, r *http.Request) {
    students := h.Repo.List()
    json.NewEncoder(w).Encode(students)
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
    updated, ok := h.Repo.Update(id, student)
    if !ok {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(updated)
}

func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    ok := h.Repo.Delete(id)
    if !ok {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}