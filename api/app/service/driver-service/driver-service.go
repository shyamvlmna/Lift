package driver

import (
	"github.com/shayamvlmna/cab-booking-app/app/database"
	"github.com/shayamvlmna/cab-booking-app/app/models"
)

func InsertDriver(newDriver *models.Driver) error {

	database.AddDriver()
	return nil
}

func GetDriver() error {
	database.FindDriver()
	return nil
}
