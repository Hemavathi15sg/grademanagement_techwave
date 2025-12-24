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

func TestCreateGrade(t *testing.T) {
    repo := repository.NewGradeRepository()
    handler := &handlers.GradeHandler{Repo: repo}

    grade := `{"student_id":1,"value":"B+","subject":"Science"}`
    req := httptest.NewRequest("POST", "/grades", strings.NewReader(grade))
    w := httptest.NewRecorder()

    handler.CreateGrade(w, req)
    res := w.Result()

    if res.StatusCode != http.StatusCreated {
        t.Fatalf("expected status 201, got %v", res.StatusCode)
    }

    var g models.Grade
    if err := json.NewDecoder(res.Body).Decode(&g); err != nil {
        t.Fatal(err)
    }
    if g.Subject != "Science" {
        t.Errorf("expected subject Science, got %s", g.Subject)
    }
}