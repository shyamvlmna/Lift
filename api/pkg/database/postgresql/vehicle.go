package database

import (
	"fmt"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"gorm.io/gorm"
)

func OpenVehicleDb() *gorm.DB {
	db, err := openDB()
	if err != nil {
		fmt.Println(err)
	}
	vehicle := models.Vehicle{}
	db.AutoMigrate(&vehicle)
	return db
}

func Insert(v *models.Vehicle) {
	db := OpenVehicleDb()
	db.Create(&v)
}
