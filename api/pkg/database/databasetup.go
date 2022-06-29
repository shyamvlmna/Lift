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

	dbHost := os.Getenv("dbHost")
	dbPort := os.Getenv("dbPort")
	dbUser := os.Getenv("dbUser")
	dbPassword := os.Getenv("dbPassword")
	dbName := os.Getenv("dbName")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)
	Db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		fmt.Println("failed to connect to Postgresql")
		return nil
	}
	fmt.Println("successfully connected to Postgresql")
	return Db
}

var Db *gorm.DB = DBSet()
