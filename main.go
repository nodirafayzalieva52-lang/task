package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"project/handlers"
	"project/internal/repository"
	"project/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		"localhost", 5432, "user123", "123", "postgres"))
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", userHandler.Register)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
