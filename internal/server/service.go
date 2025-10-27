package server

import (
	"net/http"

	"github.com/MaxGot69/url-shortener/internal/config"
	"github.com/MaxGot69/url-shortener/internal/handler"
	myMiddleware "github.com/MaxGot69/url-shortener/internal/middleware"
	"github.com/MaxGot69/url-shortener/internal/repository"
	"github.com/MaxGot69/url-shortener/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(urlService *service.URLService, userRepo repository.UserRepository, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Max"))
	})

	redirectHandler := handler.NewRedirectHandler(userRepo)
	r.Handle("/{shortCode}", redirectHandler)
	
	r.Get("/health", http.HandlerFunc(healthHandler))
	
	shortenHandler := handler.NewPostHandler(urlService)
	r.Handle("/api/v1/shorten", shortenHandler)
	
	registerHandler := handler.NewRegisterHandler(userRepo, cfg)
	r.Handle("/api/v1/register", registerHandler)
	
	loginHandler := handler.NewLoginHandler(userRepo, cfg)
	r.Handle("/api/v1/login", loginHandler)
	
	authMw := myMiddleware.NewAuthMiddleware(cfg.JWTSecret)
	statsHandler := handler.NewStatisticHandler(userRepo)
	r.With(authMw.Middleware).Handle("/api/v1/stats/{shortCode}", statsHandler)

	r.Handle("/metrics", promhttp.Handler())

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}