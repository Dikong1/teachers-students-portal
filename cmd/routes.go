package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
)

func setupRoutes(router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/teachlog", teachLoginHandler)
	router.HandleFunc("/teachreg", teachRegHandler)
	router.HandleFunc("/studlog", studLogHandler)
	router.HandleFunc("/studreg", studRegHandler)
	router.HandleFunc("/logout", logoutHandler).Methods("GET")

	router.HandleFunc("/teach/{id}", verifyToken(teachPersonalPageHandler)).Methods("GET")
	router.HandleFunc("/stud/{id}", verifyToken(studPersonalPageHandler)).Methods("GET")
	router.HandleFunc("/api/courses", verifyToken((getDataFromDatabase))).Methods("GET")

	router.HandleFunc("/verify", verifyHandler).Methods("GET")
	router.HandleFunc("/verification-failure", verifyFailureHandler).Methods("GET")

	router.HandleFunc("/admin", AdminPanelHandler).Methods("GET")
	router.HandleFunc("/admin/add-course", AdminAddCourseHandler).Methods("POST")
	router.HandleFunc("/admin/teachers", getTeacherPage).Methods("GET")
	router.HandleFunc("/admin/add-teacher", addTeacherHandler).Methods("POST")
	router.HandleFunc("/admin/delete-teacher", deleteTeacherHandler).Methods("POST")
	router.HandleFunc("/admin/students", getStudentPage).Methods("GET")
	router.HandleFunc("/admin/add-student", addStudentHandler).Methods("POST")
	router.HandleFunc("/admin/delete-student", deleteStudentHandler).Methods("POST")
}
