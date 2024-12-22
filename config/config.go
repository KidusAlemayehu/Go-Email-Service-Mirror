package config

import (
	"email-service/migrate"
	"email-service/utils/log"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadENV() {
	err := godotenv.Load()
	if err != nil {
		log.Logger.Error("Error loading .env file")
	}
}

func InitDB() (*gorm.DB, error) {
	// Connect to your database here
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	log.Logger.Info("Successfully connected to the database!")

	migrate.ModelsAutoMigrate(db)

	return db, nil
}
