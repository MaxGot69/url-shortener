// Создать структуру для конфигурвции
package config

import (
	"os"
)

type Config struct {
	Port string
}

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	return Config{Port: port}
}
