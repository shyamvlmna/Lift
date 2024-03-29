package driver

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/shayamvlmna/lift/pkg/service/auth"

	"github.com/shayamvlmna/lift/pkg/models"
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
func GetAllDrivers() ([]models.Driver, error) {

	drivers, err := d.GetAll()

	if err != nil {
		return nil, err
	}
	return *drivers, nil
}

func DriverRequests() (*[]models.Driver, error) {
	return models.DriverRequests()
}

func Payout(amount string, driverId uint) error {

	return models.AddPayoutRequest(amount, driverId)

}

func PayoutRequests(driverid uint) ([]models.PayoutResponse, error) {
	return models.GetPayoutStatus(driverid)
}

func PayoutHistory(id uint) []models.PayoutResponse {
	return models.PayoutHistory(id)
}

func ApproveDriver(id uint) error {
	return d.ApproveToDrive(id)
}

func BlockDriver(id uint) error {
	return d.BlockUnblock(id)
}
func UnBlockDriver(id uint) error {
	return d.BlockUnblock(id)
}

func GetBankDetails(id uint) (*models.Bank, error) {
	return models.GetBankDetails(id)
}

func UpdateBankDetails(id uint, bank *models.Bank) error {
	return bank.UpdateBank(id, bank)
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

//delete driver from the database by the id
func DeleteDriver(id uint64) error {
	return d.Delete(id)
}
