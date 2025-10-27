package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MaxGot69/url-shortener/internal/config"
	"github.com/MaxGot69/url-shortener/internal/models"
	"github.com/MaxGot69/url-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	appLogger "github.com/MaxGot69/url-shortener/pkg/logger"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

type RegisterHandler struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewRegisterHandler(repo repository.UserRepository, cfg *config.Config) *RegisterHandler {
	return &RegisterHandler{
		repo:      repo,
		jwtSecret: cfg.JWTSecret,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appLogger.Logger.Error("Failed to decode register request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	appLogger.Logger.Info("Registration attempt", "email", req.Email)

	checkUser, err := h.repo.FindByEmail(req.Email)
	if err == nil && checkUser != nil {
		appLogger.Logger.Warn("User already exists", "email", req.Email)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	passwordBytes := []byte(req.Password)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		appLogger.Logger.Error("Failed to hash password", "error", err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	newUser := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	err = h.repo.CreateUser(newUser)
	if err != nil {
		appLogger.Logger.Error("Failed to create user", "error", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})

	appLogger.Logger.Info("User registered", "email", req.Email, "user_id", newUser.ID)
}

type LoginHandler struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewLoginHandler(repo repository.UserRepository, cfg *config.Config) *LoginHandler {
	return &LoginHandler{
		repo:      repo,
		jwtSecret: cfg.JWTSecret,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appLogger.Logger.Error("Failed to decode login request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	appLogger.Logger.Info("Login attempt", "email", req.Email)

	user, err := h.repo.FindByEmail(req.Email)
	if err != nil {
		appLogger.Logger.Error("Database error during login", "error", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		appLogger.Logger.Warn("User not found", "email", req.Email)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		appLogger.Logger.Warn("Invalid password", "email", req.Email)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"ID":    user.ID,
		"Email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		appLogger.Logger.Error("Failed to generate token", "error", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: signedToken,
	}

	json.NewEncoder(w).Encode(response)
	appLogger.Logger.Info("Login successful", "email", req.Email, "user_id", user.ID)
}

type StatisticHandler struct {
	repo repository.UserRepository
}

func NewStatisticHandler(repo repository.UserRepository) *StatisticHandler {
	return &StatisticHandler{
		repo: repo,
	}
}

func (h *StatisticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")

	appLogger.Logger.Info("Statistics request", "short_code", shortCode)

	urlRepo, ok := h.repo.(repository.URLRepository)
	if !ok {
		appLogger.Logger.Error("Repository does not implement URLRepository")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	url, err := urlRepo.GetURLByShortCode(shortCode)
	if err != nil {
		appLogger.Logger.Warn("URL not found for statistics", "short_code", shortCode, "error", err)
		http.Error(w, "URL not Found", http.StatusNotFound)
		return
	}

	response := struct {
		ShortCode   string `json:"short_code"`
		OriginalURL string `json:"original_url"`
		ClickCount  int    `json:"clicks"`
	}{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
		ClickCount:  url.ClickCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	appLogger.Logger.Info("Statistics retrieved", "short_code", shortCode)
}