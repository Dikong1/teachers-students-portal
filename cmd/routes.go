package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
)

func setupRoutes(router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/teach/{id}", teachPersonalPageHandler)
	router.HandleFunc("/stud/{id}", studPersonalPageHandler)
	router.HandleFunc("/teachlog", teachLoginHandler)
	router.HandleFunc("/teachreg", teachRegHandler)
	router.HandleFunc("/studlog", studLogHandler)
	router.HandleFunc("/studreg", studRegHandler)
	router.HandleFunc("/api/courses", getDataFromDatabase)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})
}
