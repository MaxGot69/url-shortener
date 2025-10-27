package service

import (
	"crypto/rand"
	"math/big"
	"net/url"

	"github.com/MaxGot69/url-shortener/internal/models"
	"github.com/MaxGot69/url-shortener/internal/repository"
)

type URLServiceInterface interface {
	GenerateRandomShortString(n int) (string, error)
	IsValidURL(u string) bool
	SaveURL(url models.URL) error
}

type URLService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) GenerateRandomShortString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}

func (s *URLService) IsValidURL(u string) bool {
	parseURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	if parseURL.Scheme == "" || parseURL.Host == "" {
		return false
	}
	return true
}

func (s *URLService) SaveURL(url models.URL) error {
	return s.repo.SaveURL(url)
}