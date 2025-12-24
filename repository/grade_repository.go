package repository

import (
    "techwave/models"
    "sync"
)

type GradeRepository struct {
    data map[int]models.Grade
    mu   sync.RWMutex
    nextID int
}

func NewGradeRepository() *GradeRepository {
    return &GradeRepository{
        data: make(map[int]models.Grade),
        nextID: 1,
    }
}

func (r *GradeRepository) Create(g models.Grade) models.Grade {
    r.mu.Lock()
    defer r.mu.Unlock()
    g.ID = r.nextID
    r.nextID++
    r.data[g.ID] = g
    return g
}

func (r *GradeRepository) Get(id int) (models.Grade, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    g, ok := r.data[id]
    return g, ok
}

func (r *GradeRepository) List() []models.Grade {
    r.mu.RLock()
    defer r.mu.RUnlock()
    grades := make([]models.Grade, 0, len(r.data))
    for _, g := range r.data {
        grades = append(grades, g)
    }
    return grades
}

func (r *GradeRepository) Update(id int, g models.Grade) (models.Grade, bool) {
    r.mu.Lock()
    defer r.mu.Unlock()
    _, exists := r.data[id]
    if !exists {
        return models.Grade{}, false
    }
    g.ID = id
    r.data[id] = g
    return g, true
}

func (r *GradeRepository) Delete(id int) bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, exists := r.data[id]; exists {
        delete(r.data, id)
        return true
    }
    return false
}