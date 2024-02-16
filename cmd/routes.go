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

	router.HandleFunc("/teach/{id}", verifyToken(teachPersonalPageHandler)).Methods("GET")
	router.HandleFunc("/stud/{id}", verifyToken(studPersonalPageHandler)).Methods("GET")
	router.HandleFunc("/api/courses", verifyToken((getDataFromDatabase))).Methods("GET")

	router.HandleFunc("/add-student", addStudentHandler).Methods("POST")
	router.HandleFunc("/delete-student", deleteStudentHandler).Methods("POST")
}
