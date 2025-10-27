package main

import (
	"log"
	"net/http"

	"github.com/MaxGot69/url-shortener/internal/config"
	"github.com/MaxGot69/url-shortener/internal/metrics"
	"github.com/MaxGot69/url-shortener/internal/models"
	"github.com/MaxGot69/url-shortener/internal/repository"
	"github.com/MaxGot69/url-shortener/internal/server"
	"github.com/MaxGot69/url-shortener/internal/service"
	"github.com/MaxGot69/url-shortener/pkg/database"
	appLogger "github.com/MaxGot69/url-shortener/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	
	appLogger.Init(cfg.LogLevel)
	appLogger.Logger.Info("Initializing application")
	
	metrics.Init()

	db, err := database.PostgresConnection(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		appLogger.Logger.Error("Failed to connect to database", "error", err)
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&models.URL{}, &models.User{})
	if err != nil {
		appLogger.Logger.Error("Failed to run migrations", "error", err)
		log.Fatal("Failed to run migrations:", err)
	}

	repo := &repository.PostgresRepository{DB: db}
	urlService := service.NewURLService(repo)
	
	router := server.NewRouter(urlService, repo, &cfg)

	appLogger.Logger.Info("Starting server", "port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
