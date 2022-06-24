package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
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

func UserData(db *gorm.DB, table string) *gorm.DB {
	user := &models.User{}
	err := Db.AutoMigrate(&user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("user db opened")
	return Db

}
func DriverData(db *gorm.DB) *gorm.DB {
	driver := &models.Driver{}
	err := Db.AutoMigrate(&driver)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("driver db opened")
	return Db
}
func TripData(db *gorm.DB, table string) *gorm.DB {
	driver := &models.Trip{}
	err := Db.AutoMigrate(&driver)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("driver db opened")
	return Db
}

// func VehicleData(db *gorm.DB, table string) *gorm.DB {

// }
