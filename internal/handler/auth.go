package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MaxGot69/url-shortener/internal/models"
	"github.com/MaxGot69/url-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

var secretKey = []byte("your-secret-key")
var userRepo repository.UserReposytory
var repo repository.PostgresRepository

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Декодировать JSON из тела запроса
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Printf("Email: %s, Password: %s\n", req.Email, req.Password)

	// 2. Проверить есть ли пользователь в БД
	checkUser, err := userRepo.FindByEmail(req.Email)
	if err == nil && checkUser != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// 3. Захешировать пароль с bcrypt
	passwordBytes := []byte(req.Password)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Hashing password: %s\n", string(hash))

	// 4. Сохранить пользователя в БД
	newUser := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}
	err = userRepo.CreateUser(newUser)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// 5. Вернуть HTTP 201 Created
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Декодировать JSON
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Найти пользователя в БД по email
	checkUser, err := userRepo.FindByEmail(req.Email)
	if err != nil {
		http.Error(w, "DataBase error", http.StatusInternalServerError)
		return
	}
	if checkUser == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// 3. Сравнить пароль с bcrypt.CompareHashAndPassword()
	err = bcrypt.CompareHashAndPassword([]byte(checkUser.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
	} else {
		fmt.Printf("Susses password")

		//4.Генерация токена
		claims := jwt.MapClaims{
			"ID":    checkUser.ID,
			"Email": checkUser.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(), // Срок действия — 24 часа
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString(secretKey)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		fmt.Println("JWT:", token)
		response := struct {
			Token string `json:"token"`
		}{
			Token: signedToken,
		}
		// Отправляем JSON-ответ клиенту
		json.NewEncoder(w).Encode(response)
	}
}

// Статистика по ссылкам с jwt middleware
func StatisticHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode") // Получаем короткую из url
	// поиcк в бд
	url, err := repo.GetURLByShortCode(shortCode)

	// JSON
	if err != nil {
		http.Error(w, "URL not Found", http.StatusBadRequest)
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
	json.NewEncoder(w).Encode(response)

}
