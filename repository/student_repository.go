package repository

import (
    "techwave/models"
    "sync"
)

type StudentRepository struct {
    data map[int]models.Student
    mu   sync.RWMutex
    nextID int
}

func NewStudentRepository() *StudentRepository {
    return &StudentRepository{
        data: make(map[int]models.Student),
        nextID: 1,
    }
}

func (r *StudentRepository) Create(s models.Student) models.Student {
    r.mu.Lock()
    defer r.mu.Unlock()
    s.ID = r.nextID
    r.nextID++
    r.data[s.ID] = s
    return s
}

func (r *StudentRepository) Get(id int) (models.Student, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    s, ok := r.data[id]
    return s, ok
}

func (r *StudentRepository) List() []models.Student {
    r.mu.RLock()
    defer r.mu.RUnlock()
    students := make([]models.Student, 0, len(r.data))
    for _, s := range r.data {
        students = append(students, s)
    }
    return students
}

func (r *StudentRepository) Update(id int, s models.Student) (models.Student, bool) {
    r.mu.Lock()
    defer r.mu.Unlock()
    _, exists := r.data[id]
    if !exists {
        return models.Student{}, false
    }
    s.ID = id
    r.data[id] = s
    return s, true
}

func (r *StudentRepository) Delete(id int) bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, exists := r.data[id]; exists {
        delete(r.data, id)
        return true
    }
    return false
}