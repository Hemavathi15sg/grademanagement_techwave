package routes

import (
    "github.com/gorilla/mux"
    "techwave/handlers"
    "techwave/repository"
)

func RegisterRoutes(r *mux.Router) {
    studentRepo := repository.NewStudentRepository()
    studentHandler := &handlers.StudentHandler{Repo: studentRepo}

    r.HandleFunc("/students", studentHandler.CreateStudent).Methods("POST")
    r.HandleFunc("/students", studentHandler.ListStudents).Methods("GET")
    r.HandleFunc("/students/{id}", studentHandler.GetStudent).Methods("GET")
    r.HandleFunc("/students/{id}", studentHandler.UpdateStudent).Methods("PUT")
    r.HandleFunc("/students/{id}", studentHandler.DeleteStudent).Methods("DELETE")

    gradeRepo := repository.NewGradeRepository()
    gradeHandler := &handlers.GradeHandler{Repo: gradeRepo}

    r.HandleFunc("/grades", gradeHandler.CreateGrade).Methods("POST")
    r.HandleFunc("/grades", gradeHandler.ListGrades).Methods("GET")
    r.HandleFunc("/grades/{id}", gradeHandler.GetGrade).Methods("GET")
    r.HandleFunc("/grades/{id}", gradeHandler.UpdateGrade).Methods("PUT")
    r.HandleFunc("/grades/{id}", gradeHandler.DeleteGrade).Methods("DELETE")
}