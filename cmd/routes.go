package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
)

func setupRoutes(router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/vol/{id}", teachPersonalPageHandler)
	router.HandleFunc("/chil/{id}", studPersonalPageHandler)
	router.HandleFunc("/vollogin", teachLoginHandler)
	router.HandleFunc("/volreg", teachRegHandler)
	router.HandleFunc("/chilog", studLogHandler)
	router.HandleFunc("/chireg", studRegHandler)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})
}
