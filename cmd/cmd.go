package cmd

import (
	"Platform/db"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	err1 := godotenv.Load(".env")
	if err1 != nil {
		fmt.Println("Error loading .env file:", err)
	}

	router := mux.NewRouter()
	router.Use(RateLimitMiddleware)
	setupRoutes(router)

	srv := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log_.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log_.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log_.Fatalf("Server forced to shutdown: %v", err)
	}
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
