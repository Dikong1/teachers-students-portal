package cmd

import (
	"Platform/db"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

var log_ = logrus.New()

func RunServer() {
	err := db.DbConnection()
	if err != nil {
		log_.WithFields(logrus.Fields{
			"action": "database_connection",
			"status": "failed",
		}).Fatal("Database connection failed: ", err)
	}

	router := mux.NewRouter()

	router.Use(RateLimitMiddleware)

	setupRoutes(router)

	port := ":3000"
    log_.WithFields(logrus.Fields{
        "action": "server_start",
        "status": "running",
        "port":   port,
    }).Info("Starting server...")

    log_.Fatal(http.ListenAndServe(port, router))
}

var limiter = rate.NewLimiter(1, 3) // 1 request per second, burst of 3

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
