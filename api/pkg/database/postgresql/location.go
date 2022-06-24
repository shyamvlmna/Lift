package database

import (
	"fmt"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
	"gorm.io/gorm"
)

func OpenLocationDb() (*gorm.DB, error) {
	Db, err := openDB()
	if err != nil {
		return nil, err
	}
	location := &models.Location{}
	err = Db.AutoMigrate(&location)
	if err != nil {
		return nil, err
	}
	fmt.Println("trip db opened")
	return Db, nil
}
