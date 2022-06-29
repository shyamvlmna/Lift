package driver

import (
	"encoding/json"

	"github.com/shayamvlmna/cab-booking-app/pkg/database/redis"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
)

var d = &models.Driver{}

//return boolean to check if the driver exist or not
func IsDriverExists(key, value string) bool {
	_, err := d.Get(key, value)
	return err
}

//accepts druver models and pass to the user database to insert
//retun error if any
func AddDriver(newDriver *models.Driver) error {
	return newDriver.Add()
}

//returns a driver model by accepting a key and a value
//eg:if searching using id, key is "id" and value is the id of the driver to search
func GetDriver(key, value string) models.Driver {

	p, err := redis.GetData("data")
	if err != nil {
		driver, _ := d.Get(key, value)
		return driver
	}

	driver := models.Driver{}

	json.Unmarshal([]byte(p), &driver)

	return driver
}

//return all drivers in the database
func GetAllDrivers() []models.Driver {

	return *d.GetAll()
}

//update the driver by accepting the updated driver fields
//only update fields with null values
func UpdateDriver(driver *models.Driver) {
	driver.Update()
}

//delete driver from the database by the id
func DeleteDriver(id uint64) error {
	return d.Delete(id)
}

func ApproveDriver(id uint64) error {
	return d.BlockUnblock(id)
}

func BlockDriver(id uint64) error {
	return d.BlockUnblock(id)
}
func UnBlockDriver(id uint64) error {
	return d.BlockUnblock(id)
}
