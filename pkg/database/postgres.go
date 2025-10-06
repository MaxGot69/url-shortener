package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresConnection() (*gorm.DB, error) {
	dsn := "host=localhost user=maxim password=secret dbname=urlshortener port=5432 sslmode=disable TimeZone=Asia/Almaty"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
