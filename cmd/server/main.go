// Создать базовый сервер на на порту 8080
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MaxGot69/url-shortener/internal/config"
	"github.com/MaxGot69/url-shortener/internal/models"
	"github.com/MaxGot69/url-shortener/internal/repository"
	"github.com/MaxGot69/url-shortener/internal/server"
	"github.com/MaxGot69/url-shortener/internal/service"
	"github.com/MaxGot69/url-shortener/pkg/database"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.PostgresConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	repo := &repository.PostgresRepository{DB: db}
	urlService := service.NewURLService(repo) //  создаём сервис

	router := server.NewRouter(urlService) // передаём его в NewRouter

	fmt.Println("Starting server on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))

	if err := db.AutoMigrate(&models.URL{}); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
}
