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

	router.HandleFunc("/admin", AdminPanelHandler).Methods("GET")
	router.HandleFunc("/admin/add-teacher", getAddingTeacherPage).Methods("GET")
	router.HandleFunc("/admin/add-teacher", addTeacherHandler).Methods("POST")
    router.HandleFunc("/admin/add-student", getAddingStudentPage).Methods("GET")
    router.HandleFunc("/admin/add-student", addStudentHandler).Methods("POST")
}
