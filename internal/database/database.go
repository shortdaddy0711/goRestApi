package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Failed to load env file")
	}

	return os.Getenv(key)
}

func NewDatabase() (*gorm.DB, error) {
	log.Info("Setting up new database connection")

	dbUsername := goDotEnvVariable("DB_USERNAME")
	dbPassword := goDotEnvVariable("DB_PASSWORD")
	dbHost := goDotEnvVariable("DB_HOST")
	dbTable := goDotEnvVariable("DB_TABLE")
	dbPort := goDotEnvVariable("DB_PORT")
	sslmode := goDotEnvVariable("SSL_MODE")

	connectString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbHost, dbPort, dbUsername, dbTable, dbPassword, sslmode)

	db, err := gorm.Open("postgres", connectString)
	if err != nil {
		return db, err
	}

	if err := db.DB().Ping(); err != nil {
		return db, err
	}

	return db, nil
}
