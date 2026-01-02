package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"techwave/handlers"
	"techwave/models"
	"techwave/repository"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupDepartmentTest() (*handlers.DepartmentHandler, *mux.Router) {
	repo := repository.NewDepartmentRepository(nil) // No Redis for tests
	handler := &handlers.DepartmentHandler{Repo: repo}
	router := mux.NewRouter()

	router.HandleFunc("/api/departments", handler.CreateDepartment).Methods("POST")
	router.HandleFunc("/api/departments", handler.GetAllDepartments).Methods("GET")
	router.HandleFunc("/api/departments/{id}", handler.GetDepartment).Methods("GET")
	router.HandleFunc("/api/departments/{id}", handler.UpdateDepartment).Methods("PUT")
	router.HandleFunc("/api/departments/{id}", handler.DeleteDepartment).Methods("DELETE")

	return handler, router
}

func TestCreateDepartment_Success(t *testing.T) {
	_, router := setupDepartmentTest()

	department := NewValidDepartment()
	department.ID = 0 // Clear ID for creation

	body, _ := json.Marshal(department)
	req := httptest.NewRequest("POST", "/api/departments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var created models.Department
	json.NewDecoder(w.Body).Decode(&created)
	assert.NotZero(t, created.ID)
	assert.Equal(t, department.DepartmentName, created.DepartmentName)
	assert.Equal(t, department.DepartmentCode, created.DepartmentCode)
	assert.NotZero(t, created.CreatedAt)
	assert.NotZero(t, created.UpdatedAt)
}

func TestCreateDepartment_WithoutHeadName(t *testing.T) {
	_, router := setupDepartmentTest()

	department := NewDepartmentWithoutHeadName()
	department.ID = 0

	body, _ := json.Marshal(department)
	req := httptest.NewRequest("POST", "/api/departments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var created models.Department
	json.NewDecoder(w.Body).Decode(&created)
	assert.Empty(t, created.DepartmentHead)
}

func TestCreateDepartment_EmptyName(t *testing.T) {
	_, router := setupDepartmentTest()

	department := NewDepartmentWithEmptyName()
	department.ID = 0

	body, _ := json.Marshal(department)
	req := httptest.NewRequest("POST", "/api/departments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "name is required")
}

func TestCreateDepartment_EmptyCode(t *testing.T) {
	_, router := setupDepartmentTest()

	department := NewDepartmentWithEmptyCode()
	department.ID = 0

	body, _ := json.Marshal(department)
	req := httptest.NewRequest("POST", "/api/departments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "code is required")
}

func TestCreateDepartment_DuplicateCode(t *testing.T) {
	_, router := setupDepartmentTest()

	department := NewValidDepartment()
	department.ID = 0

	// Create first department
	body, _ := json.Marshal(department)
	req := httptest.NewRequest("POST", "/api/departments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Try to create duplicate
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/departments", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "already exists")
}

func TestGetDepartment_Success(t *testing.T) {
	handler, router := setupDepartmentTest()

	// Create department
	department := NewValidDepartment()
	department.ID = 0
	created, _ := handler.Repo.Create(department)

	req := httptest.NewRequest("GET", "/api/departments/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrieved models.Department
	json.NewDecoder(w.Body).Decode(&retrieved)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.DepartmentName, retrieved.DepartmentName)
	assert.Equal(t, created.DepartmentCode, retrieved.DepartmentCode)
}

func TestGetDepartment_NotFound(t *testing.T) {
	_, router := setupDepartmentTest()

	req := httptest.NewRequest("GET", "/api/departments/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAllDepartments_Empty(t *testing.T) {
	_, router := setupDepartmentTest()

	req := httptest.NewRequest("GET", "/api/departments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var departments []models.Department
	json.NewDecoder(w.Body).Decode(&departments)
	assert.Empty(t, departments)
}

func TestGetAllDepartments_Multiple(t *testing.T) {
	handler, router := setupDepartmentTest()

	// Create multiple departments
	departments := NewCommonDepartments()
	for _, dept := range departments {
		dept.ID = 0
		handler.Repo.Create(dept)
	}

	req := httptest.NewRequest("GET", "/api/departments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrieved []models.Department
	json.NewDecoder(w.Body).Decode(&retrieved)
	assert.Len(t, retrieved, len(departments))
}

func TestGetDepartmentByCode_Success(t *testing.T) {
	handler, router := setupDepartmentTest()

	// Create department
	department := NewDepartmentWithCode("MATH")
	department.ID = 0
	handler.Repo.Create(department)

	req := httptest.NewRequest("GET", "/api/departments?code=MATH", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrieved models.Department
	json.NewDecoder(w.Body).Decode(&retrieved)
	assert.Equal(t, "MATH", retrieved.DepartmentCode)
}

func TestGetDepartmentByCode_NotFound(t *testing.T) {
	_, router := setupDepartmentTest()

	req := httptest.NewRequest("GET", "/api/departments?code=NOTFOUND", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateDepartment_Success(t *testing.T) {
	handler, router := setupDepartmentTest()

	// Create department
	department := NewValidDepartment()
	department.ID = 0
	handler.Repo.Create(department)

	// Update
	updates := map[string]interface{}{
		"name":      "Updated Computer Science",
		"head_name": "Dr. Jane Doe",
	}
	body, _ := json.Marshal(updates)
	req := httptest.NewRequest("PUT", "/api/departments/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updated models.Department
	json.NewDecoder(w.Body).Decode(&updated)
	assert.Equal(t, "Updated Computer Science", updated.DepartmentName)
	assert.Equal(t, "Dr. Jane Doe", updated.DepartmentHead)
}

func TestUpdateDepartment_NotFound(t *testing.T) {
	_, router := setupDepartmentTest()

	updates := map[string]interface{}{
		"name": "Updated Name",
	}
	body, _ := json.Marshal(updates)
	req := httptest.NewRequest("PUT", "/api/departments/999", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateDepartment_DuplicateCode(t *testing.T) {
	handler, router := setupDepartmentTest()

	// Create two departments
	dept1 := NewDepartmentWithCode("CS")
	dept1.ID = 0
	handler.Repo.Create(dept1)

	dept2 := NewDepartmentWithCode("MATH")
	dept2.ID = 0
	handler.Repo.Create(dept2)

	// Try to update dept2 with dept1's code
	updates := map[string]interface{}{
		"code": "CS",
	}
	body, _ := json.Marshal(updates)
	req := httptest.NewRequest("PUT", "/api/departments/2", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "already exists")
}

func TestDeleteDepartment_Success(t *testing.T) {
	handler, router := setupDepartmentTest()

	// Create department
	department := NewValidDepartment()
	department.ID = 0
	handler.Repo.Create(department)

	req := httptest.NewRequest("DELETE", "/api/departments/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify deletion
	_, exists := handler.Repo.Get(1)
	assert.False(t, exists)
}

func TestDeleteDepartment_NotFound(t *testing.T) {
	_, router := setupDepartmentTest()

	req := httptest.NewRequest("DELETE", "/api/departments/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDepartmentValidation_ShortName(t *testing.T) {
	_, router := setupDepartmentTest()

	department := NewDepartmentWithShortName()
	department.ID = 0

	body, _ := json.Marshal(department)
	req := httptest.NewRequest("POST", "/api/departments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "between 2 and 100 characters")
}

func TestDepartmentValidation_LongName(t *testing.T) {
	_, router := setupDepartmentTest()

	department := NewDepartmentWithLongName()
	department.ID = 0

	body, _ := json.Marshal(department)
	req := httptest.NewRequest("POST", "/api/departments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "between 2 and 100 characters")
}

func TestDepartmentValidation_ShortCode(t *testing.T) {
	_, router := setupDepartmentTest()

	department := NewDepartmentWithShortCode()
	department.ID = 0

	body, _ := json.Marshal(department)
	req := httptest.NewRequest("POST", "/api/departments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "between 2 and 10 characters")
}

func TestDepartmentValidation_LongCode(t *testing.T) {
	_, router := setupDepartmentTest()

	department := NewDepartmentWithLongCode()
	department.ID = 0

	body, _ := json.Marshal(department)
	req := httptest.NewRequest("POST", "/api/departments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "uppercase letters")
}
