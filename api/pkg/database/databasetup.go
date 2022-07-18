package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBSet() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// dbHost := os.Getenv("DB_HOST")
	dbHost := "sculift_pg"
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)
	Db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		fmt.Println("failed to connect to Postgresql")
		return nil
	}
	fmt.Println("successfully connected to Postgresql")

	return Db
}

var Db *gorm.DB = DBSet()
