package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openDB() (*gorm.DB, error) {
	godotenv.Load()

	dbHost := os.Getenv("dbHost")
	dbPort := os.Getenv("dbPort")
	dbUser := os.Getenv("dbUser")
	dbPassword := os.Getenv("dbPassword")
	dbName := os.Getenv("dbName")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)
	Db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	

	return Db, nil
}
