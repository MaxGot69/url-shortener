package service

import (
	"crypto/rand"
	"math/big"
	"net/url"
)

// Генерация случайной строки для короткой ссылки
func GenerateRandomShortString(n int) (string, error) {
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

// функция для валидации url
func IsValidURL(u string) bool {
	parseURL, err := url.Parse(u) // парсим строку в url
	if err != nil {
		return false // ошибка - неваллидный url
	}
	if parseURL.Scheme == "" || parseURL.Host == "" {
		return false // неполный или некорректный url
	}
	return true // вадидный
}
