package driver

import (
	database "github.com/shayamvlmna/cab-booking-app/pkg/database/postgresql"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
)

//accepts druver models and pass to the user database to insert
//retun error if any
func AddDriver(newDriver *models.Driver) error {
	return database.InsertDriver(newDriver)
}

//returns a driver model by accepting a key and a value
//eg:if searching using id, key is "id" and value is the id of the driver to search
func GetDriver(key, value string) models.Driver {
	driver, _ := database.FindDriver(key, value)
	return driver
}

//return all drivers in the database
func GetDrivers() []models.Driver {

	return *database.GetDrivers()
}

//update the driver by accepting the updated driver fields
//only update fields with null values
func UpdateDriver(driver *models.Driver) {
	database.UpdateDriver(driver)
}

//delete driver from the database by the id
func DeleteDriver(id string) {

}

//return boolean to check if the driver exist or not
func IsDriverExists(key, value string) bool {
	_, err := database.FindDriver(key, value)
	return err
}
