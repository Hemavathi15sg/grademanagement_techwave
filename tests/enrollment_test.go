package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"techwave/handlers"
	"techwave/models"
	"techwave/repository"

	"github.com/gorilla/mux"
)

func TestCreateEnrollment(t *testing.T) {
	repo := repository.NewEnrollmentRepository()
	handler := &handlers.EnrollmentHandler{Repo: repo}

	enrollment := `{"student_id":123,"course_id":456,"status":"pending"}`
	req := httptest.NewRequest("POST", "/api/enrollments", strings.NewReader(enrollment))
	w := httptest.NewRecorder()

	handler.CreateEnrollment(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %v", res.StatusCode)
	}

	var e models.Enrollment
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		t.Fatal(err)
	}
	if e.StudentID != 123 {
		t.Errorf("expected student_id 123, got %d", e.StudentID)
	}
	if e.Status != models.StatusPending {
		t.Errorf("expected status pending, got %s", e.Status)
	}
	if e.ID == 0 {
		t.Error("expected ID to be generated")
	}
}

func TestCreateEnrollmentValidation(t *testing.T) {
	repo := repository.NewEnrollmentRepository()
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Test missing student_id
	enrollment := `{"course_id":456,"status":"pending"}`
	req := httptest.NewRequest("POST", "/api/enrollments", strings.NewReader(enrollment))
	w := httptest.NewRecorder()

	handler.CreateEnrollment(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for missing student_id, got %v", res.StatusCode)
	}

	// Test invalid status
	enrollment = `{"student_id":123,"course_id":456,"status":"invalid"}`
	req = httptest.NewRequest("POST", "/api/enrollments", strings.NewReader(enrollment))
	w = httptest.NewRecorder()

	handler.CreateEnrollment(w, req)
	res = w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid status, got %v", res.StatusCode)
	}
}

func TestGetEnrollment(t *testing.T) {
	repo := repository.NewEnrollmentRepository()
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create an enrollment first
	enrollment := models.Enrollment{
		StudentID: 123,
		CourseID:  456,
		Status:    models.StatusActive,
	}
	created, err := repo.Create(enrollment)
	if err != nil {
		t.Fatal(err)
	}

	// Get the enrollment
	req := httptest.NewRequest("GET", "/api/enrollments/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetEnrollment(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", res.StatusCode)
	}

	var e models.Enrollment
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		t.Fatal(err)
	}
	if e.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, e.ID)
	}
}

func TestGetEnrollmentNotFound(t *testing.T) {
	repo := repository.NewEnrollmentRepository()
	handler := &handlers.EnrollmentHandler{Repo: repo}

	req := httptest.NewRequest("GET", "/api/enrollments/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.GetEnrollment(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %v", res.StatusCode)
	}
}

func TestGetAllEnrollments(t *testing.T) {
	repo := repository.NewEnrollmentRepository()
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create test enrollments
	enrollments := []models.Enrollment{
		{StudentID: 1, CourseID: 1, Status: models.StatusActive},
		{StudentID: 2, CourseID: 2, Status: models.StatusPending},
	}

	for _, e := range enrollments {
		if _, err := repo.Create(e); err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest("GET", "/api/enrollments", nil)
	w := httptest.NewRecorder()

	handler.GetAllEnrollments(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", res.StatusCode)
	}

	var result []models.Enrollment
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 enrollments, got %d", len(result))
	}
}

func TestGetEnrollmentsByStudentID(t *testing.T) {
	repo := repository.NewEnrollmentRepository()
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create enrollments for different students
	enrollments := []models.Enrollment{
		{StudentID: 1, CourseID: 1, Status: models.StatusActive},
		{StudentID: 1, CourseID: 2, Status: models.StatusPending},
		{StudentID: 2, CourseID: 1, Status: models.StatusActive},
	}

	for _, e := range enrollments {
		if _, err := repo.Create(e); err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest("GET", "/api/enrollments?student_id=1", nil)
	w := httptest.NewRecorder()

	handler.GetAllEnrollments(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", res.StatusCode)
	}

	var result []models.Enrollment
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 enrollments for student 1, got %d", len(result))
	}
}

func TestUpdateEnrollment(t *testing.T) {
	repo := repository.NewEnrollmentRepository()
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create an enrollment
	enrollment := models.Enrollment{
		StudentID: 123,
		CourseID:  456,
		Status:    models.StatusPending,
	}
	created, err := repo.Create(enrollment)
	if err != nil {
		t.Fatal(err)
	}

	// Update to active status
	update := `{"status":"active"}`
	req := httptest.NewRequest("PUT", "/api/enrollments/1", strings.NewReader(update))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateEnrollment(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", res.StatusCode)
	}

	var e models.Enrollment
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		t.Fatal(err)
	}
	if e.Status != models.StatusActive {
		t.Errorf("expected status active, got %s", e.Status)
	}
	if e.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, e.ID)
	}
}

func TestUpdateEnrollmentInvalidTransition(t *testing.T) {
	repo := repository.NewEnrollmentRepository()
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create a completed enrollment
	enrollment := models.Enrollment{
		StudentID: 123,
		CourseID:  456,
		Status:    models.StatusCompleted,
	}
	_, err := repo.Create(enrollment)
	if err != nil {
		t.Fatal(err)
	}

	// Try to update to active (invalid transition from completed)
	update := `{"status":"active"}`
	req := httptest.NewRequest("PUT", "/api/enrollments/1", strings.NewReader(update))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateEnrollment(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid transition, got %v", res.StatusCode)
	}
}

func TestDeleteEnrollment(t *testing.T) {
	repo := repository.NewEnrollmentRepository()
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create an enrollment
	enrollment := models.Enrollment{
		StudentID: 123,
		CourseID:  456,
		Status:    models.StatusActive,
	}
	created, err := repo.Create(enrollment)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the enrollment
	req := httptest.NewRequest("DELETE", "/api/enrollments/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteEnrollment(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %v", res.StatusCode)
	}

	// Verify it's deleted
	_, ok := repo.Get(created.ID)
	if ok {
		t.Error("expected enrollment to be deleted")
	}
}

func TestGetEnrollmentStats(t *testing.T) {
	repo := repository.NewEnrollmentRepository()
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create enrollments with different statuses
	enrollments := []models.Enrollment{
		{StudentID: 1, CourseID: 1, Status: models.StatusPending},
		{StudentID: 2, CourseID: 2, Status: models.StatusPending},
		{StudentID: 3, CourseID: 3, Status: models.StatusActive},
		{StudentID: 4, CourseID: 4, Status: models.StatusCompleted},
	}

	for _, e := range enrollments {
		if _, err := repo.Create(e); err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest("GET", "/api/enrollments/stats", nil)
	w := httptest.NewRecorder()

	handler.GetEnrollmentStats(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", res.StatusCode)
	}

	var stats repository.EnrollmentStats
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		t.Fatal(err)
	}

	if stats.TotalEnrollments != 4 {
		t.Errorf("expected 4 total enrollments, got %d", stats.TotalEnrollments)
	}
	if stats.ByStatus["pending"] != 2 {
		t.Errorf("expected 2 pending enrollments, got %d", stats.ByStatus["pending"])
	}
	if stats.ByStatus["active"] != 1 {
		t.Errorf("expected 1 active enrollment, got %d", stats.ByStatus["active"])
	}
	if stats.ByStatus["completed"] != 1 {
		t.Errorf("expected 1 completed enrollment, got %d", stats.ByStatus["completed"])
	}
}
