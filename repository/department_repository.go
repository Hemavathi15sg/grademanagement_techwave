package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"techwave/models"
	"time"

	"github.com/redis/go-redis/v9"
)

// DepartmentRepository manages department data with Redis caching
type DepartmentRepository struct {
	data        map[int]models.Department
	codeIndex   map[string]int // Map department code to ID for fast lookup
	mu          sync.RWMutex
	nextID      int
	redisClient *redis.Client
	ctx         context.Context
}

// NewDepartmentRepository creates a new department repository with optional Redis caching
func NewDepartmentRepository(redisClient *redis.Client) *DepartmentRepository {
	return &DepartmentRepository{
		data:        make(map[int]models.Department),
		codeIndex:   make(map[string]int),
		nextID:      1,
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// Create stores a new department with auto-generated ID and timestamps
func (r *DepartmentRepository) Create(department models.Department) (models.Department, error) {
	// Validate department before creating
	if err := department.Validate(); err != nil {
		return models.Department{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if code already exists
	if _, exists := r.codeIndex[department.DepartmentCode]; exists {
		return models.Department{}, fmt.Errorf("department with code %s already exists", department.DepartmentCode)
	}

	// Generate ID
	department.ID = r.nextID
	r.nextID++

	// Set timestamps
	now := time.Now()
	department.CreatedAt = now
	department.UpdatedAt = now

	r.data[department.ID] = department
	r.codeIndex[department.DepartmentCode] = department.ID

	// Cache in Redis if available
	if r.redisClient != nil {
		r.cacheSet(department)
	}

	return department, nil
}

// Get retrieves a department by its ID with Redis caching
func (r *DepartmentRepository) Get(id int) (models.Department, bool) {
	// Try Redis cache first if available
	if r.redisClient != nil {
		if cached, found := r.cacheGet(id); found {
			return cached, true
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	department, ok := r.data[id]

	// Cache the result if found and Redis is available
	if ok && r.redisClient != nil {
		r.cacheSet(department)
	}

	return department, ok
}

// List retrieves all departments
func (r *DepartmentRepository) List() []models.Department {
	r.mu.RLock()
	defer r.mu.RUnlock()
	departments := make([]models.Department, 0, len(r.data))
	for _, d := range r.data {
		departments = append(departments, d)
	}
	return departments
}

// GetByCode retrieves a department by its unique code with Redis caching
func (r *DepartmentRepository) GetByCode(code string) (models.Department, bool) {
	// Try Redis cache first if available
	if r.redisClient != nil {
		if cached, found := r.cacheGetByCode(code); found {
			return cached, true
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.codeIndex[code]
	if !exists {
		return models.Department{}, false
	}

	department := r.data[id]

	// Cache the result if Redis is available
	if r.redisClient != nil {
		r.cacheSet(department)
	}

	return department, true
}

// Update updates an existing department with validation
func (r *DepartmentRepository) Update(id int, updates models.Department) (models.Department, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get existing department
	existing, exists := r.data[id]
	if !exists {
		return models.Department{}, nil
	}

	// Update fields if provided
	if updates.DepartmentName != "" {
		existing.DepartmentName = updates.DepartmentName
	}
	if updates.DepartmentCode != "" && updates.DepartmentCode != existing.DepartmentCode {
		// Check if new code already exists
		if existingID, codeExists := r.codeIndex[updates.DepartmentCode]; codeExists && existingID != id {
			return models.Department{}, fmt.Errorf("department with code %s already exists", updates.DepartmentCode)
		}
		// Remove old code index and add new one
		delete(r.codeIndex, existing.DepartmentCode)
		r.codeIndex[updates.DepartmentCode] = id
		existing.DepartmentCode = updates.DepartmentCode
	}
	if updates.DepartmentHead != "" {
		existing.DepartmentHead = updates.DepartmentHead
	}
	if updates.AnnualBudget > 0 {
		existing.AnnualBudget = updates.AnnualBudget
	}
	if updates.Status != "" {
		existing.Status = updates.Status
	}

	// Validate updated department
	if err := existing.Validate(); err != nil {
		return models.Department{}, err
	}

	// Update timestamp
	existing.UpdatedAt = time.Now()

	r.data[id] = existing

	// Invalidate cache if Redis is available
	if r.redisClient != nil {
		r.cacheDelete(id)
		r.cacheSet(existing)
	}

	return existing, nil
}

// Delete removes a department
func (r *DepartmentRepository) Delete(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	department, exists := r.data[id]
	if exists {
		delete(r.data, id)
		delete(r.codeIndex, department.DepartmentCode)

		// Invalidate cache if Redis is available
		if r.redisClient != nil {
			r.cacheDelete(id)
		}
		return true
	}
	return false
}

// Redis caching helper methods

func (r *DepartmentRepository) cacheKey(id int) string {
	return fmt.Sprintf("department:%d", id)
}

func (r *DepartmentRepository) cacheKeyByCode(code string) string {
	return fmt.Sprintf("department:code:%s", code)
}

func (r *DepartmentRepository) cacheSet(department models.Department) {
	data, err := json.Marshal(department)
	if err != nil {
		return
	}

	// Cache by ID
	r.redisClient.Set(r.ctx, r.cacheKey(department.ID), data, 5*time.Minute)

	// Cache by code
	r.redisClient.Set(r.ctx, r.cacheKeyByCode(department.DepartmentCode), data, 5*time.Minute)
}

func (r *DepartmentRepository) cacheGet(id int) (models.Department, bool) {
	val, err := r.redisClient.Get(r.ctx, r.cacheKey(id)).Result()
	if err != nil {
		return models.Department{}, false
	}

	var department models.Department
	if err := json.Unmarshal([]byte(val), &department); err != nil {
		return models.Department{}, false
	}

	return department, true
}

func (r *DepartmentRepository) cacheGetByCode(code string) (models.Department, bool) {
	val, err := r.redisClient.Get(r.ctx, r.cacheKeyByCode(code)).Result()
	if err != nil {
		return models.Department{}, false
	}

	var department models.Department
	if err := json.Unmarshal([]byte(val), &department); err != nil {
		return models.Department{}, false
	}

	return department, true
}

func (r *DepartmentRepository) cacheDelete(id int) {
	// Get department first to delete code cache
	if department, ok := r.data[id]; ok {
		r.redisClient.Del(r.ctx, r.cacheKeyByCode(department.DepartmentCode))
	}
	r.redisClient.Del(r.ctx, r.cacheKey(id))
}
