package driver

import (
	"github.com/shayamvlmna/cab-booking-app/app/database"
	"github.com/shayamvlmna/cab-booking-app/app/models"
)

func AddDriver(newDriver *models.Driver) error {
	return database.InsertDriver(newDriver)
}

func GetDriver(key string) error {
	database.FindDriver(key)
	return nil
}
func IsDriverExists(key string) bool {
	return database.CheckDriver(key)
}
