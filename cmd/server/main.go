// Создать базовый сервер на на порту 8080
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MaxGot69/url-shortener/internal/config"
	"github.com/MaxGot69/url-shortener/internal/server"
	"github.com/MaxGot69/url-shortener/models"
	"github.com/MaxGot69/url-shortener/pkg/database"
)

func main() {
	// Тут конфиг регистрирую
	cfg := config.LoadConfig()
	router := server.NewRouter()

	fmt.Println("Starting server on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))

	db, err := database.PostgresConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.URL{}); err != nil {
		log.Fatal("Failed to ran migratioins:", err)
	}

}
