package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaxGot69/url-shortener/internal/metrics"
	"github.com/MaxGot69/url-shortener/internal/models"
	"github.com/MaxGot69/url-shortener/internal/service"
	appLogger "github.com/MaxGot69/url-shortener/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func init() {
	appLogger.Init("error")
	metrics.Init()
}

type MockURLRepository struct{}

func (m *MockURLRepository) SaveURL(url models.URL) error {
	return nil
}

func (m *MockURLRepository) GetURLByShortCode(shortCode string) (*models.URL, error) {
	return &models.URL{ShortCode: shortCode, OriginalURL: "https://example.com", ClickCount: 0}, nil
}

func (m *MockURLRepository) UpdateClickCount(url models.URL) error {
	return nil
}

func TestNewPostHandler(t *testing.T) {
	mockRepo := &MockURLRepository{}
	urlService := service.NewURLService(mockRepo)
	
	handler := NewPostHandler(urlService)
	
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.urlService)
}

func TestPostHandler_Success(t *testing.T) {
	mockRepo := &MockURLRepository{}
	urlService := service.NewURLService(mockRepo)
	handler := NewPostHandler(urlService)

	reqBody := URLRequest{OriginalURL: "https://example.com"}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	
	var response ShortenResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NotEmpty(t, response.ShortURL)
}

func TestPostHandler_InvalidURL(t *testing.T) {
	mockRepo := &MockURLRepository{}
	urlService := service.NewURLService(mockRepo)
	handler := NewPostHandler(urlService)

	reqBody := URLRequest{OriginalURL: "not-a-valid-url"}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPostHandler_InvalidJSON(t *testing.T) {
	mockRepo := &MockURLRepository{}
	urlService := service.NewURLService(mockRepo)
	handler := NewPostHandler(urlService)

	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer([]byte("invalid json")))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPostHandler_InvalidMethod(t *testing.T) {
	mockRepo := &MockURLRepository{}
	urlService := service.NewURLService(mockRepo)
	handler := NewPostHandler(urlService)

	req, _ := http.NewRequest("GET", "/shorten", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}
