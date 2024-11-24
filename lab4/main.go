package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"lab4/handlers"
	"lab4/middleware"
	"lab4/monitoring"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// Create new router
	r := mux.NewRouter()

	// Security Middleware (CSRF Protection, Security Headers)
	r.Use(middleware.CSRFMiddleware)
	r.Use(middleware.SecurityHeadersMiddleware)

	// Routes
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/protected", handlers.ProtectedHandler).Methods("GET")

	// JWT Authentication Middleware for protected routes
	r.Use(middleware.JWTMiddleware)

	// Logging Middleware
	r.Use(middleware.LoggingMiddleware(logger))

	// Monitoring (Prometheus) Middleware
	r.Handle("/metrics", monitoring.MetricsHandler())

	useHTTPS := false
	if useHTTPS {
		// Start HTTPS server with TLS (for HTTPS)
		log.Println("Starting server with HTTPS on port 8443")
		srv := &http.Server{
			Handler:      r,
			Addr:         ":8443",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
	} else {
		// Start HTTP server (default)
		log.Println("Starting server with HTTP on port 8080")
		srv := &http.Server{
			Handler:      r,
			Addr:         ":8080",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.Fatal(srv.ListenAndServe()) // Default to HTTP
	}

	log.Println("Starting server on port 8000")
}
