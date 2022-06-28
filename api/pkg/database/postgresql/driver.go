package database

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/shayamvlmna/cab-booking-app/pkg/database"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
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
	fmt.Println("driver db opened")
	return Db, nil
}

// func closeDriverdb(db *gorm.DB) {

// 	sqlDb, err := db.DB()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	sqlDb.Close()
// 	fmt.Println("driver db closed")
// }

//insert driver model to the driver database
//return error if any
func InsertDriver(driver *models.Driver) error {
	db := database.DriverData(database.Db)

	result := db.Create(&driver)

	return result.Error
}

func FindDriver(key, value string) (models.Driver, bool) {

	db := database.DriverData(database.Db)

	driver := &models.Driver{}
	result := db.Where(key+"=?", value).First(&driver)

	if result.Error == gorm.ErrRecordNotFound {
		return *driver, false
	} else {
		return *driver, true
	}
}

//get and return all drivers from the driver database
func GetDrivers() *[]models.Driver {
	db := database.DriverData(database.Db)

	drivers := &[]models.Driver{}
	db.Find(&drivers)

	return drivers
}

//update a driver by getting updated driver fields
//only update the not null driver fields
func UpdateDriver(updatedDriver *models.Driver) error {
	db := database.DriverData(database.Db)

	driver := &models.Driver{}

	id := strconv.Itoa(int(updatedDriver.ID))

	db.Where("id=?", id).First(&driver)

	result := db.Model(&driver).Updates(models.Driver{

		FirstName:   updatedDriver.FirstName,
		LastName:    updatedDriver.LastName,
		PhoneNumber: updatedDriver.PhoneNumber,
		Email:       updatedDriver.Email,
		Password:    updatedDriver.Password,
		City:        updatedDriver.City,
		Active:      false,
		Cab:         models.Vehicle{},
	})
	return result.Error
}

func ActiveStatus(id uint64) error {
	db := database.DriverData(database.Db)

	driver := &models.Driver{}

	db.Where("DriverId=?", id).First(&driver)

	if !driver.Approved {
		return errors.New("accesDenied")
	}

	if driver.Active {
		driver.Active = false
		result := db.Save(&driver)
		return result.Error
	}

	driver.Active = true
	result := db.Save(&driver)
	return result.Error
}

func ApproveDriver(id uint64) error {
	db := database.DriverData(database.Db)

	driver := &models.Driver{}

	db.Where("id=?", id).First(&driver)

	if !driver.Approved {
		driver.Approved = true
		result := db.Save(&driver)
		return result.Error
	}
	driver.Approved = false
	result := db.Save(&driver)
	return result.Error
}

//delete driver by id
//returns err if any
func DeleteDriver(id uint64) error {
	db := database.DriverData(database.Db)

	driver := &models.Driver{}

	result := db.Delete(&driver, id)

	return result.Error
}
