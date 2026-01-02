package routes

import (
	"techwave/handlers"
	"techwave/repository"

	"github.com/gorilla/mux"
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

	// Enrollment routes with in-memory repository
	enrollmentRepo := repository.NewEnrollmentRepository()
	enrollmentHandler := &handlers.EnrollmentHandler{Repo: enrollmentRepo}

	// Use /api prefix for enrollment routes as specified in requirements
	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/enrollments", enrollmentHandler.CreateEnrollment).Methods("POST")
	apiRouter.HandleFunc("/enrollments", enrollmentHandler.GetAllEnrollments).Methods("GET")
	apiRouter.HandleFunc("/enrollments/stats", enrollmentHandler.GetEnrollmentStats).Methods("GET")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.GetEnrollment).Methods("GET")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.UpdateEnrollment).Methods("PUT")
	apiRouter.HandleFunc("/enrollments/{id}", enrollmentHandler.DeleteEnrollment).Methods("DELETE")

	// Department routes with Redis caching support
	departmentRepo := repository.NewDepartmentRepository(nil) // Redis client can be passed here
	departmentHandler := &handlers.DepartmentHandler{Repo: departmentRepo}

	apiRouter.HandleFunc("/departments", departmentHandler.CreateDepartment).Methods("POST")
	apiRouter.HandleFunc("/departments", departmentHandler.GetAllDepartments).Methods("GET")
	apiRouter.HandleFunc("/departments/{id}", departmentHandler.GetDepartment).Methods("GET")
	apiRouter.HandleFunc("/departments/{id}", departmentHandler.UpdateDepartment).Methods("PUT")
	apiRouter.HandleFunc("/departments/{id}", departmentHandler.DeleteDepartment).Methods("DELETE")
}
