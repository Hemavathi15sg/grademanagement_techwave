package routes

import (
	"techwave/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	// No repositories - handlers do everything (messy!)
	studentHandler := &handlers.StudentHandler{}

	r.HandleFunc("/students", studentHandler.CreateStudent).Methods("POST")
	r.HandleFunc("/students", studentHandler.ListStudents).Methods("GET")
	r.HandleFunc("/students/{id}", studentHandler.GetStudent).Methods("GET")
	r.HandleFunc("/students/{id}", studentHandler.UpdateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", studentHandler.DeleteStudent).Methods("DELETE")

	gradeHandler := &handlers.GradeHandler{}

	r.HandleFunc("/grades", gradeHandler.CreateGrade).Methods("POST")
	r.HandleFunc("/grades", gradeHandler.ListGrades).Methods("GET")
	r.HandleFunc("/grades/{id}", gradeHandler.GetGrade).Methods("GET")
	r.HandleFunc("/grades/{id}", gradeHandler.UpdateGrade).Methods("PUT")
	r.HandleFunc("/grades/{id}", gradeHandler.DeleteGrade).Methods("DELETE")
}
