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
	cache2 "project/pkg/cache"
	smtp2 "project/pkg/smtp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		"localhost", 5432, "postgres", "20102010", "mydb"))
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	cache := cache2.NewMemoryCache()
	smtp := smtp2.NewSMTP("smtp.gmail.com", "587", "rfsu vjub fnmm ditg", "golang.tester1974@gmail.com")

	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository, cache, smtp)
	userHandler := handler.NewUserHandler(userService)

	middleware := handler.New(userRepository)

	mux := http.NewServeMux()
	mux.Handle("POST /register", http.HandlerFunc(userHandler.Register))
	mux.Handle("POST /verify", http.HandlerFunc(userHandler.Verify))
	mux.Handle("POST /login", http.HandlerFunc(userHandler.Login))
	mux.Handle("POST /get/me", middleware.Auth(http.HandlerFunc(userHandler.GetMe)))
	mux.Handle("PUT /update/me", middleware.Auth(http.HandlerFunc(userHandler.UpdateMe)))
	mux.Handle("DELETE /delete/me", middleware.Auth(http.HandlerFunc(userHandler.DeleteMe)))
	mux.Handle("POST /orders", middleware.Auth(http.HandlerFunc(userHandler.CreateOrder)))
	mux.Handle("GET /orders", middleware.Auth(http.HandlerFunc(userHandler.GetMyOrders)))
	mux.Handle("GET /orders/{id}", middleware.Auth(http.HandlerFunc(userHandler.GetOrderByID)))
	mux.Handle("PUT /orders/{id}", middleware.Auth(http.HandlerFunc(userHandler.UpdateOrder)))
	mux.Handle("DELETE /orders/{id}", middleware.Auth(http.HandlerFunc(userHandler.DeleteOrder)))
	mux.Handle("GET /admin/users", middleware.AuthAdmin(http.HandlerFunc(userHandler.AdminGetAllUsers)))
	mux.Handle("GET /admin/orders", middleware.AuthAdmin(http.HandlerFunc(userHandler.AdminGetAllOrders)))
	mux.Handle("PATCH /admin/users/{id}/role", middleware.AuthAdmin(http.HandlerFunc(userHandler.AdminUpdateUserRole)))
	mux.Handle("PUT /users/change-password", middleware.Auth(http.HandlerFunc(userHandler.UpdateUserPassword)))
	mux.Handle("PATCH /orders/{id}/cancel", middleware.Auth(http.HandlerFunc(userHandler.CancelOrder)))
	mux.Handle("GET /users/profile", middleware.Auth(http.HandlerFunc(userHandler.GetUserAndOrders)))
	mux.Handle("POST /auth/refresh", http.HandlerFunc(userHandler.Refresh))
	mux.Handle("POST /auth/logout", http.HandlerFunc(userHandler.Logout))
	
	
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
