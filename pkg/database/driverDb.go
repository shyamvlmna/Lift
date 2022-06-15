package database

import (
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"gorm.io/gorm"
)

func OpenDriverDb() (*gorm.DB, error) {
	Db, err := openDB()
	if err != nil {
		return nil, err
	}
	driver := &models.Driver{}
	err = Db.AutoMigrate(&driver)
	if err != nil {
		return nil, err
	}
	return Db, nil
}

func AddDriver(){

}
func FindDriver(){

}
func UpdateDriver() {

}