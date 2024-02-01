package cmd

import (
	"SantaWeb/db"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RunServer() {
	err := db.DbConnection()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	router := mux.NewRouter()
	setupRoutes(router)

	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
