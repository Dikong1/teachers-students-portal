package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
)

func setupRoutes(router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/vol/{id}", volunteerPersonalPageHandler)
	router.HandleFunc("/chil/{id}", childPersonalPageHandler)
	router.HandleFunc("/vollogin", volLoginHandler)
	router.HandleFunc("/volreg", volRegHandler)
	router.HandleFunc("/chilog", chiLogHandler)
	router.HandleFunc("/chireg", chiRegHandler)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})
	router.HandleFunc("/update-wishes", updateWishesHandler).Methods("POST")
}
