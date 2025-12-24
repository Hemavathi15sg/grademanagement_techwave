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
)

func TestCreateStudent(t *testing.T) {
    repo := repository.NewStudentRepository()
    handler := &handlers.StudentHandler{Repo: repo}

    student := `{"name":"Alice","email":"alice@example.com","grade":"A"}`
    req := httptest.NewRequest("POST", "/students", strings.NewReader(student))
    w := httptest.NewRecorder()

    handler.CreateStudent(w, req)
    res := w.Result()

    if res.StatusCode != http.StatusCreated {
        t.Fatalf("expected status 201, got %v", res.StatusCode)
    }

    var s models.Student
    if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
        t.Fatal(err)
    }
    if s.Name != "Alice" {
        t.Errorf("expected name Alice, got %s", s.Name)
    }
}