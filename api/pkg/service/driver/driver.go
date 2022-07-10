package driver

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/shayamvlmna/cab-booking-app/pkg/service/auth"

	"github.com/shayamvlmna/cab-booking-app/pkg/models"
)

var d = &models.Driver{}

func RegisterDriver(newDriver *models.Driver) error {
	newDriver.PhoneNumber = auth.GetPhone()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newDriver.Password), bcrypt.DefaultCost)
	newDriver.Password = string(hashPassword)

	if err := AddDriver(newDriver); err != nil {
		return err
	}
	auth.StorePhone(newDriver.PhoneNumber)
	return nil
}

// IsDriverExists return boolean to check if the driver exist or not
func IsDriverExists(key, value string) bool {
	_, err := d.Get(key, value)
	return err
}

// AddDriver accepts druver models and pass to the user database to insert
//retun error if any
func AddDriver(newDriver *models.Driver) error {
	return newDriver.Add()
}

// GetDriver returns a driver model by accepting a key and a value
//eg:if searching using id, key is "id" and value is the id of the driver to search
func GetDriver(key, value string) *models.Driver {

	//p, err := redis.GetData("data")
	//if err != nil {
	driver, _ := d.Get(key, value)
	return &driver
	//}

	//driver := &models.Driver{}
	//
	//err = json.Unmarshal([]byte(p), &driver)
	//if err != nil {
	//	return nil
	//}
	//
	//return driver
}

// GetAllDrivers return all drivers in the database
func GetAllDrivers() []models.Driver {

	return *d.GetAll()
}

func RegisterToDrive() {

}

// UpdateDriver update the driver by accepting the updated driver fields
//only update fields with null values
func UpdateDriver(driver models.Driver) {
	err := d.Update(driver)
	if err != nil {
		return
	}
}

func Payout(amount string, driverId uint) error {

	err := models.AddPayout(amount, driverId)
	if err != nil {
		return err
	}
	return nil
}

func PayoutRequests(driverid uint) *models.Payout {
	return models.GetPayoutStatus(driverid)
}

//delete driver from the database by the id
func DeleteDriver(id uint64) error {
	return d.Delete(id)
}

func ApproveDriver(id uint) error {
	return d.BlockUnblock(id)
}

func BlockDriver(id uint) error {
	return d.BlockUnblock(id)
}
func UnBlockDriver(id uint) error {
	return d.BlockUnblock(id)
}
