package main

import (
	"log"
	"net/http"
	"os"

	"auth-service/internal/app"
	"auth-service/internal/domain"
	"auth-service/internal/interfaces/http"
)

func main() {
	// Get JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Initialize services
	tokenService := domain.NewJWTService(jwtSecret)
	userService := app.NewUserService(nil) // TODO: Implement user repository
	authHandler := http.NewAuthHandler(userService, tokenService)

	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/auth/validate", authHandler.ValidateToken)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting auth service on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
