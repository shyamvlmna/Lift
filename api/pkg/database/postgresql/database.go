package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// dbHost := os.Getenv("dbHost")
	dbHost := "sculift_pg"
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)
	Db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return Db, nil
}
