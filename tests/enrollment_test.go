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

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// getTestRedisClient returns a Redis client configured for testing
// Uses database 15 to avoid conflicts with production data
func getTestRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use a separate DB for testing
	})
}

// cleanupTestRedis clears all data from the test Redis database
func cleanupTestRedis(client *redis.Client) {
	client.FlushDB(client.Context())
}

func TestCreateEnrollment(t *testing.T) {
	client := getTestRedisClient()
	defer cleanupTestRedis(client)
	
	repo := repository.NewEnrollmentRepository(client)
	handler := &handlers.EnrollmentHandler{Repo: repo}

	enrollment := `{"student_id":"student-123","course_id":"course-456","status":"pending"}`
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
	if e.StudentID != "student-123" {
		t.Errorf("expected student_id student-123, got %s", e.StudentID)
	}
	if e.Status != models.StatusPending {
		t.Errorf("expected status pending, got %s", e.Status)
	}
	if e.ID == "" {
		t.Error("expected ID to be generated")
	}
}

func TestCreateEnrollmentValidation(t *testing.T) {
	client := getTestRedisClient()
	defer cleanupTestRedis(client)
	
	repo := repository.NewEnrollmentRepository(client)
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Test missing student_id
	enrollment := `{"course_id":"course-456","status":"pending"}`
	req := httptest.NewRequest("POST", "/api/enrollments", strings.NewReader(enrollment))
	w := httptest.NewRecorder()

	handler.CreateEnrollment(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for missing student_id, got %v", res.StatusCode)
	}

	// Test invalid status
	enrollment = `{"student_id":"student-123","course_id":"course-456","status":"invalid"}`
	req = httptest.NewRequest("POST", "/api/enrollments", strings.NewReader(enrollment))
	w = httptest.NewRecorder()

	handler.CreateEnrollment(w, req)
	res = w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid status, got %v", res.StatusCode)
	}
}

func TestGetEnrollment(t *testing.T) {
	client := getTestRedisClient()
	defer cleanupTestRedis(client)
	
	repo := repository.NewEnrollmentRepository(client)
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create an enrollment first
	enrollment := models.Enrollment{
		StudentID: "student-123",
		CourseID:  "course-456",
		Status:    models.StatusActive,
	}
	created, err := repo.Create(&enrollment)
	if err != nil {
		t.Fatal(err)
	}

	// Get the enrollment
	req := httptest.NewRequest("GET", "/api/enrollments/"+created.ID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": created.ID})
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
		t.Errorf("expected ID %s, got %s", created.ID, e.ID)
	}
}

func TestGetEnrollmentNotFound(t *testing.T) {
	client := getTestRedisClient()
	defer cleanupTestRedis(client)
	
	repo := repository.NewEnrollmentRepository(client)
	handler := &handlers.EnrollmentHandler{Repo: repo}

	req := httptest.NewRequest("GET", "/api/enrollments/nonexistent-id", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "nonexistent-id"})
	w := httptest.NewRecorder()

	handler.GetEnrollment(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %v", res.StatusCode)
	}
}

func TestGetAllEnrollments(t *testing.T) {
	client := getTestRedisClient()
	defer cleanupTestRedis(client)
	
	repo := repository.NewEnrollmentRepository(client)
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create test enrollments
	enrollments := []models.Enrollment{
		{StudentID: "student-1", CourseID: "course-1", Status: models.StatusActive},
		{StudentID: "student-2", CourseID: "course-2", Status: models.StatusPending},
	}

	for _, e := range enrollments {
		if _, err := repo.Create(&e); err != nil {
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

	var result []*models.Enrollment
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 enrollments, got %d", len(result))
	}
}

func TestGetEnrollmentsByStudentID(t *testing.T) {
	client := getTestRedisClient()
	defer cleanupTestRedis(client)
	
	repo := repository.NewEnrollmentRepository(client)
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create enrollments for different students
	enrollments := []models.Enrollment{
		{StudentID: "student-1", CourseID: "course-1", Status: models.StatusActive},
		{StudentID: "student-1", CourseID: "course-2", Status: models.StatusPending},
		{StudentID: "student-2", CourseID: "course-1", Status: models.StatusActive},
	}

	for _, e := range enrollments {
		if _, err := repo.Create(&e); err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest("GET", "/api/enrollments?student_id=student-1", nil)
	w := httptest.NewRecorder()

	handler.GetAllEnrollments(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", res.StatusCode)
	}

	var result []*models.Enrollment
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 enrollments for student-1, got %d", len(result))
	}
}

func TestUpdateEnrollment(t *testing.T) {
	client := getTestRedisClient()
	defer cleanupTestRedis(client)
	
	repo := repository.NewEnrollmentRepository(client)
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create an enrollment
	enrollment := models.Enrollment{
		StudentID: "student-123",
		CourseID:  "course-456",
		Status:    models.StatusPending,
	}
	created, err := repo.Create(&enrollment)
	if err != nil {
		t.Fatal(err)
	}

	// Update to active status
	update := `{"status":"active"}`
	req := httptest.NewRequest("PUT", "/api/enrollments/"+created.ID, strings.NewReader(update))
	req = mux.SetURLVars(req, map[string]string{"id": created.ID})
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
}

func TestUpdateEnrollmentInvalidTransition(t *testing.T) {
	client := getTestRedisClient()
	defer cleanupTestRedis(client)
	
	repo := repository.NewEnrollmentRepository(client)
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create a completed enrollment
	enrollment := models.Enrollment{
		StudentID: "student-123",
		CourseID:  "course-456",
		Status:    models.StatusCompleted,
	}
	created, err := repo.Create(&enrollment)
	if err != nil {
		t.Fatal(err)
	}

	// Try to update to active (invalid transition from completed)
	update := `{"status":"active"}`
	req := httptest.NewRequest("PUT", "/api/enrollments/"+created.ID, strings.NewReader(update))
	req = mux.SetURLVars(req, map[string]string{"id": created.ID})
	w := httptest.NewRecorder()

	handler.UpdateEnrollment(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid transition, got %v", res.StatusCode)
	}
}

func TestDeleteEnrollment(t *testing.T) {
	client := getTestRedisClient()
	defer cleanupTestRedis(client)
	
	repo := repository.NewEnrollmentRepository(client)
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create an enrollment
	enrollment := models.Enrollment{
		StudentID: "student-123",
		CourseID:  "course-456",
		Status:    models.StatusActive,
	}
	created, err := repo.Create(&enrollment)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the enrollment
	req := httptest.NewRequest("DELETE", "/api/enrollments/"+created.ID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": created.ID})
	w := httptest.NewRecorder()

	handler.DeleteEnrollment(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %v", res.StatusCode)
	}

	// Verify it's deleted
	deleted, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if deleted != nil {
		t.Error("expected enrollment to be deleted")
	}
}

func TestGetEnrollmentStats(t *testing.T) {
	client := getTestRedisClient()
	defer cleanupTestRedis(client)
	
	repo := repository.NewEnrollmentRepository(client)
	handler := &handlers.EnrollmentHandler{Repo: repo}

	// Create enrollments with different statuses
	enrollments := []models.Enrollment{
		{StudentID: "student-1", CourseID: "course-1", Status: models.StatusPending},
		{StudentID: "student-2", CourseID: "course-2", Status: models.StatusPending},
		{StudentID: "student-3", CourseID: "course-3", Status: models.StatusActive},
		{StudentID: "student-4", CourseID: "course-4", Status: models.StatusCompleted},
	}

	for _, e := range enrollments {
		if _, err := repo.Create(&e); err != nil {
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
