package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/api"
	"github.com/Lacky1234union/OsintTeleBot/internal/app/repositories"
	"github.com/Lacky1234union/OsintTeleBot/internal/app/services"
	"github.com/Lacky1234union/OsintTeleBot/pkg/database"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Get database configuration from environment variables with defaults
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "host.docker.internal"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "appdb"),
	}

	// Инициализация БД
	db, err := database.NewPostgresDB(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация слоев
	personRepo := repositories.NewPersonRepository(db)
	personService := services.NewPersonService(personRepo)
	personAPI := api.NewPersonAPI(personService)

	// Настройка HTTP маршрутов
	http.HandleFunc("/api/person", personAPI.RegisterUserHandler)
	http.HandleFunc("/api/person/email", personAPI.FindUserByEmailHandler)
	http.HandleFunc("/api/person/name", personAPI.FindUserByNameHandler)
	http.HandleFunc("/api/person/phone", personAPI.FindUserByPhoneHandler)

	// Запуск HTTP сервера
	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	log.Printf("Server started on :8080 with database connection to %s:%s", dbConfig.Host, dbConfig.Port)

	// Ожидание сигнала завершения
	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
