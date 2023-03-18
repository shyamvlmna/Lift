package mysql

import (
	"fmt"
	"github.com/shayamvlmna/lift/pkg/models"

	"gorm.io/gorm"
)

func OpenTripDb() (*gorm.DB, error) {
	Db, err := openDB()
	if err != nil {
		return nil, err
	}
	trip := &models.Trip{}
	err = Db.AutoMigrate(&trip)
	if err != nil {
		return nil, err
	}
	fmt.Println("trip db opened")
	return Db, nil
}

func GetTrips(role string, id uint) *[]models.Trip {

	db, _ := OpenTripDb()
	rides := []models.Trip{}
	db.Where(role, id).Find(&rides)

	return &rides
}
