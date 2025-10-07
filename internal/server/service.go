package server

import (
	"net/http"

	"github.com/MaxGot69/url-shortener/internal/handler"
	"github.com/MaxGot69/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(urlService *service.URLService) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger) //автоматически логирует информацию о каждом HTTP-запросе
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Max"))
	})

	r.Get("/{shortCode}", handler.RedirectHandler) // Обработчик GET запросов для редиректа (из handler/redirect.go)
	r.Get("/health", healthHeandler)
	r.Route("/api/v1/", func(r chi.Router) {
		r.Post("/shorten", handler.PostHandler(urlService)) // Обработчик POST запросов для сокращения url (из handler/shorten.go)
	})

	return r

}
func healthHeandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
