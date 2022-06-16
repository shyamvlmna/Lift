package driver

import (
	"github.com/shayamvlmna/cab-booking-app/app/database"
	"github.com/shayamvlmna/cab-booking-app/app/models"
)

func AddDriver(newDriver *models.Driver) error {
	return database.InsertDriver(newDriver)
}

func GetDriver(key string) models.Driver {
	driver, _ := database.FindDriver(key)
	return driver
}
func IsDriverExists(key string) bool {
	_, err := database.FindDriver(key)
	return err
}
