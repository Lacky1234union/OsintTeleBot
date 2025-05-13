package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/repositories"
	"github.com/Lacky1234union/OsintTeleBot/internal/app/services"
	"github.com/Lacky1234union/OsintTeleBot/pkg/database"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Инициализация БД
	db, err := database.NewPostgresDB(ctx, database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		DBName:   "appdb",
	})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Инициализация слоев
	userRepo := repositories.NewRepository(db)
	PersonService := services.NewPersonService(userRepo)

	// Далее инициализация транспорта (HTTP, gRPC и т.д.)
	// ...
}
